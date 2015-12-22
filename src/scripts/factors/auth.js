import Firebase from 'firebase'
import {send, match, set, patch} from '../core'

let rootRef
let foliantRef

export function getRef (path) {
  return path ? foliantRef.child(path) : null
}

export const defaultRefs = {
  get lang () {return foliantRef.child('defaults/langs/eng')},
  get names () {return foliantRef.child('defaults/names/eng')},
  get words () {return foliantRef.child('defaults/words/eng')}
}

match('init', () => {
  rootRef = new Firebase('https://incandescent-torch-3438.firebaseio.com')
  foliantRef = rootRef.child('foliant')

  rootRef.onAuth(authData => {
    if (authData) {
      // Establish personal data refs.
      // send({
      //   type: 'set',
      //   path: ['refPaths'],
      //   value: {
      //     names: `personal/${authData.uid}/names/eng`,
      //     words: `personal/${authData.uid}/words/eng`
      //   }
      // })
      set('refPaths', {
        names: `personal/${authData.uid}/names/eng`,
        words: `personal/${authData.uid}/words/eng`
      })
    } else {
      // Clear personal data.
      // send({type: 'patch', value: {refPaths: null, words: null, names: null}})
      patch({refPaths: null, words: null, names: null})
      // When deauthed, auth anonymously.
      rootRef.authAnonymously(err => {
        if (err) console.error(err)
      })
    }

    // send({type: 'set', path: ['auth'], value: authData})
    set('auth', authData)
  })
})

match('auth/logout', () => {
  rootRef.unauth()
})

match('auth/loginTwitter', () => {
  rootRef.authWithOAuthRedirect('twitter', err => {
    if (err) console.error(err)
  })
})
