FROM golang:1.17-alpine3.15 as builder
WORKDIR  /home/email-sending-system
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o ./email-sending-system


FROM alpine:3
RUN apk add --update ca-certificates
RUN apk add --no-cache tzdata && \
    cp -f /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime && \
    apk del tzdata
WORKDIR /app
COPY testdata testdata
COPY --from=builder /home/email-sending-system .

ENTRYPOINT ["./email-sending-system"]
