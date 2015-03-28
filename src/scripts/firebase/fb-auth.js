/**
 * Authentication.
 */

angular.module('astil.firebase')
.factory('auth', function(fbRoot, $firebaseAuth) {

  var auth = $firebaseAuth(fbRoot)

  // Check auth state and login anonymously whenever we're logged out.
  auth.$onAuth(function(authData) {
    if (authData === null) auth.$authAnonymously()
  })

  return auth

})
