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
		尊敬的用户：
		<br>
		<br>
		您好！欢迎您使用云收银 APP！
		<br>
		<br>
		您只需点击下方的链接，即可执行密码重置的请求：
		<br>
		<a href="%s">%s</a>

		<br>
		<br>
		如果以上链接无效，请复制此网址，并将其粘贴到新的浏览器窗口中。
		<br>
		<br>
		感谢您使用云收银 APP，祝您使用愉快！
		<br>
		<br>
		这只是一封系统发送的邮件。我们并不回答对此邮件的回复。
		</body>
	</html>
	`
)

// 邮件模板消息体
var (
	activation    = emailTemplate{Title: "云收银帐号激活", Body: actTemplate}
	promote       = emailTemplate{Title: "申请提升限额", Body: promoteTemplate}
	open          = emailTemplate{Title: "【感谢您注册云收银】", Body: openTemplate}
	resetPassword = emailTemplate{Title: "云收银客户端密码重置", Body: resetTemplate}
)

type emailTemplate struct {
	Title string
	Body  string
}
