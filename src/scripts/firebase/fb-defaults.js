/**
 * Default public datasets.
 */

var module = angular.module('astil.firebase')

/**
 * Default language.
 */
module.factory('defaultLang', function(fbRoot, $firebaseObject, Traits) {

  var ref = fbRoot.child('/defaults/langs/eng')

  var Lang = $firebaseObject.$extend({
    // Produces a customised traits object.
    $traits(): Traits {
      var traits = new Traits()
      if (this.knownSounds && this.knownSounds.length) {
        traits.knownSounds = new Traits.StringSet(this.knownSounds)
      }
      if (this.knownVowels && this.knownVowels.length) {
        traits.knownVowels = new Traits.StringSet(this.knownVowels)
      }
      return traits
    }
  })

  return new Lang(ref)

})

/**
 * Default names array.
 */
module.factory('defaultNames', function(fbRoot, FBArray) {

  var ref = fbRoot.child('/defaults/names/eng')
  var names = new FBArray(ref)
  names.$title = 'Names'

  return names

})

/**
 * Default words array.
 */
module.factory('defaultWords', function(fbRoot, FBArray) {

  var ref = fbRoot.child('/defaults/words/eng')
  var words = new FBArray(ref)
  words.$title = 'Words'

  return words

})
