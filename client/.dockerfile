FROM golang:1.21.1

WORKDIR /app

COPY ./client .

RUN go build .

ENTRYPOINT ["tail", "-f", "/dev/null"]