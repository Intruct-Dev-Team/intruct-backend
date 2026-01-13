package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"go.uber.org/zap"
)

const (
	maxLogLenValue   = 100
	maxBodySizeToLog = 10 * 1024
	filteredValue    = "[FILTERED]"
	truncatedValue   = "[TRUNCATED]"
	fileTooLarge     = "[FILE_TOO_LARGE]"
	binaryContent    = "[BINARY_CONTENT]"
)

// list of sensitive fields, which should be masked
var sensitiveFields = []string{
	"password",
	"token",
	"jwt",
	"authorization",
	"api_key",
	"access_token",
	"refresh_token",
	"credit_card",
	"card_number",
}

type responseWriter struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		body:           bytes.NewBuffer([]byte{}),
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(body []byte) (int, error) {
	rw.body.Write(body)
	return rw.ResponseWriter.Write(body)
}

// isMultipartFormData checks if content type is multipart/form-data
func isMultipartFormData(contentType string) bool {
	return strings.HasPrefix(strings.ToLower(contentType), "multipart/form-data")
}

// isBinaryContent checks if content appears to be binary
func isBinaryContent(data []byte) bool {
	// Check first 512 bytes for non-printable characters
	checkLen := len(data)
	if checkLen > 512 {
		checkLen = 512
	}

	for i := 0; i < checkLen; i++ {
		b := data[i]
		// If byte is not printable ASCII and not common whitespace
		if b < 0x20 && b != '\n' && b != '\r' && b != '\t' {
			return true
		}
		if b > 0x7E {
			return true
		}
	}
	return false
}

// maskSensitiveAndBigData mask sensitive data in JSON
func maskSensitiveAndBigData(data []byte) string {
	// Check if data is binary
	if isBinaryContent(data) {
		return binaryContent
	}

	// try parse JSON
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		// if it is not JSON, just check on sensitive fields
		bodyStr := string(data)
		for _, field := range sensitiveFields {
			if strings.Contains(strings.ToLower(bodyStr), field) {
				return filteredValue
			}
		}

		// Truncate if too long
		if len(bodyStr) > maxLogLenValue {
			return bodyStr[:maxLogLenValue] + "..." + truncatedValue
		}
		return bodyStr
	}

	maskJSONFields(jsonMap)

	maskedJSON, err := json.Marshal(jsonMap)
	if err != nil {
		return filteredValue
	}
	return string(maskedJSON)
}

// maskJSONFields recursively masks sensitive and long fields in JSON structure
func maskJSONFields(data map[string]interface{}) {
	for key, value := range data {
		lowerKey := strings.ToLower(key)

		// check if field is sensitive
		for _, sensitive := range sensitiveFields {
			if strings.Contains(lowerKey, sensitive) {
				data[key] = filteredValue
				break
			}
		}

		// mask long strings
		if str, ok := value.(string); ok && len(str) > maxLogLenValue {
			data[key] = str[:maxLogLenValue] + "..." + truncatedValue
		}

		// recursively process nested objects
		switch v := value.(type) {
		case map[string]interface{}:
			maskJSONFields(v)
		case []interface{}:
			for i, item := range v {
				if mapItem, ok := item.(map[string]interface{}); ok {
					maskJSONFields(mapItem)
				} else if str, ok := item.(string); ok && len(str) > maxLogLenValue {
					v[i] = str[:maxLogLenValue] + "..." + truncatedValue
				}
			}
		}
	}
}

func LoggerMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// mask sensitive headers
			headers := make(map[string]string)
			for key, values := range r.Header {
				lowerKey := strings.ToLower(key)
				if slices.Contains(sensitiveFields, lowerKey) {
					headers[key] = filteredValue
				} else if len(strings.Join(values, ", ")) >= maxLogLenValue {
					headers[key] = truncatedValue
				} else {
					headers[key] = strings.Join(values, ", ")
				}
			}

			logger.Info("request started",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
			)

			// Check content type - skip body reading for multipart/form-data
			contentType := r.Header.Get("Content-Type")
			var requestBody []byte
			shouldLogBody := true

			if isMultipartFormData(contentType) {
				// Don't read multipart form data body - it contains files
				shouldLogBody = false
			} else if r.Body != nil && r.ContentLength <= maxBodySizeToLog {
				// Only read body if it's not too large
				requestBody, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
			} else if r.ContentLength > maxBodySizeToLog {
				shouldLogBody = false
			}

			rw := newResponseWriter(w)
			next.ServeHTTP(rw, r)

			duration := time.Since(start)

			// form log with masked data
			logFields := []zap.Field{
				zap.Duration("duration", duration),
				zap.Int("status", rw.status),
			}

			// log request body / response body only in case of error
			if rw.status >= 400 {
				if shouldLogBody && len(requestBody) > 0 {
					maskedRequest := maskSensitiveAndBigData(requestBody)
					logFields = append(logFields,
						zap.String("error_code", http.StatusText(rw.status)),
						zap.String("request_body", maskedRequest),
					)
				} else if isMultipartFormData(contentType) {
					logFields = append(logFields,
						zap.String("error_code", http.StatusText(rw.status)),
						zap.String("request_body", "[MULTIPART_FORM_DATA]"),
					)
				} else if r.ContentLength > maxBodySizeToLog {
					logFields = append(logFields,
						zap.String("error_code", http.StatusText(rw.status)),
						zap.String("request_body", fileTooLarge),
					)
				}

				// mask response body
				if rw.body.Len() > 0 && rw.body.Len() <= maxBodySizeToLog {
					maskedResponse := maskSensitiveAndBigData(rw.body.Bytes())
					logFields = append(logFields, zap.String("response_body", maskedResponse))
				}
			}

			logger.Info("request completed", logFields...)
		})
	}
}
