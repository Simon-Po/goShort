# build
FROM golang:1.24.6-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./src .
ENV CGO_ENABLED=0
RUN go build -o /app/main

# run
FROM alpine:latest
ARG INSTALL_SQLITE=false
WORKDIR /app

RUN if [ "$INSTALL_SQLITE" = "true" ] ; then apk add --no-cache sqlite ; fi

RUN mkdir -p /app/data
COPY --from=builder /app/main .
COPY ./site ./site
EXPOSE 8000
ENTRYPOINT ["./main"]