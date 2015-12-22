import _ from 'lodash'
import {renderTo} from '../utils'
import {read, set} from '../core'
import {login} from './login'
import {words} from './words'
import {footnote} from './footnote'

renderTo('[data-render-foliant]', function foliant () {
  const kinds = read('kinds')
  const current = read('state', 'kind')

  return (
    ['div', null,
      ['div', {className: 'sf-navbar sf-navbar-tabs'},
        kinds.map(kind => (
          ['a', {className: `interactive ${kind === current ? 'active' : ''}`,
                 onclick () {set('state', 'kind', kind)},
                 key: kind},
            ['h3', null, _.capitalize(kind)]]
        ))],

      [words, {kind: current}],

      ['br'],

      [login],

      [footnote]]
  )
})
