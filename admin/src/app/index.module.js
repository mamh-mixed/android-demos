import "babel-polyfill";

import { config } from './index.config';
import { routerConfig } from './index.route';
import { runBlock } from './index.run';
////////
//控制器 //
////////
import { MainController } from './index.controller';
import { PasswordResetController } from './views/appuser/password/reset.controller';

////////
//服务 //
////////
import { PasswordResetService } from './views/appuser/password/reset.service';

///////
//指令 //
///////
import { ToggleSubmenuDirective } from '../app/components/sidebar/sidebar.directive';
import { FgLineDirective } from './components/fgLine/fgLine.directive';
import { WavesEffectDirective } from './components/wavesEffect/wavesEffect.directive';

angular.module('quickpay', ['ngAnimate', 'ngCookies', 'ngSanitize', 'ngMessages', 'ngAria', 'ui.router', 'toastr', 'ngTable'])
  .config(config)
  .config(routerConfig)
  .run(runBlock)
  .service('passwordResetService', PasswordResetService)
  .controller('MainController', MainController)
  .controller('PasswordResetController', PasswordResetController)
  .directive('toggleSubmenu', ToggleSubmenuDirective)
  .directive('fgLine', FgLineDirective)
  .directive('wavesEffect', WavesEffectDirective);
