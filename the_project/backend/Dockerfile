FROM golang:1.22

WORKDIR /app

COPY . .

RUN go build -o todo_backend .

CMD ["./todo_backend"]
