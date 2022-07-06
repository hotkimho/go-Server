FROM golang:alpine AS builder

WORKDIR /build

COPY go.mod go.sum main.go *.go ./

RUN go mod download

RUN go build -o main .

WORKDIR /dist

RUN cp /build/main .

FROM scratch

COPY --from=builder /dist/main .

EXPOSE 8000
ENTRYPOINT ["/main"]