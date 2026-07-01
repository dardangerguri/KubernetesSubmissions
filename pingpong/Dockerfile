FROM golang:1.22

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

WORKDIR /app

RUN go build -o main .

CMD ["./main"]
