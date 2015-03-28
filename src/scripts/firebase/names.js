/**
 * Personal names array.
 */

var module = angular.module('astil.firebase')

/**
 * Returns a function that generates a new names array reference for the
 * current user, sets up a $firebaseArray object, and returns it. If the user
 * is not logged in, this returns null.
 */
module.factory('namesFactory', function(fbRoot, auth, FBArray, defaultNames) {
  return function() {
    var authData = auth.$getAuth()
    if (!authData) return null

    var ref = fbRoot.child(`/personal/${authData.uid}/names/eng`)

    // When no names are left, sync with default names.
    ref.on('value', function(snap) {
      if (_.isEmpty(snap.val())) {
        defaultNames.$loaded().then(() => {
          defaultNames.$each(ref.push, ref)
        })
      }
    })

    var names = new FBArray(ref)
    names.$title = 'Names'

    return names
  }
})
