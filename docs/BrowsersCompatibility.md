# 浏览器兼容性
![aaa](PerformanceNowNotSupportInChrome22.jpg)

以上图为例，图中的浏览器是`Chrome 22`，报错的函数是`Performance.now()`。<br>该函数的浏览器兼容性如下表所列，详情[戳此](https://developer.mozilla.org/en-US/docs/Web/API/Performance/now#Browser_compatibility)。

Feature       | Chrome | Firefox(Gecko) | IE   | Opera | Safari
:------------ | -----: | -------------: | ---: | ----: | -----:
Basic support |   24.0 |           15.0 | 10.0 |  15.0 |    8.0

管理平台使用的框架是`Polymer`，原生支持该框架的目前只有`Chrome`。但是引用一个`Web Components`的polyfills，可以使`Chrome`之外的大部分浏览器支持。<br>详情如下。

Polyfill        | IE10 | IE11+ | Chrome* | Firefox* | Safari 7+* | Chrome Android* | Mobile Safari*
--------------- | :--: | :---: | :-----: | :------: | :--------: | :-------------: | :------------:
Custom Elements | ~    | ✓     | ✓       | ✓        | ✓          | ✓               | ✓
HTML Imports    | ~    | ✓     | ✓       | ✓        | ✓          | ✓               | ✓
Shadow DOM      | ✓    | ✓     | ✓       | ✓        | ✓          | ✓               | ✓
Templates       | ✓    | ✓     | ✓       | ✓        | ✓          | ✓               | ✓

综合以上两张表格，浏览器兼容要求如下

浏览器            | 最低版本要求
:------------- | -----:
Chrome         |   24.0
Firefox(Gecko) |   15.0
Safari         |    8.0
Opera          |   15.0
IE             |    不要求
