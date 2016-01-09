export function runBlock($log, $rootScope) {
	'ngInject';
	$log.debug('runBlock end');
	let rootScope = $rootScope;
	rootScope.$on('$stateChangeStart', (event, toState, toParams, fromState, fromParams) => {
		// event.preventDefault();
		$log.debug(event, toState, toParams, fromState, fromParams);
		// transitionTo() promise will be rejected with
		// a 'transition prevented' error
	});
}
