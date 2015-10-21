(function(document) {
	'use strict';

	var notAuthPath = [
		'/404',
		'/'
	];

	// See https://github.com/Polymer/polymer/issues/1381
	window.addEventListener('WebComponentsReady', function() {
		// imports are loaded and elements have been registered
		var router = document.querySelector('#appRouter');
		var agentRoutes = ['/trade/list', '/trade/reports'];

		router.addEventListener('state-change', function(e) {
			console.log(e.type, e.detail.path);
			if (notAuthPath.indexOf(e.detail.path) >= 0) {
				return;
			}

			var userType = window.localStorage.getItem('USERTYPE');
			if (userType && userType !== '') {
				userType = userType.substr(1, userType.length - 2);
			}
			var isAccess = false;

			if (userType === 'admin') {
				return;
			}

			if (userType !== 'agent' && userType !== 'group' && userType !== 'merchant') {
				document.querySelector('#appRouter').go('/');
				return;
			}

			// TODO 按照不同的用户类型给定不同的菜单权限
			for (var i = 0; i < agentRoutes.length; i++) {
				if (agentRoutes[i] === e.detail.path) {
					isAccess = true;
					break;
				}
			}
			if (!isAccess) {
				router.go('/404');
				// window.location.href = '#/trade/list';
				e.preventDefault();
			}

			// }
		});

		router.init();
	});

})(document);
