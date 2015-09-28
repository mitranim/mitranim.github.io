import React from 'react'
import {Component, authData, root} from './data'

export class LoginButton extends Component {
  getState () {
    return {
      authData: authData.get()
    }
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
