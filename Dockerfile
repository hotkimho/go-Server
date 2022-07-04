FROM golang:alpine

WORKDIR /app

COPY . .
COPY go.mod ./
COPY go.sum ./
#COPY userAuthhandler.go ./
##COPY boardHandler.go ./
#COPY auth/users/*.go ./
#COPY config/global.go ./
#COPY model/*.go ./

RUN go mod download

RUN go build -o goserver .

EXPOSE 8000
ENTRYPOINT ["./goserver"]