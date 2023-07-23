FROM golang:1.20-alpine

WORKDIR /usr/local/barracuda

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /usr/local/bin/barracuda ./cmd/...

CMD ["barracuda"]
