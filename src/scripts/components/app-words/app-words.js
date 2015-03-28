var module = angular.module('astil.components.appWords', [
  'foliant',
  'astil.firebase',
  'astil.attributes',
  'astil.controllers.generic'
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

module.factory('appWordsCtrl', function(Traits, CtrlGeneric, lang, names, words) {

  return class Controller extends CtrlGeneric {

    constructor(scope) {
      super()

      /**
       * @type Lang
       */
      this.lang = lang

      /**
       * @type string[][]
       */
      this.stores = [names, words]

      /**
       * After each store has loaded, generate words and mark readiness status.
       */
      this.stores.forEach(store => {
        store.$loaded().then(() => {
          store.$ready = true
          this.generate(store)
        })
      })
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
      var traits = this.lang.$traits()
      traits.examine(_.map(store, '$value'))
      return traits
    }

    // /**
    //  * Adds the given word to the given example store or displays an error
    //  * message.
    //  */
    // add(store, word) {
    //   if (typeof word !== 'string') word = ''
    //   word = word.toLowerCase().trim()

    //   if (!word) {
    //     store.$error = 'Please input a word.'
    //     return
    //   }

    //   if (word.length < 2) {
    //     store.$error = 'The word is too short.'
    //     return
    //   }

    //   if (~store.Words.indexOf(word)) {
    //     store.$error = 'This word is already in the set.'
    //     return
    //   }

    //   try {
    //     this.getTraits(store).examine([word])
    //   } catch (err) {
    //     console.error('-- word parsing error:', err)
    //     store.$error = 'Some of these characters are not allowed in a word.'
    //     return
    //   }

    //   store.$error = ''
    //   store.$input = ''
    //   store.Words.push(word)

    //   // Refresh the generator.
    //   store.$gen = this.getTraits(store).generator()
    // }

    /**
     * Generates a group of words for the given example store.
     */
    generate(store) {
      if (!store.$gen) store.$gen = this.getTraits(store).generator()
      var words = []

      while (words.length < this.limit) {
        var word = store.$gen()
        if (!word) break
        // Skip source words.
        if (store.$has(word)) continue
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
      if (store.$has(word)) return
      store.$add(word).then(() => {
        _.pull(store.$results, word)
      })
    }

    /**
     * Removes the given word from the given example store and refreshes the
     * generator.
     */
    drop(store, item) {
      store.$remove(item).then(() => {
        store.$gen = this.getTraits(store).generator()
      })
    }

    /**
     * Returns the appropriate text class for the given example store.
     */
    textClass(title) {
      return title === 'Names' ? 'text-capitalise' : 'text-lowercase'
    }

  }

})
