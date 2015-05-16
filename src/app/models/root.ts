/**
 * Shared references and classes.
 */

import 'angular'
import Firebase from 'firebase'
import Fireproof from 'fireproof'
import {app} from 'app'
import {config} from 'utils/all'

// Bless Fireproof with $q to integrate it with our digest cycle.
var $q: ng.IQService = angular.injector(['ng']).get('$q')
Fireproof.bless($q)

// Refresh with proper $q after app bootstrap.
app.run(['$q', function($q) {
    Fireproof.bless($q)
}])

// Root reference for our firebase.
export var fbRoot = new Firebase(config.fbRootUrl)
export var root = new Fireproof(fbRoot)

// Utils for debugging.
if (config.dev) {
  window.root = root
  window.log = function log(value) {console.log(value)}
  window.val = function log(value) {console.log(value.val())}
  window.warn = function warn(value) {console.warn(value)}
}
