FROM golang:1.23-alpine AS builder

WORKDIR /go/src/app
COPY . .

RUN go mod download

RUN go build -o ./bin/dependabotnotifier ./cmd/dependabotnotifier

FROM gcr.io/distroless/static-debian11

COPY --from=builder /go/src/app/bin/dependabotnotifier /dependabotnotifier
CMD ["/dependabotnotifier"]
