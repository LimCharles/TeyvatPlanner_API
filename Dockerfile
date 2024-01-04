FROM golang:latest as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gqlgen-server .

FROM alpine:latest  
WORKDIR /root/
COPY --from=builder /app/gqlgen-server .
CMD ["./gqlgen-server"]
