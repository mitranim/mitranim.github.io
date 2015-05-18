/**
 * Disables angular-driven actions on the given element or its children when
 * the given condition is true, based on the presence of the attributes
 * ng-model, ng-click or ng-change. Example:
 *
 *   <div disable-all="myCondition">
 *       <a ng-click="do stuff"></a>
 *       <a ng-click="do stuff" ng-disabled="false"></a>
 *   </div>
 *
 * Becomes:
 *
 *   <div disable-all="myCondition">
 *       <a ng-click="do stuff" ng-disabled="myCondition"></a>
 *       <a ng-click="do stuff" ng-disabled="false || myCondition"></a>
 *   </div>
 *
 * The condition is treated as a string.
 */

import _ from 'lodash'
import {Attribute} from 'ng-decorate'

@Attribute({
  selector: '[disable-all]'
})
class VM {
  static compile($elem: ng.IAugmentedJQuery) {
    // Use the native element.
    var elem: HTMLElement = $elem[0]

    // Find elements where [ng-disabled] is relevant.
    var attr = elem.getAttribute('disable-all')
    var elems = elem.querySelectorAll('[ng-model],[ng-click],[ng-change],[ng-disabled],button[type="submit"]')

    // Add the current element.
    var list: HTMLElement[] = [].concat.apply([elem], elems)

    // Add or modify the [ng-disabled] attribute.
    _.each(list, node => {
      if (node.getAttribute('ng-disabled')) {
        var value = node.getAttribute('ng-disabled') + ' || ' + attr
        node.setAttribute('ng-disabled', value)
      } else {
        node.setAttribute('ng-disabled', attr)
      }
    })
  }
}
