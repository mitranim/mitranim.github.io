import React, {PropTypes} from 'react'
import _ from 'lodash'
import {reactiveRender} from 'symphony'
import {LoginButton} from './login'
import {renderTo, Spinner, pure} from './utils'
import {read, signals, on} from './flow'

const SourceWord = pure(props => (
  <div className='word'>
    <span>{props.text}</span>
    {props.handler ?
    <button className='fa fa-times sf-button-flat fade interactive'
            onClick={props.handler} type='button' /> : null}
  </div>
))

const GeneratedWord = pure(props => (
  <div className='word'>
    {props.handler ?
    <button className='fa fa-arrow-left sf-button-flat fade interactive'
            onClick={props.handler} type='button' /> : null}
    <span className='flex-1 text-center'>{props.text}</span>
  </div>
))

@reactiveRender
class WordsTab extends React.Component {
  static propTypes = {
    type: PropTypes.string
  }

  @on.login
  init () {
    this.sig = signals[this.props.type]
    if (this.sig) this.sig.init()
  }

  componentWillMount () {
    this.init()
  }

  componentDidUpdate (props) {
    if (this.props.type !== props.type) this.init()
  }

  render () {
    const auth = read('auth')
    const selected = read(this.props.type, 'selected')
    const inited = read(this.props.type, 'inited')
    const generated = read(this.props.type, 'generated')
    const depleted = read(this.props.type, 'depleted')
    const error = read(this.props.type, 'error')

    if (!auth || !inited) {
      return (
        <Spinner size='large' style={{minHeight: '3em'}} />
      )
    }

    return (
      <div className='widget-words'>
        {/* Left column: source words */}
        <div className='flex-1 container'>
          <h3>Source {_.capitalize(this.props.type)}</h3>
          <form onSubmit={::this.add} className='sf-label-row sf-label-dense'
                data-sf-tooltip={error} data-sf-trigger='focus' style={{height: '2.5rem'}}>
            <input name='word' ref='input' autoFocus className={`flex-11 theme-text-primary ${this.textStyle}`} placeholder='add...' />
            <button className='flex-1 fa fa-plus theme-primary' tabIndex='-1' />
          </form>
          <div className={`sm-grid-1 md-grid-2 ${this.textStyle}`}>
            {_.map(selected, (word, key) => (
              <SourceWord text={word} handler={() => this.sig.drop(key)} key={key} />
            ))}
          </div>
        </div>

        {/* Right column: generated results */}
        <div className='flex-1 container'>
          <h3>Generated {_.capitalize(this.props.type)}</h3>
          <form onSubmit={::this.generate} className='sf-label-row sf-label-dense' style={{height: '2.5rem'}}>
            <button className='flex-1 theme-accent fa fa-refresh' tabIndex='-1' />
            <button className='flex-11 theme-accent row-center-center text-center'>Generate</button>
          </form>
          <div className={`sm-grid-1 md-grid-2 ${this.textStyle}`}>
            {_.map(generated, word => (
              <GeneratedWord text={word} handler={() => {this.sig.pick(word)}} key={word} />
            ))}

            {depleted ?
            <GeneratedWord text='(depleted)'/> : null}
          </div>
        </div>
      </div>
    )
  }

  add (event) {
    event.preventDefault()
    const value = this.refs.input.value.trim()
    this.sig.add(value)
  }

  @on.didAdd
  onAdd () {
    this.refs.input.value = ''
  }

  generate (event) {
    event.preventDefault()
    this.sig.generate()
  }

  get textStyle () {return this.props.type === 'names' ? 'text-capitalise' : 'text-lowercase'}
}

@renderTo('[data-render-foliant]')
export class WordsPage extends React.Component {
  tabs = ['names', 'words']
  state = {tab: 'names'}

  render () {
    return (
      <div>
        <div className='sf-navbar sf-navbar-tabs'>
          {this.tabs.map(tab => (
            <a className={`interactive ${tab === this.state.tab ? 'active' : ''}`}
               onClick={() => {this.setState({tab})}} key={tab}>
              <h3>{_.capitalize(tab)}</h3>
            </a>
          ))}
        </div>

        <WordsTab type={this.state.tab} />

        <br />

        <LoginButton />
      </div>
    )
  }
}
