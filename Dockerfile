FROM golang:1.10

# TODO - will passing this in docker-compose override this?
ENV THIS_QUEUE 0

RUN mkdir /app
RUN mkdir -p /go/src/github.com/waltermblair/brain
COPY . /go/src/github.com/waltermblair/brain/
WORKDIR /go/src/github.com/waltermblair/brain

RUN go get -u github.com/golang/dep/...
RUN dep ensure -vendor-only
RUN go build -o /app/main .

ENTRYPOINT ["/app/main"]