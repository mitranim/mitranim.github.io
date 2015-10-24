import React from 'react'
import {render, unmountComponentAtNode} from 'react-dom'
import _ from 'lodash'
import {autorun, stop} from 'rapt'
import {setUpDataLoad} from './data'

const unmountQueue = []

export function renderTo (selector: string) {
  return (Component: typeof React.Component) => {
    onload(() => {
      const mountPoints = document.querySelectorAll(selector)
      if (mountPoints.length) setUpDataLoad()
      _.each(mountPoints, element => {
        unmountQueue.push(element)
        render(<Component />, element)
      })
    })
  }
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

/**
 * Component method decorator for reactive updates. Usage:
 *   class X extends React.Component {
 *     @reactive
 *     updateMe () {
 *       ...
 *     }
 *   }
 */
export function reactive (prototype, name, {value: reactiveFunc}) {
  if (typeof reactiveFunc !== 'function') return
  const {componentWillMount: pre, componentWillUnmount: post} = prototype

  prototype.componentWillMount = function () {
    if (typeof pre === 'function') pre.call(this)
    this[name] = reactiveFunc.bind(this)
    autorun(this[name])
  }

  prototype.componentWillUnmount = function () {
    stop(this[name])
    if (typeof post === 'function') post.call(this)
  }
}

export const Spinner = props => (
  <div className={`spinner-container ${props.size ? `size-${props.size}` : ''}`}
       style={props.style || null}>
    <div className='spinner' />
  </div>
)
