import _ from 'lodash'
import React from 'react'
import {render, unmountComponentAtNode} from 'react-dom'
import {createPure, createEmitter} from 'symphony'

// Creates a reactively updating component from a pure function.
export const pure = createPure(React.Component)

// Application-wide event emitter.
export const emit = createEmitter('init', 'loginSuccess', 'genAddSuccess')
export const on = emit.decorator

const unmountQueue = []

export function renderTo (selector: string, renderFunc: ?Function) {
  function init (Component: typeof React.Component) {
    onload(() => {
      const mountPoints = document.querySelectorAll(selector)
      if (mountPoints.length) emit('init')
      _.each(mountPoints, element => {
        unmountQueue.push(element)
        render(<Component />, element)
      })
    })
  }

  if (typeof renderFunc === 'function') init(pure(renderFunc))
  else return init
}

document.addEventListener('simple-pjax:before-transition', () => {
  unmountQueue.splice(0).forEach(unmountComponentAtNode)
})

function onload (callback: () => void): void {
  if (/loaded|complete|interactive/.test(document.readyState)) {
    callback()
  } else {
    document.addEventListener('DOMContentLoaded', function cb () {
      document.removeEventListener(cb)
      callback()
    })
  }
  document.addEventListener('simple-pjax:after-transition', callback)
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
