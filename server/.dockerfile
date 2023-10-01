FROM golang:1.21.1

WORKDIR /app

COPY ./server .

RUN go build .

EXPOSE 80

CMD ["./server", "-difficulty=4", "-time=60"]