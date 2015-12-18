import React from 'react'
import _ from 'lodash'
import {render, unmountComponentAtNode} from 'react-dom'
import {send, auto} from './core'

const initOnce = _.once(() => {send('init')})

const unmountQueue = []

export function renderTo (selector: string, renderFunc: ?Function) {
  function init (Component) {
    onload(() => {
      const mountPoints = document.querySelectorAll(selector)
      if (mountPoints.length) initOnce()
      _.each(mountPoints, element => {
        unmountQueue.push(element)
        render(<Component />, element)
      })
    })
  }

  if (typeof renderFunc === 'function') init(auto(renderFunc))
  else return init
}

document.addEventListener('simple-pjax-before-transition', () => {
  unmountQueue.splice(0).forEach(unmountComponentAtNode)
})

function onload (callback: () => void): void {
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
export const Spinner = props => {
  const {size, ...other} = props

  return (
    <div className={`spinner-container ${size ? `size-${size}` : ''}`} {...other}>
      <div className='spinner' />
    </div>
  )
}
