FROM golang:1.8

RUN apt-get update && \
apt-get install -y unzip && \
apt-get install -y wget && \
apt-get install -y libx11-6

WORKDIR /lotherry

RUN wget https://chromedriver.storage.googleapis.com/2.33/chromedriver_linux64.zip
RUN unzip chromedriver_linux64.zip

WORKDIR /go/src/app

RUN go-wrapper download github.com/fedesog/webdriver
RUN go-wrapper install github.com/fedesog/webdriver

COPY go/src/app/lotherry.go .

RUN go-wrapper download
RUN go-wrapper install

RUN apt-get install -y libnss3-dev
RUN apt-get install -y libgconf-2-4
RUN apt-get install -y libfontconfig-1

CMD ["go-wrapper", "run", "app"]
