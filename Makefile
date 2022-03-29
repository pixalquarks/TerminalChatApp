gen:
	protoc --go_out=. --go-grpc_out=. proto/*.proto

clean:
	rm pb/*.go

runserver:
	go run server.go

runclient:
	go run client.go