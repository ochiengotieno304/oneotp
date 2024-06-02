FROM golang:1.22.1-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o app cmd/main.go

FROM scratch

COPY --from=builder ["/build/app", "/"]

EXPOSE 6000

CMD ["./app"]