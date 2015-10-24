import React from 'react'
import {authData, root} from './data'
import {reactive} from './utils'

export class LoginButton extends React.Component {
  @reactive
  updateState () {
    this.setState({
      authData: authData.read()
    })
  }

  render () {
    return (
      <div className='container text-right'>
        {/* Anonymous */}
        {this.state.authData && this.state.authData.provider === 'anonymous' ?
        <div>
          <p>Anonymous session.</p>
          <p>
            <button className='sf-button-flat' onClick={::this.loginWithTwitter}>
              <span>Sign in with Twitter.</span>
              <span className='fa fa-twitter inline'></span>
            </button>
          </p>
        </div> : null}

        {/* Twitter */}
        {this.state.authData && this.state.authData.twitter ?
        <div>
          <p>Signed in as {this.state.authData.twitter.displayName}.</p>
          <p>
            <button className='sf-button-flat' onClick={::this.logout}>
              <span>Sign out</span>
              <span className='fa fa-sign-out inline'></span>
            </button>
          </p>
        </div> : null}
      </div>
    )
  }

  logout () {
    root.unauth()
  }

  loginWithTwitter () {
    root.authWithOAuthRedirect('twitter', err => {
      if (err) console.warn(err)
    })
  }
}
