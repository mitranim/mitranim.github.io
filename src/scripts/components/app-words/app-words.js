/**
 * Lets the user input sample words and loads derived synthetic words from the
 * server.
 */

angular.module('astil.components.appWords', [
  'astil.attributes',
  'astil.mixins.generic',
  'astil.models.Mode',
  'astil.models.Word',
  'astil.stores.Langs',
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

.factory('appWordsCtrl', function(mixinGeneric, Mode, Word, Langs) {

  return function() {
    var self = this

    // Use generic controller mixin.
    mixinGeneric(self)

    /**
     * Available languages.
     * @type Langs
     */
    self.langs = Langs

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
        generated: Word.readAll({params: self.params(mode)})
      }).finally(self.ready)
    }

    /**
     * Submit first request on page load. This makes an assumption that the
     * view logic (see template) also selects the first lang and the first
     * mode.
     */
    self.submit(self.langs[0].modes[0])
  }

})
