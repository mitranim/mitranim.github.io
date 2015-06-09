import {Component, autoinject} from 'ng-decorate'
import {BaseVM} from 'utils/all'
import {root, authData, defaultLang, makeNames, makeWords} from 'models/all'

@Component({
  selector: 'app-words-page'
})
class VM extends BaseVM {
  @autoinject $timeout: ng.ITimeoutService

  loading: boolean = true
  $ready: boolean = false
  authData = authData
  names: Fireproof
  words: Fireproof

  /**
   * Lang configuration.
   */
  lang = defaultLang

  constructor() {
    super()

    /**
     * On successful authentication, regenerate names and words references
     * and render the app-words component.
     */
    root.onAuth(authData => {
      // Mark as loading.
      this.loading = true
      this.$ready = false
      if (!authData) return

      this.names = makeNames(authData.uid)
      this.words = makeWords(authData.uid)

      this.$timeout(this.ready)
    })
  }

  authWithTwitter() {
    root.authWithOAuthRedirect('twitter')
  }
}
