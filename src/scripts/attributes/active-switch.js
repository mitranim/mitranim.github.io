/**
 * Directive for switching tabs in an <sf-tabset> or any other component where
 * one of several siblings needs to be set .active on click.
 */

angular.module('astil.attributes.activeSwitch', [])
.directive('activeSwitch', function() {

  return {
    restrict: 'A',
    scope: false,
    link: function(scope, $elem) {
      // Use the native DOM element.
      var elem = $elem[0]

      // Register a listener to remove the .active class from self.
      scope.$on('$active-switch', function() {
        elem.classList.remove('active')
      })

      // Register an onclick listener to remove the .active class from self and
      // all siblings, them add this class to self.
      elem.onclick = function() {
        scope.$parent.$broadcast('$active-switch')
        elem.classList.add('active')
        scope.$digest()
      }

      // If this is a first sibling in a repeater, activate by default.
      if (scope.$index === 0) elem.classList.add('active')
    }
  }

})
