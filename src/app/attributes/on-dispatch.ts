/**
 * Emulates angular2 event dispatching. Should be used on isolated directives
 * whose controllers inherit from BaseVM.
 */

import _ from 'lodash'
import {app} from 'app'
import {defaults} from 'ng-decorate'

;['ready', 'select'].forEach(eventName => {
  app.directive(`on${_.capitalize(eventName)}`, ['$parse', function($parse) {
    return {
      restrict: 'A',
      scope: false,
      link: function(scope, $elem) {
        var elem: Element = $elem[0]

        var statement = $parse(elem.getAttribute('on-' + eventName))

        /**
         * Find the isolated scope controller and assign the element to it.
         */
        var isolatedScope = $elem.isolateScope()
        if (isolatedScope) {
          var ctrl = isolatedScope[defaults.controllerAs]
          if (ctrl) ctrl.element = elem
        }

        /**
         * Add listener.
         */
        elem.addEventListener(eventName, function(event: Event) {
          statement(scope, {$event: event})
          event.stopPropagation()
        })
      }
    }
  }])
})
