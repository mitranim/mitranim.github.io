import _ from 'lodash'
import {renderTo} from '../utils'
import {set} from '../core'
import {login} from './login'
import {words} from './words'
import {footnote} from './footnote'

renderTo('[data-render-foliant]', function foliant (props, read) {
  const kinds = read('kinds')
  const current = read('state', 'kind')

  return (
    ['div', null,
      ['nav', {className: 'nav-h'},
        kinds.map(kind => (
          ['button', {className: `flat ${kind === current ? 'active' : ''}`,
                 onclick () { set(['state', 'kind'], kind) },
                 key: kind},
            ['h3', null, _.capitalize(kind)]]
        ))],

      [words, {kind: current}],

      ['br'],

      [login],

      [footnote]]
  )
})