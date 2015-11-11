## 多版本维护
### 云收银

大版本为2.0。

发布的下载链接：
http://download.cardinfolink.net/android/showmoney.apk

### 云收银银讯版 vs. 云收银

* LOGO
* App Name
* 友盟分发渠道为yinxun

注1：产品经理要求，两个AppId相同。因此，两个App不能同时安装在一台Android设备上。

注2：共享友盟key
### 云收银日文版 vs. 云收银

* Login页面没有了注册按钮
* 计算器主面板scanCodeView减少了很多按钮
* 菜单精简为4个

发布的下载链接（可修改）：
http://qrcode.cardinfolink.net/app/download/showmoney_jp-release.apk

注1：独立的友盟key

注2：3个应用使用同一个key签名。

## 使用的第三方UI代码：

* 侧滑菜单 https://github.com/jfeinstein10/SlidingMenu
* 下拉刷新 https://github.com/chrisbanes/Android-PullToRefresh

## 约定

* 所有项目使用Android Studio开发。
* 所有项目（包括使用的第三方库和SDK）使用SDK 23编译，最低兼容到8。
