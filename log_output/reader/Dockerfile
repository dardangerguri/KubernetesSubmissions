FROM golang:1.25

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

WORKDIR /app/reader

RUN go build -o main .

CMD ["./main"]
