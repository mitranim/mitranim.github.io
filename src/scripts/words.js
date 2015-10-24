import React from 'react'
import _ from 'lodash'
import Traits from 'foliant'
import {LoginButton} from './login'
import {refs, store} from './data'
import {renderTo, reactive, Spinner} from './utils'

@renderTo('[data-render-foliant]')
export class Words extends React.Component {
  @reactive
  updateState () {
    this.setState({
      defaultLang: store.defaultLang,
      defaultNames: store.defaultNames,
      defaultWords: store.defaultWords,
      names: store.names,
      words: store.words
    })
  }

  isReady () {
    return !_.any(this.state, _.isEmpty)
  }

  render () {
    if (!this.isReady()) {
      return (
        <Spinner size='large' style={{minHeight: '3em'}} />
      )
    }

    return (
      <div>
        <div className='sf-tabset'>
          <div className='sf-tabset-head'>
            <div className='sf-tab-head active'><h3>Names</h3></div>
            <div className='sf-tab-head'><h3>Words</h3></div>
          </div>

          <div className='sf-tabset-body'>
            <div className='sf-tab-body active' style={{padding: 0}}>
              <WordsTab title='Names' lang={this.state.defaultLang}
                        default={this.state.defaultNames} current={this.state.names} />
            </div>
            <div className='sf-tab-body' style={{padding: 0}}>
              <WordsTab title='Words' lang={this.state.defaultLang}
                        default={this.state.defaultWords} current={this.state.words} />
            </div>
          </div>
        </div>

        <br />

        <LoginButton />
      </div>
    )
  }
}

const limit = 12

class WordsTab extends React.Component {
  state = {
    results: [],
    error: null,
    depleted: false
  }

  generator = null
  lastCurrent = null

  componentWillReceiveProps (props) {
    this.props = props
    this.refresh()
  }

  componentWillMount () {
    this.refresh()
  }

  render () {
    return (
      <div className='widget-words'>
        {/* Left column: source words */}
        <div className='flex-1 container'>
          <h3>Source {this.props.title}</h3>
          <form onSubmit={::this.add} className='sf-label-row sf-label-dense'
                data-sf-tooltip={this.state.error} data-sf-trigger='focus' style={{height: '2.5rem'}}>
            <input name='word' autoFocus className={`flex-11 theme-text-primary ${this.textStyle}`} placeholder='add...' />
            <button className='flex-1 fa fa-plus theme-primary' tabIndex='-1'></button>
          </form>
          <div className={`sm-grid-1 md-grid-2 ${this.textStyle}`}>
            {_.map(this.props.current, (word, key) => (
              <SourceWord text={word} handler={() => this.drop(key)} key={key} />
            ))}
          </div>
        </div>

        {/* Right column: generated results */}
        <div className='flex-1 container'>
          <h3>Generated {this.props.title}</h3>
          <form onSubmit={::this.generate} className='sf-label-row sf-label-dense' style={{height: '2.5rem'}}>
            <button className='flex-1 theme-accent fa fa-refresh' tabIndex='-1'></button>
            <button className='flex-11 theme-accent row-center-center text-center'>Generate</button>
          </form>
          <div className={`sm-grid-1 md-grid-2 ${this.textStyle}`}>
            {_.map(this.state.results, word => (
              <GeneratedWord text={word} handler={() => this.pick(word)} key={word} />
            ))}
            {this.state.depleted ?
              <GeneratedWord text='(depleted)'/> : null}
          </div>
        </div>
      </div>
    )
  }

  refresh () {
    if (_.isEqual(this.props.current, this.lastCurrent)) return
    this.generator = getGenerator(this.props.lang, this.props.current, limit)
    this.state.results = this.generator()
    this.state.depleted = this.state.results.length < limit
    this.lastCurrent = this.props.current
  }

  generate (event) {
    event.preventDefault()
    const results = this.generator()
    this.setState({
      results: results,
      depleted: results.length < limit
    })
  }

  // Adds the word or displays an error message.
  add (event) {
    event.preventDefault()
    const word = event.target.word.value.toLowerCase()

    if (!word) {
      this.setState({error: 'Please input a word'})
      return
    }

    if (word.length < 2) {
      this.setState({error: 'The word is too short'})
      return
    }

    if (_.contains(this.props.current, word)) {
      this.setState({error: 'This word is already in the set'})
      return
    }

    if (!isWordValid(this.props.lang, word)) {
      this.setState({error: 'Some of these characters are not allowed in a word'})
      return
    }

    this.setState({error: null})
    this.ref.push(word)
    event.target.word.value = ''
  }

  // Adds the given word to the store, removing it from the generated results.
  // Expects the "new" set of source words to become a union of current+word to
  // avoid refreshing the generator and the generated words, because adding a
  // previously generated word to the same source set has no effect on the total
  // output.
  pick (word) {
    if (_.contains(this.props.current, word)) return
    const oldResults = this.state.results

    const ref = this.ref.push(word, err => {
      if (err) {
        delete this.lastCurrent[ref.key()]
        this.setState({results: oldResults})
      }
    })

    this.lastCurrent[ref.key()] = word
    this.setState({
      results: _.without(this.state.results, word)
    })
  }

  // Removes the given word from the store and implicitly refreshes the
  // generator and the results.
  drop (key: string) {
    this.ref.child(key).remove(err => {
      if (!err) this.ref.once('value', snap => {
        if (_.isEmpty(snap.val())) this.ref.set(this.props.default)
      })
    })
  }

  get textStyle () {return this.props.title === 'Names' ? 'text-capitalise' : 'text-lowercase'}
  get ref () {return refs[this.props.title.toLowerCase()]}
}

class SourceWord extends React.Component {
  render () {
    return (
      <div className='word'>
        <span>{this.props.text}</span>
        <button className='fa fa-times sf-button-flat fade interactive'
                onClick={this.props.handler} type='button'></button>
      </div>
    )
  }
}

class GeneratedWord extends React.Component {
  render () {
    return (
      <div className='word'>
        {this.props.handler ?
        <button className='fa fa-arrow-left sf-button-flat fade interactive'
                onClick={this.props.handler} type='button'></button> : null}
        <span className='flex-1 text-center'>{this.props.text}</span>
      </div>
    )
  }
}

function getGenerator (lang, words, limit) {
  const traits = getTraits(lang)
  traits.examine(_.toArray(words))
  const generator = traits.generator()

  return function () {
    const results = []
    let word = ''
    while ((word = generator()) && results.length < limit) {
      if (!_.contains(words, word)) results.push(word)
    }
    return results
  }
}

function getTraits (lang) {
  const traits = new Traits()
  if (lang.knownSounds && lang.knownSounds.length) {
    traits.knownSounds = new Traits.StringSet(lang.knownSounds)
  }
  if (lang.knownVowels && lang.knownVowels.length) {
    traits.knownVowels = new Traits.StringSet(lang.knownVowels)
  }
  return traits
}

function isWordValid (lang, word) {
  try {
    const traits = getTraits(lang)
    traits.examine([word])
    return true
  } catch (err) {
    console.warn(err)
    return false
  }
}
