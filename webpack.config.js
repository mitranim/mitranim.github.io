'use strict'

const pt = require('path')
const webpack = require('webpack')
const prod = process.env.NODE_ENV === 'production'

module.exports = {
  entry: {
    main: pt.resolve('src/scripts/main.js'),
  },

  output: {
    path: pt.resolve('dist/scripts'),
    filename: '[name].js',
    // For dev middleware
    publicPath: '/scripts/',
  },

  module: {
    rules: [
      {
        test: /\.js$/,
        include: pt.resolve('src/scripts'),
        use: {loader: 'babel-loader'},
      }
    ]
  },

  plugins: !prod ? [
    new webpack.HotModuleReplacementPlugin()
  ] : [
    new webpack.optimize.UglifyJsPlugin({
      minimize: true,
      compress: {warnings: false, screw_ie8: true},
      mangle: true
    })
  ],

  devtool: prod ? 'source-map' : false,

  stats: {
    colors: true,
    chunks: false,
    version: false,
    hash: false,
    assets: false
  }
}
