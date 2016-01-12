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
});
