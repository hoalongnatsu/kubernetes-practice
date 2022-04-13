FROM golang:1.17-alpine AS build
RUN apk add --no-cache git
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o run .

FROM alpine:3.9
WORKDIR /app
RUN apk add ca-certificates
COPY --from=build /build/run .
CMD ["/app/run"]