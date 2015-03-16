package cfca

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/omigo/g"
)

var requestURL = "https://test.china-clearing.com/Gateway/InterfaceII"

const (
	cfca_ev_oca_crt = `-----BEGIN CERTIFICATE-----
MIIFTjCCAzagAwIBAgIGALTPlDJmMA0GCSqGSIb3DQEBCwUAMFYxCzAJBgNVBAYTAkNOMTAwLgYD
VQQKDCdDaGluYSBGaW5hbmNpYWwgQ2VydGlmaWNhdGlvbiBBdXRob3JpdHkxFTATBgNVBAMMDENG
Q0EgRVYgUk9PVDAeFw0xMjA4MDgwNjA2MzFaFw0yOTEyMjkwNjA2MzFaMFUxCzAJBgNVBAYTAkNO
MTAwLgYDVQQKDCdDaGluYSBGaW5hbmNpYWwgQ2VydGlmaWNhdGlvbiBBdXRob3JpdHkxFDASBgNV
BAMMC0NGQ0EgRVYgT0NBMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA02OMsGFxFQIP
MKVPoRaO9rHNX41xbq8jhnbdK0MDVbxfGa3b8QTKxMcmxlRlULfsaie0cIlaRl0AUcJPQH9ftekz
h4T287xqsEAydYQHf77arWQ5nY3fR9RcoBq9pTCQbqw49S6/jHA5oPQaEoKbF0G8zfVKp5PrcKSu
fHMQyKo/Ez2UYT+gut36j4GYpAABuV6PbusPpjufsN9Br9+xqgyz8ubSp1Wl1qSlvQUQBhAJAH+a
3NMhD0illaGfTdWbF485a5NilMFGqJBa/kLVEYwG4aoKdV9vG/NFS0LKz3QVnB7bkrLjTkuGN/zQ
JP0daJ3CGAzmN+Cr2ujtXOfAYwIDAQABo4IBITCCAR0wOAYIKwYBBQUHAQEELDAqMCgGCCsGAQUF
BzABhhxodHRwOi8vb2NzcC5jZmNhLmNvbS5jbi9vY3NwMB8GA1UdIwQYMBaAFOP+Lf0o0Au1urai
xL8GqgWMk/svMA8GA1UdEwEB/wQFMAMBAf8wRAYDVR0gBD0wOzA5BgRVHSAAMDEwLwYIKwYBBQUH
AgEWI2h0dHA6Ly93d3cuY2ZjYS5jb20uY24vdXMvdXMtMTIuaHRtMDoGA1UdHwQzMDEwL6AtoCuG
KWh0dHA6Ly9jcmwuY2ZjYS5jb20uY24vZXZyY2EvUlNBL2NybDEuY3JsMA4GA1UdDwEB/wQEAwIB
BjAdBgNVHQ4EFgQUVQji3MyVbR9d3rNH6OkWxsBFd8QwDQYJKoZIhvcNAQELBQADggIBAMmFEIoC
E9UNmb2BYYhTRV12kNVucP6t683BaFTgJizIJw/ebvvTdWNTycyP5MQFlHKrIYwjvFO9Rfw8+yIs
sT3JFYiqsLBswvaMr3AIuA2mTnmasvZFe6P19qitzTRkz+TL6TFailrtnzudsvn2SeVbRiX+6Csy
NNMoPsRHTeZAEpkB7J3vh+ZAiv3gsIXtjtz5Y1iWWRZipemJ/qEfW2hDONB+T6lGcEXHDi9dIkWc
C/jFT4XPM64pagAz9gEGZg1PzFBE8QMxiwaDAOeaG0l0e/HW4wJlo4ZzOELqZGJLlYhQ8AkBYR95
NEtR9j5bWK98Lznykldk2MDLBD2mrIfMkVjMwEj4A8ElMXsLnWXXg41NN6gjUm2/IudKOaGqniPs
5SZrN36O4B3NzsaZdLznHH5H0+aksurjgme8RAG0A2OAnRG3VXBWrxud7t0KDINLs+mxY7IR+xVZ
2cw6Cer8HnAVfKPJrbdq7vyJJkIpCll+mLHaGgvv3IqiU4rrrllE3NYjKG4Fk2MiYvZg10KXA8tl
YsLt8I/RcNmC2TvjZHYVE3tanbGw53TRGFk2Vq68XOkvooOardihwRkgqcOgUvouORuvSqTlkQiz
TFH6FTUt3xuuED4dnn5N/1ijcDt0N3l5ovoyHOVcYiO4drCN96LHiUoiSfYODmpXG2tl
-----END CERTIFICATE-----`

	cfca_ev_root_crt = `-----BEGIN CERTIFICATE-----
MIIFjTCCA3WgAwIBAgIEGErM1jANBgkqhkiG9w0BAQsFADBWMQswCQYDVQQGEwJDTjEwMC4GA1UE
CgwnQ2hpbmEgRmluYW5jaWFsIENlcnRpZmljYXRpb24gQXV0aG9yaXR5MRUwEwYDVQQDDAxDRkNB
IEVWIFJPT1QwHhcNMTIwODA4MDMwNzAxWhcNMjkxMjMxMDMwNzAxWjBWMQswCQYDVQQGEwJDTjEw
MC4GA1UECgwnQ2hpbmEgRmluYW5jaWFsIENlcnRpZmljYXRpb24gQXV0aG9yaXR5MRUwEwYDVQQD
DAxDRkNBIEVWIFJPT1QwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDXXWvNED8fBVnV
BU03sQ7smCuOFR36k0sXgiFxEFLXUWRwFsJVaU2OFW2fvwwbwuCjZ9YMrM8irq93VCpLTIpTUnrD
7i7es3ElweldPe6hL6P3KjzJIx1qqx2hp/Hz7KDVRM8Vz3IvHWOX6Jn5/ZOkVIBMUtRSqy5J35DN
uF++P96hyk0g1CXohClTt7GIH//62pCfCqktQT+x8Rgp7hZZLDRJGqgG16iI0gNyejLi6mhNbiyW
ZXvKWfry4t3uMCz7zEasxGPrb382KzRzEpR/38wmnvFyXVBlWY9ps4deMm/DGIq1lY+wejfeWkU7
xzbh72fROdOXW3NiGUgthxwG+3SYIElz8AXSG7Ggo7cbcNOIabla1jj0Ytwli3i/+Oh+uFzJlU9f
py25IGvPa931DfSCt/SyZi4QKPaXWnuWFo8BGS1sbn85WAZkgwGDg8NNkt0yxoekN+kWzqotaK8K
gWU6cMGbrU1tVMoqLUuFG7OA5nBFDWteNfB/O7ic5ARwiRIlk9oKmSJgamNgTnYGmE69g60dWIol
hdLHZR4tjsbftsbhf4oEIRUpdPA+nJCdDC7xij5aqgwJHsfVPKPtl8MeNPo4+QgO48BdK4PRVmrJ
tqhUUy54Mmc9gn900PvhtgVguXDbjgv5E1hvcWAQUhC5wUEJ73IfZzF4/5YFjQIDAQABo2MwYTAf
BgNVHSMEGDAWgBTj/i39KNALtbq2osS/BqoFjJP7LzAPBgNVHRMBAf8EBTADAQH/MA4GA1UdDwEB
/wQEAwIBBjAdBgNVHQ4EFgQU4/4t/SjQC7W6tqLEvwaqBYyT+y8wDQYJKoZIhvcNAQELBQADggIB
ACXGumvrh8vegjmWPfBEp2uEcwPenStPuiB/vHiyz5ewG5zz13ku9Ui20vsXiObTej/tUxPQ4i9q
ecsAIyjmHjdXNYmEwnZPNDatZ8POQQaIxffu2Bq41gt/UP+TqhdLjOztUmCypAbqTuv0axn96/Ua
4CUqmtzHQTb3yHQFhDmVOdYLO6Qn+gjYXB74BGBSESgoA//vU2YApUo0FmZ8/Qmkrp5nGm9BC2sG
E5uPhnEFtC+NiWYzKXZUmhH4J/qyP5Hgzg0b8zAarb8iXRvTvyUFTeGSGn+ZnzxEk8rUQElsgIfX
BDrDMlI1Dlb4pd19xIsNER9Tyx6yF7Zod1rg1MvIB671Oi6ON7fQAUtDKXeMOZePglr4UeWJoBjn
aH9dCi77o0cOPaYjesYBx4/IXr9tgFa+iiS6M+qf4TIRnvHST4D2G0CvOJ4RUHlzEhLN5mydLIhy
PDCBBpEi6lmt2hkuIsKNuYyH4Ga8cyNfIWRjgEj1oDwYPZTISEEdQLpe/v5WOaHIz16eGWRGENoX
kbcFgKyLmZJ956LYBws2J+dIeWCKw9cTXPhyQN9Ky8+ZAAoACxGV2lZFA4gKn2fQ1XmxqI1AbQ3C
ekD6819kR5LLU7m7Wc5P/dAVUwHY3+vZ5nbv0CO7O6l5s9UCKc2Jo5YPSjXnTkLAdc0Hz+Ys63su
-----END CERTIFICATE-----`
)

