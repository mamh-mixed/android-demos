package unionlive

import "testing"

func TestEncryptAndSign(t *testing.T) {
	src := `{"header":{"version":"1.0","transType":"W412","transDirect":"Q","sessionId":"747f6cf7-dadf-46ef-83e9-d3c0a87b3dbf","merchantId":"182000001000000","submitTime":"20130501201012","clientTraceNo":"497540"},"body":{"couponsNo":"1809706004000705","termId":"00000667","termSn":"9e908a255b3e5989","amount":"1"}}`
	body, err := encryptAndSign([]byte(src))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("post body: %s", body)
}

func TestCheckSignAndDecrypt(t *testing.T) {
	// src := `{"header":{"version":"1.0","transType":"W412","transDirect":"Q","sessionId":"747f6cf7-dadf-46ef-83e9-d3c0a87b3dbf","merchantId":"182000001000000","submitTime":"20130501201012","clientTraceNo":"497540"},"body":{"couponsNo":"1809706004000705","termId":"00000667","termSn":"9e908a255b3e5989","amount":"1"}}`
	src := `{"header":{"version":null,"transType":null,"submitTime":null,"sessionId":null,"clientTraceNo":null,"hostTime":"20151027100450","hostTraceNo":"2f48b7a1-8a78-4979-95fd-ad9f40379912","returnCode":"0096","merchantId":null,"terminalId":null,"returnMessage":"系统错误","transDirect":null,"signMessage":""},"body":null}`
	message, err := encryptAndSign([]byte(src))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("post body: %s", message)

	jsonBytes, err := checkSignAndDecrypt(message)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("json content: %s", jsonBytes)
}

func TestCheckSignAndDecrypr(t *testing.T) {
	message := "e2ab0bab5cfd7cb84e22405158474e35o27t0qJE7OlvQ8H1UA32JSGwQt7dZlq89rjUrAi4xGSUrQLJh3L38NDzhLhDA+cPTp1W7x7UEjJQDnbIjBAaRuv7DeTKorA/lnZyrvxtT8+j+s/TPfYviKZtI9y4uiuTGsJ1HfqCfZmx5D9Uq/2FYWM3ithPBCTpEOAKua0we38KEEhF17QB6EjKOr3gYRpF9+TWvhRi1SefWtRw2M3RMx/hQem+kPEhmHFEGodTPQMm745y0OxAf0NJm1oTVDNd8xWpwEJmZg1IKWUTYCIoYOQXFjhJ7QgQl/A1YHXVAYtFAIeVJ3dI4J7ds/caaUpzyeYAwkTT8MyWOi8dAG+3opiMaY2Jc4aFALGxTlebxHnu7KUAzTNFmqubFzS9BTydGIqvHJrQvdDzrAAxwK2YsJiWbSUiwavMDPD306TnrU0="

	jsonBytes, err := checkSignAndDecrypt([]byte(message))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("json content: %s", jsonBytes)
}
