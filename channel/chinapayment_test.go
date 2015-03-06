package channel

import (
	"testing"
)

func TestConnectToChinaPay(t *testing.T) {

	request := Request{}
	request.Head.TxCode = "9999"
	request.Body.InstitutionID = "001405"

	ChinaPayRequestHandler(request)

}
