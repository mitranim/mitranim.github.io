import {Component, bindOneWay} from 'ng-decorate'
import {BaseVM} from 'utils/all'
import {defaultLang, defaultNames, defaultWords} from 'models/all'

@Component({
  selector: 'app-words'
})
class VM extends BaseVM {
  @bindOneWay() names: Fireproof
  @bindOneWay() words: Fireproof

  // Fields
  lang = defaultLang
  defaultNames = defaultNames
  defaultWords = defaultWords
}
