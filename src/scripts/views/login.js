import {send, auto} from '../core'

export const login = auto(function login (props, read) {
  const auth = read('auth')
  if (!auth) return null

  return (
    ['div', {className: 'container text-right'},
      // Anonymous
      auth.provider === 'anonymous' ?
      ['div', null,
        ['p', null, 'Anonymous session'],
        ['p', null,
          ['button', {className: 'sf-button-flat', onclick () {send('auth/loginTwitter')}},
            'Sign in with Twitter',
            ['span', {className: 'fa fa-twitter inline'}]]]] : null,

      // Twitter
      auth.twitter ?
      ['div', null,
        ['p', null, `Signed in as ${auth.twitter.displayName}`],
        ['p', null,
          ['button', {className: 'sf-button-flat', onclick () {send('auth/logout')}},
            'Sign out',
            ['span', {className: 'fa fa-sign-out inline'}]]]] : null]
  )
})
