FROM golang:1.24.6-alpine3.22 AS builder

LABEL org.opencontainers.image.description="scs-user service"


# Set up the correct module path for Go
WORKDIR /app/scs-user

# COPY .env .env

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN apk add --no-cache --update gcc g++

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o /build/scs-user ./cmd/server

FROM alpine:3.22

RUN apk add -U tzdata
ENV TZ=Asia/Singapore
RUN cp /usr/share/zoneinfo/Asia/Singapore /etc/localtime

WORKDIR /app

COPY --from=builder /build/scs-user .

CMD ["/app/scs-user"]
