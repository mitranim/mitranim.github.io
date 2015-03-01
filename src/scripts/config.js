/**
 * Basic configuration.
 */

var module = angular.module('astil.config', ['Datacore'])

module.constant('config', {
  dev: window.astilEnvironment === 'development'
  // dev: false
})

/**
 * Configure Datacore.
 */
module.run(function(Record, config) {

  if (config.dev && window.recordBaseUrl) {
    Record.baseUrl = window.recordBaseUrl
  } else {
    Record.baseUrl = 'http://api.mitranim.com'
  }

  // Match backend id key style.
  Record.prototype.$idKey = 'Id'

})
