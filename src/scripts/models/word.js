/**
 * @class Word
 */

angular.module('astil.models.Word', ['Datacore'])
.factory('Word', function($q, Record) {

  return Record.derive({

    path: 'words'

  }, {

    $name: 'Word',

    $schema: {
      Value: ''
    }

  })

})
