package httputils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const MaxSize = 3 << 20

func GetIntParamFromRequestPath(r *http.Request, paramName string) (int, error) {
	var number int

	param := r.PathValue(paramName)
	if param == "" {
		return 0, fmt.Errorf("missing {%s} param in url", paramName)
	}

	number, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}

	return number, nil
}

func DecodeImage(base64Img string) ([]byte, error) {
	imgBytes, err := base64.StdEncoding.DecodeString(base64Img)
	if err != nil {
		return nil, errors.New("invalid image format")
	}

	// check for weight
	if len(imgBytes) > MaxSize {
		return nil, errors.New("file too large")
	}

	// check is MIME-type (image)
	mimeType := http.DetectContentType(imgBytes)
	if !strings.HasPrefix(mimeType, "image/") {
		return nil, errors.New("file is not an image")
	}

	return imgBytes, nil
}
