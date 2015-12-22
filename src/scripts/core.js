import _ from 'lodash'
// Core utilities.
import {createAtom, createMb} from 'prax'
// Immutability utilities.
import {immute, replaceAtPath, mergeAtPath} from 'prax'
import {toAsync} from 'prax/async'

/**
 * State
 */

export const atom = toAsync(createAtom(immute({
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
})))

export const {read, watch, stop} = atom

/**
 * Message Bus
 */

const mb = createMb(
  // {type: 'set', path: x => x instanceof Array}, ({value, path}) => {
  //   atom.write(replaceAtPath(read(), value, path))
  // },

  // {type: 'patch'}, ({value, path}) => {
  //   atom.write(mergeAtPath(read(), value, path || []))
  // }
)

export const {match} = mb

// Hack to make `send` safe to use during a `watch` call.
export function send (msg) {
  watch(_.once(() => {mb.send(msg)}))
}

export function set (...path) {
  // send({type: 'set', path, value: path.pop()})
  atom.write(replaceAtPath(read(), path.pop(), path))
}

export function patch (...path) {
  // send({type: 'patch', path, value: path.pop()})
  atom.write(mergeAtPath(read(), path.pop(), path))
}

require('./factors/auth')
require('./factors/generate')

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
  window.atom = atom
  window.read = read
  window.send = send
}
