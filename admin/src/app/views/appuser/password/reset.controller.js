export class PasswordResetController {
  constructor($stateParams, passwordResetService) {
    'ngInject';
    let prc = this;
    prc.params = $stateParams;
    prc.passwordResetService = passwordResetService;
    prc.activate();
    prc.request = {
      username: 'fnghwsj@qq.com',
      password: 'Yun#1016',
      passwordRepeat: 'Yun#1016'
    };
  }

  activate() {
    this.passwordResetService.validateCheckCode(this.params.code);
  }

  sendRequest() {
    this.request.checkCode = this.params.code;
    this.passwordResetService.sendRequest(this.request);
  }
}
