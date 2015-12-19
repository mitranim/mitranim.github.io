import _ from 'lodash'
import {renderAt, unlink} from 'alder'
import {send, auto} from './core'

const initOnce = _.once(() => {send('init')})

export function renderTo (selector: string, view: Function) {
  const component = auto(view)

  onload(() => {
    const mountPoints = document.querySelectorAll(selector)
    if (mountPoints.length) initOnce()
    _.each(mountPoints, element => {
      renderAt(element, component)
    })
  })
}

document.addEventListener('simple-pjax-before-transition', () => {
  unlink(document.body)
})

function onload (callback) {
  if (/loaded|complete|interactive/.test(document.readyState)) {
    callback()
  } else {
    document.addEventListener('DOMContentLoaded', function cb () {
      document.removeEventListener('DOMContentLoaded', cb)
      callback()
    })
  }
  document.addEventListener('simple-pjax-after-transition', callback)
}

// Loading indicator.
export const spinner = auto(function spinner (props) {
  const {size, ...other} = props

  return (
    ['div', {className: `spinner-container ${size ? `size-${size}` : ''}`, ...other},
      ['div', {className: 'spinner'}]]
  )
})
