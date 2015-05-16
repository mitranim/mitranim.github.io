import 'angular'
import 'angular-sanitize'

/**
 * Application module.
 */
export var app = angular.module('app', ['ng', 'ngSanitize'])

/**
 * Digest hack.
 */
export function digest() {
  if (rootScope && !rootScope.$$phase) rootScope.$digest()
}
var rootScope: ng.IRootScopeService = null
app.run(['$rootScope', function($rootScope) {rootScope = $rootScope}])
