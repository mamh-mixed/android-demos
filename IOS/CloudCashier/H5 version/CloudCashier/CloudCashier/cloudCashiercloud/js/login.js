var Login = function() {
    var username, password, isAutoLogin, device;
    var init = function() {

        username = Util.getUrlParam("username");
        password = Util.getUrlParam("password");
        isAutoLogin = Util.getUrlParam("autologin");
        device = Util.getUrlParam("device");
        Util.setUsername(username);
        Util.setPassword(password);
        Util.setIsAutoLogin(isAutoLogin);
        Util.setDevice(device);
        initLayout();
        events();
    };
    var initLayout = function() {
        document.getElementById('email').value = username;
        document.getElementById('password').value = password;
        var autoLogin = (isAutoLogin === 'true');
        document.getElementById('autoLoginCheckbox').checked = autoLogin;
    };
    var events = function() {
        document.getElementById('loginBtn').onclick = login;
        document.getElementById('registerBtn').onclick = register;
    };
    // 默认数据
    var setParameter = function(data) {
        username = data.username;
        password = data.password;
        isAutoLogin = data.autologin;
        device = data.device;

        Util.setUsername(username);
        Util.setPassword(password);
        Util.setIsAutoLogin(isAutoLogin);
        Util.setDevice(device);

        initLayout();
    };
    // 登入事件
    var login = function() {
        var email = document.getElementById('email').value,
            password = document.getElementById('password').value,
            ischeck = document.getElementById('autoLoginCheckbox').checked;
        if (email.length === 0) {
            window.alert('邮箱不能为空!');
            return;
        }

        // if (!Util.isEmail(email)) {
        //     window.alert('邮箱格式不正确!');
        //     return;
        // }

        if (password.length === 0) {
            window.alert('密码不能为空!');
            return;
        }

        if (password.length < 6) {
            window.alert('密码长度不能小于六位!');
            return;
        }

        var json = {
            "username": email,
            "password": password,
            "autologin": '' + ischeck
        };
        var jsonStr = JSON.stringify(json);

        Util.sendRequestToDevice('login', jsonStr);
    };
    var register = function() {
        Util.sendRequestToDevice('jump_register', '');
    };
    return {
        init: init,
        setParameter: setParameter
    };
}(window);

Login.init();
