import React, {Component} from 'react'
import _ from 'lodash'
import {renderTo} from '../utils'
import {read, send, reactiveRender, reactiveMethod} from '../core'
import {LoginButton} from './login'
import {Words} from './words'

@renderTo('[data-render-foliant]')
@reactiveRender
export class WordsPage extends Component {
  @reactiveMethod
  init () {
    const kind = read('state', 'kind')
    if (read('auth')) send({type: 'gen/init', kind})
  }

  render () {
    const kinds = read('kinds')
    const current = read('state', 'kind')

    return (
      <div>
        <div className='sf-navbar sf-navbar-tabs'>
          {kinds.map(kind => (
            <a className={`interactive ${kind === current ? 'active' : ''}`}
               onClick={() => {select(kind)}} key={kind}>
              <h3>{_.capitalize(kind)}</h3>
            </a>
          ))}
        </div>

        <Words kind={current} />

        <br />

        <LoginButton />
      </div>
    )
  }
}

function select (kind) {
  send({type: 'patch', value: {state: {kind}}})
}
