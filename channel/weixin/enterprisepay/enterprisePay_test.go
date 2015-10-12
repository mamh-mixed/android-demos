package enterprisepay

import (
	"testing"

	"github.com/CardInfoLink/quickpay/model"
	// "github.com/CardInfoLink/quickpay/util"
)

var clientCert = []byte(`-----BEGIN CERTIFICATE-----
MIIEaDCCA9GgAwIBAgIDAYRnMA0GCSqGSIb3DQEBBQUAMIGKMQswCQYDVQQGEwJD
TjESMBAGA1UECBMJR3Vhbmdkb25nMREwDwYDVQQHEwhTaGVuemhlbjEQMA4GA1UE
ChMHVGVuY2VudDEMMAoGA1UECxMDV1hHMRMwEQYDVQQDEwpNbXBheW1jaENBMR8w
HQYJKoZIhvcNAQkBFhBtbXBheW1jaEB0ZW5jZW50MB4XDTE1MDIxMDA5NTIwNVoX
DTI1MDIwNzA5NTIwNVowgZgxCzAJBgNVBAYTAkNOMRIwEAYDVQQIEwlHdWFuZ2Rv
bmcxETAPBgNVBAcTCFNoZW56aGVuMRAwDgYDVQQKEwdUZW5jZW50MQ4wDAYDVQQL
EwVNTVBheTEtMCsGA1UEAxQk5LiK5rW36K6v6IGU5pWw5o2u5pyN5Yqh5pyJ6ZmQ
5YWs5Y+4MREwDwYDVQQEEwgxMDA4Njc2NjCCASIwDQYJKoZIhvcNAQEBBQADggEP
ADCCAQoCggEBALJAQYqloHovS+fh8+jVDfHrGqgB3Qb+gTM/E8Q0aSikyXvIAON5
9Nie1EY4nPrXf+soJigImUjeFb9jExA4pnky7s0h24RXPlTJKO4t7+qpS3IGK9Is
3UFDF19IVqqjFFVodY2ZXyyyL4wUJcHWiIuNF5YaEv7lAvXFQI45xAMD5a5H/+js
FEcvNVz8HcbH5mB5+f4igCZW9LQh17/LEmg9zcoSRNIeVeWf69hDs5u3//ntjA7u
TZXWVxj8A4tpB6i0VnX0GMW9FABBWiLZNpYxeEtdPoculmpjwxYErpEaFBuSJdlM
aXHcVNtsx04S/HxBP4x2yFhMY4ZOPScRJvcCAwEAAaOCAUYwggFCMAkGA1UdEwQC
MAAwLAYJYIZIAYb4QgENBB8WHSJDRVMtQ0EgR2VuZXJhdGUgQ2VydGlmaWNhdGUi
MB0GA1UdDgQWBBTz34uDT3mCAXESdt6AQzTpfCeyqTCBvwYDVR0jBIG3MIG0gBQ+
BSb2ImK0FVuIzWR+sNRip+WGdKGBkKSBjTCBijELMAkGA1UEBhMCQ04xEjAQBgNV
BAgTCUd1YW5nZG9uZzERMA8GA1UEBxMIU2hlbnpoZW4xEDAOBgNVBAoTB1RlbmNl
bnQxDDAKBgNVBAsTA1dYRzETMBEGA1UEAxMKTW1wYXltY2hDQTEfMB0GCSqGSIb3
DQEJARYQbW1wYXltY2hAdGVuY2VudIIJALtUlyu8AOhXMA4GA1UdDwEB/wQEAwIG
wDAWBgNVHSUBAf8EDDAKBggrBgEFBQcDAjANBgkqhkiG9w0BAQUFAAOBgQCf9iA5
tcB/DqOTdugbAHL47WZNem+9lZnhUL1bppmOnqeKu7uslg1OLNm5KKRFXWX1s2aK
SMYkeEvTrW7dtF82deP7fQ1DM3KqTIrizFsvIKxHdiIhZk1wv/MgQURWTw8Q0ILV
2OPKal8lbrkhHcAmThB8994mVzgSvwfVgADU9g==
-----END CERTIFICATE-----
`)
var clientKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAskBBiqWgei9L5+Hz6NUN8esaqAHdBv6BMz8TxDRpKKTJe8gA
43n02J7URjic+td/6ygmKAiZSN4Vv2MTEDimeTLuzSHbhFc+VMko7i3v6qlLcgYr
0izdQUMXX0hWqqMUVWh1jZlfLLIvjBQlwdaIi40XlhoS/uUC9cVAjjnEAwPlrkf/
6OwURy81XPwdxsfmYHn5/iKAJlb0tCHXv8sSaD3NyhJE0h5V5Z/r2EOzm7f/+e2M
Du5NldZXGPwDi2kHqLRWdfQYxb0UAEFaItk2ljF4S10+hy6WamPDFgSukRoUG5Il
2UxpcdxU22zHThL8fEE/jHbIWExjhk49JxEm9wIDAQABAoIBAHTfQJoBgrjaLfC5
68yrdNs7hVzG6/7b/CZ3oyQwIQrvENRCDKMZoXoumYv5LqQhi9PJnUr+aFKomqXD
9gnauvwYcw64tk+NTGcXBek04WuA2ODIPw8tL1zM+pQUwA5dosVlGj5fY1HQ+u/j
feYcHacyOVbHfdD2ovw1+t/F7Ej0uoWx0uRZRHk4fkdw/JxiOt5ZFNVVy+YLL2Na
AADZ3uaHuSNMyrr99qIHo6jlSwCz/Eg4fKpaWmM1BTm9c5xXC1pKZvQLs9FU0Xcs
9qRjl9hTOHH8lgHl+MMago4GHt8230dRmkt7gHrgW2WUYv38KHCaCUayE9FnRFBo
xWiVE9ECgYEA42KrDJ5XZQ5j0+e18Bbs5Vvg+SzHq/aBEHOThxK56kzPs75wGP+J
nAaH0RUENQQXKdaKDWEn/BrbDuWcgox6+9Q4jYyD4RKbDZmOS4prMx/OTP6nQX0c
jV4Agec0gHtM1mIitKvmHHjy6MPSJyD9x19sN/z8HjgaNcfdJzl3NBsCgYEAyK6y
1dLkCQ9rNyKoa3KYDnFrmVGwzV9MsMBHkOT2X2E55BMzOk0KzAGuwFuyz7ASku9N
Zgt+iUIsSXulrUlD9RBY7+ViHFYI81Mn4zQsT6xSUPsDdjDtJES3MyrUWFQuBqyx
T+dEdiMrnCu34uAsucjr/bPAzoq/e9m+wCdgLlUCgYBenmU9B/qn85f4yre7o16K
hnQUW9zuotHMDbv6/gDdDX90dS9iR5t0kIcdqtgoU35sC3lA3gfscSRsi4FYFarr
dcDerfUGyF47B4Xdy0iWaorHIURqDOy/qrkdVR9Uw3oSz51PfyRu/qld0HZ3j9Pq
jbuThLNIw+GsNXHCa7g9twKBgQCe6K8V4CvP/NSiUSBaDODZNvjD7Er7JQZ+Q5On
6ZYpyrxjnMyI0u7EwmRVT4dVLBBZJ4L7Vgi5uZuGCJIVmJlmWa6DL/kzhPELdIJ2
SB76a/K2yz1ffriZaoxCyRxiYS4c/Oxe4Dt27LygqoGu9mKULsSoHYmEQ5wC/1Yr
pEJmDQKBgQCQYrzCabObH3Yx0N2Uxc/rqcIXqPuW5lR6BHZkN8EzLzykxX068zfL
rwERwybTKRoGin2PuStBdbibuXt0sf8JrTCBkFPzJI/BL15ORnbG/rYe35gQrNhd
6ePbNPU4+UhogBJ6dq+fRrFUGGwuW2GIbKwea78ptA+Lq1C+siAqpg==
-----END RSA PRIVATE KEY-----
`)

func TestEnterprisePay(t *testing.T) {

	clientCert = []byte(`-----BEGIN CERTIFICATE-----\nMIIEazCCA9SgAwIBAgIDBPUGMA0GCSqGSIb3DQEBBQUAMIGKMQswCQYDVQQGEwJD\nTjESMBAGA1UECBMJR3Vhbmdkb25nMREwDwYDVQQHEwhTaGVuemhlbjEQMA4GA1UE\nChMHVGVuY2VudDEMMAoGA1UECxMDV1hHMRMwEQYDVQQDEwpNbXBheW1jaENBMR8w\nHQYJKoZIhvcNAQkBFhBtbXBheW1jaEB0ZW5jZW50MB4XDTE1MDgxMzA4NDAwMVoX\nDTI1MDgxMDA4NDAwMVowgZsxCzAJBgNVBAYTAkNOMRIwEAYDVQQIEwlHdWFuZ2Rv\nbmcxETAPBgNVBAcTCFNoZW56aGVuMRAwDgYDVQQKEwdUZW5jZW50MQ4wDAYDVQQL\nEwVNTVBheTEwMC4GA1UEAxQn5rGf6KW/5b2p5aW9576O5oqV6LWE566h55CG5pyJ\n6ZmQ5YWs5Y+4MREwDwYDVQQEEwgxMDQ0MzEzODCCASIwDQYJKoZIhvcNAQEBBQAD\nggEPADCCAQoCggEBALzyFI3D9IF/5ZHJEEqA4BfnKCFoJReGi0C/sAwT/6EAu4l3\n7vihL63J9+0fj2qkZX3gQWyP8A/YM/8XQm6igKaWTKYCcwwzJh6XKOQEyd/ITTRg\nTGL/LAY1a7ckdmE835yDCvq7r8OR9/3s8vxtGwVIT8oEmqS9o05Z+YNyfirFmYA9\nyvSSxJNcKlpZ6GT72JeR8KV/6kPb67ypH5fEPVhg+fOnw/WzM5VZOMox1weNdNUF\nJ47fAtIs9P6GbluUuf0FELva53+8PxxqTBC3Gf2ypyh7HINb/03oo3pgGak+1alC\nG17//nNFLCeFRWPEuU2zIiZpOD19emFebum54wcCAwEAAaOCAUYwggFCMAkGA1Ud\nEwQCMAAwLAYJYIZIAYb4QgENBB8WHSJDRVMtQ0EgR2VuZXJhdGUgQ2VydGlmaWNh\ndGUiMB0GA1UdDgQWBBS4nUW48Uiw9tPxlfUEdx0jEDAgRTCBvwYDVR0jBIG3MIG0\ngBQ+BSb2ImK0FVuIzWR+sNRip+WGdKGBkKSBjTCBijELMAkGA1UEBhMCQ04xEjAQ\nBgNVBAgTCUd1YW5nZG9uZzERMA8GA1UEBxMIU2hlbnpoZW4xEDAOBgNVBAoTB1Rl\nbmNlbnQxDDAKBgNVBAsTA1dYRzETMBEGA1UEAxMKTW1wYXltY2hDQTEfMB0GCSqG\nSIb3DQEJARYQbW1wYXltY2hAdGVuY2VudIIJALtUlyu8AOhXMA4GA1UdDwEB/wQE\nAwIGwDAWBgNVHSUBAf8EDDAKBggrBgEFBQcDAjANBgkqhkiG9w0BAQUFAAOBgQCp\nupR3B2GwQMcwAbS9WzppsKYr/Ovd3rVkqS3uXkA8VRI99oBFdBKkhzlTYeIYt174\nvHeeO2v0oe8ZKelQLERWgshNCVXTV8T+tyap900ISVdBdbwz9eC5lv/fmnY0Z1cy\nWVsQrZWoemNxZ+elmeDNC/rTdktaXrdBXr5AR2FV2g==\n-----END CERTIFICATE-----\n`)
	clientKey = []byte(`-----BEGIN PRIVATE KEY-----\nMIIEugIBADANBgkqhkiG9w0BAQEFAASCBKQwggSgAgEAAoIBAQC88hSNw/SBf+WR\nyRBKgOAX5yghaCUXhotAv7AME/+hALuJd+74oS+tyfftH49qpGV94EFsj/AP2DP/\nF0JuooCmlkymAnMMMyYelyjkBMnfyE00YExi/ywGNWu3JHZhPN+cgwr6u6/Dkff9\n7PL8bRsFSE/KBJqkvaNOWfmDcn4qxZmAPcr0ksSTXCpaWehk+9iXkfClf+pD2+u8\nqR+XxD1YYPnzp8P1szOVWTjKMdcHjXTVBSeO3wLSLPT+hm5blLn9BRC72ud/vD8c\nakwQtxn9sqcoexyDW/9N6KN6YBmpPtWpQhte//5zRSwnhUVjxLlNsyImaTg9fXph\nXm7pueMHAgMBAAECggEAKNCuafVTgwnqwHRLhZyTS/aOL7E8mflhaWo/EEzdopzy\n5f43bBP9pbAEU3/GzkWW+vsPhvEM7Y9JpCgVHKsT2WiMRCwxSAhgiqkilycFqMav\nDspildwDWY/+pcBFpjmtIDUakREsJbYJeRvPXYIAHUMCoWQfN8kTVuCOyrtXz8si\nfKaR/2DIof6WT2pZQhbJoYKo4262gqM4yrtgnsV2XymSRMTyDHnkTxoddzdcD5e+\nmYh3tXuFAOXIzXNyuIA9jVnqpzSupWJrVVz72+Nb75hloHgvxH2Iqx5QSxrW4MuZ\nM7Mf6DUObO26dSY5Pahd1CV9wBj2ilBDAc3h8bvneQKBgQDm9YRQWsxRs1ZZmMBS\nE9EPPWBQ1rCaYwH7hOE2MEJvCMGJmKbqidx3++vM8xlg9VtWvXJSPR5C0fvpGTby\nfRgtU0gmguSLHGtfBWlq8O1MnFS+5s2A45c63Nk95d6ngG7qRNhJGapTZd6wyos2\n7nPxBZN/kp9iYJZtpz+nYRlHAwKBgQDRbnEZbIkhG6yf/I50OQsuwYhTHdJCtSZU\nuUBo9Wb8gGlc9uRxvoPYNOWScPdh7XRslcXpxaRsYLQ5trArkFCSIafLQyexnJC2\nk0y7toFrY8iu090pZmwV9zxckTZPylaeiM7yveTo4M4XDzU3fZYoYhpbUeZADxYT\n3lxEE+eirQKBgBi+Ir+tCoiSKuUMXUYtw07bp27hoSfZBYRZlvsELonQqVNBXFhy\nDoF4JqndPVHK108ymoW+8Hf+IPu27NELn8RzUJ53lV5l3AbhAIspZnK3qMiO12NA\nfpoawNdFwhW1x6wnVfh23G1002ejO0jWQOuISmX3YkXazSLQMRyQZXHbAn9w5RgS\nnG9PaZukooAoJpfgnHLgWzL8wqnnAfW2npLlilNBydVbe1eXNOyMBFlWcKpR3mrR\nmbkIsxh0BAbzdNf57iqFZ63+EGfyA3VZuwgI28FSfOK7bnrVXHEtdRnR6UDINPdZ\n4wOChaySwZ/uScSoADo0//EelwNPso4KJLC5AoGAKqVVrxfAtWs57ciF8zPZsh1+\ne+KfGbEkj8I4YppJcjw9d8UeAltFnZSATjkPiiO6F5JBDNrBwBFYs94/xnL4GUzx\njLQXop3yte0Tr3nfN2vb+/pPLVHP9FQH7WQ9DBZ2t9Mvie05X4OwSTGRaEpUU0ts\niuhniRoP1U9oFA5xmKo=\n-----END PRIVATE KEY-----\n`)
	req := &model.ScanPayRequest{}
	req.ActTxamt = "500"
	req.AppID = "wx889539c239f029ee"
	req.ChanMerId = "1261318101"
	req.CheckName = "NO_CHECK"
	req.Desc = "您的中奖金额已成功兑付！"
	req.OpenId = "opM-kwBHOXEqdwTY9S9lHtXjgevo"
	req.OrderNum = "14433639524984836448"
	req.SignKey = "b0f05a24dd405df9c4a82dd2cbac697a"

	req.WeixinClientCert = clientCert
	req.WeixinClientKey = clientKey

	ret, _ := DefaultClient.ProcessPay(req)
	t.Logf("%+v", ret)
}

