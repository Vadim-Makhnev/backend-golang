FROM golang:1.24.1

WORKDIR /app

COPY go.sum go.mod ./

RUN go mod download

COPY . .

RUN go build -o myapp ./cmd

CMD ["./myapp"]