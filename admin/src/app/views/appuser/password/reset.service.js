// var SHA256 = require("crypto-js/sha256");
// import { SHA256 } from 'crypto-js';
export class PasswordResetService {
	constructor($http, toastr, $log, $state) {
		'ngInject';

		this.toastr = toastr;
		this.log = $log;
		this.state = $state;
		this.$http = $http;
		this.apiHost = '/master/user/app/password/reset';
	}

	/**
	 * validateCheckCode 校验验证码
	 * 如果验证码不存在或者空，跳转到404页面
	 */
	validateCheckCode(checkCode) {
		if (!checkCode || checkCode === '') {
			this.state.go('notFound');
			return false;
		}

		return true;
	}

	/**
	 * 保存数据
	 */
	sendRequest(params = {}) {
		if (!this.validate(params)) {
			return {status: 7, message: "DATA_VALIDATE_FAIL"};
		}

		// 密码加密
		let data = {
			username: params.username,
			checkCode: params.checkCode,
			password: CryptoJS.MD5(params.password).toString()
		};

		return this.$http.post(this.apiHost, angular.toJson(data))
			.then((response) => {
				let body = response.data;
				if (body.status === 0) {
					this.toastr.success('保存成功');
					return body;
				}

				switch (body.message) {
					case 'MISS_REQUIRED_PARAMETER':
						this.toastr.error('缺失必要参数：' + response.data);
						break;
					case 'PASSWORD_MUST_BE_COMPLEX':
						this.toastr.error('密码必须包含大小写字母和数字，并且长度不小于8位字符');
						break;
					case 'USERNAME_NOT_MATCH':
						this.toastr.error('输入的用户名未注册');
						break;
					case 'INVALID_CHECK_CODE':
						this.toastr.error('没有权限访问');
						break;
					case 'OPERATION_OUT_OF_DATE':
						this.toastr.error('操作超时，请重新申请密码修改');
						break;
					case 'ALREADY_OPERATED':
						this.toastr.info('您已经变更过密码了，请勿重新操作');
						break;
					default:
						this.toastr.error('系统错误');
				}

				return body;
			})
			.catch((error) => {
				this.log.error('XHR Failed for sendRequest: ' + error.data);
				return error;
			});
	}

	/**
	 * 参数验证,
	 * return true if params is valid;false is invalidate;
	 */
	validate(params) {
		if (!params) {
			this.log.error('parmas is null');
			return false;
		}

		if (!(angular.isObject(params) && !angular.isArray(params))) {
			this.log.error('parmas is not object');
			return false;
		}

		if (params.username === '') {
			this.log.error('username is required');
			this.toastr.error('用户名不能为空', 'ERROR');
			return false;
		}

		if (params.password === '') {
			this.log.error('password is required');
			this.toastr.error('密码不能为空', 'ERROR');
			return false;
		}

		if (params.passwordRepeat === '') {
			this.log.error('passwordRepeat is required');
			this.toastr.error('重复密码不能为空', 'ERROR');
			return false;
		}

		if (params.password !== params.passwordRepeat) {
			this.log.error('the two password is not equal');
			this.toastr.error('两次输入的密码不一致', 'ERROR');
			return false;
		}

		if (!/^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])[a-zA-Z\d#\$@\.\_]{8,50}$/.test(params.password)) {
			this.log.warn('the two password is not equal');
			this.toastr.warning('新密码必须包含大小写字母和数字', 'WARNING');
			return false;
		}

		return true;
	}
}
