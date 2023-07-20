dependabotnotifier:
	go build -o bin/dependabotnotifier cmd/dependabotnotifier/*.go

test:
	go test ./...


