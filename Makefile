all: api

api:
	$(MAKE) -C api/hello/v1

run: api
	go run cmd/gateway/main.go

.PHONY: api
