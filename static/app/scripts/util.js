;(function(window) {
  var Util = (function() {
    var init = function() {};
    var toast = function(text, duration) {
      if (!text || typeof text !== 'string' || text === '') {
        return;
      }
      if (!duration || typeof duration !== 'number' || duration <= 0) {
        return;
      }
      var toast = Polymer.Base.$$('#toast');
      toast.text = text;
      toast.duration = duration;
      toast.show();
    };
    return {
      init: init,
      toast: toast
    };
  }());

  window.Util = Util;

})(window);
