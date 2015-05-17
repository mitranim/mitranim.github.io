/**
 * Generic viewmodel.
 */

import _ from 'lodash'
import {Service} from 'ng-decorate'
import {digest} from 'app'

@Service({
  inject: ['$q'],
  serviceName: 'BaseVM'
})
export class BaseVM {
  // Angular services.
  $q: ng.IQService

  // Loading status.
  loading: boolean = false
  // Ready method.
  ready = () => {
    this.loading = false
    this.dispatch('ready')
  }
  // Wrapped element.
  $element: ng.IAugmentedJQueryStatic
  // DOM element.
  element: HTMLElement

  // Static properties.
  static requireAuth: boolean = false

  constructor(...args) {
    // Assign injected arguments to self under matching key names.
    var inject = this.constructor.$inject
    if (inject && inject.length === args.length) {
      for (let index in args) this[inject[index]] = args[index]
    }
    if (this.$element) this.element = this.$element[0]
  }

  // Sets the given refs to sync values to self under the given keys. Returns a
  // hash of Fireproof promises returned from .on calls.
  sync(hash: {[key: string]: Fireproof}): {[key: string]: ng.IPromise<any>} {
    return _.mapValues(hash, (ref, key) => {
      return ref.on('value', snap => {
        this[key] = snap.val()
        digest()
      })
    })
  }

  /**
   * Dispatches the given custom event with the given data, if the
   * associated DOM element is available.
   */
  dispatch(eventName: string, data?: any) {
    if (this.element instanceof HTMLElement) {
      var event = new CustomEvent(eventName, {bubbles: true, detail: data})
      this.element.dispatchEvent(event)
    }
  }

  /**
   * Loads the given values and assigns them to self.
   */
  load(hash: any): ng.IPromise<any> {
    this.loading = true
    return this.$q.all(hash).then(data => {
      _.assign(this, data)
    }).finally(this.ready)
  }
}
