import {Component} from 'ng-decorate'
import {defaultLang, defaultNames, defaultWords} from 'models/all'

@Component({
  moduleName: 'app',
  selector: 'app-words',
  scope: {
    names: '=',
    words: '='
  }
})
class VM {
  // Bindable
  names: Fireproof
  words: Fireproof

  // Fields
  lang = defaultLang
  defaultNames = defaultNames
  defaultWords = defaultWords
}
