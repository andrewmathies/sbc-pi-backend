FROM golang:1.8

WORKDIR /go/src/sbc-pi-backend
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["sbc-pi-backend"]
