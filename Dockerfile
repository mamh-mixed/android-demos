FROM golang

RUN go get -u -v github.com/omigo/log
RUN go get -u -v gopkg.in/mgo.v2
RUN go get -u -v gopkg.in/mgo.v2/bson

RUN go install github.com/CardInfoLink/quickpay

ENTRYPOINT /go/bin/quickpay

EXPOSE 3009
