FROM golang:1.21 as builder

WORKDIR /mnt

COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o main .
EXPOSE 9090

CMD ["./main"]