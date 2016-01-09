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
		.state('notFound', {
			url: '/404',
			templateUrl: 'app/views/common/404.html'
		});


	$urlRouterProvider.otherwise('/404');
}
