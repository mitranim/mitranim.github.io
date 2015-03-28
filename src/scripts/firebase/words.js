/**
 * Personal words array.
 */

var module = angular.module('astil.firebase')

/**
 * Returns a function that generates a new words array reference for the
 * current user, sets up a $firebaseArray object, and returns it. If the user
 * is not logged in, this returns null.
 */
module.factory('wordsFactory', function(fbRoot, auth, FBArray, defaultWords) {
  return function() {
    var authData = auth.$getAuth()
    if (!authData) return null

    var ref = fbRoot.child(`/personal/${authData.uid}/words/eng`)

    // When no words are left, sync with default words.
    ref.on('value', function(snap) {
      if (_.isEmpty(snap.val())) {
        defaultWords.$loaded().then(() => {
          defaultWords.$each(ref.push, ref)
        })
      }
    })

    var words = new FBArray(ref)
    words.$title = 'Words'

    return words
  }
})
