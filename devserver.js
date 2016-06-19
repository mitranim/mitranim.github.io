'use strict'

const bs = require('browser-sync').create()
const config = require('./webpack.config')
const prod = process.env.NODE_ENV === 'production'

bs.init({
  startPath: '/',
  server: {
    baseDir: 'dist',
    middleware: !prod ? hmr() : null
  },
  port: 11204,
  files: 'dist',
  open: false,
  online: false,
  ui: false,
  ghostMode: false,
  notify: false
})

function hmr () {
  const compiler = require('webpack')(extend(config, {
    entry: ['webpack-hot-middleware/client', config.entry]
  }))

  return [
    require('webpack-dev-middleware')(compiler, {
      publicPath: '/',
      noInfo: true
    }),
    require('webpack-hot-middleware')(compiler)
  ]
}

function extend () {
  return [].reduce.call(arguments, Object.assign, {})
}
