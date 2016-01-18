describe('PasswordResetController', () => {
  let prc;

  beforeEach(angular.mock.module('quickpay'));

  beforeEach(inject(($controller, $stateParams, passwordResetService) => {
    spyOn(passwordResetService, 'validateCheckCode').and.callThrough();
    spyOn(passwordResetService, 'sendRequest').and.callThrough();

    prc = $controller('PasswordResetController');
  }));

  it('should be an object', () => {
    expect(prc.request).toEqual(jasmine.any(Object));
  });

  describe('togglePasswordShow function', () => {
    it('the default value of inputType should be password', () => {
      expect(prc.inputType === 'password').toBeTruthy();
    });

    it('should change to be "text" after first click and to be "password" after second click', () => {
      prc.togglePasswordShow();
      expect(prc.inputType === 'text').toBeTruthy();

      prc.togglePasswordShow();
      expect(prc.inputType === 'password').toBeTruthy();
    });
  });
});
