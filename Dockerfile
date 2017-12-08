FROM golang:1.8

RUN apt-get update && \
apt-get install -y unzip && \
apt-get install -y wget

WORKDIR /lotherry

RUN wget https://chromedriver.storage.googleapis.com/2.33/chromedriver_mac64.zip
RUN unzip chromedriver_mac64.zip

WORKDIR /go/src/app

RUN go-wrapper download github.com/fedesog/webdriver
RUN go-wrapper install github.com/fedesog/webdriver

COPY go/src/app/lotherry.go .

RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run", "app"]
