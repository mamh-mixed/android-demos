export class PasswordResetService {
	constructor($http, toastr, $log, $state) {
		'ngInject';

		this.toastr = toastr;
		this.log = $log;
		this.state = $state;
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
			return;
		}

    this.toastr.success('保存成功');
    // TODO 成功后跳转
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
