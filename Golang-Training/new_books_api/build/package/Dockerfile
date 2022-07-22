FROM golang:1.17-alpine AS builder

# Install all required dependencies for building
RUN apk update
RUN apk add git make build-base

WORKDIR /service

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make compile

FROM alpine:3.14.2
COPY --from=builder ./service/build/out/new_books_api new_books_api
COPY --from=builder ./service/configs/default.yaml default.yaml
ENTRYPOINT ["./new_books_api", "-configFile", "default.yaml"]