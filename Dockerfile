FROM golang:1.22

COPY . /app
WORKDIR /app

RUN mkdir -p bin/
RUN curl http://chromedriver.storage.googleapis.com/$(curl http://chromedriver.storage.googleapis.com/LATEST_RELEASE)/chromedriver_linux64.zip -o ./bin/chromedriver
RUN chmod a+x ./bin/chromedriver

RUN go mod download
RUN GOOS=linux make bin/sync

CMD ["./bin/sync"]
