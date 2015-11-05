import {createStore} from 'redux'
import {immute, replaceAtPath, mergeAtRoot, createReader} from 'symphony'

/**
 * Store
 */

const store = createStore((state, action) => {
  switch (action.type) {
    case 'set': {
      state = replaceAtPath(state, action.value, action.path)
      break
    }
    case 'patch': {
      state = mergeAtRoot(state, action.value)
      break
    }
  }
  return state
}, immute({
  refPaths: {
    names: null,
    words: null
  },

  auth: null,

  defaults: {
    names: null,
    words: null
  },

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
  }
}))

export const dispatch = store.dispatch
export const read = createReader(store)

/**
 * Utils
 */

if (window.developmentMode) {
  window.store = store
  window.read = read
}
