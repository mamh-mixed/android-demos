var CartDemo = function() {
    var dataMap = {};
    var cart = {};
    var init = function() {
        _initialItems();
        _events();
    };
    var _events = function() {
        $('#container').on('click', '.buy-btn', function(e) {
            var $this = $(this),
                idx = $this.data('id');
            if (!cart[idx]) {
                cart[idx] = 1;
            } else {
                cart[idx]++;
            }
            var sum = 0;
            for (var i in cart) {
                sum += cart[i];
            }
            $('#cartCount').text(sum);
            _saveToStorage();
        });
        $('#cartIcon').click(function() {
            window.location.href = "./cart.html";
        });
    };
    var _saveToStorage = function() {
        var storage = [];
        for (var i in dataMap) {
            if (cart[i]) {
                var item = dataMap[i];
                item.count = cart[i];
                storage[storage.length] = item;
            }
        }
        window.localStorage.setItem('cartItem', JSON.stringify(storage));
    };
    var _initialItems = function() {
        for (var i = 0, l = datas.length; i < l; i++) {
            dataMap[i + 1] = datas[i];
        }
        $('#container').empty();
        var htmls = [];
        for (var i = 0, l = datas.length; i < l; i++) {
            htmls[htmls.length] = _fillItemToTemplate(datas[i]);
        }
        var htmlss = [];
        for (var i = 0, l = htmls.length; i < l; i++) {
            var idx = parseInt(i / 2);
            if (i % 2 == 0) {
                htmlss[idx] = '';
            }
            htmlss[idx] += htmls[i];
        }
        $('#container').html('<div class="row">' + htmlss.join('</div><div class="row">') + '</div>')

    };
    //使用handlebars填充一个预订单的信息到html模板中
    var _fillItemToTemplate = function(item) {
        var source = $("#itemTemplate").html();
        var template = Handlebars.compile(source);
        return template(item);
    };
    var datas = [{
        id: 1,
        pic: 'img/n1.jpg',
        price: 0.01,
        name: '衣服1'
    }, {
        id: 2,
        pic: 'img/n2.jpg',
        price: 0.1,
        name: '衣服2'
    }, {
        id: 3,
        pic: 'img/n3.jpg',
        price: 1.01,
        name: '衣服3'
    }, {
        id: 4,
        pic: 'img/n4.jpg',
        price: 2.11,
        name: '衣服4'
    }, {
        id: 5,
        pic: 'img/n5.jpg',
        price: 3.01,
        name: '衣服5'
    }, {
        id: 6,
        pic: 'img/n6.jpg',
        price: 10.00,
        name: '衣服6'
    }, {
        id: 7,
        pic: 'img/n7.jpg',
        price: 15.00,
        name: '衣服7'
    }, {
        id: 8,
        pic: 'img/n8.jpg',
        price: 100.00,
        name: '衣服8'
    }];
    return {
        init: init
    }
}();
$(function() {
    CartDemo.init();
});
