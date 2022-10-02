FROM golang:latest

ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update

# build go app
RUN go build -o gnp ./cmd/main.go

# Add docker-compose-wait tool -------------------
ENV WAIT_VERSION 2.7.2
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait

CMD ["./gnp"]