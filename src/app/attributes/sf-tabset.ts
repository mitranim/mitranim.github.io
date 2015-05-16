/**
 * Enhances <sf-tabset> with JS-driven tab switching.
 */

import _ from 'lodash'
import {Attribute} from 'ng-decorate'

@Attribute({
  moduleName: 'app',
  selector: 'sf-tabset',
  restrict: 'E'
})
class VM {
  element: HTMLElement

  static $inject = ['$element']
  constructor($element) {
    this.element = $element[0]

    // If no label is marked as active, activate the first.
    var labels = this.getLabels()
    if (labels.length && !_.any(labels, label => label.classList.contains('active'))) {
      labels[0].classList.add('active')
    }

    // Aggregate clicks on tab labels and make the clicked one active.
    this.element.addEventListener('click', event => {
      // Ignore if the target is not among the tab labels.
      var labels = this.getLabels()
      if (!_.contains(labels, event.target)) return

      // Deactivate all and activate the target.
      _.each(labels, label => {
        label.classList.remove('active')
      });
      (<HTMLElement>event.target).classList.add('active')
    })
  }

  getLabels(): HTMLElement[] {
    return <HTMLElement[]>_.filter(this.element.childNodes, node => {
      return node instanceof Element && (<Element>node).tagName === 'LABEL'
    })
  }
}
