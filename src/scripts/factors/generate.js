import _ from 'lodash'
import Traits from 'foliant'
import {match, multimatch, pipe} from 'prax'
// Circular dependency, TODO reconsider.
import {read, send} from '../core'
import {getRef, defaultRefs} from './auth'

const limit = 12
let generators = {}
const unsubs = []

export default () => pipe(
  multimatch('auth/logout', next => msg => {
    generators = {}
    while (unsubs.length) unsubs.shift()()
    next(msg)
  }),

  match({type: 'gen/init', kind: Boolean}, ({kind}) => {
    if (read(kind, 'inited') || !read('auth')) return

    const ref = getRef(read('refPaths', kind))
    if (!ref) throw Error(`ref not found for path: ${read('refPaths', kind)}`)

    loadSelected(kind, ref)
      .then(() => {
        send({type: 'patch', value: {[kind]: {inited: true}}})
      })
      .then(() => checkEmpty(kind))
      .then(() => {
        send({type: 'gen/generate', kind, onlyEmpty: true})
      })
  }),

  match({type: 'gen/generate', kind: Boolean}, ({kind, onlyEmpty}) => {
    if (onlyEmpty && read(kind, 'generated')) return

    const selected = read(kind, 'selected')
    if (!selected) return

    if (!generators[kind]) {
      generators[kind] = getGenerator(selected)
    }
    const generated = generators[kind]()

    send({
      type: 'patch',
      value: {
        [kind]: {
          generated,
          depleted: generated.length < limit
        }
      }
    })
  }),

  match({type: 'gen/add', kind: Boolean, word: _.isString}, ({kind, word}) => {
    if (typeof word !== 'string') {
      return send(err(kind, 'The word must be a string'))
    }

    word = word.toLowerCase().trim()

    if (!word) {
      return send(err(kind, 'Please input a word'))
    }

    if (word.length < 2) {
      return send(err(kind, 'The word is too short'))
    }

    const selected = read(kind, 'selected')

    if (_.contains(selected, word)) {
      return send(err(kind, 'This word is already in the set'))
    }

    if (!isWordValid(word)) {
      return send(err(kind, 'Some of these characters are not allowed in a word'))
    }

    const ref = getRef(read('refPaths', kind))

    if (ref) {
      send(err(null))

      ref.push(word, err => {
        if (!err) {
          send('gen/clearWord')
          send({type: 'gen/resetGenerator', kind})
          send({type: 'gen/generate', kind})
        }
      })
    }
  }),

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

    send({
      type: 'set',
      path: [kind, 'generated'],
      value: generated
    })

    ref.push(word, err => {
      if (err) {
        // Roll back the optimistic change.
        send({
          type: 'set',
          path: [kind, 'generated'],
          value: prevGenerated
        })
      }
    })
  }),

  match({type: 'gen/drop', kind: Boolean, key: Boolean}, ({kind, key}) => {
    getRef(read('refPaths', kind)).child(key).remove(err => {
      if (err) {
        console.error(err)
      } else {
        checkEmpty(kind).then(() => {
          send({type: 'gen/resetGenerator', kind})
          send({type: 'gen/generate', kind})
        })
      }
    })
  }),

  match('gen/clearWord', () => {
    send({type: 'patch', value: {state: {word: ''}}})
  }),

  match({type: 'gen/resetGenerator', kind: Boolean}, ({kind}) => {
    generators[kind] = null
  })
)

/**
 * Utils
 */

// If the selected words are empty, fills them with the defaults.
function checkEmpty (kind) {
  const selected = read(kind, 'selected')
  const ref = getRef(read('refPaths', kind))
  const def = defaultRefs[kind]
  if (selected || !ref || !def) return Promise.resolve()

  return new Promise(resolve => {
    const defaults = read('defaults', kind)

    if (defaults) {
      ref.set(defaults)
      resolve()
    } else {
      def.once('value', snap => {
        ref.set(snap.val())
        send({
          type: 'set',
          path: ['defaults', kind],
          value: snap.val()
        })
        resolve()
      }, () => {resolve()})
    }
  })
}

function loadSelected (kind, ref) {
  return new Promise(resolve => {
    const handler = ref.on('value', snap => {
      send({
        type: 'set',
        path: [kind, 'selected'],
        value: snap.val()
      })
      resolve()
    }, () => {
      resolve()
      ref.off('value', handler)
    })

    // Completely paranoid cleanup.
    unsubs.push(() => {ref.off('value', handler)})
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

function err (kind, value) {
  return {type: 'set', path: [kind, 'error'], value}
}
