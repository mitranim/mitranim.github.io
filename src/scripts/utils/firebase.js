var module = angular.module('astil.firebase', [
  'firebase', 'astil.config', 'foliant'
])

module.factory('lang', function(config, $firebaseObject, Traits) {

  var ref = new Firebase(config.fbRootUrl + '/defaults/langs/eng')

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

module.factory('names', function(config, baseArray) {

  var ref = new Firebase(config.fbRootUrl + '/defaults/names/eng')
  var names = new baseArray(ref)
  names.$title = 'Names'

  console.log("-- names:", names);
  window.names = names

  return names

})

module.factory('words', function(config, baseArray) {

  var ref = new Firebase(config.fbRootUrl + '/defaults/words/eng')
  var words = new baseArray(ref)
  words.$title = 'Words'

  console.log("-- words:", words);
  window.words = words

  return words

})

module.factory('baseArray', function($firebaseArray) {

  return $firebaseArray.$extend({

    /**
     * Iterates over primitive values or object values.
     */
    $each: function(iterator) {
      this.$list.forEach(val => {
        if (val.$value != null) iterator(val.$value, val.$id)
        else iterator(val, val.$id)
      })
    },

    /**
     * Reports whether the array has the given primitive value.
     */
    $has: function(value): boolean {
      var yes = false
      try {
        this.$each(val => {if (yes = val === value) throw null})
      } catch (err) {
        if (err !== null) throw err
      }
      return yes
    }

  })

})
