/**
 * Generic controller mixin.
 */

angular.module('astil.mixins.generic', [])
.factory('mixinGeneric', function($q: $Q) {

  return function(self) {

    /**
     * Loads records with the given promise hash and assigns them to self.
     */
    self.load = function(qHash: any): Promise {
      self.loading = true
      return $q.all(qHash).then(_.curry(_.assign, 2)(self))
    }

    /**
     * Loads records with the given promise hash and assigns them to the given
     * object. If the destination is not an object, returns a rejected
     * promise.
     */
    self.loadTo = function(destination: {}, qHash: any): Promise {
      if (!_.isObject(destination)) {
        return $q.reject('Destination must be an object.')
      }
      self.loading = true
      return $q.all(qHash).then(_.curry(_.assign, 2)(destination))
    }

    /**
     * Marks end of loading.
     */
    self.ready = function() {
      self.loading = false
    }

  }

})
