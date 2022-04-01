FROM golang:1.18 as builder
WORKDIR /build
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o urlshort cmd/main.go

# alpine does not work with CGO required by Sqlite
FROM ubuntu:latest
COPY --from=builder /build/urlshort .
EXPOSE 8080
EXPOSE 50051
RUN ls
ENTRYPOINT ["./urlshort"]