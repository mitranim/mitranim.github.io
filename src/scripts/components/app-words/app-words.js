/**
 * Lets the user input sample words and loads derived synthetic words from the
 * server.
 */

angular.module('astil.components.appWords', [
  'astil.attributes',
  'astil.mixins.generic',

  'astil.models.Mode',
  'astil.models.Word',

  'astil.stores.Lang',
  'astil.stores.Mode',
  'astil.stores.Word',
])

.directive('appWords', function(appWordsCtrl) {
  return {
    restrict: 'E',
    scope: {},
    templateUrl: 'components/app-words/app-words.html',
    controllerAs: 'self',
    bindToController: true,
    controller: appWordsCtrl
  }
})

.factory('appWordsCtrl', function(mixinGeneric, Mode, Word, Langs, Modes, Words) {

  return function() {
    var self = this

    // Use generic controller mixin.
    mixinGeneric(self)

    /**
     * Available languages.
     * @type Langs
     */
    self.langs = Langs.records

    // Removes the given word from the given mode and from the word store.
    self.drop = function(mode: Mode, word: Word): void {
      _.pull(mode.$source, word)
      _.pull(mode.$generated, word)
      _.pull(Words.records, word)
      Words.$saveLS()
    }

    // Moves the given word from a mode's $generated to $source.
    self.pick = function(mode: Mode, word: Word): void {
      _.pull(mode.$generated, word)
      self.add(mode, word.Value)
    }

    // Converts the given string to a word and adds it to the given mode's
    // $source, as well as to the word store.
    self.add = function(mode: Mode, string: string): void {
      var value = string.toLowerCase().trim()
      if (!value) {
        mode.$error = 'Please input a word.'
        return
      }

      var word = new Word({Value: value, ModeId: mode.Id})

      if (word.Value.length < 2) {
        mode.$error = 'The word is too short.'
        return
      }

      if (!word.$valid()) {
        mode.$error = 'Some of these characters are not allowed in a word.'
        return
      }

      if (_.some(mode.$source, {Value: word.Value})) {
        mode.$error = 'This word is already in the set.'
        return
      }

      // Add to mode.
      mode.$error = ''
      mode.$source.push(word)
      mode.$word = ''

      // Add to store.
      Words.records.push(word)
      Words.$saveLS()
    }

    /**
     * Generates request parameters.
     */
    self.params = function(mode: Mode): {} {
      return {
        words:    mode.words(),
        soundset: mode.soundset || null
      }
    }

    /**
     * Loads words for the given mode.
     */
    self.submit = function(mode: Mode) {
      // Ignore if we're already making a request.
      if (self.loading) return

      self.loadTo(mode, {
        $generated: Word.readAll({params: self.params(mode)})
      }).finally(self.ready)
    }

    /**
     * Submit first request on page load. This makes an assumption that the
     * view logic also selects the first lang and the first mode.
     */
    self.submit(self.langs[0].$modes[0])
  }

})
