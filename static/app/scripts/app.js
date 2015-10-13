(function(document) {
  'use strict';

  // See https://github.com/Polymer/polymer/issues/1381
  window.addEventListener('WebComponentsReady', function() {
    // imports are loaded and elements have been registered
    var router = document.querySelector('#appRouter');
    var agentRoutes = ['/trade/list','/trade/reports'];

    router.addEventListener('state-change', function(e) {
      // TODO 权限校验
      console.log(e.type, e.detail.path);
      var userType = window.sessionStorage.getItem('USERTYPE');
      var isAccess=false;
      if (userType === 'agent' || userType === 'group' || userType === 'merchant') {
        for (var i = 0; i < agentRoutes.length; i++) {
           if(agentRoutes[i]===e.detail.path){
               isAccess=true;
               break;
           }
        }
        if (!isAccess){
            window.location.href = '#/trade/list';
            e.preventDefault();
        }

      }
    });

    router.init();
  });

})(document);
