generate_protoc:
	protoc \
    --go_out=models \
    --go_opt=paths=source_relative \
    --go-grpc_out=models \
    --go-grpc_opt=paths=source_relative \
    **/*.proto