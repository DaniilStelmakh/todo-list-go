FROM golang:1.23.2

WORKDIR /usr/src/app

COPY . .

RUN go mod download

RUN go build -o /todo_app/main todo_app/main.go

CMD [ "/todo_app/main" ]