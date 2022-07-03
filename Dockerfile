FROM golang:alpine AS builder

WORKDIR /app

COPY . .
#COPY go.mod ./
#COPY *.go ./
#COPY auth/users/*.go ./
#COPY config/*.go ./
#COPY model/*.go ./

RUN go mod tidy

RUN go build main -o goserver


EXPOSE 8000
ENTRYPOINT ["/goserver"]