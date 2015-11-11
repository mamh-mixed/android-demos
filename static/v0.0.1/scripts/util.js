;
(function(window) {
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
			if (!toast) {
				toast = document.getElementById('toast');
			}

			if (!toast) {
				window.alert(text);
				return;
			}
			toast.text = text;
			toast.duration = duration;
			toast.show();
		};
		var showLoginDialog = function() {
			var abc = document.getElementsByTagName('paper-dialog'),
			a = null;
			for(var i = 0, l = abc.length; i < l; i++){
				a = abc[i];
				if (typeof a.close === 'function'){
					abc[i].close();
				}
			}
			Polymer.Base.fire('open-dialog-please', '', {
				node: document.getElementById('reloginDialog')
			});
		};
		var hideLoginDialog = function() {
			Polymer.Base.fire('close-dialog-please', '', {
				node: document.getElementById('reloginDialog')
			});
		};
		var fire = function(type, detail, node) {
			// create a CustomEvent the old way for IE9/10 support
			var event = document.createEvent('CustomEvent');
			// initCustomEvent(type, bubbles, cancelable, detail)
			event.initCustomEvent(type, false, true, detail);
			// returns false when event.preventDefault() is called, true otherwise
			return node.dispatchEvent(event);
		};
		var query = function(obj) {
			var q = '';
			for (var k in obj) {
				if (!obj[k]) {
					continue;
				}
				var v = encodeURIComponent(obj[k]);
				q += '&' + k + '=' + v;
			}
			return q.substring(1);
		};
		return {
			init: init,
			fire: fire,
			toast: toast,
			query: query,
			showLoginDialog: showLoginDialog,
			hideLoginDialog: hideLoginDialog
		};
	}());

	window.Util = Util;

	Date.prototype.format = function(format) {
		var o = {
			'M+': this.getMonth() + 1, //month
			'd+': this.getDate(), //day
			'h+': this.getHours(), //hour
			'm+': this.getMinutes(), //minute
			's+': this.getSeconds(), //second
			'q+': Math.floor((this.getMonth() + 3) / 3), //quarter
			'S': this.getMilliseconds() //millisecond
		};

		if (/(y+)/.test(format)) {
			format = format.replace(RegExp.$1, (this.getFullYear() + '').substr(4 - RegExp.$1.length));
		}

		for (var k in o) {
			if (new RegExp('(' + k + ')').test(format)) {
				format = format.replace(RegExp.$1, RegExp.$1.length === 1 ? o[k] : ('00' + o[k]).substr(('' + o[k]).length));
			}
		}
		return format;
	};

})(window);
