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

  var baseUrl: string = 'http://api.mitranim.com'
  if (config.dev && typeof window.recordBaseUrl === 'string') {
    baseUrl = window.recordBaseUrl
  }

  Record.prototype.$id = function() {return this.Id}

  Record.prototype.$path = function() {return baseUrl}

})
