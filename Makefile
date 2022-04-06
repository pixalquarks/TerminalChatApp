gen:
	protoc --go_out=. --go-grpc_out=. proto/*.proto

gents:
	protoc -I=".\Client\src\proto" --ts_out=".\Client\src\proto" chat.proto

clean:
	rm pb/*.go

server:
	go run server.go

goClient:
	go run client.go helper.go