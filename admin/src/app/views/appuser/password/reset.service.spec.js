describe('service resetPassword', () => {
	beforeEach(angular.mock.module('quickpay'));

	beforeEach(inject((toastr, $state, $log) => {
		spyOn(toastr, 'info').and.callThrough();
		spyOn(toastr, 'error').and.callThrough();
		spyOn(toastr, 'warning').and.callThrough();
		spyOn(toastr, 'success').and.callThrough();
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

		it('should success when send request', inject((passwordResetService, $state, $httpBackend) => {
			let response = {
				status: 0,
				message: 'SUCCESS'
			}
			$httpBackend.when('POST', passwordResetService.apiHost).respond(200, response);

			let params = {
				password: 'Yun#1016',
				checkCode: 'f3fbe685172f46a768e9aab29cba6134'
			};
			let data;

			// 测试成功案例
			passwordResetService.sendRequest(params).then(body => {
				data = body;
			});
			$httpBackend.flush();
			expect(data).toEqual(jasmine.any(Object));
			expect(data.status).toEqual(0);
			expect(data.message).toEqual('SUCCESS');
			expect($state.go).toHaveBeenCalled();
		}));

		it('should error when send request', inject((passwordResetService, toastr, $httpBackend) => {
			let response = {
				status: 1,
				message: 'PASSWORD_MUST_BE_COMPLEX'
			}
			$httpBackend.when('POST', passwordResetService.apiHost).respond(200, response);

			let params = {
				password: 'Yun#1016',
				checkCode: 'f3fbe685172f46a768e9aab29cba6134'
			};
			let data;

			passwordResetService.sendRequest(params).then(body => {
				data = body;
			});
			$httpBackend.flush();
			expect(data).toEqual(jasmine.any(Object));
			expect(data.status).toEqual(1);
			expect(data.message).toEqual('PASSWORD_MUST_BE_COMPLEX');
			expect(toastr.error).toHaveBeenCalled();
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
				password: ''
			}, false];

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
				password: '120943629'
			}, false];

			result = passwordResetService.validate(params);
			expect(result).toEqual(false);
			expect($log.warn).toHaveBeenCalled();
			expect(toastr.warning).toHaveBeenCalled();
		}));
	});
});
