import _ from 'lodash'
import Traits from 'foliant'
import {Component} from 'ng-decorate'
import {BaseVM} from 'utils/all'

@Component({
  moduleName: 'app',
  selector: 'app-words-tab',
  scope: {
    title: '@',
    lang: '=',
    words: '=',
    defaultWords: '='
  }
})
class VM extends BaseVM {
  /**
   * Bindable
   */
  title: string
  lang: Fireproof
  words: Fireproof
  defaultWords: Fireproof

  /**
   * Fields
   */
  loading: boolean = true
  langVal: Lang
  wordsVal: {[key: string]: string}
  defaultWordsVal: {[key: string]: string}
  word: string = ''
  error: string = null
  depleted: boolean = false
  // Words generator
  gen: () => string
  // Generated words
  results: string[] = []

  constructor() {
    super()

    this.$q.all(this.sync({
      langVal: this.lang,
      wordsVal: this.words,
      defaultWordsVal: this.defaultWords
    }))
    // First generation. Produces nothing if the words set is empty.
    .then(this.generate.bind(this))
    .then(() => this.words.on('value', snap => {
      // When the stored words are depleted, refresh from the default collection.
      if (_.isEmpty(snap.val())) {
        return this.words.set(this.defaultWordsVal).then(this.generate.bind(this))
      }
    }))
    .then(this.ready)
    .catch(console.warn.bind(console))
  }

  /**
   * Word count limit.
   */
  get limit(): number {return 12}

  /**
   * Produces a traits object.
   */
  getTraits(): Traits {
    // Account for lang properties.
    var traits = new Traits()
    if (this.langVal.knownSounds && this.langVal.knownSounds.length) {
      traits.knownSounds = new Traits.StringSet(this.langVal.knownSounds)
    }
    if (this.langVal.knownVowels && this.langVal.knownVowels.length) {
      traits.knownVowels = new Traits.StringSet(this.langVal.knownVowels)
    }

    // Account for word characteristics.
    traits.examine(_.toArray(this.wordsVal))

    return traits
  }

  /**
   * Adds the word or displays an error message.
   */
  add(): void {
    this.word = this.word.toLowerCase()

    if (!this.word) {
      this.error = 'Please input a word.'
      return
    }

    if (this.word.length < 2) {
      this.error = 'The word is too short.'
      return
    }

    if (_.contains(this.wordsVal, this.word)) {
      this.error = 'This word is already in the set.'
      return
    }

    try {
      this.getTraits().examine([this.word])
    } catch (err) {
      console.warn('-- word parsing error:', err)
      this.error = 'Some of these characters are not allowed in a word.'
      return
    }

    this.error = ''
    var ref = this.words.push(this.word)
    this.wordsVal[ref.key()] = this.word
    ref.then(_.noop, () => {delete this.wordsVal[ref.key()]})
    this.word = ''

    // Refresh the generator.
    this.gen = this.getTraits().generator()
  }

  /**
   * Generates a group of words.
   */
  generate() {
    // Remove error, if any.
    this.error = ''

    // Regenerate the generator, if necessary.
    if (!this.gen) this.gen = this.getTraits().generator()
    var words = []

    while (words.length < this.limit) {
      var word = this.gen()
      if (!word) break
      // Skip source words.
      if (_.contains(this.wordsVal, word)) continue
      words.push(word)
    }

    if (words.length < this.limit) this.depleted = true
    else this.depleted = false

    this.results = words
  }

  /**
   * Adds the given word to the store, removing it from the generated results.
   * Doesn't refresh the generator because adding a previously generated word
   * to the same source set has no effect on the total output.
   */
  pick(word) {
    if (_.contains(this.wordsVal, word)) return
    var ref = this.words.push(word)
    // Move from results to words.
    this.wordsVal[ref.key()] = word
    _.pull(this.results, word)
    // On failure, move back (disregarding order).
    ref.then(_.noop, () => {
      delete this.wordsVal[ref.key()]
      this.results.push(word)
    })
  }

  /**
   * Removes the given word from the store and refreshes the generator.
   */
  drop(key: string) {
    this.words.child(key).remove().then(() => {
      this.gen = this.getTraits().generator()
    })
  }
}
