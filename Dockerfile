FROM golang:1.23.3

WORKDIR /merch-store
COPY . .

RUN go build -o /build ./internal/cmd


EXPOSE 8080

CMD ["/build"]