all: gen

.PHONY : gen

gen:
	protoc iam.proto --go_out=./iam_pb --go-grpc_out=./iam_pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative
	protoc filer.proto --go_out=./filer_pb --go-grpc_out=./filer_pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative
	protoc master.proto --go_out=./master_pb --go-grpc_out=./master_pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative