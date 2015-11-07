import _ from 'lodash'
import Traits from 'foliant'
import {emit} from './utils'

// Used for lazily exported Firebase refs.
const auth = require('./auth')

const limit = 12
const kinds = ['names', 'words']
let generators = {}
const flush = []

export function transducer (action, dispatch, read) {
  if (action.type === 'logout') {
    generators = {}
    while (flush.length) flush.shift()()
    return action
  }

  const kind = action.kind
  if (!kind || !~kinds.indexOf(kind)) return action

  switch (action.type) {
    case 'genInit': {
      if (read(kind, 'inited') || !read('auth')) return

      return [
        {type: 'genLoadSelected', kind},
        {type: 'patch', value: {[kind]: {inited: true}}},
        {type: 'genCheckEmpty', kind},
        {type: 'genGenerate', kind, onlyEmpty: true}
      ]
    }

    case 'genLoadSelected': {
      const ref = auth.getRef(read('refPaths', kind))
      if (!ref) return Promise.reject(`ref not found for path: ${read('refPaths', kind)}`)
      return loadSelected(kind, dispatch, ref)
    }

    case 'genCheckEmpty': {
      return checkForEmpty(kind, dispatch, read)
    }

    case 'genGenerate': {
      if (action.onlyEmpty && read(kind, 'generated')) return

      const selected = read(kind, 'selected')
      if (!selected) return

      if (!generators[kind]) {
        generators[kind] = getGenerator(selected, limit)
      }
      const generated = generators[kind]()

      return {
        type: 'patch',
        value: {
          [kind]: {
            generated,
            depleted: generated.length < limit
          }
        }
      }
    }

    case 'genAdd': {
      let word = action.value

      if (typeof word !== 'string') {
        return err(kind, 'The word must be a string')
      }

      word = word.toLowerCase().trim()

      if (!word) {
        return err(kind, 'Please input a word')
      }

      if (word.length < 2) {
        return err(kind, 'The word is too short')
      }

      const selected = read(kind, 'selected')

      if (_.contains(selected, word)) {
        return err(kind, 'This word is already in the set')
      }

      if (!isWordValid(word)) {
        return err(kind, 'Some of these characters are not allowed in a word')
      }

      const ref = auth.getRef(read('refPaths', kind))

      if (ref) {
        dispatch(err(null))
        emit('genAddSuccess')

        ref.push(word, err => {
          if (!err) {
            generators[kind] = null
            dispatch({type: 'genGenerate', kind})
          }
        })
      }

      break
    }

    // Adds the given word to the store, removing it from the generated results.
    // We don't need to refresh the generator and the generated words, because
    // adding a previously generated word to the same source set has no effect
    // on the total output from this sample.
    case 'genPick': {
      const word = action.value
      const selected = read(kind, 'selected')
      if (_.contains(selected, word)) return

      const ref = auth.getRef(read('refPaths', kind))
      if (!ref) {
        throw Error(`no ref found at path: ${read('refPaths', kind)}`)
      }

      // Optimistically remove from results.
      const prevGenerated = read(kind, 'generated')
      const generated = _.without(prevGenerated, word)

      dispatch({
        type: 'set',
        path: [kind, 'generated'],
        value: generated
      })

      ref.push(word, err => {
        if (err) {
          // Roll back the optimistic change.
          dispatch({
            type: 'set',
            path: [kind, 'generated'],
            value: prevGenerated
          })
        }
      })

      // Babel gets confused by its own generated code if I use `break` here.
      return
    }

    // Removes the given word from the selected group.
    case 'genDrop': {
      auth.getRef(read('refPaths', kind)).child(action.value).remove(err => {
        if (err) {
          console.error(err)
        } else {
          dispatch([
            {type: 'genCheckEmpty', kind},
            {type: 'genResetGenerator', kind},
            {type: 'genGenerate', kind}
          ])
        }
      })

      break
    }

    case 'genResetGenerator': {
      generators[kind] = null
      break
    }
  }

  return action
}

/**
 * Utils
 */

// If the selected words are empty, fills them with the defaults.
function checkForEmpty (kind, dispatch, read) {
  const selected = read(kind, 'selected')
  const ref = auth.getRef(read('refPaths', kind))
  const def = auth.defaultRefs[kind]
  if (selected || !ref || !def) return Promise.resolve()

  return new Promise(resolve => {
    const defaults = read('defaults', kind)

    if (defaults) {
      ref.set(defaults)
      resolve()
    } else {
      def.once('value', snap => {
        ref.set(snap.val())
        dispatch({
          type: 'set',
          path: ['defaults', kind],
          value: snap.val()
        })
        resolve()
      }, () => {resolve()})
    }
  })
}

function loadSelected (kind, dispatch, ref) {
  return new Promise(resolve => {
    const handler = ref.on('value', snap => {
      dispatch({
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
    flush.push(() => {ref.off('value', handler)})
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

function err (kind, value) {
  return {type: 'set', path: [kind, 'error'], value}
}