func TestEnterpriseQuery(t *testing.T) {
	clientCert = []byte(`-----BEGIN CERTIFICATE-----
MIIEazCCA9SgAwIBAgIDBPUGMA0GCSqGSIb3DQEBBQUAMIGKMQswCQYDVQQGEwJD
TjESMBAGA1UECBMJR3Vhbmdkb25nMREwDwYDVQQHEwhTaGVuemhlbjEQMA4GA1UE
ChMHVGVuY2VudDEMMAoGA1UECxMDV1hHMRMwEQYDVQQDEwpNbXBheW1jaENBMR8w
HQYJKoZIhvcNAQkBFhBtbXBheW1jaEB0ZW5jZW50MB4XDTE1MDgxMzA4NDAwMVoX
DTI1MDgxMDA4NDAwMVowgZsxCzAJBgNVBAYTAkNOMRIwEAYDVQQIEwlHdWFuZ2Rv
bmcxETAPBgNVBAcTCFNoZW56aGVuMRAwDgYDVQQKEwdUZW5jZW50MQ4wDAYDVQQL
EwVNTVBheTEwMC4GA1UEAxQn5rGf6KW/5b2p5aW9576O5oqV6LWE566h55CG5pyJ
6ZmQ5YWs5Y+4MREwDwYDVQQEEwgxMDQ0MzEzODCCASIwDQYJKoZIhvcNAQEBBQAD
ggEPADCCAQoCggEBALzyFI3D9IF/5ZHJEEqA4BfnKCFoJReGi0C/sAwT/6EAu4l3
7vihL63J9+0fj2qkZX3gQWyP8A/YM/8XQm6igKaWTKYCcwwzJh6XKOQEyd/ITTRg
TGL/LAY1a7ckdmE835yDCvq7r8OR9/3s8vxtGwVIT8oEmqS9o05Z+YNyfirFmYA9
yvSSxJNcKlpZ6GT72JeR8KV/6kPb67ypH5fEPVhg+fOnw/WzM5VZOMox1weNdNUF
J47fAtIs9P6GbluUuf0FELva53+8PxxqTBC3Gf2ypyh7HINb/03oo3pgGak+1alC
G17//nNFLCeFRWPEuU2zIiZpOD19emFebum54wcCAwEAAaOCAUYwggFCMAkGA1Ud
EwQCMAAwLAYJYIZIAYb4QgENBB8WHSJDRVMtQ0EgR2VuZXJhdGUgQ2VydGlmaWNh
dGUiMB0GA1UdDgQWBBS4nUW48Uiw9tPxlfUEdx0jEDAgRTCBvwYDVR0jBIG3MIG0
gBQ+BSb2ImK0FVuIzWR+sNRip+WGdKGBkKSBjTCBijELMAkGA1UEBhMCQ04xEjAQ
BgNVBAgTCUd1YW5nZG9uZzERMA8GA1UEBxMIU2hlbnpoZW4xEDAOBgNVBAoTB1Rl
bmNlbnQxDDAKBgNVBAsTA1dYRzETMBEGA1UEAxMKTW1wYXltY2hDQTEfMB0GCSqG
SIb3DQEJARYQbW1wYXltY2hAdGVuY2VudIIJALtUlyu8AOhXMA4GA1UdDwEB/wQE
AwIGwDAWBgNVHSUBAf8EDDAKBggrBgEFBQcDAjANBgkqhkiG9w0BAQUFAAOBgQCp
upR3B2GwQMcwAbS9WzppsKYr/Ovd3rVkqS3uXkA8VRI99oBFdBKkhzlTYeIYt174
vHeeO2v0oe8ZKelQLERWgshNCVXTV8T+tyap900ISVdBdbwz9eC5lv/fmnY0Z1cy
WVsQrZWoemNxZ+elmeDNC/rTdktaXrdBXr5AR2FV2g==
-----END CERTIFICATE-----
`)
	clientKey = []byte(`-----BEGIN PRIVATE KEY-----
MIIEugIBADANBgkqhkiG9w0BAQEFAASCBKQwggSgAgEAAoIBAQC88hSNw/SBf+WR
yRBKgOAX5yghaCUXhotAv7AME/+hALuJd+74oS+tyfftH49qpGV94EFsj/AP2DP/
F0JuooCmlkymAnMMMyYelyjkBMnfyE00YExi/ywGNWu3JHZhPN+cgwr6u6/Dkff9
7PL8bRsFSE/KBJqkvaNOWfmDcn4qxZmAPcr0ksSTXCpaWehk+9iXkfClf+pD2+u8
qR+XxD1YYPnzp8P1szOVWTjKMdcHjXTVBSeO3wLSLPT+hm5blLn9BRC72ud/vD8c
akwQtxn9sqcoexyDW/9N6KN6YBmpPtWpQhte//5zRSwnhUVjxLlNsyImaTg9fXph
Xm7pueMHAgMBAAECggEAKNCuafVTgwnqwHRLhZyTS/aOL7E8mflhaWo/EEzdopzy
5f43bBP9pbAEU3/GzkWW+vsPhvEM7Y9JpCgVHKsT2WiMRCwxSAhgiqkilycFqMav
DspildwDWY/+pcBFpjmtIDUakREsJbYJeRvPXYIAHUMCoWQfN8kTVuCOyrtXz8si
fKaR/2DIof6WT2pZQhbJoYKo4262gqM4yrtgnsV2XymSRMTyDHnkTxoddzdcD5e+
mYh3tXuFAOXIzXNyuIA9jVnqpzSupWJrVVz72+Nb75hloHgvxH2Iqx5QSxrW4MuZ
M7Mf6DUObO26dSY5Pahd1CV9wBj2ilBDAc3h8bvneQKBgQDm9YRQWsxRs1ZZmMBS
E9EPPWBQ1rCaYwH7hOE2MEJvCMGJmKbqidx3++vM8xlg9VtWvXJSPR5C0fvpGTby
fRgtU0gmguSLHGtfBWlq8O1MnFS+5s2A45c63Nk95d6ngG7qRNhJGapTZd6wyos2
7nPxBZN/kp9iYJZtpz+nYRlHAwKBgQDRbnEZbIkhG6yf/I50OQsuwYhTHdJCtSZU
uUBo9Wb8gGlc9uRxvoPYNOWScPdh7XRslcXpxaRsYLQ5trArkFCSIafLQyexnJC2
k0y7toFrY8iu090pZmwV9zxckTZPylaeiM7yveTo4M4XDzU3fZYoYhpbUeZADxYT
3lxEE+eirQKBgBi+Ir+tCoiSKuUMXUYtw07bp27hoSfZBYRZlvsELonQqVNBXFhy
DoF4JqndPVHK108ymoW+8Hf+IPu27NELn8RzUJ53lV5l3AbhAIspZnK3qMiO12NA
fpoawNdFwhW1x6wnVfh23G1002ejO0jWQOuISmX3YkXazSLQMRyQZXHbAn9w5RgS
nG9PaZukooAoJpfgnHLgWzL8wqnnAfW2npLlilNBydVbe1eXNOyMBFlWcKpR3mrR
mbkIsxh0BAbzdNf57iqFZ63+EGfyA3VZuwgI28FSfOK7bnrVXHEtdRnR6UDINPdZ
4wOChaySwZ/uScSoADo0//EelwNPso4KJLC5AoGAKqVVrxfAtWs57ciF8zPZsh1+
e+KfGbEkj8I4YppJcjw9d8UeAltFnZSATjkPiiO6F5JBDNrBwBFYs94/xnL4GUzx
jLQXop3yte0Tr3nfN2vb+/pPLVHP9FQH7WQ9DBZ2t9Mvie05X4OwSTGRaEpUU0ts
iuhniRoP1U9oFA5xmKo=
-----END PRIVATE KEY-----
`)
	req := &model.ScanPayRequest{}
	req.AppID = "wx889539c239f029ee"
	req.ChanMerId = "1261318101"
	req.SignKey = "b0f05a24dd405df9c4a82dd2cbac697a"
	req.OrigOrderNum = "14445733185618269533"
	req.WeixinClientCert = clientCert
	req.WeixinClientKey = clientKey

	DefaultClient.ProcessEnquiry(req)
}
