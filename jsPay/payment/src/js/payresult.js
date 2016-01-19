'use strict';

require('../css/index.css');
require('../css/fail.css');
require('./vendor/util');

let state = Util.getUrlParam('state');
if (state === '0') {
	window.location.replace('result.html?code=' + Util.getUrlParam('code') + '&orderAmount=' + Util.getNormalTxamt(Util.getUrlParam('txamt')));
} else {
	window.location.replace('fail.html?errorDetail=' + Util.getUrlParam('errorDetail'));
}
