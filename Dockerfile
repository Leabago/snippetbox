FROM golang:latest
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o snippetbox ./cmd/web/
RUN ls -l
CMD ["/app/snippetbox"]