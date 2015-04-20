package cil

// only for test
var (
	applePayMerId  = "123456" // apple pay 测试用商户号
	testAccountNO  = "6225220100740059"
	testIdentNum   = "130412"
	testChcd       = "00000050"
	testMchntId    = "050310058120002"
	testTerminalId = "00000001"
	testCVV2       = "111"
	testValiDate   = "1605"
	testPhoneNum   = "13611111111"
)

// apple pay测试增加了如下路由
// {
//     "merId": "123456",
//     "cardBrand": "VIS",
//     "chanCode": "APT",
//     "chanMerId": "APT123456"
// }

// apple pay测试增加了渠道商户
// {
// 	"chanCode" : "APT",
//     "chanMerId" : "APT123456",
//     "chanMerName" : "Apple Pay测试渠道商户",
// 	"terminalId": "TID123456789",
// 	"insCode" : "99667788"
// }
