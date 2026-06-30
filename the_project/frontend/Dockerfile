FROM golang:1.22

WORKDIR /app

COPY . .

RUN go build -o todo_frontend .

CMD ["./todo_frontend"]
