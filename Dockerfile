FROM golang:1.20 AS build
WORKDIR /app
ENV CGO_ENABLED=0
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags "-w -s" -o server cmd/main.go

FROM alpine:3.17 AS runtime
COPY --from=build /app ./
EXPOSE 8080/tcp
ENTRYPOINT ["/server"]
