var updatepassword = function() {
    var oldpwd = document.getElementById('oldpwd').value;
    var newpwd1 = document.getElementById('newpwd1').value;
    var newpwd2 = document.getElementById('newpwd2').value;

    if (oldpwd.length === 0) {
        alert("原密码不能为空!");
        return;
    }

    if (newpwd1.length === 0) {
        alert("新密码不能为空!");
        return;
    }

    if (newpwd1.length < 6) {
        alert("密码长度不能小于六位!");
        return;
    }

    if (newpwd1 != newpwd2) {
        alert("确认密码不一致!");
        return;
    }


    var device = Util._getUrlParam("device");
    var json = {
        "oldpwd": oldpwd,
        "newpwd": newpwd1
    };

    var jsonStr = JSON.stringify(json);
    if (device == "android") {
        window.couldCashier.OnDo("updatepassword", jsonStr);

    } else if (device == "ios") {
        var result = "couldcashier://updatepassword:/" + jsonStr;
        document.location = result;
    }


}


var getaccount = function() {

    var XMLHttpReq, resultURL;
    var url = "http://211.147.72.70:10003/getinfo";
    var password = Util._getUrlParam("password");
    var username = Util._getUrlParam("username");
    var key = Util._getUrlParam("key");
    var clientid = Util._getUrlParam("clientid");

    var now = new Date();
    var transtime = now.Format("YYYYMMddHHmmss");
    password = hex_md5(password);

    var signatureStr = "clientid=" + clientid +
        "&password=" + password +
        "&transtime=" + transtime +
        "&username=" + username + key;


    var sign = hex_sha1(signatureStr);

    var data = "clientid=" + clientid +
        "&password=" + password +
        "&transtime=" + transtime +
        "&username=" + username +
        "&sign=" + sign;

    try {
        XMLHttpReq = new ActiveXObject("Msxml2.XMLHTTP");
    } catch (E) {
        try {
            XMLHttpReq = new ActiveXObject("Microsoft.XMLHTTP");
        } catch (E) {
            XMLHttpReq = new XMLHttpRequest();
        }
    }

    XMLHttpReq.open("post", url, false);
    XMLHttpReq.onreadystatechange = function() {
        if (XMLHttpReq.readyState == 4) {
            if (XMLHttpReq.status == 200) {
                var text = XMLHttpReq.responseText;
                var jsonobj = eval('(' + text + ')');
                var info = jsonobj.info;
                var bank_open = info.bank_open;
                var payee = info.payee;
                var payee_card = info.payee_card;
                var phone_num = info.phone_num;
                document.getElementById('bank_open').value = bank_open;
                document.getElementById('payee').value = payee;
                document.getElementById('payee_card').value = payee_card;
                document.getElementById('phone_num').value = phone_num;

            }
        }
    };
    XMLHttpReq.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    XMLHttpReq.send(data);
}

var updateaccount = function() {
    var bank_open = document.getElementById('bank_open').value;
    var payee = document.getElementById('payee').value;
    var payee_card = document.getElementById('payee_card').value;
    var phone_num = document.getElementById('phone_num').value;

    if (bank_open.length === 0) {
        alert("开户行不能为空!");
        return;
    }

    if (payee.length === 0) {
        alert("姓名不能为空!");
        return;
    }

    if (payee_card.length === 0) {
        alert("银行卡号不能为空!");
        return;
    }

    if (!Util.luhmCheck(payee_card)) {
        alert("银行卡号验证未通过!");
        return;
    }

    if (!Util.validatemobile(phone_num)) {
        return;
    }


    var device = Util._getUrlParam("device");
    var json = {
        "bank_open": bank_open,
        "payee": payee,
        "payee_card": payee_card,
        "phone_num": phone_num
    };
    var jsonStr = JSON.stringify(json);
    if (device == "android") {
        window.couldCashier.OnDo("updateaccount", jsonStr);

    } else if (device == "ios") {
        var result = "couldcashier://updateaccount:/" + jsonStr;
        document.location = result;
    }

}


var limitincrease = function() {
    var limitincrease_phone_num = document.getElementById('limitincrease_phone_num').value;
    var limitincrease_payee = document.getElementById('limitincrease_payee').value;
    var limitincrease_email = document.getElementById('limitincrease_email').value;


    if (!Util.validatemobile(limitincrease_phone_num)) {
        return;
    }

    if (limitincrease_payee.length === 0) {
        alert("姓名不能为空!");
        return;
    }

    if (limitincrease_email.length === 0) {
        alert("邮箱不能为空!");
        return;
    }

    if (!Util.isEmail(limitincrease_email)) {
        alert("邮箱格式不正确!");
        return;
    }

    var device = Util._getUrlParam("device");
    var json = {
        "payee": limitincrease_payee,
        "phone_num": limitincrease_phone_num,
        "email": limitincrease_email
    };
    var jsonStr = JSON.stringify(json);
    if (device == "android") {
        window.couldCashier.OnDo("limitincrease", jsonStr);

    } else if (device == "ios") {
        var result = "couldcashier://limitincrease:/" + jsonStr;
        document.location = result;
    }


}
