/**
 * Shared references and classes.
 */

var module = angular.module('astil.firebase')

module.factory('fbRoot', function(config) {
  var fbRoot = new Firebase(config.fbRootUrl)
  if (config.dev) {
    window.fbRoot = fbRoot
  }
  return fbRoot
})

module.factory('FBArray', function($firebaseArray) {

  return $firebaseArray.$extend({

    /**
     * Iterates over primitive values or object values.
     */
    $each: function(iterator, thisArg) {
      this.$list.forEach(val => {
        if (val.$value != null) val = val.$value
        if (thisArg) iterator.call(thisArg, val, val.$id)
        else iterator(val, val.$id)
      })
    },

    /**
     * Reports whether the array has the given primitive value.
     */
    $has: function(value): boolean {
      var yes = false
      try {
        this.$each(val => {if (yes = val === value) throw null})
      } catch (err) {
        if (err !== null) throw err
      }
      return yes
    }

  })

})

if (window.astilEnvironment === 'development') {
  window.log = function log(value) {console.log(value)}
  window.warn = function warn(value) {console.warn(value)}
}
