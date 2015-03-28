/**
 * Personal words array.
 */

angular.module('astil.firebase')
.factory('wordsPromise', function(fbRoot, auth, FBArray, defaultWords) {
  return auth.$waitForAuth().then(function recur(authData) {

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

  })
})
