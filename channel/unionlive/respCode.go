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
	SysRespCode["C0"] = "字段不能为空"
	SysRespCode["C1"] = "卡券已被核销"
	SysRespCode["C2"] = "卡券已过期"
	SysRespCode["C3"] = "无效的卡券"
	SysRespCode["C4"] = "券状态异常"
	SysRespCode["C5"] = "未到卡券使用时间"
	SysRespCode["C6"] = "商户不能使用该卡券"
	SysRespCode["C7"] = "金额达不到满足优惠条件的最小金额"
	SysRespCode["C8"] = "不能用此支付类型"
	// SysRespCode["58"] = "未知应答"
	SysRespCode["09"] = "刷卡券"

	ChanSysRespCode = make(map[string]string)
	ChanSysRespCode["0000"] = "00" //成功
	ChanSysRespCode["0096"] = "91" //系统错误
	ChanSysRespCode["9999"] = "91" //系统错误
	ChanSysRespCode["3040"] = "C3" //没有找到券信息
	ChanSysRespCode["1548"] = "25" //未找到该终端的消费交易记录
	ChanSysRespCode["1556"] = "96" //交易方向不正确,无法交易
	ChanSysRespCode["1557"] = "H5" //消费次数必须大于等于0
	ChanSysRespCode["1457"] = "96" //请传入终端硬件序列号
	ChanSysRespCode["1352"] = "96" //终端硬件序列号不能为空
	ChanSysRespCode["1347"] = "03" //商户信息不存在
	ChanSysRespCode["1348"] = "H5" //使用次数错误
	ChanSysRespCode["1109"] = "C3" //券码不存在
	ChanSysRespCode["1123"] = "C3" //不支持的券码类型
	ChanSysRespCode["1103"] = "C2" //该券已于XX过期
	ChanSysRespCode["1117"] = "C4" //券已作废
	ChanSysRespCode["1104"] = "C4" //券号状态异常
	ChanSysRespCode["1107"] = "C5" //该券在XX后才能使用
	ChanSysRespCode["1105"] = "C6" //商户不能使用该码
	ChanSysRespCode["1102"] = "C1" //可用次数不足
	ChanSysRespCode["3115"] = "C1" //此类型的券不能再次验券
	ChanSysRespCode["3116"] = "C0" //原交易后台交易流水号不能为空
	ChanSysRespCode["3117"] = "C0" //原验证交易提交时间格式不正确
	ChanSysRespCode["3118"] = "C7" //金额达不到满足优惠条件的最小金额
	ChanSysRespCode["3119"] = "C0" //没有对应的礼包券券号信息
	ChanSysRespCode["3120"] = "C5" //此时间段不能验证
	ChanSysRespCode["3121"] = "C0" //原验证交易客户端流水号不能为空
	ChanSysRespCode["3122"] = "96" //原验证交易验证的次数不能为空
	ChanSysRespCode["3123"] = "C0" //未找到原验证记录
	ChanSysRespCode["3124"] = "C8" //不能用此支付类型
	ChanSysRespCode["1101"] = "C1"
	ChanSysRespCode["1563"] = "96" //无效的终端
	ChanSysRespCode["36"] = "09"   // 刷卡活动券
}
