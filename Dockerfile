FROM chromedp/headless-shell

RUN apt-get update && apt-get -y upgrade && apt-get -y install \
  golang-go make gcc g++ ca-certificates chromium xvfb

COPY . /app
WORKDIR /app

RUN go mod download
RUN GOOS=linux make bin/sync

ENTRYPOINT ["/bin/sh"]
CMD ["-c", "./bin/sync"]
