FROM golang:1.11

VOLUME ["/app"]
WORKDIR /app

# Get sendgrid-go
RUN mkdir -p /go/src/github.com/sendgrid && \
    cd /go/src/github.com/sendgrid && \
    git clone https://www.github.com/sendgrid/rest && \
    cd rest && \
    go get

ENTRYPOINT ["tail", "-f", "/dev/null"]
