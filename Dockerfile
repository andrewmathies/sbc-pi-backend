FROM golang:1.8

WORKDIR /go/src/app
COPY . .

RUN apt-get update && apt-get -y upgrade

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["sbc-pi-backend"]
