FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY go.* ./
RUN go mod download
COPY . .

RUN go test -v ./...
RUN go run honnef.co/go/tools/cmd/staticcheck@latest ./...
RUN go run golang.org/x/vuln/cmd/govulncheck@latest ./...
RUN go run golang.org/x/tools/cmd/deadcode@latest -test ./...
RUN go build -o ./bin/dependabotnotifier ./cmd/dependabotnotifier

FROM gcr.io/distroless/static-debian11
WORKDIR /app
COPY --from=builder /src/bin/dependabotnotifier /app/dependabotnotifier
CMD ["/app/dependabotnotifier"]
