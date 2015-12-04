/*
Copyright (c) 2010,2011,2012,2013,2014,2015 CardInfoLink http://www.show.money
License: MIT - http://mrgnrdrck.mit-license.org
*/
(function(root, factory) {
	'use strict';

	if (typeof define === 'function' && define.amd) {
		// AMD. Register as an anonymous module.
		define(['exports'], factory);

	} else if (typeof exports === 'object') {
		// CommonJS
		factory(exports);

	}

	// Browser globals
	var Util = {};
	root.Util = Util;
	factory(Util);

}((typeof window === 'object' && window) || this, function(Util) {
	'use strict';
	Date.prototype.Format = function(formatStr) {
		var str = formatStr;
		str = str.replace(/yyyy|YYYY/, this.getFullYear());
		str = str.replace(/yy|YY/, (this.getYear() % 100) > 9 ? (this.getYear() % 100).toString() : '0' + (this.getYear() % 100));
		var month = this.getMonth() + 1;
		str = str.replace(/MM/, month > 9 ? month.toString() : '0' + month);
		str = str.replace(/M/g, month);
		str = str.replace(/dd|DD/, this.getDate() > 9 ? this.getDate().toString() : '0' + this.getDate());
		str = str.replace(/d|D/g, this.getDate());

		str = str.replace(/hh|HH/, this.getHours() > 9 ? this.getHours().toString() : '0' + this.getHours());
		str = str.replace(/h|H/g, this.getHours());
		str = str.replace(/mm/, this.getMinutes() > 9 ? this.getMinutes().toString() : '0' + this.getMinutes());
		str = str.replace(/m/g, this.getMinutes());

		str = str.replace(/ss|SS/, this.getSeconds() > 9 ? this.getSeconds().toString() : '0' + this.getSeconds());
		str = str.replace(/s|S/g, this.getSeconds());
		return str;
	};
	/**
	 *
	 */
	Util.isEmail = (email) => {
		var reg = /^([a-zA-Z0-9._-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+/;
		return reg.test(email);
	};

	Util.getUrlParam = (key) => {
		var reg = new RegExp('(^|&)' + key + '=([^&]*)(&|$)', 'i');
		var r = window.location.search.substr(1).match(reg);
		if (r !== null) return window.unescape(r[2]);
		return null;
	};

	Util.validatemobile = (mobile) => {
		if (mobile.length === 0) {
			window.alert('请输入手机号码！');

			return false;
		}
		if (mobile.length != 11) {
			window.alert('请输入有效的手机号码！');

			return false;
		}

		var myreg = /^0?1[3|4|5|8][0-9]\d{8}$/;
		if (!myreg.test(mobile)) {
			window.alert('请输入有效的手机号码！');

			return false;
		}
		return true;
	};

	Util.getTxamt = (txamt) => {
		var str = txamt;
		var i = parseFloat(str);
		var j = i.toFixed(2);
		j = j * 100;
		var num = j;
		str = '' + num;
		var k = 12 - str.length;
		var sum = '';
		for (var l = 0; l < k; l++) {
			sum = sum + '0';
		}
		sum = sum + str;
		return sum;
	};

	Util.getNormalTxamt = (txamt) => {
		var str = txamt;
		if (str !== undefined) {
			var sum = '';
			var index = 0;
			var c = str.charCodeAt(index);
			while (c == '0') {
				index++;
				c = str.charCodeAt(index);
			}
			sum = str.substring(index);
			var i = parseFloat(sum);
			i = i / 100;
			var j = i.toFixed(2);
			sum = '' + j;
			return sum;

		}
	};

	Util.getServer = () => {
		// 测试环境地址
		var server = 'http://test.quick.ipay.so';
		// 生产环境地址
		//  var server='https://api.shou.money';
		// 开发环境地址
		// var server = 'http://192.168.199.193:6800';
		return server;
	};

	Util.luhmCheck = (bankno) => {
		var lastNum = bankno.substr(bankno.length - 1, 1); //取出最后一位（与luhm进行比较）

		var first15Num = bankno.substr(0, bankno.length - 1); //前15或18位
		var newArr = [];
		for (var i = first15Num.length - 1; i > -1; i--) { //前15或18位倒序存进数组
			newArr.push(first15Num.substr(i, 1));
		}
		var arrJiShu = []; //奇数位*2的积 <9
		var arrJiShu2 = []; //奇数位*2的积 >9

		var arrOuShu = []; //偶数位数组
		for (var j = 0; j < newArr.length; j++) {
			if ((j + 1) % 2 == 1) { //奇数位
				if (parseInt(newArr[j]) * 2 < 9)
					arrJiShu.push(parseInt(newArr[j]) * 2);
				else
					arrJiShu2.push(parseInt(newArr[j]) * 2);
			} else //偶数位
				arrOuShu.push(newArr[j]);
		}

		var jishuChild1 = []; //奇数位*2 >9 的分割之后的数组个位数
		var jishuChild2 = []; //奇数位*2 >9 的分割之后的数组十位数
		for (var h = 0; h < arrJiShu2.length; h++) {
			jishuChild1.push(parseInt(arrJiShu2[h]) % 10);
			jishuChild2.push(parseInt(arrJiShu2[h]) / 10);
		}

		var sumJiShu = 0; //奇数位*2 < 9 的数组之和
		var sumOuShu = 0; //偶数位数组之和
		var sumJiShuChild1 = 0; //奇数位*2 >9 的分割之后的数组个位数之和
		var sumJiShuChild2 = 0; //奇数位*2 >9 的分割之后的数组十位数之和
		var sumTotal = 0;
		for (var m = 0; m < arrJiShu.length; m++) {
			sumJiShu = sumJiShu + parseInt(arrJiShu[m]);
		}

		for (var n = 0; n < arrOuShu.length; n++) {
			sumOuShu = sumOuShu + parseInt(arrOuShu[n]);
		}

		for (var p = 0; p < jishuChild1.length; p++) {
			sumJiShuChild1 = sumJiShuChild1 + parseInt(jishuChild1[p]);
			sumJiShuChild2 = sumJiShuChild2 + parseInt(jishuChild2[p]);
		}
		//计算总和
		sumTotal = parseInt(sumJiShu) + parseInt(sumOuShu) + parseInt(sumJiShuChild1) + parseInt(sumJiShuChild2);

		//计算Luhm值
		var k = parseInt(sumTotal) % 10 === 0 ? 10 : parseInt(sumTotal) % 10;
		var luhm = 10 - k;

		if (lastNum == luhm) {
			return true;
		} else {
			return false;
		}
	};
}));
