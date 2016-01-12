describe('service resetPassword', () => {
	beforeEach(angular.mock.module('quickpay'));

	beforeEach(inject((toastr, $state, $log) => {
		spyOn(toastr, 'info').and.callThrough();
		spyOn(toastr, 'error').and.callThrough();
		spyOn(toastr, 'warning').and.callThrough();
		spyOn($state, 'go').and.callThrough();
		spyOn($log, 'error').and.callThrough();
		spyOn($log, 'warn').and.callThrough();
	}));

	it('should be registered', inject(passwordResetService => {
		expect(passwordResetService).not.toEqual(null);
	}));

	describe('apiHost variable', () => {
		it('should exist', inject(passwordResetService => {
			expect(passwordResetService.apiHost).not.toEqual(null);
			expect(passwordResetService.apiHost).toEqual(jasmine.any(String));
		}));
	});

	describe('toastr,log, state variable', () => {
		it('should exist', inject(passwordResetService => {
			expect(passwordResetService.toastr).not.toEqual(null);
			expect(passwordResetService.log).not.toEqual(null);
			expect(passwordResetService.state).not.toEqual(null);
		}));
	});

	describe('validateCheckCode function', () => {
		it('should exist', inject(passwordResetService => {
			expect(passwordResetService.validateCheckCode).not.toEqual(null);
		}));

		it('should return boolean value and $state have been called', inject((passwordResetService, $state) => {
			let result = passwordResetService.validateCheckCode('123');
			expect(result).not.toEqual(undefined);
			expect(result).toBeTruthy();

			result = passwordResetService.validateCheckCode();
			expect(result).not.toBeTruthy();
			expect($state.go).toHaveBeenCalled();
		}));
	});

	describe('sendRequest function', () => {
		it('should exist', inject(passwordResetService => {
			expect(passwordResetService.sendRequest).not.toEqual(null);
		}));

		// TODO 重写
		it('should return data', inject((passwordResetService, toastr) => {
			// let data = passwordResetService.sendRequest({});
			// expect(data).not.toEqual(null);
			// expect(data).toEqual(jasmine.any(Number));
			// expect(data === 1).toBeTruthy();
			//
			// expect(toastr.info).toHaveBeenCalled();
		}));
	});

	describe('validate function', () => {
		it('should exist', inject(passwordResetService => {
			expect(passwordResetService.validate).not.toEqual(null);
		}));

		it('should retrun false when params is null or not an object', inject((passwordResetService, $log) => {
			let result = passwordResetService.validate();
			expect(result).toEqual(false);
			expect($log.error).toHaveBeenCalled();

			result = passwordResetService.validate([]);
			expect(result).toEqual(false);
			expect($log.error).toHaveBeenCalled();

			result = passwordResetService.validate('');
			expect(result).toEqual(false);
			expect($log.error).toHaveBeenCalled();
		}));

		it('should return false when required params is missing', inject((passwordResetService, $log, toastr) => {
			let [params, result] = [{
				username: '',
				password: '',
				passwordRepeat: ''
			}, false];

			result = passwordResetService.validate(params);
			expect(result).toEqual(false);
			expect($log.error).toHaveBeenCalled();
			expect(toastr.error).toHaveBeenCalled();

			params.username = '1234';
			result = passwordResetService.validate(params);
			expect(result).toEqual(false);
			expect($log.error).toHaveBeenCalled();
			expect(toastr.error).toHaveBeenCalled();

			params.password = '1234';
			result = passwordResetService.validate(params);
			expect(result).toEqual(false);
			expect($log.error).toHaveBeenCalled();
			expect(toastr.error).toHaveBeenCalled();

		}));

		it('password checing', inject((passwordResetService, $log, toastr) => {
			let [params, result] = [{
				username: 'wonsikin',
				password: '120943629',
				passwordRepeat: '123'
			}, false];

			result = passwordResetService.validate(params);
			expect(result).toEqual(false);
			expect($log.error).toHaveBeenCalled();
			expect(toastr.error).toHaveBeenCalled();

			params.passwordRepeat = '120943629';
			result = passwordResetService.validate(params);
			expect(result).toEqual(false);
			expect($log.warn).toHaveBeenCalled();
			expect(toastr.warning).toHaveBeenCalled();

			params.password = '$Wsj123456';
			params.passwordRepeat = '$Wsj123456';
			result = passwordResetService.validate(params);
			expect(result).toBeTruthy();
		}));
	});
});
