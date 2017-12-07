FROM golang:1.8

WORKDIR /go/src/app

RUN go-wrapper download github.com/PuerkitoBio/fetchbot
RUN go-wrapper install github.com/PuerkitoBio/fetchbot

COPY go/src/app/lotherry.go .

RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run", "app"]
