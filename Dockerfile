FROM golang:1.20 as builder

WORKDIR /go/src/app
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 make dependabotnotifier

RUN make test

FROM gcr.io/distroless/static-debian11

COPY --from=builder /go/src/app/bin/dependabotnotifier /dependabotnotifier
CMD ["/dependabotnotifier"]
