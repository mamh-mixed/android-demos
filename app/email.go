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
			%s，申请提升限额。邮箱：%s，手机：%s，商户号：%s。
		</body>
	</html>
	`

	openTemplate = `
	<html>
		<body>
			您好，您申请的商户参数信息如下：
			<h3>%s</h3>
			<h4>app登录信息：</h4>
			注册邮箱：%s</br>
			<h4>桌面版信息：</h4>
			商户号： 	%s </br>
			密钥：	%s
			<h4>网页版信息：</h4>
			<div id="code">
			%s
			</div>
		</body>
	</html>
	`

	resetTemplate = `
	<html>
		<body>
			<h3>
				点击以下链接重置密码</br>
				<a href="%s">%s</a>
			</h3>
		</body>
	</html>
	`
)

// 邮件模板消息体
var (
	activation    = emailTemplate{Title: "云收银帐号激活", Body: actTemplate}
	promote       = emailTemplate{Title: "申请提升限额", Body: promoteTemplate}
	open          = emailTemplate{Title: "【感谢您注册云收银】", Body: openTemplate}
	resetPassword = emailTemplate{Title: "重置密码", Body: resetTemplate}
)

type emailTemplate struct {
	Title string
	Body  string
}
