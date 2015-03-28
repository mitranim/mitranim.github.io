/**
 * Personal names array.
 */

angular.module('astil.firebase')
.factory('namesPromise', function(fbRoot, auth, FBArray, defaultNames) {
  return auth.$waitForAuth().then(function recur(authData) {

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

  })
})
