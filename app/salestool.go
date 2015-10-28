// 销售工具接口
package app

import (
	"archive/zip"
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/CardInfoLink/quickpay/email"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/quickpay/qiniu"
	"github.com/CardInfoLink/quickpay/util"
	"github.com/omigo/log"
	"github.com/tealeg/xlsx"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"time"
)

var tokenMap = make(map[string]*model.Session)
var qrImage = "tools/qr/image/%s/%s.jpg"

// CompanyLogin 销售人员-公司级别登录
func CompanyLogin(w http.ResponseWriter, r *http.Request) {

	if !checkSign(r) {
		w.Write(jsonMarshal(model.SIGN_FAIL))
		return
	}

	username := r.FormValue("username")
	user, err := mongo.UserColl.FindOneUser(username, "", "")
	if err != nil {
		w.Write(jsonMarshal(model.USERNAME_NO_EXIST))
		return
	}

	// 不是公司这一级的
	if user.UserType != model.UserTypeCompany {
		w.Write(jsonMarshal(model.USERNAME_NO_EXIST))
		return
	}

	password := r.FormValue("password")
	ps := fmt.Sprintf("%x", sha1.Sum([]byte((model.RAND_PWD + "{" + username + "}" + password))))
	if user.Password != ps {
		w.Write(jsonMarshal(model.USERNAME_PASSWORD_ERROR))
		return
	}

	if user.SubAgentCode == "" {
		log.Errorf("userType is company,but can not find subAgentCode, username=%s", username)
		w.Write(jsonMarshal(model.USER_DATA_ERROR))
		return
	}

	result := model.SUCCESS1
	result.AccessToken = genAccessToken(user)

	w.Write(jsonMarshal(result))
}

// UserList 用户列表
func UserList(w http.ResponseWriter, r *http.Request) {

	var agentUser *model.User
	var ok bool

	// 验证token
	if agentUser, ok = checkAccessToken(r.FormValue("accessToken")); !ok {
		w.Write(jsonMarshal(model.TOKEN_ERROR))
		return
	}

	users, err := mongo.AppUserCol.Find(&model.AppUserContiditon{SubAgentCode: agentUser.SubAgentCode})
	if err != nil {
		w.Write(jsonMarshal(model.SYSTEM_ERROR))
		return
	}

	var merIds []string
	var userMap = make(map[string]*model.AppUser)
	for _, user := range users {
		merIds = append(merIds, user.MerId)
		userMap[user.MerId] = user
	}

	// 关联商户信息
	mers, err := mongo.MerchantColl.FuzzyFind(&model.QueryCondition{
		MerIds: merIds,
	})
	if err != nil {
		w.Write(jsonMarshal(model.SYSTEM_ERROR))
		return
	}

	for _, m := range mers {
		if user, ok := userMap[m.MerId]; ok {
			user.BankOpen = m.Detail.OpenBankName
			user.Payee = m.Detail.AcctName
			user.PayeeCard = m.Detail.AcctNum
			user.PhoneNum = m.Detail.ContactTel
			user.SignKey = m.SignKey
			user.AgentCode = m.AgentCode
			user.UniqueId = m.UniqueId
			user.MerName = m.Detail.MerName
			user.Images = m.Detail.Images
			user.Password = "" // 不显示
		}
	}

	// 成功返回
	result := model.SUCCESS1
	result.Users = users

	w.Write(jsonMarshal(result))
}

// UserRegister 用户注册
func UserRegister(w http.ResponseWriter, r *http.Request) {

	var agentUser *model.User
	var ok bool

	// 验证token
	if agentUser, ok = checkAccessToken(r.FormValue("accessToken")); !ok {
		w.Write(jsonMarshal(model.TOKEN_ERROR))
		return
	}

	req := &reqParams{
		UserName:     r.FormValue("username"),
		Password:     r.FormValue("password"),
		Transtime:    time.Now().Format("20060102150405"),
		Remark:       "company_register",
		UserFrom:     model.SalesToolsRegister,
		BelongsTo:    agentUser.UserName,
		SubAgentCode: agentUser.SubAgentCode,
	}
	// 注册

	w.Write(jsonMarshal(User.register(req)))
}

