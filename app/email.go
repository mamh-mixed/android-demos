package app

var (
	andyLi      = "andy.li@cardinfolink.com"
	actTemplate = `	
	<html>
		<body>
			<h3>
				点击以下链接以激活账户</br>
				<a href="%s">%s</a>
			</h3>
		</body>
	</html>
	`

	promoteTemplate = `
	<html>
		<body>
			<h3>
				Hello,Andy.li:
			</h3>
			%s，申请提升限额。邮箱：%s，手机：%s。
		</body>
	</html>
	`
)

// 邮件模板消息体
var (
	activation = emailTemplate{Title: "云收银帐号激活", Body: actTemplate}
	promote    = emailTemplate{Title: "申请提升限额", Body: promoteTemplate}
)

type emailTemplate struct {
	Title string
	Body  string
}
