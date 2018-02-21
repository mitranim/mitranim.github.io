'use strict'

const pt = require('path')
const webpack = require('webpack')

const PROD = process.env.NODE_ENV === 'production'

const SRC_DIR = pt.resolve('src/scripts')
const TEMPLATE_SRC_DIR = pt.resolve('src/templates')
const CMS_SRC_DIR = pt.resolve('cms/scripts')
const PUBLIC_DIR = pt.resolve('public/scripts')

module.exports = {
  entry: {
    main: pt.join(SRC_DIR, 'main.js'),
    cms: pt.join(CMS_SRC_DIR, 'cms.js'),
  },

  output: {
    path: PUBLIC_DIR,
    filename: '[name].js',
  },

  module: {
    rules: [
      {
        test: /\.jsx?$/,
        include: [SRC_DIR, TEMPLATE_SRC_DIR, CMS_SRC_DIR],
        use: {loader: 'babel-loader'},
      },
      {
        test: /\.svg$/,
        include: pt.resolve('src'),
        use: {loader: 'raw-loader'},
      },
      ...(!PROD ? [] : [
        // disable dev features and warnings in React and related libs
        {
          test: /react.*\.jsx?$/,
          include: /node_modules/,
          use: {loader: 'transform-loader', options: {envify: true}},
        },
      ]),
    ],
  },

  plugins: [
    new webpack.ProvidePlugin({
      React: 'react',
    }),
    // ...(PROD ? [
    //   new webpack.optimize.UglifyJsPlugin({
    //     mangle: {toplevel: true},
    //     compress: {warnings: false},
    //     output: {comments: false},
    //   }),
    // ] : []),
  ],

  devtool: PROD ? 'source-map' : false,

  stats: {
    assets: false,
    colors: true,
    hash: false,
    modules: false,
    timings: true,
    version: false,
  },
}
