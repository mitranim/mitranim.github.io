/**
 * Lets the user input sample words and loads derived synthetic words from the
 * server.
 */

angular.module('astil.components.appWords', [
  'astil.attributes',
  'astil.mixins.generic',
  'astil.models.Word'
])
.directive('appWords', function(mixinGeneric, Word) {

  return {
    restrict: 'E',
    scope: {},
    templateUrl: 'components/app-words/app-words.html',
    controllerAs: 'self',
    bindToController: true,
    controller: [Controller]
  }

  function Controller() {
    var self = this

    // Use generic controller mixin.
    mixinGeneric(self)

    /**
     * Available languages.
     * @type [Hash]
     */
    self.langs = [
      {
        /** @type String */
        title: 'English',

        /** @type String */
        soundset: 'eng',

        /** @type String */
        source: [
          'jasmine', 'katie', 'nariko', 'nebula', 'aurora', 'theron',
          'quasar', 'graphene', 'nanite', 'orchestra', 'eridium',
        ].join(' '),

        /** @type [Word] */
        records: null,
      },
      // {
      //   /** @type String */
      //   title: 'Russian',

      //   /** @type String */
      //   soundset: 'cyr',

      //   /** @type String */
      //   source: [
      //     'дмитрий', 'владимир', 'степан', 'перун', 'хорс',
      //     'род', 'симаргл', 'велес', 'сварог',
      //   ].join(' '),

      //   /** @type [Word] */
      //   records: null,
      // },
    ]

    /**
     * Generates request parameters for the given language.
     * @param   Hash
     * @returns Hash
     */
    self.params = function(lang) {
      if (!_.isObject(lang)) return null
      if (typeof lang.source !== 'string') return null
      var words = _.invoke(lang.source.split(/\s+|,/), 'trim')
      return {
        words:    words,
        soundset: lang.soundset || null
      }
    }

    /**
     * Loads words from server using the given source.
     * @returns Promise
     */
    self.submit = function(lang) {
      self.loading = true
      return Word.readAll({params: self.params(lang)}).then(function(records) {
        // console.log("-- records:", self.records);
        // console.log("-- words:", _.map(self.records, 'Value'));
        lang.records = records
      }).finally(self.ready)
    }

    /**
     * Run first submit on page load.
     */
    self.submit(self.langs[0])
  }
})
