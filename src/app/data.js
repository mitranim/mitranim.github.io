import Firebase from 'firebase'
import _ from 'lodash'
import {autorun, Source} from 'rapt'

const fbRootUrl = 'https://incandescent-torch-3438.firebaseio.com'

/**
 * References.
 */

export const root = new Firebase(fbRootUrl)

const refMappers = {
  defaultLang: authData => root.child('foliant/defaults/langs/eng'),
  defaultNames: authData => root.child('foliant/defaults/names/eng'),
  defaultWords: authData => root.child('foliant/defaults/words/eng'),
  names: authData => authData ? root.child(`foliant/personal/${authData.uid}/names/eng`) : null,
  words: authData => authData ? root.child(`foliant/personal/${authData.uid}/words/eng`) : null
}

/**
 * Reactive values.
 */

export const authData = new Source(null)

export const refs = Object.create(null)
export const store = Object.create(null)

Object.keys(refMappers).forEach(key => {
  ;[refs, store].forEach(object => {
    const source = new Source(null)

    Object.defineProperty(object, key, {
      get () {
        return source.read()
      },
      set (value) {
        source.write(value)
      },
      enumerable: true,
      configurable: false
    })
  })
})

/**
 * Lazily set up data loading. This function should be run when the data is
 * required for the first time.
 */
export const setUpDataLoad = _.once(function () {

  /**
   * Auth handlers.
   */

  root.onAuth(newAuthData => {
    // When deauthed, auth anonymously.
    if (!newAuthData) root.authAnonymously(::console.error)

    /**
     * Refresh all reactive variables.
     */

    authData.write(newAuthData)

    _.each(refMappers, (mapper, key) => {
      // Refresh ref.
      const ref = refs[key] = mapper(newAuthData)

      // Refresh value.
      if (ref) {
        const handler = ref.on('value', snap => {
          store[key] = snap.val()
        }, () => {
          ref.off('value', handler)
        })
      }
    })
  })

  // Reactively refresh names and words.
  autorun(function () {
    const namesRef = refs.names
    if (namesRef) {
      namesRef.on('value', snap => {
        if (!snap.val()) {
          const defNamesRef = refs.defaultNames
          const handler = defNamesRef.once('value', snap => {
            namesRef.set(snap.val())
          }, () => {
            namesRef.off('value', handler)
          })
        }
      })
    }

    const wordsRef = refs.words
    if (wordsRef) {
      wordsRef.on('value', snap => {
        if (!snap.val()) {
          const defWordsRef = refs.defaultWords
          const handler = defWordsRef.once('value', snap => {
            wordsRef.set(snap.val())
          }, () => {
            wordsRef.off('value', handler)
          })
        }
      })
    }
  })
})
