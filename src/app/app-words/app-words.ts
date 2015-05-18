import {Component} from 'ng-decorate'
import {BaseVM} from 'utils/all'
import {defaultLang, defaultNames, defaultWords} from 'models/all'

@Component({
  selector: 'app-words',
  scope: {
    names: '=',
    words: '='
  }
})
class VM extends BaseVM {
  // Bindable
  names: Fireproof
  words: Fireproof

  // Fields
  lang = defaultLang
  defaultNames = defaultNames
  defaultWords = defaultWords
}
