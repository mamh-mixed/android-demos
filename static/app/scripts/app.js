(function(document) {
  'use strict';

  var app = document.querySelector('#app');

  app.addEventListener('dom-change', function() {
    console.log('Our app is ready to rock!');
  });

  // See https://github.com/Polymer/polymer/issues/1381
  window.addEventListener('WebComponentsReady', function() {
    // imports are loaded and elements have been registered
    app.checkIsLogin();
  });

  // Main area's paper-scroll-header-panel custom condensing transformation of
  // the appName in the middle-container and the bottom title in the bottom-container.
  // The appName is moved to top and shrunk on condensing. The bottom sub title
  // is shrunk to nothing on condensing.
  // app.paperHeaderTransform = function(e) {
  //   var appName = document.querySelector('#mainToolbar .app-name');
  //   var middleContainer = document.querySelector('#mainToolbar .middle-container');
  //   var bottomContainer = document.querySelector('#mainToolbar .bottom-container');
  //   var detail = e.detail;
  //   var heightDiff = detail.height - detail.condensedHeight;
  //   var yRatio = Math.min(1, detail.y / heightDiff);
  //   var maxMiddleScale = 0.50; // appName max size when condensed. The smaller the number the smaller the condensed size.
  //   var scaleMiddle = Math.max(maxMiddleScale, (heightDiff - detail.y) / (heightDiff / (1 - maxMiddleScale)) + maxMiddleScale);
  //   var scaleBottom = 1 - yRatio;
  //
  //   // Move/translate middleContainer
  //   Polymer.Base.transform('translate3d(0,' + yRatio * 100 + '%,0)', middleContainer);
  //
  //   // Scale bottomContainer and bottom sub title to nothing and back
  //   Polymer.Base.transform('scale(' + scaleBottom + ') translateZ(0)', bottomContainer);
  //
  //   // Scale middleContainer appName
  //   Polymer.Base.transform('scale(' + scaleMiddle + ') translateZ(0)', appName);
  // };

  // Close drawer after menu item is selected if drawerPanel is narrow
  app.onDataRouteClick = function() {
    var drawerPanel = document.querySelector('#paperDrawerPanel');
    if (drawerPanel.narrow) {
      drawerPanel.closeDrawer();
    }
  };

  // Scroll page to top and expand header
  app.scrollPageToTop = function() {
    document.getElementById('mainContainer').scrollTop = 0;
  };

  // app.route = {
  //     path: "/system"
  // };
  // 初始化菜单
  app.adminMenus = [{
    "nameCN": "系统管理",
    "nameEN": "",
    "icon": "bug-report",
    "url": "",
    "route": "/system",
    "children": [{
      "nameCN": "用户",
      "nameEN": "",
      "icon": "account-box",
      "url": "#/system/users",
    }]
  }, {
    "nameCN": "扫码支付",
    "nameEN": "",
    "icon": "build",
    "url": "",
    "route": "/config",
    "children": [{
      "nameCN": "代理",
      "nameEN": "",
      "icon": "assignment-ind",
      "url": "#/config/agents",
      "route": "config.agents"
    }, {
      "nameCN": "公司",
      "nameEN": "",
      "icon": "account-circle",
      "url": "#/config/subAgents",
      "route": "config.subAgents"
    }, {
      "nameCN": "商户",
      "nameEN": "",
      "icon": "home",
      "url": "#/config/groups",
      "route": "config.groups"
    }, {
      "nameCN": "门店录入",
      "nameEN": "",
      "icon": "create",
      "url": "#/config/configMer",
      "route": "config.configMer"
    }, {
      "nameCN": "批量导入",
      "nameEN": "",
      "icon": "cloud-upload",
      "url": "#/config/upload",
      "route": "config.upload"
    }, {
      "nameCN": "门店",
      "nameEN": "",
      "icon": "store",
      "url": "#/config/merchants",
      "route": "config.merchants"
    }, {
      "nameCN": "渠道商户",
      "nameEN": "",
      "icon": "account-balance",
      "url": "#/config/chanMers",
      "route": "config.chanMers"
    }, {
      "nameCN": "路由策略",
      "nameEN": "",
      "icon": "settings-input-component",
      "url": "#/config/routes",
      "route": "config.routes"
    }]
  }, {
    "nameCN": "快捷支付",
    "nameEN": "",
    "icon": "swap-horiz",
    "url": "",
    "route": "/quickPay",
    "children": [{
      "nameCN": "机构商户",
      "nameEN": "",
      "icon": "store",
      "url": "#/quickPay/merchants",
      "route": "quickPay.merchants"
    }, {
      "nameCN": "渠道商户",
      "nameEN": "",
      "icon": "account-balance",
      "url": "#/quickPay/chanMers",
      "route": "quickPay.chanMers"
    }, {
      "nameCN": "路由策略",
      "nameEN": "",
      "icon": "settings-input-component",
      "url": "#/quickPay/routes",
      "route": "quickPay.routes"
    }]
  }, {
    "nameCN": "交易查询",
    "nameEN": "",
    "icon": "swap-horiz",
    "url": "",
    "route": "/trade",
    "children": [{
      "nameCN": "扫码交易",
      "nameEN": "",
      "icon": "payment",
      "url": "#/trade/list",
      "route": "trade.list"
    }, {
      "nameCN": "绑定交易",
      "nameEN": "",
      "icon": "payment",
      "url": "#/trade/bindingpay/list",
      "route": "trade.bindingpay.list"
    }, {
      "nameCN": "报表下载",
      "nameEN": "",
      "icon": "file-download",
      "url": "#/trade/reports",
      "route": "trade.reports"
    }, {
      "nameCN": "旧系统报表下载",
      "nameEN": "",
      "icon": "cloud-download",
      "url": "#/trade/old/reports",
      "route": "trade.old.reports"
    }]
  }, {
    "nameCN": "内部测试",
    "nameEN": "",
    "icon": "bug-report",
    "url": "",
    "route": "/test",
    "children": [{
      "nameCN": "扫码支付",
      "nameEN": "",
      "icon": "pageview",
      "url": "#/test/scanpay",
      "route": "test.scanpay"
    }, {
      "nameCN": "绑定支付",
      "nameEN": "",
      "icon": "icons:speaker-notes",
      "url": "#/test/bindingpay",
      "route": "test.bindingpay"
    }, {
      "nameCN": "七牛 Ajax 上传文件",
      "nameEN": "",
      "icon": "cloud-upload",
      "url": "#/test/qiniu",
      "route": "test.qiniu"
    }, {
      "nameCN": "HTML5读取文件内容",
      "nameEN": "",
      "icon": "get-app",
      "url": "#/test/filereader",
      "route": "test.filereader"
    }]
  }];
  app.agentMenus = [{
    "nameCN": "交易查询",
    "nameEN": "",
    "icon": "swap-horiz",
    "url": "",
    "route": "trade",
    "children": [{
      "nameCN": "扫码交易",
      "nameEN": "",
      "icon": "payment",
      "url": "#/trade/list",
      "route": "trade.list"
    }, {
      "nameCN": "报表下载",
      "nameEN": "",
      "icon": "file-download",
      "url": "#/trade/reports",
      "route": "trade.reports"
    }]
  }];

  app.showLoginDialog = false;
  // 核对是否登入
  app.checkIsLogin = function() {
    this.$.findSessionAjax.generateRequest();
  };
  app.handleFindSessionError = function(e) {
    this.showLoginDialog = true;
    this.showMenu = true;
  };
  app.handleFindSessionResponse = function(e) {
    var response = e.detail.response;

    // 如果没有登录
    if (response.status !== 0) {
      this.showLoginDialog = true;
      this.showMenu = true;
      return
    }

    // 登录了，但当前页不是登录页，渲染菜单
    this.user = response.data;
    this.userType = this.user.userType;
    if (this.userType === 'admin') {
      this.menus = this.adminMenus;
    } else if (this.userType === 'agent' || this.userType === 'group' || this.userType === 'merchant') {
      this.menus = this.agentMenus;
    }
  };
  // 当没有登入的时候，跳转到登入页面
  app.handleNotLogin = function() {
    document.querySelector('#appRouter').go('/');
    return;
  };
  app.startWith = function(path, prefix) {
    cosole.log(path, prefix);
    if (!path || !prefix) {
      return false;
    }
    return path.substr(0, prefix.length) === prefix;
  };
  app.handleDeleteSessionResponse = function() {};
  // 退出
  app.logoutHandle = function() {
    // 清除客户端 SessionID
    document.cookie = 'QUICKMASTERID=; path=/master/; expires=Thu, 01 Jan 1970 00:00:01 GMT;';

    // 删除后台sesseion信息
    this.$.deleteSessionAjax.params = {
      'sessionId': this.sessionId,
    };
    this.$.deleteSessionAjax.generateRequest();

    // 清除localStorage中的信息
    this.sessionId = '';
    this.userType = '';
    this.user = {};
    this.$.logoutBtn.hidden = true;

    // 跳转至登录页面
    document.querySelector('#appRouter').go('/');
  };
  // 进入全屏模式
  app.enterFullscreenMode = function() {
    this.isFullscreen = true;
  };
  // 退出全屏模式
  app.exitFullscreenMode = function() {
    this.isFullscreen = false;
  };
  app.paperResponsiveChange = function(e, detail) {
    if (!this.isFullscreen) {
      this.isNarrow = detail.narrow;
    }
  };


})(document);
