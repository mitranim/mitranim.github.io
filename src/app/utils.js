import React from 'react'
import {render} from 'react-dom'
import _ from 'lodash'
import {setUpDataLoad} from './data'

export function renderTo (selector: string) {
  return (Component: typeof React.Component) => {
    onload(() => {
      setUpDataLoad()
      _.each(document.querySelectorAll(selector), element => {
        render(<Component />, element)
      })
    })
  }
}

function onload (callback: () => void): void {
  if (/loaded|complete|interactive/.test(document.readyState)) callback()
  document.addEventListener('DOMContentLoaded', callback)
}

export const Spinner = props => (
  <div className={`spinner-container ${props.size ? `size-${props.size}` : ''}`}
       style={props.style || null}>
    <div className='spinner' />
  </div>
)
