FROM golang

RUN go get github.com/omigo/log
RUN go get gopkg.in/mgo.v2
RUN go get gopkg.in/mgo.v2/bson
RUN go get github.com/CardInfoLink/quickpay

RUN go install github.com/CardInfoLink/quickpay

ENTRYPOINT /go/bin/quickpay

EXPOSE 3009
