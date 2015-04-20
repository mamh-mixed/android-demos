package cil

// only for test
var (
	applePayMerId  = "123456" // apple pay 测试用商户号
	testChcd       = "00000050"
	testMchntId    = "050310058120002"
	testTerminalId = "00000001"

	// 万事达卡测试数据
	testMSCCard       = "5457210001000019"
	testMSCCVV2       = "300"
	testMSCValidDate  = "1412"
	testMSCTrackdata2 = "5457210001000019=1412101080080748"

	// VISA卡测试数据
	testVISCard       = "4761340000000019"
	testVISCVV2       = "830"
	testVISValidDate  = "1712"
	testVISTrackdata2 = "4761340000000019=171210114991787"

	// 银联卡测试数据
	testCUPCard      = "6225220100740059"
	testCUPCVV2      = "111"
	testCUPValidDate = "1605"
	testCUPPhone     = "13611111111"
	testCUPIdentNum  = "130412"

	// Apple Pay测试数据
	testAPPCard       = "5180841200282463"
	testAPPExpireDate = "180531"
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
