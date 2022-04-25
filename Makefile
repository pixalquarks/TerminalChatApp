gen:
	protoc --go_out=. --go-grpc_out=. proto/*.proto

clean:
	rm pb/*.go

server:
	go run server.go args.go

goClient:
	go run ./Client

buildClient:
	go build -o ./build ./Client

buildServer:
	go build -o ./build server.go args.go