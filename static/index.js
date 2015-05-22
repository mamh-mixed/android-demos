(function() {
    "use strict";

    var DEFAULT_ROUTE = 'bindingPayment';

    var template = document.querySelector('#t');

    // 菜单列表在此配置
    var menus = [{
        name: "测试",
        icon: "bug-report",
        submenus: [
            {name: '创建绑定关系', icon:"credit-card", hash: 'bindingCreate', url: 'bindingPayment/bindingCreate.html'},
            {name: '解除绑定关系', icon:"credit-card", hash: 'bindingRemove', url: 'bindingPayment/bindingRemove.html'},
            {name: '查询绑定关系', icon:"credit-card", hash: 'bindingEnquiry', url: 'bindingPayment/bindingEnquiry.html'},
            {name: '绑定支付', icon:"credit-card", hash: 'bindingPayment', url: 'bindingPayment/bindingPayment.html'},
            {name: '退款', icon:"credit-card", hash: 'refund', url: 'bindingPayment/refund.html'},
            {name: '交易对账汇总', icon:"credit-card", hash: 'billingSummary', url: 'bindingPayment/billingSummary.html'},
            {name: '交易对账明细', icon:"credit-card", hash: 'billingDetails', url: 'bindingPayment/billingDetails.html'},
            {name: '查询订单状态', icon:"credit-card", hash: 'orderEnquiry', url: 'bindingPayment/orderEnquiry.html'},
            {name: '无卡直接支付', icon:"credit-card", hash: 'noTrackPayment', url: 'bindingPayment/noTrackPayment.html'},
            {name: 'Apple Pay', icon:"credit-card", hash: 'applePay', url: 'bindingPayment/applePay.html'}
        ]
    }, {
        name: "配置",
        icon: "settings",
        submenus: [
            {name: '商户配置', icon:"store", hash: 'config/merchant', url: 'config/merchant.html'},
            {name: '渠道商户配置', icon:"account-balance", hash: 'config/channelMerchant', url: 'config/channelMerchant.html'},
            {name: '路由配置', icon:"settings-input-component", hash: 'config/routerPolicy', url: 'config/routerPolicy.html'},
        ]
    }];

    template.addEventListener('template-bound', function(e) {
        var keys = document.querySelector('#keys');

        this.menus = menus;
        this.route = this.route || DEFAULT_ROUTE; // Select initial route.
    });

    template.keyHandler = function(e, detail, sender) {
        // Select page by num key.
        var num = parseInt(detail.key);
        if (!isNaN(num) && num <= this.pages.length) {
            pages.selectIndex(num - 1);
            return;
        }

        switch (detail.key) {
            case 'left':
            case 'up':
                pages.selectPrevious();
                break;
            case 'right':
            case 'down':
                pages.selectNext();
                break;
            case 'space':
                detail.shift ? pages.selectPrevious() : pages.selectNext();
                break;
        }
    };

    // 切换商户事件
    template.chooseMerchantHandler = function (e, detail, sender) {
        document.querySelector("#merchantChangeOverlay").toggle();
    };

})();
