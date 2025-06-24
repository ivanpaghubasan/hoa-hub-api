FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o hoa-hub-api .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/hoa-hub-api .

EXPOSE 9000

CMD [ "./app" ]