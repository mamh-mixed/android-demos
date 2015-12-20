{
	"Currency":"CNY",
	"ChanCode":{
		"ALP":"支付宝",
		"WXP":"微信",
		"Unknown":"未知"
	},
	"ReportName":{
		"SettleJournal":"对账单",
		"SettleSummary":"汇总清算文件"
	},
	"Role":{
		"Company" :"公司",
		"Agent"   :"代理",
		"Group"   :"商户",
		"Mer"     :"门店"
	},
	"BusicdType":{
		"Purc":"下单并支付",
		"Paut":"预下单",
		"Inqy":"查询",
		"Refd":"退款",
		"Void":"撤销",
		"Canc":"取消",
		"Qyzf":"企业付款",
		"Jszf":"公众号支付",
		"Veri":"卡券核销",
		"Unknown":"未知"
	},
	"TransStatus":{
		"TransHandling":"交易处理中",
		"TransFail":"交易失败",
		"TransSuccess":"交易成功",
		"TransClosed":"交易已退款",
		"Unknown":"未知"
	},
	"StatReport":{
		"Title"					:"商户交易报表汇总",
		"Total" 				:"总计：",
		"StartDate"				:"开始日期：",
		"EndDate" 				:"结束日期：",
		"Remark"				:"注：手续费为每笔单笔计算后四舍五入精确到分，跟总额计算手续费略有误差。因本表仅统计了讯联数据系统的数据，数据仅供参考",
		"MerId"					:"商户号",
		"MerName"				:"商户名称",
		"Summary"				:"汇总",
		"ALP" 					:"支付宝",
		"WXP"					:"微信",
		"AgentName"				:"代理名称",
		"TotalCount"			:"总笔数",
		"TotalAmt" 				:"总金额",
		"Count"					:"笔数",
		"Amt"					:"金额",
		"Fee"					:"手续费"
	},
	"TransReport":{
		"SheetName"      			:"商户交易报表",
		"MerId"          			:"商户号",
		"MerName"        			:"商户名称",
		"OrderNum"       			:"订单号",
		"TransAmt"       			:"交易金额",
		"TransCurr"					:"交易币种",
		"MerFee"					:"商户手续费",
		"ChanCode"       			:"渠道",
		"TransTime"      			:"交易时间",
		"PayTime"					:"支付时间",
		"TransStatus"    			:"交易状态",
		"ChanMerId"     			:"渠道商户号",
		"AgentCode"      			:"机构",
		"TerminalId"     			:"终端号",
		"Busicd"         			:"交易类型",
		"OrigOrderNum"   			:"原订单号",
		"Remark"					:"备注",
		"IsSettled"					:"是否参与清算",
		"RefundAmt"      			:"退款金额",
		"Fee"            			:"手续费",
		"SettAmt"        			:"清算金额",
		"TotalTransAmt"  			:"交易总额",
		"TotalRefundAmt" 			:"退款总额",
		"TotalFee"       			:"手续费总额",
		"TotalSettAmt"   			:"清算总额",
		"Yes"						:"是",
		"No"						:"否"
	},
	"ImportMessage":{
		"Yes"					:"是",
		"No"					:"否",
		"SysErr"   				:"系统错误，请重新上传。",
		"EmptyErr"				:"上传表格为空，请检查。",
		"FileErr"				:"无法获取文件，请重新上传。",
		"CellMapErr"			:"%d 行：%s",
		"ColNumErr"          	:"列数有误，实际应为 %d 行，读取到 %d 行。具体信息为：%s",
		"ImportSuccess"			:"成功处理 %d 行数据，耗时 %s。",
		"MerIdRepeat"        	:"门店号(%s)重复",
		"DataHandleErr":{
			"NotSupportOperation":"第 %d 行，暂不支持 %s 操作。",
			"NoMerId"            :"第 %d 行，门店号为空",
			"MerIdFormatErr"     :"第 %d 行，门店号 %s 格式错误",
			"MerIdExist"         :"门店：%s 已存在",
			"MerIdNotExist"      :"门店：%s 不存在",
			"ALPMerchantErr"     :"支付宝商户(%s): %s",
			"WXPMerchantErr"     :"微信商户(%s): %s",
			"UsernameExist"      :"门店：%s 账户信息-用户名：%s 已存在",
			"UsernameNotExist"   :"门店：%s 账户信息-用户名：%s 不存在",
			"AgentNotExist"  	 :"门店：%s 代理代码(%s)不存在",
			"CompanyNotExist"	 :"门店：%s 公司代码(%s)不存在",
			"CompanyBelongsErr"  :"门店：%s 公司代码(%s)不属于该代理",
			"GroupNotExist"      :"门店：%s 商户代码(%s)不存在",
			"GroupBelongsErr"    :"门店：%s 商户代码(%s)不属于该代理",
			"NoALPKey"           :"支付宝商户：%s 密钥为空",
			"NoALPRouteToUdpSf"  :"没找到门店：%s，对应的支付宝商户，无法变更清算标识。",
			"NoWXPRouteToUdpSf"  :"没找到门店：%s，对应的微信商户，无法变更清算标识。",
			"WXPNotAgentMode"    :"微信商户：%s 并不是受理商模式",
			"SysConfigErr"       :"系统错误配置，请联系管理员。",
			"AgentMerInfoErr"    :"微信商户：%s 代理商商户号填写错误，应为 %s，实际为 %s",
			"AgentModeNotMatch"  :"微信商户：%s 为受理商模式",
			"NoSuchAgentMer"     :"微信商户：%s 系统中没有代码为 %s 的代理商商户",
			"NoWXPKey"           :"微信商户：%s 密钥为空",
			"CILFeeErr"          :"讯联跟渠道费率格式错误(%s)",
			"MerFeeErr"          :"商户跟讯联费率格式错误(%s)",
			"CILFeeOverMax"      :"讯联跟渠道费率超过最大值 3% (%s)",
			"MerFeeOverMax"      :"商户跟讯联费率超过最大值 3% (%s)"
		},
		"ValidateErr":{
			"NoMerName"       :"门店：%s 门店名称为空",
			"NoSignKey"		  :"门店：%s 开启验签需要填写签名密钥",
			"NoAgentCode"     :"门店：%s 代理代码为空",
			"OpenSignValueErr":"是否开启验签：%s 取值错误，应为【是】或【否】",
			"AddAcctValueErr" :"是否新增账户信息：%s 取值错误，应为【是】或【否】",
			"UNOrPWDEmptyErr" :"门店：%s 新增账户信息用户名或密码为空",
			"SignLengthErr"   :"门店：%s 签名密钥长度错误(%s)",
			"NoCommodityName" :"门店：%s 商品名称为空",
			"IsAgentStrErr"   :"是否代理商模式：%s 取值错误，应为【是】或【否】",
			"NoWXPMer"        :"门店：%s 代理商模式需要填写微信商户号",
			"WXPSettFlagErr"  :"微信商户清算标识：%s 取值错误，应为[CIL,CHANNEL,AGENT,COMPANY,GROUP]",
			"ALPSettFlagErr"  :"支付宝商户清算标识：%s 取值错误，应为[CIL,CHANNEL,AGENT,COMPANY,GROUP]",
			"NoOverseasChanMer" :"门店：%s 支付宝境外商户必须填写 merchant_name 和 merchant_no",
			"NoSchemeType"		:"门店：%s 支付宝境外商户必须填写计费方案"
		}
	},
	"MerchantExport":{
		"Title":"商户表",
		"MerId":"商户编号",
		"MerName":"商户名",
		"IsNeedSign":"是否验签",
		"SignKey":"签名密钥",
		"BillUrl":"账单链接",
		"PayUrl":"支付链接",
		"Yes":"是",
		"No":"否"
	}
}
