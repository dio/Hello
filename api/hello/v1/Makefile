generate:
	protoc -I/usr/local/include -I. \
		-I$(GOPATH)/src \
		-I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway \
		-I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=paths=source_relative,plugins=grpc:. \
		--grpc-gateway_out=paths=source_relative:. \
		--swagger_out=json_names_for_fields=true,allow_merge=true,merge_file_name=core,fqn_for_swagger_name=true:. \
		hello.proto
