FROM golang:1.25

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o todo_frontend .

CMD ["./todo_frontend"]
