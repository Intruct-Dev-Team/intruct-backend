GOPATH := $(shell go env GOPATH)
export PATH := $(PATH):$(GOPATH)/bin

.PHONY: generate-swagger
generate-swagger:
	@echo "Generating Swagger docs..."
	@swag init -g ./cmd/main/main.go -o ./_docs/swagger

# migrate create -ext sql -dir migrations -seq create_teachers_table
