import _ from 'lodash'
import Firebase from 'firebase'
import {createSignals, createDecorator} from 'symphony'
import {dispatch} from './store'
export * from './store'

/**
 * Signals
 */

export const signals = createSignals(dispatch, {  // eslint-disable-line
  init: {},
  login: {twitter: {}, facebook: {}},
  logout: {},
  names: {init: {}, generate: {}, add: {}, pick: {}, drop: {}},
  words: {init: {}, generate: {}, add: {}, pick: {}, drop: {}},
  didAdd: {}
})

/**
 * Subscriptions
 */

// Delaying initialisation spares us from establishing a Firebase connection
// too early (i.e. on most pages).
signals.init.subscribe(_.once(() => {

  /**
   * Refs
   */

  const rootRef = new Firebase('https://incandescent-torch-3438.firebaseio.com')
  const foliantRef = rootRef.child('foliant')

  exports.getRef = path => path ? foliantRef.child(path) : null

  exports.defaultRefs = {
    lang: foliantRef.child('defaults/langs/eng'),
    names: foliantRef.child('defaults/names/eng'),
    words: foliantRef.child('defaults/words/eng')
  }

  /**
   * Auth
   */

  rootRef.onAuth(authData => {
    dispatch({
      type: 'set',
      path: ['auth'],
      value: authData
    })

    if (authData) {
      // Establish personal data refs.
      dispatch({
        type: 'set',
        path: ['refPaths'],
        value: {
          names: `personal/${authData.uid}/names/eng`,
          words: `personal/${authData.uid}/words/eng`
        }
      })
    } else {
      clearPersonalData()
      // When deauthed, auth anonymously.
      rootRef.authAnonymously(err => {
        if (err) console.error(err)
      })
    }

    signals.login()
  })

  signals.login.twitter.subscribe((__, out) => {
    out(new Promise((resolve, reject) => {
      rootRef.authWithOAuthRedirect('twitter', err => {
        if (err) reject(err)
        else resolve()
      })
    }))
  })

  signals.logout.subscribe(() => {rootRef.unauth()})

  // Set up word generation.
  require('./generate')
}))

/**
 * Decorators
 */

export const on = createDecorator('subscribe', signals)
// export const done = createDecorator('done', signals)

/**
 * Utils
 */

function clearPersonalData () {
  dispatch({
    type: 'patch',
    value: {
      refPaths: null,
      words: null,
      names: null
    }
  })
}

if (window.developmentMode) {
  window.signals = signals
}
