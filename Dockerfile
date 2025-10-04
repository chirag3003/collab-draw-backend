FROM golang

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o main .

EXPOSE 5000

CMD ["./main"]