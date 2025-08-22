FROM cgr.dev/chainguard/go:latest AS builder
ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /src
COPY go.* ./
RUN go mod download
COPY . .

RUN go test -v ./...
RUN go run honnef.co/go/tools/cmd/staticcheck@latest ./...
RUN go run golang.org/x/vuln/cmd/govulncheck@latest ./...
RUN go run golang.org/x/tools/cmd/deadcode@latest -test ./...
RUN go build -a -installsuffix cgo -o ./bin/dependabotnotifier ./cmd/dependabotnotifier

FROM cgr.dev/chainguard/static
WORKDIR /app
COPY --from=builder /src/bin/dependabotnotifier /app/dependabotnotifier
ENTRYPOINT ["/app/dependabotnotifier"]
