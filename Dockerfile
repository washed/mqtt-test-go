FROM golang:alpine

WORKDIR /usr/src/app/
COPY . /usr/src/app/

RUN go mod download

ENTRYPOINT ["go", "run"]
