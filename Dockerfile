FROM golang:1.22-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

WORKDIR /app

COPY go.mod go.sum ./
COPY . .

RUN go mod tidy


RUN go build -o autodeployment .

CMD ["./autodeployment"]
