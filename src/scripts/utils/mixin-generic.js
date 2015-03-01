/**
 * Generic controller mixin.
 */

angular.module('astil.mixins.generic', [])
.factory('mixinGeneric', function($q) {
  return function(self) {

    /**
     * Loads records with the given promise hash and assigns them to self.
     */
    self.load = function(qHash) {
      self.loading = true
      return $q.all(qHash).then(_.curry(_.assign, 2)(self))
    }

    /**
     * Marks end of loading.
     */
    self.ready = function() {
      self.loading = false
    }

  }
})
