FROM golang:1.22

WORKDIR /app

COPY pingpong/go.mod ./
RUN go mod download

COPY . .

WORKDIR /app/pingpong

RUN go build -o main .

CMD ["./main"]
