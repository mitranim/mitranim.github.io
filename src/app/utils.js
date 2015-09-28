import React from 'react'
import {setUpDataLoad} from './data'

export function renderTo (selector: string) {
  return (Component: typeof React.Component) => {
    onload(() => {
      const elements = document.querySelectorAll(selector)
      for (let i = 0; i < elements.length; ++i) {
        setUpDataLoad()
        React.render(<Component/>, elements[i])
      }
    })
  }
}

function asapOnce (callback: () => void): void {
  if (/loaded|complete|interactive/.test(document.readyState)) callback()
  else document.addEventListener('DOMContentLoaded', function cb () {
    document.removeEventListener('DOMContentLoaded', cb)
    callback()
  })
}

function onload (callback: () => void): void {
  if (/loaded|complete|interactive/.test(document.readyState)) callback()
  document.addEventListener('DOMContentLoaded', callback)
}

export class Spinner extends React.Component {
  render () {
    return (
      <div className={`spinner-container ${this.props.size ? `size-${this.props.size}` : ''}`}
           style={this.props.style || null}>
        <div className='spinner' />
      </div>
    )
  }
}
