/**
 * Basic configuration.
 */

var module = angular.module('astil.config', ['foliant'])

module.constant('config', {
  dev: window.astilEnvironment === 'development',

  baseUrl: window.astilEnvironment === 'development' && typeof window.recordBaseUrl === 'string' ?
           window.recordBaseUrl : 'http://api.mitranim.com',

  fbRootUrl: 'https://incandescent-torch-3438.firebaseio.com'
})

/**
 * Hack foliant's Pair prototype to compare pairs to strings more easily.
 */
module.run(function(Traits) {
  var traits = new Traits(['blah'])
  var pair = traits.pairSet[0]

  // Get the Pair constructor.
  var Pair = pair.constructor

  // This lets us compare pairs to strings like so: pair == 'bl'
  Pair.prototype.toString = function() {
    return this[0] + this[1]
  }
})
