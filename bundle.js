'use strict';

var path    = require('path');
var flags   = require('yargs').argv;
var Builder = require('systemjs-builder');
var builder = new Builder();

function prod() {
  return flags.prod === true || flags.prod === 'true';
}

builder.loadConfig('./system.config.js')
  .then(function() {
    builder.config({
      defaultJSExtensions: true
    });

    if (prod()) {
      builder.config({
        'paths': {
          '*': 'mitranim-master/app/*',
          'npm:*': 'node_modules/*'
        },
        map: {
          'react': 'npm:react/dist/react.min'
        }
      });
    }

    console.info('Starting bundling...');

    return builder.buildSFX('app', './mitranim-master/app.js', {
      runtime: false,
      minify: prod(),
      sourceMaps: false
    })
    .then(function() {
      console.info('Finished bundling.');
    })
    .catch(function(err) {
      console.warn(err);
    });
  });
