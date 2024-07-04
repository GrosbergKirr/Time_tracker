FROM golang:1.22
ADD ./bin/main /main
WORKDIR /app
COPY . .
COPY .env.docker .env
CMD ["/main"]