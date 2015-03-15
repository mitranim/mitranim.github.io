'use strict'

describe('app-words', function() {
  var self

  /*--------------------------------- Setup ---------------------------------*/

  /**
   * Load modules.
   */
  beforeEach(module('astil.templates', 'astil.components.appWords'))

  /**
   * Data responses.
   */
  beforeEach()

  /**
   * Compile component.
   */
  beforeEach(inject(function($rootScope, $compile) {
    this.compile = function(scope) {
      var element = angular.element('<app-words></app-words>')
      $compile(element)(scope)
      scope.$digest()
      self = element.isolateScope().self
    }

    this.scope = $rootScope.$new()
    this.compile(this.scope)
  }))

  /*--------------------------------- Spec ----------------------------------*/

  /*------------------------------- Constants -------------------------------*/

})
