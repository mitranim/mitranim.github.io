import _ from 'lodash'
import {spinner} from '../utils'
import {read, send, auto, set} from '../core'

/**
 * Misc
 */

const sourceWord = auto(({text, action}) => (
  ['div', {className: 'foliant--word'}, text,
    action
    ? ['button', {type: 'button',
                  className: 'fa fa-times flat text-gray',
                  onclick: action}] : null]
))

const generatedWord = auto(({text, action}) => (
  ['div', {className: 'foliant--word'},
    action
    ? ['button', {type: 'button',
                  className: 'fa fa-arrow-left flat text-gray',
                  onclick: action}] : null,
    ['span', {className: 'flex-1 text-center'}, text]]
))

/**
 * words
 */

export const words = auto(function words (props, read) {
  const auth = read('auth')
  const kind = read('state', 'kind')
  const inited = read(kind, 'inited')

  if (!auth || !inited) {
    return (
      [spinner, {size: 'large', style: {minHeight: '3em', width: '100%'}}]
    )
  }

  const {selected, generated, depleted, error} = read(kind)
  const word = read('state', 'word')
  const textStyle = kind === 'names' ? 'text-capitalize' : 'text-lowercase'

  return (
    ['div', {className: 'foliant'},
      // Left column: source words
      ['div', {className: 'flex-1 space-out', style: 'padding-right: 0.5rem'},
        ['h3', null, `Source ${_.capitalize(kind)}`],
        ['form', {onsubmit: addWord,
                  className: `layout-row`,
                  'data-tooltip': error,
                  style: 'height: 2.5rem'},
          ['input', {className: `flex-11 ${textStyle}`,
                     placeholder: 'add...',
                     value: word,
                     oninput: changeWord,
                     onblur: clearError}],
          ['button', {className: 'flex-1 fa fa-plus', tabindex: -1}]],
        ['div', {className: `sm-grid-1 md-grid-2 ${textStyle}`},
          _.map(selected, (word, key) => (
            [sourceWord, {text: word, action () { dropWord(key) }, key}]
          ))]],

      // Right column: generated results
      ['div', {className: 'flex-1 space-out', style: 'padding-left: 0.5rem'},
        ['h3', null, `Generated ${_.capitalize(kind)}`],
        ['button', {onclick: generate, className: 'text-center',
                    style: 'height: 2.5rem; width: 100%'},
          ['span', {className: 'fa fa-refresh inline', style: 'float: left'}],
          'Generate'],
        ['div', {className: `sm-grid-1 md-grid-2 ${textStyle}`},
          _.map(generated, word => (
            [generatedWord, {text: word, action () { pickWord(word) }, key: word}]
          )),
          depleted
          ? [generatedWord, {text: '(depleted)'}] : null]]]
  )
})

/**
 * Utils
 */

function changeWord (event) {
  event.preventDefault()
  const input = event.target
  const value = input.value.trim()
  if (input.value !== value) input.value = value
  set(['state', 'word'], value)
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

function generate () {
  const kind = read('state', 'kind')
  send({type: 'gen/generate', kind})
}

function pickWord (word) {
  const kind = read('state', 'kind')
  send({type: 'gen/pick', kind, word})
}

function clearError () {
  const kind = read('state', 'kind')
  set([kind, 'error'], null)
}
