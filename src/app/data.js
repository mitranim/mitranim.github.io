import Firebase from 'firebase'
import React from 'react'
import _ from 'lodash'
import {Tracker} from './tracker'
import {ReactiveVar} from './reactive-var'
// import {ReactiveDict} from './reactive-dict'

const fbRootUrl = 'https://incandescent-torch-3438.firebaseio.com'

/**
 * References.
 */

export const root = new Firebase(fbRootUrl)

const RefMappers = {
  defaultLang: authData => root.child('foliant/defaults/langs/eng'),
  defaultNames: authData => root.child('foliant/defaults/names/eng'),
  defaultWords: authData => root.child('foliant/defaults/words/eng'),
  names: authData => authData ? root.child(`foliant/personal/${authData.uid}/names/eng`) : null,
  words: authData => authData ? root.child(`foliant/personal/${authData.uid}/words/eng`) : null
}

/**
 * Reactive values.
 */

export const authData = new ReactiveVar(null)
export const Refs = _.mapValues(RefMappers, () => new ReactiveVar(null))
export const Values = _.mapValues(RefMappers, () => new ReactiveVar(null))

/**
 * Set up lazy data load. This function should be run when the data is
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

    authData.set(newAuthData)

    _.each(Refs, (refVar, key) => {
      // Refresh ref.
      const ref = RefMappers[key](newAuthData)
      refVar.set(ref)

      // Refresh value.
      if (ref) {
        const handler = ref.on('value', snap => {
          Values[key].set(snap.val())
        }, () => {
          ref.off('value', handler)
        })
      }
    })
  })

  // Reactively refresh names and words.
  Tracker.autorun(function () {
    const namesRef = Refs.names.get()
    if (namesRef) {
      namesRef.on('value', snap => {
        if (!snap.val()) {
          const defNamesRef = Refs.defaultNames.get()
          const handler = defNamesRef.once('value', snap => {
            namesRef.set(snap.val())
          }, () => {
            namesRef.off('value', handler)
          })
        }
      })
    }

    const wordsRef = Refs.words.get()
    if (wordsRef) {
      wordsRef.on('value', snap => {
        if (!snap.val()) {
          const defWordsRef = Refs.defaultWords.get()
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

/**
 * Component extension.
 */

export class Component extends React.Component {
  componentWillMount () {
    if (typeof this.getState === 'function') {
      Tracker.autorun(() => {
        // Assuming `this.getState()` accesses reactive data sources.
        this.setState(this.getState())
      })
    }
  }

  componentWillUnmount () {
    // ... TODO cleanup?
  }
}
