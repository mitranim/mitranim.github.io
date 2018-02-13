'use strict'

const bs = require('browser-sync').create()
const {log} = require('gulp-util')
const {mapVals} = require('fpx')
const config = require('./webpack.config')

const prod = process.env.NODE_ENV === 'production'

const middleware = []

if (prod) {
  require('webpack')(config).watch({}, (err, stats) => {
    log('[webpack]', stats.toString(config.stats))
    if (err) log('[webpack]', err.message)
  })
}
else {
  const compiler = require('webpack')(extend(config, {
    entry: mapVals(
      fsPath => ['webpack-hot-middleware/client?noInfo=true', fsPath],
      config.entry
    ),
  }))
  middleware.push(
    require('webpack-dev-middleware')(compiler, {
      publicPath: config.output.publicPath,
      noInfo: true,
    })
  )
  middleware.push(
    require('webpack-hot-middleware')(compiler)
  )
}

bs.init({
  startPath: '/',
  server: {
    baseDir: 'dist',
    middleware,
  },
  port: 11204,
  files: 'dist',
  open: false,
  online: false,
  ui: false,
  ghostMode: false,
  notify: false,
})

function extend () {
  return Object.assign({}, ...arguments)
}
