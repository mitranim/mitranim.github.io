/**
 * Authentication object. Encapsulates auth state, user data, and provides
 * custom methods for operating on them.
 */

import {digest} from 'app'
import {root, fbRoot} from './root'

/**
 * When deauthed, auth anonymously. Using a vanilla Firebase instance because
 * Fireproof bugs out on this callback.
 */
fbRoot.onAuth(authData => {
  if (!authData) fbRoot.authAnonymously(err => {if (err) throw err})
})

/**
 * Authentication state class.
 */
export class AuthData {

  uid: string
  provider: string
  twitter: {}

  constructor() {
    // Refresh state on auth events.
    root.onAuth(authData => {
      for (let key in this) {
        delete this[key]
      }

      if (!authData) return

      for (let key in authData) {
        this[key] = authData[key]
      }

      digest()
    })
  }

  /**
   * Checks if we're logged in as a registered user.
   */
  isAuthed(): boolean {
    return this.provider && this.provider !== 'anonymous'
  }

  /**
   * Logs the user out.
   */
  unauth() {
    root.unauth()
  }
}

/**
 * Wrapped auth state.
 */
export var authData = new AuthData()
