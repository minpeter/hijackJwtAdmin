FROM golang:1.19-alpine

WORKDIR /app


ENV BACKEND_PORT=4000

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY data/data.go ./data/
COPY .env.* ./

RUN go build

EXPOSE $BACKEND_PORT

CMD [ "./hijackJwtAdmin" ]