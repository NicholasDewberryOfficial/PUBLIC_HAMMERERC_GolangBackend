FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o application

RUN chmod +x application

EXPOSE 8080

CMD ["./application"]
#ENTRYPOINT [".app"]