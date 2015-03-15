var module = angular.module('astil.components.appWords', [
  'foliant',
  'astil.attributes',
  'astil.mixins.generic',

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

module.factory('appWordsCtrl', function(mixinGeneric, Traits, Lang, NamesExample, WordsExample) {

  return function Controller($scope) {
    // Use generic controller mixin.
    mixinGeneric(this)

    /**
     * Word count limit.
     * @type Number
     */
    this.limit = 12

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

    /******************************** Methods ********************************/

    /**
     * Takes a store object and produces a traits object based on its lang
     * and with its words.
     */
    this.getTraits = function(store) {
      var lang = _.find(this.langs, {Id: store.LangId})
      var traits = lang.$traits()
      traits.examine(store.Words)
      return traits
    }

    /**
     * Adds the given word to the given example store or displays an error
     * message.
     */
    this.add = function(store, word) {
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
      store.$gen = this.getTraits(store).generator()
    }

    /**
     * Generates a group of words for the given example store.
     */
    this.generate = function(store) {
      if (!store.$gen) store.$gen = this.getTraits(store).generator()
      var words = []

      while (words.length < this.limit) {
        var word = store.$gen()
        if (!word) break
        if (~store.Words.indexOf(word)) continue
        words.push(word)
      }

      if (words.length < this.limit) store.$depleted = true
      else delete store.$depleted

      store.$results = words
    }

    /**
     * Adds the given word to the given example store, removing it from the
     * generated results.
     */
    this.pick = function(store, word) {
      if (~store.Words.indexOf(word)) return
      store.Words.push(word)
      _.pull(store.$results, word)
      store.$gen = this.getTraits(store).generator()
    }

    /**
     * Removes the given word from the given example store.
     */
    this.drop = function(store, word) {
      _.pull(store.Words, word)
      store.$gen = this.getTraits(store).generator()
    }

    /**
     * Returns the appropriate text class for the given example store.
     */
    this.textClass = function(store) {
      return store.Title === 'Names' ? 'text-capitalise' : 'text-lowercase'
    }

    /********************************* Load **********************************/

    /**
     * Data load.
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
    .then(() => {this.lang = this.langs[0]})
    // Watch words in each example and save on change.
    .then(() => {
      _.each(this.nameExamples, record => {
        $scope.$watch(() => record.Words, () => record.$saveLS(), true)
      })
      _.each(this.wordExamples, record => {
        $scope.$watch(() => record.Words, () => record.$saveLS(), true)
      })
    })
    .then(this.ready)
  }

})
