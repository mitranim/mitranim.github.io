import _ from 'lodash'
import Traits from 'foliant'
import {read, signals, dispatch} from './flow'

// Used for lazily exported Firebase refs.
const flow = require('./flow')

const limit = 12

;['names', 'words'].forEach(type => {
  const sig = signals[type]
  let generator = null

  signals.logout.subscribe(() => {
    generator = null
  })

  sig.init.subscribe((__, out) => {
    const inited = read(type, 'inited')
    if (inited) return

    // Create a subscription for personal selected words, if it doesn't exist
    // yet.
    const ref = flow.getRef(read('refPaths', type))
    if (!ref) return

    out(new Promise(resolve => {
      const handler = ref.on('value', snap => {
        dispatch({
          type: 'set',
          path: [type, 'selected'],
          value: snap.val()
        })
        resolve()
      }, () => {
        resolve()
        ref.off('value', handler)
      })

      // Completely paranoid cleanup.
      const off = signals.logout.subscribe(() => {
        ref.off('value', handler)
        off()
      })
    }))

    out({
      type: 'set',
      path: [type, 'inited'],
      value: true
    })

    out(() => checkForEmpty(type))

    out(() => {
      if (!read(type, 'generated')) sig.generate()
    })
  })

  sig.generate.subscribe(() => {
    const selected = read(type, 'selected')
    if (!selected) return

    if (!generator) {
      generator = getGenerator(selected, limit)
    }
    const generated = generator()

    dispatch({
      type: 'patch',
      value: {
        [type]: {
          generated,
          depleted: generated.length < limit
        }
      }
    })
  })

  sig.add.subscribe(word => {
    word = word.toLowerCase()

    if (!word) {
      dispatch(err(type, 'Please input a word'))
      return
    }

    if (word.length < 2) {
      dispatch(err(type, 'The word is too short'))
      return
    }

    const selected = read(type, 'selected')

    if (_.contains(selected, word)) {
      dispatch(err(type, 'This word is already in the set'))
      return
    }

    if (!isWordValid(word)) {
      dispatch(err(type, 'Some of these characters are not allowed in a word'))
      return
    }

    const ref = flow.getRef(read('refPaths', type))
    if (ref) {
      dispatch(err(null))
      signals.didAdd()
      ref.push(word, err => {
        if (!err) {
          generator = null
          sig.generate()
        }
      })
    }
  })

  // Adds the given word to the store, removing it from the generated results.
  // We don't need to refresh the generator and the generated words, because
  // adding a previously generated word to the same source set has no effect on
  // the total output from this sample.
  sig.pick.subscribe((word, out) => {
    const selected = read(type, 'selected')
    if (_.contains(selected, word)) return

    const ref = flow.getRef(read('refPaths', type))
    if (!ref) {
      throw Error(`no ref found at path: ${read('refPaths', type)}`)
    }

    // Optimistically remove from results.
    const prevGenerated = read(type, 'generated')
    const generated = _.without(prevGenerated, word)
    out({
      type: 'set',
      path: [type, 'generated'],
      value: generated
    })

    out(new Promise((resolve, reject) => {
      ref.push(word, err => {
        if (err) {
          // Roll back the optimistic change.
          dispatch({
            type: 'set',
            path: [type, 'generated'],
            value: prevGenerated
          })
          reject(err)
        } else {
          resolve()
        }
      })
    }))
  })

  // Removes the given word from the selected group.
  sig.drop.subscribe((key, out) => {
    out(new Promise((resolve, reject) => {
      flow.getRef(read('refPaths', type)).child(key).remove(err => {
        if (err) {
          reject(err)
        } else {
          out(checkForEmpty(type))
          out(() => {generator = null})
          out(sig.generate)
          resolve()
        }
      })
    }))
  })
})

/**
 * Utils
 */

// If the selected words are empty, fills them with the defaults.
function checkForEmpty (type) {
  const selected = read(type, 'selected')
  const ref = flow.getRef(read('refPaths', type))
  const def = flow.defaultRefs[type]
  if (selected || !ref || !def) return Promise.resolve()

  return new Promise(resolve => {
    const defaults = read('defaults', type)

    if (defaults) {
      ref.set(defaults)
      resolve()
    } else {
      def.once('value', snap => {
        ref.set(snap.val())
        dispatch({
          type: 'set',
          path: ['defaults', type],
          value: snap.val()
        })
        resolve()
      }, () => {resolve()})
    }
  })
}

function getGenerator (sourceWords, limit) {
  const traits = new Traits()
  traits.examine(_.toArray(sourceWords))
  const generator = traits.generator()

  return () => {
    const results = []
    let word = ''
    while ((word = generator()) && results.length < limit) {
      if (!_.contains(sourceWords, word)) results.push(word)
    }
    return results
  }
}

function isWordValid (word) {
  try {
    const traits = new Traits()
    traits.examine([word])
    return true
  } catch (err) {
    console.warn(err)
    return false
  }
}

function err (type, value) {
  return {type: 'set', path: [type, 'error'], value}
}
