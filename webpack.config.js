'use strict'

const pt = require('path')
const webpack = require('webpack')
const prod = process.env.NODE_ENV === 'production'

module.exports = {
  entry: pt.resolve('src/scripts/app.js'),

  output: {
    path: pt.resolve('dist'),
    filename: 'app.js'
  },

  module: {
    loaders: [
      {
        test: /\.js$/,
        loader: 'babel',
        include: pt.resolve('src/scripts')
      }
    ]
  },

  plugins: !prod ? [
    new webpack.HotModuleReplacementPlugin()
  ] : [
    new webpack.optimize.UglifyJsPlugin({
      minimize: true,
      compress: {screw_ie8: true},
      mangle: true
    })
  ],

  devtool: prod ? 'source-map' : null,

  stats: {
    colors: true,
    chunks: false,
    version: false,
    hash: false,
    assets: false
  }
}
