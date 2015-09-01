var UA = window.navigator.userAgent,
    CLICK = 'click';
if (/ipad|iphone|android/.test(UA.toLowerCase())) {
    CLICK = 'click';
}
// Avoid `console` errors in browsers that lack a console.
(function() {
    var method;
    var noop = function() {};
    var methods = [
        'assert', 'clear', 'count', 'debug', 'dir', 'dirxml', 'error',
        'exception', 'group', 'groupCollapsed', 'groupEnd', 'info', 'log',
        'markTimeline', 'profile', 'profileEnd', 'table', 'time', 'timeEnd',
        'timeline', 'timelineEnd', 'timeStamp', 'trace', 'warn'
    ];
    var length = methods.length;
    var console = (window.console = window.console || {});

    while (length--) {
        method = methods[length];

        // Only stub undefined methods.
        if (!console[method]) {
            console[method] = noop;
        }
    }
}());

// Place any jQuery/helper plugins in here.
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
var Util = function() {
    var USERNAME = 'CLOUD_CASHIER_USERNAME';
    var PASSWORD = 'CLOUD_CASHIER_PASSWORD';
    var ISAUTOLOGIN = 'CLOUD_CASHIER_IS_AUTO_LOGIN';
    var DEVICE = 'CLOUD_CASHIER_DEVICE';
    var KEY = 'CLOUD_CASHIER_KEY';
    var CLIENTID = 'CLOUD_CASHIER_CLIENT_ID';
    // 验证手机号的方法
    var validatemobile = function(mobile) {
        if (mobile.length === 0) {
            window.alert('请输入手机号码！');
            return false;
        }
        if (mobile.length !== 11) {
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

    // 验证银行账号的方法
    var luhmCheck = function(bankno) {
        var lastNum = bankno.substr(bankno.length - 1, 1); //取出最后一位（与luhm进行比较）

        var first15Num = bankno.substr(0, bankno.length - 1); //前15或18位
        var newArr = new Array();
        for (var i = first15Num.length - 1; i > -1; i--) { //前15或18位倒序存进数组
            newArr.push(first15Num.substr(i, 1));
        }
        var arrJiShu = new Array(); //奇数位*2的积 <9
        var arrJiShu2 = new Array(); //奇数位*2的积 >9

        var arrOuShu = new Array(); //偶数位数组
        for (var j = 0; j < newArr.length; j++) {
            if ((j + 1) % 2 == 1) { //奇数位
                if (parseInt(newArr[j]) * 2 < 9)
                    arrJiShu.push(parseInt(newArr[j]) * 2);
                else
                    arrJiShu2.push(parseInt(newArr[j]) * 2);
            } else //偶数位
                arrOuShu.push(newArr[j]);
        }

        var jishu_child1 = new Array(); //奇数位*2 >9 的分割之后的数组个位数
        var jishu_child2 = new Array(); //奇数位*2 >9 的分割之后的数组十位数
        for (var h = 0; h < arrJiShu2.length; h++) {
            jishu_child1.push(parseInt(arrJiShu2[h]) % 10);
            jishu_child2.push(parseInt(arrJiShu2[h]) / 10);
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

        for (var p = 0; p < jishu_child1.length; p++) {
            sumJiShuChild1 = sumJiShuChild1 + parseInt(jishu_child1[p]);
            sumJiShuChild2 = sumJiShuChild2 + parseInt(jishu_child2[p]);
        }
        //计算总和
        sumTotal = parseInt(sumJiShu) + parseInt(sumOuShu) + parseInt(sumJiShuChild1) + parseInt(sumJiShuChild2);

        //计算Luhm值
        var k = parseInt(sumTotal) % 10 == 0 ? 10 : parseInt(sumTotal) % 10;
        var luhm = 10 - k;

        if (lastNum == luhm) {

            return true;
        } else {

            return false;
        }
    }
    // 验证邮箱的方法
    var isEmail = function(str) {
        var reg = /^([a-zA-Z0-9._-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+/;
        return reg.test(str);
    };

    var _getUrlParam = function(key) {
        var reg = new RegExp("(^|&)" + key + "=([^&]*)(&|$)", "i");
        var r = window.location.search.substr(1).match(reg);
        if (r !== null) return window.unescape(r[2]);
        return null;
    };

    var getTxamt = function(txamt) {
        var str = txamt;
        var i = parseFloat(str);
        var j = i.toFixed(2);
        j = j * 100;
        var num = j;
        str = "" + num;
        var k = 12 - str.length;
        var sum = "";
        for (var l = 0; l < k; l++) {
            sum = sum + "0";
        }
        sum = sum + str;
        return sum;
    };

    var getNormalTxamt = function(txamt) {
        var str = txamt;
        if (str !== undefined) {
            var sum = "";
            var index = 0;
            var c = str.charCodeAt(index);
            while (c === '0') {
                index++;
                c = str.charCodeAt(index);
            }
            sum = str.substring(index);
            var i = parseFloat(sum);
            i = i / 100;
            var j = i.toFixed(2);
            sum = "" + j;
            return sum;

        }
    };
    // 向手机发送请求，这是js和native app 交互的方法
    var sendRequestToDevice = function(action, jsonData) {
        var device = getDevice();
        if (device === '') {
            var UA = window.navigator.userAgent;
            if (/ipad|iphone/.test(UA.toLowerCase())) {
                device = 'ios';
            } else {
                device = 'android';
            }
        }

        if (device === 'android') {
            window.couldCashier.OnDo(action, jsonData);
        } else if (device === 'ios') {
            document.location = 'cloudcashier://' + action + ':/' + jsonData;
        }
    };
    var saveToLocalStorage = function(key, value) {
        window.localStorage.setItem(key, value);
    };
    var setUsername = function(value) {
        saveToLocalStorage(USERNAME, value);
    };
    var getUsername = function() {
        return window.localStorage.getItem(USERNAME) || '';
    };
    var setPassword = function(value) {
        saveToLocalStorage(PASSWORD, value);
    };
    var getPassword = function() {
        return window.localStorage.getItem(PASSWORD) || '';
    };
    var setIsAutoLogin = function(value) {
        saveToLocalStorage(ISAUTOLOGIN, value);
    };
    var getIsAutoLogin = function() {
        var tem = window.localStorage.getItem(ISAUTOLOGIN);
        return tem === 'true';
    };
    var setDevice = function(value) {
        saveToLocalStorage(DEVICE, value);
    };
    var getDevice = function() {
        return window.localStorage.getItem(DEVICE) || '';
    };
    var setKey = function(value) {
        saveToLocalStorage(KEY, value);
    };
    var getKey = function() {
        return window.localStorage.getItem(KEY) || '';
    };
    var setClientId = function(value) {
        saveToLocalStorage(CLIENTID, value);
    };
    var getClientId = function() {
        return window.localStorage.getItem(CLIENTID) || '';
    };
    return {
        getUrlParam: _getUrlParam,
        isEmail: isEmail,
        luhmCheck: luhmCheck,
        validatemobile: validatemobile,
        getTxamt: getTxamt,
        getNormalTxamt: getNormalTxamt,
        sendRequestToDevice: sendRequestToDevice,
        saveToLocalStorage: saveToLocalStorage,
        setUsername: setUsername,
        getUsername: getUsername,
        setPassword: setPassword,
        getPassword: getPassword,
        setIsAutoLogin: setIsAutoLogin,
        getIsAutoLogin: getIsAutoLogin,
        setDevice: setDevice,
        getDevice: getDevice,
        setKey: setKey,
        getKey: getKey,
        setClientId: setClientId,
        getClientId: getClientId
    };
}(window);

// ajax请求链接
//var UrlMap = {
//    getInfo: 'http://211.147.72.70:10003/getinfo',
//    bill: 'http://211.147.72.70:10003/bill',
//    getOrder: 'http://211.147.72.70:10003/getOrder'
//};

// ajax请求链接 生产
var UrlMap = {
    getInfo: 'http://211.144.213.118/getinfo',
    bill: 'http://211.144.213.118/bill',
    getOrder: 'http://211.144.213.118/getOrder'
};
