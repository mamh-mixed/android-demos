export class PasswordResetController {
  constructor($stateParams, passwordResetService) {
    'ngInject';
    let prc = this;
    prc.params = $stateParams;
    prc.service = passwordResetService;
    prc.activate();
    prc.request = {
      username: '',
      password: '',
      passwordRepeat: ''
    };
  }

  activate() {
    this.service.validateCheckCode(this.params.code);
  }

  sendRequest() {
    this.request.checkCode = this.params.code;
    this.service.sendRequest(this.request);
  }
}
