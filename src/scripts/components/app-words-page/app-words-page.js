var module = angular.module('astil.components.appWordsPage', [
  'astil.controllers.generic',
  'astil.firebase',
  'astil.components.appWords'
])

module.directive('appWordsPage', function(appWordsPageCtrl) {
  return {
    restrict: 'E',
    scope: {},
    templateUrl: 'components/app-words-page/app-words-page.html',
    controllerAs: 'self',
    bindToController: true,
    controller: [appWordsPageCtrl]
  }
})

module.factory('appWordsPageCtrl', function($q, $timeout, fbRoot, auth, CtrlGeneric, defaultLang, namesFactory, wordsFactory) {
  return class extends CtrlGeneric {

    constructor() {
      super()

      /**
       * Loading status.
       * @type Boolean
       */
      this.loading = true

      /**
       * Lang configuration.
       * @type Lang
       */
      this.lang = defaultLang

      /**
       * On successful authentication, regenerate names and words references
       * and render the app-words component.
       */
      fbRoot.onAuth(authData => {
        // Put auth data into scope.
        this.authData = authData

        // Reset data and mark as loading.
        this.reset()
        if (!authData) return

        this.names = namesFactory()
        this.words = wordsFactory()

        $timeout(this.ready)
      })
    }

    authWithTwitter() {
      auth.$authWithOAuthRedirect('twitter')
    }

    /**
     * Logs the user out.
     */
    unauth() {
      this.reset()
      auth.$unauth()
    }

    /**
     * Destroys the currently used firebase objects and sets status to
     * 'loading'.
     */
    reset() {
      if (this.names) this.names.$destroy()
      if (this.words) this.words.$destroy()
      this.loading = true
    }

  }
})