var cli *http.Client

// 初始化中金 Https 客户端
func init() {
	certs := x509.NewCertPool()
	certs.AppendCertsFromPEM([]byte(cfca_ev_oca_crt))
	certs.AppendCertsFromPEM([]byte(cfca_ev_root_crt))

	// 发送请求
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: certs,
			// InsecureSkipVerify: true, // only for testing
		},
	}
	cli = &http.Client{Transport: tr}
}

// sendRequest 对中金接口访问的统一处理
func sendRequest(req *BindingRequest) *BindingResponse {
	values := prepareRequestData(req)
	if values == nil {
		return nil
	}

	body := send(values)
	if body == nil {
		return nil
	}

	return processResponseBody(body)
}

func prepareRequestData(req *BindingRequest) (v *url.Values) {
	// xml 编组
	xmlBytes, err := xml.Marshal(req)
	if err != nil {
		g.Error("unable to marshal xml:", err)
		return nil
	}
	g.Debug("请求报文: %s", xmlBytes)

	// 对 xml 作 base64 编码
	b64Str := base64.StdEncoding.EncodeToString(xmlBytes)
	g.Trace("base64: %s", b64Str)

	// 对 xml 签名
	hexSign := signatureUseSha1WithRsa(xmlBytes, req.SignCert)
	g.Trace("请求签名: %s", hexSign)

	// 准备参数
	v = &url.Values{}
	v.Add("message", b64Str)
	v.Add("signature", hexSign)

	return v
}

