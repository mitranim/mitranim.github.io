/**
 * Basic configuration.
 */

var module = angular.module('astil.config', [])

module.constant('config', {
  dev: window.astilEnvironment === 'development',
  // dev: false,

  baseUrl: window.astilEnvironment === 'development' && typeof window.recordBaseUrl === 'string' ?
           window.recordBaseUrl : 'http://api.mitranim.com'
})
