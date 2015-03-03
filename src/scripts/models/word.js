/**
 * @class Word
 */

angular.module('astil.models.Word', ['Datacore'])
.factory('Word', function(Record) {

  return Record.derive({

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
    $valid: function() {
      return typeof this.Value === 'string' && /^[a-zа-я]{2,}$/.test(this.Value)
    }

  })

})