func send(v *url.Values) (body []byte) {
	ret, err := cli.PostForm(requestURL, *v)
	if err != nil {
		g.Error("unable to connect Cfca gratway ", err)
		return nil
	}

	// 处理返回报文
	body, err = ioutil.ReadAll(ret.Body)
	if err != nil {
		g.Error("unable to read from resp ", err)
		return nil
	}
	g.Trace("resp: [%s]", body)

	return body
}

func processResponseBody(body []byte) (resp *BindingResponse) {
	// 得到报文和签名
	result := strings.Split(string(body), ",")
	rb64Str := strings.TrimSpace(result[0])
	// 数据 base64 解码
	rxmlBytes, err := base64.StdEncoding.DecodeString(rb64Str)
	if err != nil {
		g.Error("unable to decode base64 content ", err)
	}
	g.Debug("返回报文: %s", rxmlBytes)

	// 验签
	rhexSign := strings.TrimSpace(result[1])
	g.Trace("返回签名: %s", rhexSign)
	err = checkSignatureUseSha1WithRsa(rxmlBytes, rhexSign)
	if err != nil {
		g.Error("check sign failed ", err)
		return nil
	}

	// 解编 xml
	resp = new(BindingResponse)
	err = xml.Unmarshal(rxmlBytes, resp)
	if err != nil {
		g.Error("unable to unmarshal xml ", err)
		return nil
	}

	return resp
}
