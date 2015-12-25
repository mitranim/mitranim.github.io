import {createAtom, createMb} from 'prax'
import {asyncStrategy} from 'prax/async'

/**
 * State
 */

export const atom = createAtom({
  refPaths: {
    names: null,
    words: null
  },

  auth: null,

  defaults: {
    names: null,
    words: null
  },

  kinds: ['names', 'words'],

  names: {
    inited: false,
    selected: null,
    generated: null,
    depleted: null,
    error: null
  },

  words: {
    inited: false,
    selected: null,
    generated: null,
    depleted: null,
    error: null
  },

  // Misc view states.
  state: {
    kind: 'names',
    word: ''
  }
}, asyncStrategy)

export const {read, set, patch, watch, stop} = atom

/**
 * Message Bus
 */

const mb = createMb()

export const {send, match} = mb

/**
 * App Logic
 */

require('./factors/auth')
require('./factors/generate')

/**
 * Rendering
 */

export function auto (view) {
  return function component (render, props) {
    function update (read) {render(view(props, read))}
    watch(update)
    return function unsub () {stop(update)}
  }
}

/**
 * Utils
 */

if (window.developmentMode) {
  window.atom = atom
  window.read = read
  window.send = send
}
