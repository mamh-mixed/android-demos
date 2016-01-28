export function routerConfig($stateProvider, $urlRouterProvider) {
	'ngInject';
	$stateProvider
		.state('home', {
			url: '/',
			templateUrl: 'app/views/dashboard/dashboard.html'
		})
		.state('appPwdFgt', {
			url: '/app/password/forget/{code}',
			templateUrl: 'app/views/appuser/password/reset.html',
			controller: 'PasswordResetController',
			controllerAs: 'prc'
		})
		.state('appPwdRstSucc', {
			url: '/app/password/reset/success',
			templateUrl: 'app/views/appuser/password/success.html'
		})
		.state('appAgreement', {
			url: '/app/agreement',
			templateUrl: 'app/views/appuser/agreement/index.html'
		})
		// 旧浏览器
		.state('oldBro', {
			url: '/601',
			templateUrl: 'app/views/common/601.html'
		})
		.state('notFound', {
			url: '/404',
			templateUrl: 'app/views/common/404.html'
		});


	$urlRouterProvider.otherwise('/404');
}
