FROM golang:1.22 

RUN apt-get update && apt-get -y upgrade && apt-get -y install gcc g++ ca-certificates chromium xvfb

COPY . /app
WORKDIR /app

RUN go mod download
RUN GOOS=linux make bin/sync

CMD ["./bin/sync"]
