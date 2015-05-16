import Traits from 'foliant'

export var config = {
  dev: window.astilEnvironment === 'development',
  baseUrl: window.astilEnvironment === 'development' && typeof window.recordBaseUrl === 'string' ?
           window.recordBaseUrl : 'http://api.mitranim.com',
  fbRootUrl: 'https://incandescent-torch-3438.firebaseio.com'
}
