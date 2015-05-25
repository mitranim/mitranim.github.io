/**
 * Semantic shortcut to including an SVG template with ng-include.
 */

import {Attribute} from 'ng-decorate'

@Attribute({
  selector: '[svg-icon]',
  injectStatic: ['$templateCache']
})
class VM {
  static $templateCache: ng.ITemplateCacheService

  static template($element: ng.IAugmentedJQuery) {
    var element: HTMLElement = $element[0]

    var src = 'svg/' + element.getAttribute('svg-icon') + '.svg'
    element.removeAttribute('svg-icon')
    element.classList.add('sf-icon')

    return VM.$templateCache.get(src)
  }
}
