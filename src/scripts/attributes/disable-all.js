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

angular.module('astil.attributes.disableAll', [])
.directive('disableAll', function() {

  return {
    restrict: 'A',
    scope: false,
    compile: function(elem, attrs) {
      // Find elements with relevant attributes and modify them.
      var attr = attrs.disableAll
      var elems = elem[0].querySelectorAll('[ng-model],[ng-click],[ng-change],[ng-disabled],button[type="submit"]')
      _.each(elems, function(elem) {
        var $elem = angular.element(elem)
        if ($elem.attr('ng-disabled')) {
          $elem.attr('ng-disabled', $elem.attr('ng-disabled') + ' || ' + attr)
        } else {
          $elem.attr('ng-disabled', attr)
        }
      })
    }
  }

})
