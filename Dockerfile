FROM golang

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go build -o main ./cmd/main.go

EXPOSE 8080
EXPOSE 8081
CMD [ "/app/main" ]