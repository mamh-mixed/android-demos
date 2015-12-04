'use strict';

require('../css/index.css');
require('../css/fail.css');
var $ = require('webpack-zepto');

document.title = window.localStorage.getItem('title_one');
(function() {
	$(document).ready(() => {
    $('#returnBtn').on('click', () => {
			WeixinJSBridge.call('closeWindow');
		});
	});
})($);