// GetQiniuToken
func GetQiniuToken(w http.ResponseWriter, r *http.Request) {

	// 验证token
	if _, ok := checkAccessToken(r.FormValue("accessToken")); !ok {
		w.Write(jsonMarshal(model.TOKEN_ERROR))
		return
	}

	result := model.SUCCESS1
	result.UploadToken = qiniu.GetUploadtoken()

	w.Write(jsonMarshal(result))
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(w http.ResponseWriter, r *http.Request) {

	var agentUser *model.User
	var ok bool

	// 验证token
	if agentUser, ok = checkAccessToken(r.FormValue("accessToken")); !ok {
		w.Write(jsonMarshal(model.TOKEN_ERROR))
		return
	}

	appUser, err := mongo.AppUserCol.FindOne(r.FormValue("username"))
	if err != nil {
		w.Write(jsonMarshal(model.USERNAME_NO_EXIST))
		return
	}

	req := &reqParams{
		BankOpen:   r.FormValue("bank_open"),
		Payee:      r.FormValue("payee"),
		PayeeCard:  r.FormValue("payee_card"),
		PhoneNum:   r.FormValue("phone_num"),
		Transtime:  r.FormValue("transtime"),
		Province:   r.FormValue("province"),
		City:       r.FormValue("city"),
		BranchBank: r.FormValue("branch_bank"),
		BankNo:     r.FormValue("bankNo"),
		MerName:    r.FormValue("merName"),
		Images:     r.Form["image"],
		AppUser:    appUser, // 带上user
	}

	// 默认返回
	result := model.SUCCESS1

	var merchant *model.Merchant
	// 还没申请商户
	if appUser.MerId == "" {
		merchant = &model.Merchant{
			AgentCode:    agentUser.AgentCode,
			SubAgentCode: agentUser.SubAgentCode,
			Permission:   []string{model.Paut, model.Purc, model.Canc, model.Void, model.Inqy, model.Refd, model.Jszf, model.Qyzf},
			MerStatus:    model.MerStatusNormal,
			Remark:       "user_register",
			TransCurr:    "156",
			RefundType:   model.CurrentDayRefund, // 只能当天退
			IsNeedSign:   true,
			SignKey:      fmt.Sprintf("%x", randBytes(16)),
			Detail: model.MerDetail{
				MerName:       req.MerName,
				CommodityName: req.MerName,
				Province:      req.Province,
				City:          req.City,
				OpenBankName:  req.BankOpen,
				BankName:      req.BranchBank,
				BankId:        req.BankNo,
				AcctName:      req.Payee,
				AcctNum:       req.PayeeCard,
				ContactTel:    req.PhoneNum,
				TitleTwo:      req.MerName,
				Images:        req.Images,
			},
		}

		subAgent, err := mongo.SubAgentColl.Find(agentUser.SubAgentCode)
		if err == nil {
			merchant.AgentName = subAgent.AgentName
			merchant.SubAgentName = subAgent.SubAgentName
		}

		if err := genMerId(merchant, subAgent.AgentCode+"0"); err != nil {
			w.Write(jsonMarshal(model.SYSTEM_ERROR))
			return
		}
		if err := genRouter(merchant); err != nil {
			w.Write(jsonMarshal(model.SYSTEM_ERROR))
			return
		}
	} else {
		result = User.updateSettInfo(req)
		if result.State != model.SUCCESS {
			w.Write(jsonMarshal(result))
			return
		}
		merchant = req.m
	}

	appUser.MerId = merchant.MerId
	appUser.UniqueId = merchant.UniqueId
	appUser.SignKey = merchant.SignKey
	appUser.AgentCode = merchant.AgentCode
	appUser.SubAgentCode = agentUser.SubAgentCode
	mongo.AppUserCol.Upsert(appUser)
	result.User = appUser
	appUser.Password = "" //不显示

	w.Write(jsonMarshal(result))
}

// UserActivate 用户激活
func UserActivate(w http.ResponseWriter, r *http.Request) {
	// 验证token
	if _, ok := checkAccessToken(r.FormValue("accessToken")); !ok {
		w.Write(jsonMarshal(model.TOKEN_ERROR))
		return
	}

	username := r.FormValue("username")
	appUser, err := mongo.AppUserCol.FindOne(username)
	if err != nil {
		w.Write(jsonMarshal(model.USERNAME_NO_EXIST))
		return
	}

	m, err := mongo.MerchantColl.Find(appUser.MerId)
	if err != nil {
		w.Write(jsonMarshal(model.MERID_NO_EXIST))
		return
	}

	// 设置为已激活
	appUser.Activate = "true"
	err = mongo.AppUserCol.Upsert(appUser)
	if err != nil {
		w.Write(jsonMarshal(model.SYSTEM_ERROR))
		return
	}

	// 异步处理邮件和上传图片
	go func() {

		// 生成支付设计图
		payImage := genImageWithQrCode(payImageTemplate, m.Detail.PayUrl, m.Detail.MerName)

		// 写入缓冲池
		var isEnocdeSuccess = true
		payBuf := bytes.NewBuffer([]byte{})
		err = jpeg.Encode(payBuf, payImage, nil)
		if err != nil {
			log.Errorf("jpeg encode fail: %s", err)
			isEnocdeSuccess = false
		}

		// 发email告知用户已开通成功
		email := &email.Email{To: username, Title: open.Title}
		if isEnocdeSuccess {
			jpg64 := base64.StdEncoding.EncodeToString(payBuf.Bytes())
			image := fmt.Sprintf(`<img src="data:image/jpeg;base64,%s" style=width:213px;height:300px;/>`, jpg64)
			email.Body = fmt.Sprintf(open.Body, m.Detail.MerName, username, m.MerId, m.SignKey, image)
		} else {
			email.Body = fmt.Sprintf(open.Body, m.Detail.MerName, username, m.MerId, m.SignKey, "")
		}
		err = email.Send()
		if err != nil {
			log.Errorf("send email fail: %s, To=%s, body=%s", err, username, email.Body)
		}

		// upload payImage
		if isEnocdeSuccess {
			err = qiniu.Put(fmt.Sprintf(qrImage, m.MerId, "pay"), int64(payBuf.Len()), payBuf)
			if err != nil {
				log.Errorf("fail to upload image : %s", err)
			}
		}

		// 生成账单设计图
		billImage := genImageWithQrCode(billImageTemplate, m.Detail.BillUrl, m.Detail.MerName)
		billBuf := bytes.NewBuffer([]byte{})
		err = jpeg.Encode(billBuf, billImage, nil)
		if err != nil {
			log.Errorf("jpeg encode fail: %s", err)
			isEnocdeSuccess = false
		} else {
			isEnocdeSuccess = true
		}

		// upload billImage
		if isEnocdeSuccess {
			err = qiniu.Put(fmt.Sprintf(qrImage, m.MerId, "bill"), int64(billBuf.Len()), billBuf)
			if err != nil {
				log.Errorf("fail to upload image : %s", err)
			}
		}

	}()

	w.Write(jsonMarshal(model.SUCCESS1))

}

// GetDownloadUrl 生成下载地址
func GetDownloadUrl(w http.ResponseWriter, r *http.Request) {
	// 验证token
	if _, ok := checkAccessToken(r.FormValue("accessToken")); !ok {
		w.Write(jsonMarshal(model.TOKEN_ERROR))
		return
	}

	merId := r.FormValue("merId")
	imageType := r.FormValue("imageType")

	if merId == "" || imageType == "" {
		w.Write(jsonMarshal(model.PARAMS_EMPTY))
		return
	}

	result := model.SUCCESS1
	result.DownloadUrl = qiniu.MakePrivateUrl(fmt.Sprintf(qrImage, merId, imageType))
	w.Write(jsonMarshal(result))
}

// checkAccessToken
func checkAccessToken(token string) (*model.User, bool) {

	if s, ok := tokenMap[token]; ok {
		if time.Now().After(s.Expires) {
			delete(tokenMap, token)
			return nil, false
		}
		return s.User, true
	}
	// 向数据库里查找
	s, err := mongo.SessionColl.Find(token)
	if err != nil {
		return nil, false
	}

	// 放到map里
	tokenMap[token] = s

	return s.User, true
}

func genAccessToken(user *model.User) string {
	s := &model.Session{
		SessionID:  util.SerialNumber(),
		User:       user,
		CreateTime: time.Now(),
	}
	s.UpdateTime = s.CreateTime
	s.Expires = s.CreateTime.Add(24 * time.Hour)
	mongo.SessionColl.Add(s)
	tokenMap[s.SessionID] = s
	return s.SessionID
}

// NotifySalesman 每天汇总当天用户数据给业务人员
func NotifySalesman(em string) {
	day := time.Now().Format("2006-01-02")
	all, err := mongo.AppUserCol.Find(&model.AppUserContiditon{
		RegisterFrom: model.SalesToolsRegister,
		StartTime:    day + " 00:00:00",
		EndTime:      day + " 23:59:59",
	})
	if err != nil {
		log.Errorf("find appUser error:%s", err)
		return
	}

	// 归属
	c := make(map[string][]*model.AppUser)
	for _, u := range all {
		if users, ok := c[u.BelongsTo]; ok {
			users = append(users, u)
			c[u.BelongsTo] = users
		} else {
			c[u.BelongsTo] = []*model.AppUser{u}
		}
	}

	for k, v := range c {
		// 打包成zip
		buf := new(bytes.Buffer)
		w := zip.NewWriter(buf)

		var eds []excelData
		for _, u := range v {
			// log.Debugf("k=%s c=%s u=%s", k, u, v)
			m, err := mongo.MerchantColl.Find(u.MerId)
			if err != nil {
				log.Errorf("fail to find merchant(%s): %s", u.MerId, err)
				continue
			}

			eds = append(eds, excelData{m: m, u: u})

			// 获取图片
			for index, iu := range m.Detail.Images {
				resp, err := http.Get(qiniu.MakePrivateUrl(iu))
				if err != nil {
					log.Errorf("http.get err: %s", err)
					continue
				}
				defer resp.Body.Close()
				bs, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Errorf("read body err: %s", err)
					continue
				}

				// 按商户号排
				name := "/" + m.MerId + "/" + fmt.Sprintf("%d.jpg", index+1)
				f, err := w.Create(name)
				if err != nil {
					log.Errorf("create file err: %s", err)
					continue
				}
				f.Write(bs)
			}
		}
		w.Close()

		if em == "" {
			em = k
		}

		e := email.Email{To: em, Title: "当日商户汇总", Body: "如题，见附件。"}
		e.Attach(buf, "商户.zip", "")

		// 生成excel
		excel := genExcel(eds)
		ebuf := bytes.NewBuffer([]byte{})
		err = excel.Write(ebuf)
		if err == nil {
			e.Attach(ebuf, "汇总表.xlsx", "")
		} else {
			log.Errorf("fail to gen excel: %s", err)
		}

		e.Send()
	}
}

