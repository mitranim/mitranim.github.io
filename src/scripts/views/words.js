import React, {PropTypes} from 'react'
import _ from 'lodash'
import {Spinner} from '../utils'
import {read, send, auto} from '../core'

/**
 * Misc
 */

const SourceWord = auto(({text, action}) => (
  <div className='word'>
    <span>{text}</span>
    {action ?
    <button className='fa fa-times sf-button-flat fade interactive'
            onClick={action} type='button' /> : null}
  </div>
))

SourceWord.propTypes = {
  text: PropTypes.string.isRequired,
  action: PropTypes.func
}

const GeneratedWord = auto(({text, action}) => (
  <div className='word'>
    {action ?
    <button className='fa fa-arrow-left sf-button-flat fade interactive'
            onClick={action} type='button' /> : null}
    <span className='flex-1 text-center'>{text}</span>
  </div>
))

GeneratedWord.propTypes = SourceWord.propTypes

/**
 * Words
 */

export const Words = auto(() => {
  const auth = read('auth')
  const kind = read('state', 'kind')
  const inited = read(kind, 'inited')

  if (!auth || !inited) {
    return (
      <Spinner size='large' style={{minHeight: '3em'}} />
    )
  }

  const {selected, generated, depleted, error} = read(kind)
  const word = read('state', 'word')
  const textStyle = kind === 'names' ? 'text-capitalise' : 'text-lowercase'

  return (
    <div className='widget-words'>
      {/* Left column: source words */}
      <div className='flex-1 container'>
        <h3>Source {_.capitalize(kind)}</h3>
        <form onSubmit={addWord} className='sf-label-row sf-label-dense'
              data-sf-tooltip={error} data-sf-trigger='focus' style={{height: '2.5rem'}}>
          <input autoFocus className={`flex-11 theme-text-primary ${textStyle}`} placeholder='add...'
                 value={word} onChange={changeWord} />
          <button className='flex-1 fa fa-plus theme-primary' tabIndex='-1' />
        </form>
        <div className={`sm-grid-1 md-grid-2 ${textStyle}`}>
          {_.map(selected, (word, key) => (
            <SourceWord text={word} action={() => {dropWord(key)}} key={key} />
          ))}
        </div>
      </div>

      {/* Right column: generated results */}
      <div className='flex-1 container'>
        <h3>Generated {_.capitalize(kind)}</h3>
        <form onSubmit={generate} className='sf-label-row sf-label-dense' style={{height: '2.5rem'}}>
          <button className='flex-1 theme-accent fa fa-refresh' tabIndex='-1' />
          <button className='flex-11 theme-accent row-center-center text-center'>Generate</button>
        </form>
        <div className={`sm-grid-1 md-grid-2 ${textStyle}`}>
          {_.map(generated, word => (
            <GeneratedWord text={word} action={() => {pickWord(word)}} key={word} />
          ))}

          {depleted ?
          <GeneratedWord text='(depleted)' /> : null}
        </div>
      </div>
    </div>
  )
})

/**
 * Utils
 */

function changeWord ({target: {value}}) {
  send({type: 'set', path: ['state', 'word'], value: value.trim()})
}

function addWord (event) {
  event.preventDefault()
  const kind = read('state', 'kind')
  const word = read('state', 'word')
  send({type: 'gen/add', kind, word})
}

function dropWord (key) {
  const kind = read('state', 'kind')
  send({type: 'gen/drop', kind, key})
}

function generate (event) {
  event.preventDefault()
  const kind = read('state', 'kind')
  send({type: 'gen/generate', kind})
}

function pickWord (word) {
  const kind = read('state', 'kind')
  send({type: 'gen/pick', kind, word})
}
