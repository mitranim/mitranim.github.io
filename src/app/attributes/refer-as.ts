/**
 * <my-outer-component>
 *     <my-inner-component refer-as="self.innerComponent"></my-inner-component>
 * </my-outer-component>
 */

import {Attribute} from 'ng-decorate'

@Attribute({
  selector: '[refer-as]',
  inject: ['$parse']
})
class VM {
  $parse: ng.IParseService

  static $inject = ['$scope', '$element']
  constructor($scope: ng.IScope, $element: ng.IAugmentedJQuery) {
    var element: HTMLElement = $element[0]

    /**
     * Find the isolated scope controller.
     */
    var isolatedScope = $element.isolateScope()
    if (!isolatedScope) return
    var ctrl = (<any>isolatedScope).self
    if (!ctrl) return

    /**
     * Assign it on the outer scope under the given path.
     */
    var path = element.getAttribute('refer-as')
    this.$parse(path).assign($scope, ctrl)
  }
}
