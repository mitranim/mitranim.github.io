'use strict';

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
  beforeEach(inject(function($q, Word) {
    var records = words().map(Word)
    spyOn(Word, 'readAll').and.returnValue($q.when(records))
  }))

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

  it('provides default input', function() {
    expect(_.size(self.source)).toBeGreaterThan(0)
    expect(_.every(self.source, _.isString)).toBe(true)
  })

  it('includes source words into request parameters', function() {
    self.source = 'one two three'
    expect(self.params()).toEqual({words: ['one', 'two', 'three']})
  })

  it('loads words and saves the result', inject(function(Word) {
    self.source = 'one two three'
    self.submit()
    this.scope.$digest()
    expect(Word.readAll).toHaveBeenCalled()
    expect(_.map(self.records, 'Value')).toEqual(_.map(words(), 'Value'))
  }))

  /*------------------------------- Constants -------------------------------*/

  function words() {
    return [
      {Value: 'first'},
      {Value: 'second'}
    ]
  }

})
