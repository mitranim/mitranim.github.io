import _ from 'lodash'
import Traits from 'foliant'
import {read, send, match, watch, set, patch} from '../core'
import {getRef, defaultRef} from './auth'

const limit = 12
let generators = {}
const unsubs = []

match('init', () => {
  watch(read => {
    const kind = read('state', 'kind')
    if (read('auth')) send({type: 'gen/init', kind})
  })
})

match('auth/logout', () => {
  generators = {}
  while (unsubs.length) unsubs.shift()()
})

match({type: 'gen/init', kind: Boolean}, ({kind}) => {
  if (read(kind, 'inited') || !read('auth')) return

  const ref = getRef(read('refPaths', kind))
  if (!ref) throw Error(`ref not found for path: ${read('refPaths', kind)}`)

  loadSelected(kind, ref)
    .then(() => {
      set([kind, 'inited'], true)
    })
    .then(() => checkEmpty(kind))
    .then(() => {
      send({type: 'gen/generate', kind, onlyEmpty: true})
    })
})

match({type: 'gen/generate', kind: Boolean}, ({kind, onlyEmpty}) => {
  if (onlyEmpty && read(kind, 'generated')) return

  const selected = read(kind, 'selected')
  if (!selected) return

  if (!generators[kind]) {
    generators[kind] = getGenerator(selected)
  }
  const generated = generators[kind]()

  patch([kind], {generated, depleted: generated.length < limit})
})

match({type: 'gen/add', kind: Boolean, word: _.isString}, ({kind, word}) => {
  if (typeof word !== 'string') {
    return set([kind, 'error'], 'The word must be a string')
  }

  word = word.toLowerCase().trim()

  if (!word) {
    return set([kind, 'error'], 'Please input a word')
  }

  if (word.length < 2) {
    return set([kind, 'error'], 'The word is too short')
  }

  const selected = read(kind, 'selected')

  if (_.contains(selected, word)) {
    return set([kind, 'error'], 'This word is already in the set')
  }

  if (!isWordValid(word)) {
    return set([kind, 'error'], 'Some of these characters are not allowed in a word')
  }

  const ref = getRef(read('refPaths', kind))

  if (ref) {
    set([kind, 'error'], null)

    ref.push(word, err => {
      if (!err) {
        set(['state', 'word'], '')
        generators[kind] = null
        send({type: 'gen/generate', kind})
      }
    })
  }
})

// Adds the given word to the store, removing it from the generated results.
// We don't need to refresh the generator and the generated words, because
// adding a previously generated word to the same source set has no effect
// on the total output from this sample.
match({type: 'gen/pick', kind: Boolean, word: Boolean}, ({kind, word}) => {
  const selected = read(kind, 'selected')
  if (_.contains(selected, word)) return

  const ref = getRef(read('refPaths', kind))
  if (!ref) {
    throw Error(`no ref found at path: ${read('refPaths', kind)}`)
  }

  // Optimistically remove from results.
  const prevGenerated = read(kind, 'generated')
  const generated = _.without(prevGenerated, word)

  set([kind, 'generated'], generated)

  ref.push(word, err => {
    if (err) {
      // Roll back the optimistic change.
      set([kind, 'generated'], prevGenerated)
    }
  })
})

match({type: 'gen/drop', kind: Boolean, key: Boolean}, ({kind, key}) => {
  getRef(read('refPaths', kind)).child(key).remove(err => {
    if (err) {
      console.error(err)
    } else {
      checkEmpty(kind).then(() => {
        generators[kind] = null
        send({type: 'gen/generate', kind})
      })
    }
  })
})

/**
 * Utils
 */

// If the selected words are empty, fills them with the defaults.
function checkEmpty (kind) {
  const selected = read(kind, 'selected')
  const ref = getRef(read('refPaths', kind))
  const def = defaultRef(kind)
  if (selected || !ref || !def) return Promise.resolve()

  return new Promise(resolve => {
    const defaults = read('defaults', kind)

    if (defaults) {
      ref.set(defaults)
      resolve()
    } else {
      def.once('value', snap => {
        ref.set(snap.val())
        set(['defaults', kind], snap.val())
        resolve()
      }, () => { resolve() })
    }
  })
}

function loadSelected (kind, ref) {
  return new Promise(resolve => {
    const handler = ref.on('value', snap => {
      set([kind, 'selected'], snap.val())
      resolve()
    }, () => {
      resolve()
      ref.off('value', handler)
    })

    // Completely paranoid cleanup.
    unsubs.push(() => { ref.off('value', handler) })
  })
}

function getGenerator (sourceWords) {
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
