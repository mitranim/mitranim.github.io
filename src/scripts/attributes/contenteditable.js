/**
 * Mimics [ng-model] functionality for [contenteditable].
 */

angular.module('astil.attributes.contenteditable', ['ngSanitize'])
.directive('contenteditable', function($sce) {

  return {
    restrict: 'A',
    require: '?ngModel',
    scope: false,
    link: function(scope, elem, a, ngModel) {
      if (!ngModel) return

      // Specify how UI should be updated
      ngModel.$render = function() {
        elem[0].innerText = ngModel.$viewValue
      }

      // Listen for change events to enable binding
      elem.on('blur keyup change', function() {
        scope.$evalAsync(read)
      })
      read()

      // Write data to the model
      function read() {
        ngModel.$setViewValue(elem[0].innerText)
      }
    }
  }

})
