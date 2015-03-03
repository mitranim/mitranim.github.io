/**
 * @class Lang
 */

angular.module('astil.models.Lang', [
  'Datacore', 'astil.models.Mode'
])
.factory('Lang', function(Record, Mode) {

  return Record.derive({

    path: 'langs'

  }, {

    $name: 'Lang',

    $schema: {
      /**
       * @type String
       */
      title: '',

      /**
       * @type Words
       */
      modes: function(modes) {
        modes = Mode.collection(modes)
        _.invoke(modes, '$castExtend')
        return modes
      },
    },

  })

})
