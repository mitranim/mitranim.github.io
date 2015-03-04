/**
 * @class Mode
 */

angular.module('astil.models.Mode', [
  'Datacore', 'astil.models.Word'
])
.factory('Mode', function(Record, Word) {

  var Mode = Record.derive({

    path: 'modes'

  }, {

    $name: 'Mode',

    $schema: {
      /**
       * @type String
       */
      title: '',

      /**
       * Polymorphic adapter. ToDo split into monomorphic.
       */
      source: function(words: Word[]): Word.collection {
        if (_.some(words, _.isString)) {
          return Word.collection(_.map(words, function(value) {
            return Word({Value: value})
          }))
        }
        return Word.collection(words)
      },

      /**
       * @type String
       */
      soundset: '',

      /**
       * @type String
       */
      LangId: '',
    },

    $extendedSchema: {
      /**
       * @type Words
       */
      generated: null,

      /**
       * @type String
       */
      textMode: function(): string {
        if (this.title === 'Names') return 'text-capitalise'
        return 'text-lowercase'
      },

      /**
       * @type String
       */
      word: '',
    },

    /******************************** Methods ********************************/

    /**
     * Returns the values of own source words as an array of strings.
     */
    words: function(): string {
      return _.invoke(_.map(this.source, 'Value'), 'toLowerCase')
    },

    /**
     * Removes the given word from source and generated.
     */
    drop: function(word: Word) {
      _.pull(this.source, word)
      _.pull(this.generated, word)
    },

    /**
     * Moves the given word from generated to source.
     */
    pick: function(word: Word) {
      _.pull(this.generated, word)
      this.source.push(word)
    },

    /**
     * Converts the given string to a word and adds it to source.
     */
    add: function(string: string) {
      var value = string.toLowerCase().trim()
      if (!value) {
        this.error = 'Please input a word.'
        return
      }

      var word = Word({Value: value})

      if (word.Value.length < 2) {
        this.error = 'The word is too short.'
        return
      }

      if (!word.$valid()) {
        this.error = 'Some of these characters are not allowed in a word.'
        return
      }

      if (_.some(this.source, {Value: word.Value})) {
        this.error = 'This word is already in the set.'
        return
      }

      delete this.error
      this.source.push(word)
      this.word = ''
    },

  })

  return Mode

})
