var cartPay = function() {
    var cartItems = [];
    var sumAmount = 0;
    var init = function() {
        _readData();
        _event();
    };
    var _event = function() {
        
    };
    var _readData = function() {
        var str = window.localStorage.getItem('cartItem');
        if (!str) {
            return
        }

        cartItems = JSON.parse(str);
        var htmls = [];
        sumAmount = 0;
        for (var i = 0, l = cartItems.length; i < l; i++) {
            var item = cartItems[i];
            item.total = item.price * item.count;
            sumAmount += item.total;
            item.total=toDecimal2(item.total);
            htmls[htmls.length] = _fillItemToTemplate(item);
        }
        sumAmount=toDecimal2(sumAmount);
        $('#totalAmount').text(sumAmount);
        $('#container').html(htmls.join(''));
    };
    //使用handlebars填充一个预订单的信息到html模板中
    var _fillItemToTemplate = function(item) {
        var source = $("#itemTemplate").html();
        var template = Handlebars.compile(source);
        return template(item);
    };

    var toDecimal2=function (x) {  
            var f = parseFloat(x);  
            if (isNaN(f)) {  
                return false;  
            }  
            var f = Math.round(x*100)/100;  
            var s = f.toString();  
            var rs = s.indexOf('.');  
            if (rs < 0) {  
                rs = s.length;  
                s += '.';  
            }  
            while (s.length <= rs + 2) {  
                s += '0';  
            }  
            return s;  
        }  
    return {
        init: init
    }
}();
$(function() {
    cartPay.init();
});
