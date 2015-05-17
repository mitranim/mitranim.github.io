import 'angular'
import 'angular-sanitize'
import {defaults} from 'ng-decorate'

/**
 * Application module.
 */
export var app = angular.module('app', ['ng', 'ngSanitize'])

/**
 * Use this module for all decorations.
 */
defaults.module = app

/**
 * Digest hack.
 */
export function digest() {
  if (rootScope && !rootScope.$$phase) rootScope.$digest()
}
var rootScope: ng.IRootScopeService = null
app.run(['$rootScope', function($rootScope) {rootScope = $rootScope}])
