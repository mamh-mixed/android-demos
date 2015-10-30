package unionlive

var SysRespCode map[string]string
var ChanSysRespCode map[string]string

func init() {
	go func() {
		initRespCode()
	}()

}

func initRespCode() {
	SysRespCode = make(map[string]string)
	SysRespCode["00"] = "成功"
	SysRespCode["91"] = "外部系统错误"
	SysRespCode["25"] = "订单不存在"
	SysRespCode["96"] = "系统错误"
	SysRespCode["H5"] = "格式错误"
	SysRespCode["03"] = "商户错误"
	SysRespCode["C1"] = "卡券已被核销"
	SysRespCode["C2"] = "卡券已过期"
	SysRespCode["C3"] = "无效的卡券"
	SysRespCode["C4"] = "券状态异常"
	SysRespCode["C5"] = "未到卡券使用时间"
	SysRespCode["C6"] = "商户不能使用该卡券"
	// SysRespCode["58"] = "未知应答"

	ChanSysRespCode = make(map[string]string)
	ChanSysRespCode["0000"] = "00"
	ChanSysRespCode["0096"] = "91"
	ChanSysRespCode["9999"] = "91"
	ChanSysRespCode["3040"] = "C3"
	ChanSysRespCode["1548"] = "25"
	ChanSysRespCode["1556"] = "96"
	ChanSysRespCode["1557"] = "H5"
	ChanSysRespCode["1457"] = "96"
	ChanSysRespCode["1352"] = "96"
	ChanSysRespCode["1347"] = "03"
	ChanSysRespCode["1348"] = "H5"
	ChanSysRespCode["1109"] = "C3"
	ChanSysRespCode["1123"] = "C3"
	ChanSysRespCode["1103"] = "C2"
	ChanSysRespCode["1117"] = "C4"
	ChanSysRespCode["1104"] = "C4"
	ChanSysRespCode["1107"] = "C5"
	ChanSysRespCode["1105"] = "C6"
	ChanSysRespCode["1102"] = "C1"
	ChanSysRespCode["1101"] = "C1"
}
