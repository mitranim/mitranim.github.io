import Firebase from 'firebase'
import {match, set, patch} from '../core'

let rootRef
let foliantRef

export function getRef (path) {
  return path ? foliantRef.child(path) : null
}

export function defaultRef (type) {
  return foliantRef.child(`defaults/${type}/eng`)
}

match('init', () => {
  rootRef = new Firebase('https://incandescent-torch-3438.firebaseio.com')
  foliantRef = rootRef.child('foliant')

  rootRef.onAuth(authData => {
    if (authData) {
      // Establish personal data refs.
      set(['refPaths'], {
        names: `personal/${authData.uid}/names/eng`,
        words: `personal/${authData.uid}/words/eng`
      })
    } else {
      // Clear personal data.
      patch([], {refPaths: null, words: null, names: null})
      // When deauthed, auth anonymously.
      rootRef.authAnonymously(err => {
        if (err) console.error(err)
      })
    }

    set(['auth'], authData)
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
