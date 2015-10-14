FROM golang

RUN go get ./...
ADD . /go/src/github.com/CardInfoLink/quickpay

RUN go install github.com/CardInfoLink/quickpay

WORKDIR /go/src/github.com/CardInfoLink/quickpay
ENTRYPOINT /go/bin/quickpay

EXPOSE 6800 6600 6601
