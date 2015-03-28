/**
 * Generic controller.
 */

angular.module('astil.controllers.generic', [])
.factory('CtrlGeneric', function($q: $Q) {

  return class CtrlGeneric {

    ready: Function;
    loading: boolean;

    constructor() {
      /**
       * Marks end of loading.
       */
      this.ready = () => {this.loading = false}

      /**
       * Refer self.
       */
      this.refer()
    }

    /**
     * Loads records with the given promise hash and assigns them to self.
     */
    load(qHash: any): Promise {
      this.loading = true
      return $q.all(qHash).then(_.curry(_.assign, 2)(this))
    }

    /**
     * Loads records with the given promise hash and assigns them to the given
     * object. If the destination is not an object, returns a rejected
     * promise.
     */
    loadTo(destination: {}, qHash: any): Promise {
      if (!_.isObject(destination)) {
        return $q.reject('Destination must be an object.')
      }
      this.loading = true
      return $q.all(qHash).then(_.curry(_.assign, 2)(destination))
    }

    /**
     * Sets the self-reference.
     */
    refer() {
      if (this.element && this.element.hasAttribute('reference')) {
        this.reference = this
      }
    }

  }

})
