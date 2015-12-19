import _ from 'lodash'
// Core utilities.
import {createAtom, createFq} from 'prax'
// Immutability utilities.
import {immute, replaceAtPath, mergeAtPath} from 'prax'

/**
 * State
 */

export const atom = createAtom(immute({
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
}))

export const {read, watch, stop} = atom

/**
 * FQ
 */

import auth from './factors/auth'
import generate from './factors/generate'

const writer = read => next => msg => {
  if (msg === 'init' || msg === 'auth/logout') return
  const {type, value, path} = msg

  switch (type) {
    case 'set':
      next(replaceAtPath(read(), value, path))
      break
    case 'patch':
      next(mergeAtPath(read(), value, path || []))
      break
    default:
      console.warn('Discarding unrecognised message:', msg)
  }
}

const fq = createFq(auth, generate, writer)

const fqSend = fq(atom.read, atom.write)

// Hack to make `send` safe to use during a `watch` call.
export function send (msg) {
  watch(_.once(() => {fqSend(msg)}))
}

/**
 * Rendering
 */

export const auto = view => (render, props) => {
  const update = watch(() => {render(view(props))})
  return () => {stop(update)}
}

/**
 * Utils
 */

if (window.developmentMode) {
  window.read = read
  window.send = send
}
