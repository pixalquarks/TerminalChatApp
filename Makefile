gen:
	protoc --go_out=. --go-grpc_out=. proto/*.proto

clean:
	rm pb/*.go

server:
	go run server.go

client:
	go run client.go helper.go

ui:
	go run ./Client

buildClient:
	go build -o ./build ./Client

buildServer:
	go build -o ./build server.go