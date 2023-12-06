FROM golang:1.21.4-alpine AS builder

WORKDIR /app

# 보안관제 서버에서는 local.env 대신 production.env 사용
ENV DOT_ENV=local.env
ENV BACKEND_PORT=4000

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY data/data.go ./data/

RUN go build -o server .

EXPOSE $BACKEND_PORT

CMD [ "./server" ]

FROM alpine:3.18.4

COPY --from=builder /app/server /app
COPY $DOT_ENV ./

EXPOSE 8000

CMD /app