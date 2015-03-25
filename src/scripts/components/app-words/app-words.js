var module = angular.module('astil.components.appWords', [
  'foliant',
  'astil.attributes',
  'astil.controllers.generic',

  'astil.models.Lang',
  'astil.models.NamesExample',
  'astil.models.WordsExample'
])

module.directive('appWords', function(appWordsCtrl) {
  return {
    restrict: 'E',
    scope: {},
    templateUrl: 'components/app-words/app-words.html',
    controllerAs: 'self',
    bindToController: true,
    controller: ['$scope', appWordsCtrl]
  }
})

module.factory('appWordsCtrl', function(Traits, CtrlGeneric, Lang, NamesExample, WordsExample) {

  return class Controller extends CtrlGeneric {

    constructor(scope) {
      super()
      this.scope = scope

      /**
       * Languages.
       * @type Lang[]
       */
      this.langs = null

      /**
       * Selected lang.
       * @type Lang
       */
      this.lang = null

      /**
       * Example names.
       * @type NamesExample[]
       */
      this.nameExamples = null

      /**
       * Example words.
       * @type WordsExample[]
       */
      this.wordExamples = null

      /**
       * Load data.
       */
      this.load({
        langs: Lang.readAll(),
        nameExamples: NamesExample.readAll(),
        wordExamples: WordsExample.readAll()
      })
      // Sort out examples.
      .then(() => {
        _.each(this.langs, lang => {
          lang.$names = _.find(this.nameExamples, {LangId: lang.Id})
          lang.$words = _.find(this.wordExamples, {LangId: lang.Id})
        })
      })
      // Stick to the first lang for now.
      .then(() => {this.lang = _.first(this.langs)})
      // Watch each example and save on change.
      .then(() => {
        var scope = this.scope
        function watcher(record) {
          scope.$watch(_.constant(record), () => record.$saveLS(), true)
        }
        _.each(this.nameExamples, watcher)
        _.each(this.wordExamples, watcher)
      })
      .then(this.ready)
    }

    /**
     * Word count limit.
     * @type Number
     */
    get limit() {return 12}

    /**
     * Takes a store object and produces a traits object based on its lang
     * and with its words.
     */
    getTraits(store) {
      var lang = _.find(this.langs, {Id: store.LangId})
      var traits = lang.$traits()
      traits.examine(store.Words)
      return traits
    }

    /**
     * Adds the given word to the given example store or displays an error
     * message.
     */
    add(store, word) {
      if (typeof word !== 'string') word = ''
      word = word.toLowerCase().trim()

      if (!word) {
        store.$error = 'Please input a word.'
        return
      }

      if (word.length < 2) {
        store.$error = 'The word is too short.'
        return
      }

      if (~store.Words.indexOf(word)) {
        store.$error = 'This word is already in the set.'
        return
      }

      try {
        this.getTraits(store).examine([word])
      } catch (err) {
        console.error('-- word parsing error:', err)
        store.$error = 'Some of these characters are not allowed in a word.'
        return
      }

      store.$error = ''
      store.$input = ''
      store.Words.push(word)

      // Refresh the generator.
      store.$gen = this.getTraits(store).generator()
    }

    /**
     * Generates a group of words for the given example store.
     */
    generate(store) {
      if (!store.$gen) store.$gen = this.getTraits(store).generator()
      var words = []

      // Regex filter to use with words.
      var reg = new RegExp(store.filter || '', 'gi')

      while (words.length < this.limit) {
        var word = store.$gen()
        if (!word) break
        // Skip source words.
        if (~store.Words.indexOf(word)) continue
        // Skip words matching the filter. Using String#match because
        // RegExp#test is unpredictable when used on several strings in
        // succession.
        // if (word.match(reg)) continue
        words.push(word)
      }

      if (words.length < this.limit) store.$depleted = true
      else delete store.$depleted

      store.$results = words
    }

    /**
     * Adds the given word to the given example store, removing it from the
     * generated results. Doesn't refresh the generator because adding a
     * previously generated word to the same source set has no effect on the
     * total output.
     */
    pick(store, word) {
      if (~store.Words.indexOf(word)) return
      store.Words.push(word)
      _.pull(store.$results, word)
    }

    /**
     * Removes the given word from the given example store.
     */
    drop(store, word) {
      _.pull(store.Words, word)
      // Refresh the generator.
      store.$gen = this.getTraits(store).generator()
    }

    /**
     * Returns the appropriate text class for the given example store.
     */
    textClass(store) {
      return store.Title === 'Names' ? 'text-capitalise' : 'text-lowercase'
    }

  }

})
