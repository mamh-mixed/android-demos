package weixin

import (
	"errors"
	"log"
	"os"
)

type WeixinRequestFactory struct{}

func (_ WeixinRequestFactory) createRequestData(busiType int) WeixinRequest {
	switch busiType {
	case MicroPay:
		return new(MicropayRequest)
	case OrderQuery:
		return new(OrderqueryRequest)
	default:
		log.Println(errors.New("createRequestData errors"))
		os.Exit(1)
		return nil
	}
}

type WeixinResponseFactory struct{}

func (_ WeixinResponseFactory) createResponseData(m WeixinRequest) WeixinResponse {

	switch m.(type) {

	case *MicropayRequest:
		return new(MicroPayResponse)
	case *OrderqueryRequest:
		return new(OrderqueryResponse)
	default:
		log.Println(errors.New("createResponseData errors"))
		os.Exit(1)
		return nil
	}
}
