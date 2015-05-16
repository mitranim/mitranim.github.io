import 'angular'
import {app} from 'app'

import 'views'
import 'attributes/all'
import 'app-words/app-words'
import 'app-words-page/app-words-page'
import 'app-words-tab/app-words-tab'
import 'models/all'
import 'utils/all'

// Manual bootstrapping.
angular.element(document).ready(() => {
  angular.bootstrap(document.body, [app.name], {
    strictDi: true
  })
})
