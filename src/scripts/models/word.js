/**
 * @class Word
 */

angular.module('astil.models.Word', ['Datacore'])
.factory('Word', function(Record) {

  /**
   * Class.
   */
  var Word = Record.derive({

    path: 'words'

  }, {

    $name: 'Word',

    $schema: {
      Value: ''
    },

    /**
     * Validates the word.
     * @returns Boolean
     */
    $valid: function(): boolean {
      return typeof this.Value === 'string' && /^[a-zа-я]{2,}$/.test(this.Value)
    }

  })

  /**
   * Export.
   */
  return Word

})
