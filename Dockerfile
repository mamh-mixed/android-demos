FROM golang

RUN go get github.com/omigo/log && \
    go get github.com/nu7hatch/gouuid && \
    go get github.com/axgle/mahonia && \
    go get gopkg.in/mgo.v2 && \
    go get gopkg.in/mgo.v2/bson
ADD . /go/src/github.com/CardInfoLink/quickpay

RUN go install github.com/CardInfoLink/quickpay

WORKDIR /go/src/github.com/CardInfoLink/quickpay
ENTRYPOINT /go/bin/quickpay

EXPOSE 3009
