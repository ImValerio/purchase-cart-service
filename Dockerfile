FROM golang:1.19 as builder

WORKDIR /mnt

COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o main .

FROM debian:bullseye

WORKDIR /root/

COPY --from=builder /mnt/ .
EXPOSE 9090

CMD ["./main"]