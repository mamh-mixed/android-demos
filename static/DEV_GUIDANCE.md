#开发指南
---
##项目运行
* 首先得安装Node.js，然后当你拉下代码后运行以下命令：

```sh
npm install -g gulp bower && npm install && bower install
```

如果出现错误，可能是没有权限，再执行npm命令之前加上sudo。

* 根据不同需求在项目的根目录下执行以下命令：

```sh
#开发的时候
gulp serve

#本地查看
gulp serve:dist

#发布
gulp
```
发布的文件在```dist/```目录下。

##开发流程

该项目使用[Yeoman](http://yeoman.io/)构建的。使用的生成器(Generator)是[generator-polymer](https://github.com/yeoman/generator-polymer)。所以按照以下命令安装```yeoman```和```generator-polymer```：

```sh
npm install -g yo generator-polymer
```
####可用的生成器

- [polymer (aka polymer:app)](#app)
- [polymer:element](#element-alias-el)
- [polymer:seed](#seed)
- [polymer:gh](#gh)

**注意: 生成器只能在项目根目录下执行**

### App
创建一个新的Polymer应用，生成所有项目需要的模板。

例子:
```bash
yo polymer
```

### Element (别名: El)
在`app/elements`目录下生成一个Polymer元素并且可选择是否把它加到`app/elements/elements.html`中的import列表。

例子:
```bash
yo polymer:element my-element

# 或者使用别名

yo polymer:el my-element
```

**注意: 必须传递一个元素的名称，并且元素的名称必须包含破折号 "-"**

生成元素可以包括需要导入的依赖。例如，如果你想创建一个`fancy-menu`的元素，这个元素需要导入 `paper-button` 和 `paper-checkbox` 作为依赖，你可以照下面例子生成这个文件。

```bash
yo polymer:el fancy-menu paper-button paper-checkbox
```

#### Options

```
--docs, 在你的元素中包含 iron-component-page 文档和 demo.html
--path, 替换默认目录结构, ex: --path foo/bar 将会把你的文件放置在 app/elements/foo/bar
```

### Seed
创建一个基于[seed-element workflow](https://github.com/polymerelements/seed-element)可复用的polymer元素。**这个生成器只能用于创建一个想要通过bower分享给别人的Polymer元素库。**在Polymer应用中不使用。

如果要预览创建的Polymer seed元素，你一使用[polyserve](https://github.com/PolymerLabs/polyserve)工具。

例子:
```bash
mkdir -p my-foo && cd $_
yo polymer:seed my-foo
polyserve
```

### Gh
为[seed-element](#seed)创建一个Github展示页面。

例子:
```bash
cd my-foo
yo polymer:gh
```

##代码

#####Util.toast(text, duration)
弹窗显示信息

```javascript
Util.toast('Hello world!', 100);
```
