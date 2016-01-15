const PASSWORD_INPUT = 'password';
const TEXT_INPUT = 'text';
export class PasswordResetController {
  constructor($stateParams, passwordResetService) {
    'ngInject';
    let prc = this;
    prc.params = $stateParams;
    prc.passwordResetService = passwordResetService;
    prc.activate();
    prc.request = {
      password: ''
    };
    prc.inputType = PASSWORD_INPUT;
  }

  activate() {
    this.passwordResetService.validateCheckCode(this.params.code);
  }

  sendRequest() {
    this.request.checkCode = this.params.code;
    this.passwordResetService.sendRequest(this.request);
  }

  togglePasswordShow() {
    if (this.inputType === PASSWORD_INPUT) {
      this.inputType = TEXT_INPUT;
    } else {
      this.inputType = PASSWORD_INPUT;
    }
  }
}
