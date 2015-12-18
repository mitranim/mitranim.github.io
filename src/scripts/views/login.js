import React from 'react'
import {read, send, auto} from '../core'

export const LoginButton = auto(() => {
  const auth = read('auth')
  if (!auth) return null

  return (
    <div className='container text-right'>
      {/* Anonymous */}
      {auth.provider === 'anonymous' ?
      <div>
        <p>Anonymous session.</p>
        <p>
          <button className='sf-button-flat' onClick={() => {send('auth/loginTwitter')}}>
            <span>Sign in with Twitter.</span>
            <span className='fa fa-twitter inline' />
          </button>
        </p>
      </div> : null}

      {/* Twitter */}
      {auth.twitter ?
      <div>
        <p>Signed in as {auth.twitter.displayName}.</p>
        <p>
          <button className='sf-button-flat' onClick={() => {send('auth/logout')}}>
            <span>Sign out</span>
            <span className='fa fa-sign-out inline' />
          </button>
        </p>
      </div> : null}
    </div>
  )
})
