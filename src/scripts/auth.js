import Firebase from 'firebase'
import {emit} from './utils'

let rootRef
let foliantRef

export function transducer (action, dispatch) {
  switch (action.type) {
    case 'logout': {
      rootRef.unauth()
      break
    }

    case 'loginTwitter': {
      rootRef.authWithOAuthRedirect('twitter', err => {
        if (err) console.error(err)
      })
      break
    }

    case 'init': {
      rootRef = new Firebase('https://incandescent-torch-3438.firebaseio.com')
      foliantRef = rootRef.child('foliant')

      exports.getRef = path => path ? foliantRef.child(path) : null

      exports.defaultRefs = {
        lang: foliantRef.child('defaults/langs/eng'),
        names: foliantRef.child('defaults/names/eng'),
        words: foliantRef.child('defaults/words/eng')
      }

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
          emit('loginSuccess')
        } else {
          dispatch(clearPersonalData())
          // When deauthed, auth anonymously.
          rootRef.authAnonymously(err => {
            if (err) console.error(err)
          })
        }
      })

      break
    }
  }

  return action
}

/**
 * Utils
 */

function clearPersonalData () {
  return {
    type: 'patch',
    value: {
      refPaths: null,
      words: null,
      names: null
    }
  }
}
