'use strict';
require('../css/index.css');
require('./vendor/util');


let server = Util.getServer(),
	merchantCode = Util.getUrlParam('merchantCode');
let redirectUri = server + '/scanpay/weChat/oauth2?merchantCode=' + merchantCode + '&showwxpaytitle=1';
window.location.href = redirectUri;
