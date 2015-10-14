(function(document) {
  'use strict';

  // See https://github.com/Polymer/polymer/issues/1381
  window.addEventListener('WebComponentsReady', function() {
    // imports are loaded and elements have been registered
    var router = document.querySelector('#appRouter');

    router.addEventListener('state-change', function(e) {
      // TODO 权限校验
      //   console.log(e.type, e.detail.path);
    });

    router.init();
  });

})(document);
