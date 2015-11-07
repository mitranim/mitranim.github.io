import _ from 'lodash'
import {createStore, applyMiddleware} from 'redux'
import {createMiddleware, immute, replaceAtPath, mergeAtRoot, createReader} from 'symphony'
import {emit} from './utils'

/**
 * Transducing middleware
 */

import {transducer as auth} from './auth'
import {transducer as generate} from './generate'

const create = applyMiddleware(createMiddleware(auth, generate))(createStore)

/**
 * Store
 */

const store = create((state, {type, value, path}) => {
  switch (type) {
    case 'set': {
      return replaceAtPath(state, value, path)
    }
    case 'patch': {
      return mergeAtRoot(state, value)
    }
  }
  return state
},
immute({
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

emit.on('init', _.once(() => {dispatch({type: 'init'})}))

/**
 * Utils
 */

if (window.developmentMode) {
  window.store = store
  window.read = read
}
