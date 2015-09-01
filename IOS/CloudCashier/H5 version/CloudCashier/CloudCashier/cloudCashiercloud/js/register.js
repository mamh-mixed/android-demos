var Register = function() {
    var init = function() {
    };
    var next = function() {
        var email = document.getElementById('email').value,
            password1 = document.getElementById('password1').value,
            password2 = document.getElementById('password2').value;

        if (email.length === 0) {
            window.alert('邮箱不能为空!');
            return;
        }

        if (!Util.isEmail(email)) {
            window.alert('邮箱格式不正确!');
            return;
        }

        if (password1.length === 0) {
            window.alert('密码不能为空!');
            return;
        }

        if (password1.length < 6) {
            window.alert('密码长度不能小于六位!');
            return;
        }

        if (password2 !== password1) {
            window.alert('密码不一致!');
            return;
        }

        var json = {
            "username": email,
            "password": password1
        };
        var jsonStr = JSON.stringify(json);

        Util.sendRequestToDevice('register', jsonStr);
    };
    var improveinfo = function() {
        var bank_open = document.getElementById('bank_open').value,
            payee = document.getElementById('payee').value,
            payee_card = document.getElementById('payee_card').value,
            phone_num = document.getElementById('phone_num').value;

        if (bank_open.length === 0) {
            window.alert('开户行不能为空!');
            return;
        }

        if (payee.length === 0) {
            window.alert('姓名不能为空!');
            return;
        }

        if (payee_card.length === 0) {
            window.alert('银行卡号不能为空!');
            return;
        }

        if (!Util.luhmCheck(payee_card)) {
            window.alert('银行卡号验证未通过!');
            return;
        }

        if (!Util.validatemobile(phone_num)) {
            return;
        }
        var json = {
            "bank_open": bank_open,
            "payee": payee,
            "payee_card": payee_card,
            "phone_num": phone_num
        };
        var jsonStr = JSON.stringify(json);
        Util.sendRequestToDevice('improveinfo', jsonStr);
    };
    var back = function() {
        window.history.go(-1);
        Util.sendRequestToDevice('back', '');
    };
    return {
        init: init,
        back: back,
        next: next,
        improveinfo: improveinfo
    };
}(window);

Register.init();
