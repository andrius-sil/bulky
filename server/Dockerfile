FROM golang:1.10

WORKDIR /go/src/app
COPY . .

RUN go get -d -v
RUN go install -v

EXPOSE 3000
CMD ["app"]

