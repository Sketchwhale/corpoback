FROM golang:alpine

ENV GO111MODULE=on \
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64

WORKDIR /stuff

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 10000

CMD ["/stuff/main"]
