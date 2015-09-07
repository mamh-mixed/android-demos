/*
Copyright (c) 2015 The Polymer Project Authors. All rights reserved.
This code may only be used under the BSD style license found at http://polymer.github.io/LICENSE.txt
The complete set of authors may be found at http://polymer.github.io/AUTHORS.txt
The complete set of contributors may be found at http://polymer.github.io/CONTRIBUTORS.txt
Code distributed by Google as part of the polymer project is also
subject to an additional IP rights grant found at http://polymer.github.io/PATENTS.txt
*/

(function(document) {
	'use strict';

	// Grab a reference to our auto-binding template
	// and give it some initial binding values
	// Learn more about auto-binding templates at http://goo.gl/Dx1u2g
	var app = document.querySelector('#app');

	app.displayInstalledToast = function() {
		document.querySelector('#caching-complete').show();
	};

	// Listen for template bound event to know when bindings
	// have resolved and content has been stamped to the page
	app.addEventListener('dom-change', function() {
		console.log('Our app is ready to rock!');
		window.localStorage.setItem("userType","admin");
		// 获取userType
		var userType = window.localStorage.getItem("userType");
		// userType为空跳转至登陆页面
		// if (!userType||userType===''){
		// 	window.location.href='http://www.baidu.com';
		// }
		if (userType==='admin'){
			document.querySelector('#menuAjax').generateRequest();
		}else if(userType==='agent'||userType==='group'||userType==='merchant'){
			document.querySelector('#agentAjax').generateRequest();
		}
	});


	// See https://github.com/Polymer/polymer/issues/1381
	window.addEventListener('WebComponentsReady', function() {
		// imports are loaded and elements have been registered
		console.log('ready');
	});

	// 主页面滚动的时候，顶部的title的动画效果
	addEventListener('paper-header-transform', function(e) {
		var appName = document.querySelector('.app-name');
		var middleContainer = document.querySelector('.middle-container');
		var bottomContainer = document.querySelector('.bottom-container');
		var detail = e.detail;
		var heightDiff = detail.height - detail.condensedHeight;
		var yRatio = Math.min(1, detail.y / heightDiff);
		var maxMiddleScale = 0.50; // appName max size when condensed. The smaller the number the smaller the condensed size.
		var scaleMiddle = Math.max(maxMiddleScale, (heightDiff - detail.y) / (heightDiff / (1 - maxMiddleScale)) + maxMiddleScale);
		var scaleBottom = 1 - yRatio;

		// Move/translate middleContainer
		Polymer.Base.transform('translate3d(0,' + yRatio * 100 + '%,0)', middleContainer);

		// Scale bottomContainer and bottom sub title to nothing and back
		Polymer.Base.transform('scale(' + scaleBottom + ') translateZ(0)', bottomContainer);

		// Scale middleContainer appName
		Polymer.Base.transform('scale(' + scaleMiddle + ') translateZ(0)', appName);
	});

	// Close drawer after menu item is selected if drawerPanel is narrow
	app.onMenuSelect = function() {
		var mainContainer = document.querySelector('#mainContainer');
		mainContainer.scrollTop = 140;
		var drawerPanel = document.querySelector('#paperDrawerPanel');
		if (drawerPanel.narrow) {
			drawerPanel.closeDrawer();
		}
	};

	app.startWith = function(name, prefix) {
		if (!name || !prefix) {
			return true;
		}
		return name.substring(0, prefix.length) === prefix;
	};

	// 菜单列表加载
	app.handleQueryResponse = function(e) {
		var response = e.detail.response;
		app.menus = response.slice(0);
	};

})(document);
