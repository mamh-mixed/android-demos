'use strict';

require('../css/index.css');
require('../css/fail.css');
require('./vendor/util');
var $ = require('webpack-zepto');

document.title = window.localStorage.getItem('title_one');
(function() {
	$(document).ready(() => {
		alert(Util.getUrlParam('errorDetail'));
    $('#returnBtn').on('click', () => {
			WeixinJSBridge.call('closeWindow');
		});
	});
})($);