type excelData struct {
	m *model.Merchant
	u *model.AppUser
}

func genExcel(eds []excelData) *xlsx.File {

	var sheet *xlsx.Sheet
	var row *xlsx.Row

	excel := xlsx.NewFile()
	sheet, _ = excel.AddSheet("原始-商户信息表")

	row = sheet.AddRow()

	type rowType struct {
		A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, T, U, V, W, X, Y, Z, AA, AB, AC, AD, AE, AF, AG, AH, AI, AJ, AK, AL, AM, AN, AO, AP, AQ, AR, AS, AT string
	}

	row.WriteStruct(&rowType{"商家营业简称", "公司名称", "注册地址", "营业执照注册号", "经营范围", "营业期限", "注册资本", "预计年收入", "员工人数", "营业场所面积", "证件持有人类型", "证件持有人姓名", "证件类型", "证件号码", "证件有效期限", "组织机构代码", "有效期",
		"商家简称", "售卖商品具体描述", "客服电话", "账户类型", "开户行代码", "开户银行城市", "开户名称", "开户支行", "银行账号", "主要联系人姓名", "主要联系人手机号码", "主要联系人邮箱", "联系地址", "公司传真", "营业执照影印件（资质）", "运营者证件",
		"组织机构代码证（扫描件)", "门店照片", "个户工商户营业执照扫描件", "《餐饮服务许可证》/《食品卫生许可证》", "关注公众服务号(APPID)", "支付宝账户", "申请业务范围", "商家设备数量（台）", "商户号", "商户密钥", "app注册邮箱", "app密码md5值", "收款码链接"}, -1)

	// 填充数据
	for _, ed := range eds {
		row = sheet.AddRow()
		row.WriteStruct(&rowType{
			ed.m.Detail.MerName, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
			ed.m.Detail.MerName, "", "", "个体", ed.m.Detail.BankId, ed.m.Detail.OpenBankName, ed.m.Detail.AcctName, ed.m.Detail.BankName, ed.m.Detail.AcctNum, ed.m.Detail.Contact, ed.m.Detail.ContactTel, "", "", "", "", "",
			"附件形式提供", "附件形式提供", "附件形式提供", "附件形式提供", "", "", "", "", ed.m.MerId, ed.m.SignKey, ed.u.UserName, ed.u.Password, ed.m.Detail.PayUrl,
		}, -1)
	}
	return excel
}
