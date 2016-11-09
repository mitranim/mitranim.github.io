'use strict'

/* ***************************** Dependencies ********************************/

const $ = require('gulp-load-plugins')()
const del = require('del')
const gulp = require('gulp')
const webpack = require('webpack')
const {fork} = require('child_process')

const statilOptions = require('./statil')
const webpackConfig = require('./webpack.config')

/* ******************************** Globals **********************************/

const prod = process.env.NODE_ENV === 'production'

const src = {
  static: [
    'src/static/**/*',
    'node_modules/font-awesome/@(fonts)/**/*',
  ],
  html: 'src/html/**/*',
  xml: 'src/xml/**/*',
  stylesCore: 'src/styles/app.scss',
  styles: 'src/styles/**/*.scss',
  images: 'src/images/**/*',
}

const out = {
  root: 'dist',
  xml: 'dist/**/*.xml',
  styles: 'dist/styles',
  images: 'dist/images',
}

function noop () {}

/* ********************************* Tasks ***********************************/

/* --------------------------------- Clear ---------------------------------- */

gulp.task('clear', () => (
  // Skips dotfiles like `.git` and `.gitignore`
  del(out.root + '/*').catch(noop)
))

/* -------------------------------- Static --------------------------------- */

gulp.task('static:build', () => (
  gulp.src(src.static).pipe(gulp.dest(out.root))
))

gulp.task('static:watch', () => {
  $.watch(src.static, gulp.series('static:build'))
})

/* -------------------------------- Scripts ---------------------------------*/

gulp.task('scripts:build', done => {
  webpack(webpackConfig, (err, stats) => {
    if (err) {
      throw new $.util.PluginError('webpack', err, {showProperties: false})
    }
    $.util.log('[webpack]', stats.toString(webpackConfig.stats))
    if (stats.hasErrors()) {
      throw new $.util.PluginError('webpack', 'plugin error', {showProperties: false})
    }
    done()
  })
})

/* --------------------------------- HTML -----------------------------------*/

gulp.task('html:build', () => (
  gulp.src(src.html)
    .pipe($.statil(statilOptions()))
    .pipe($.if(prod, $.minifyHtml({empty: true, loose: true})))
    .pipe(gulp.dest(out.root))
))

gulp.task('html:watch', () => {
  $.watch(src.html, gulp.series('html:build'))
})

/* ---------------------------------- XML -----------------------------------*/

gulp.task('xml:build', () => (
  gulp.src([src.html, src.xml])
    .pipe($.statil(statilOptions()))
    .pipe($.filter('*feed*'))
    .pipe($.rename('feed.xml'))
    .pipe(gulp.dest(out.root))
))

gulp.task('xml:watch', () => {
  $.watch(src.html, gulp.series('xml:build'))
  $.watch(src.xml, gulp.series('xml:build'))
})

/* -------------------------------- Styles ----------------------------------*/

gulp.task('styles:build', () => (
  gulp.src(src.stylesCore)
    .pipe($.sass())
    .pipe($.autoprefixer())
    .pipe($.cleanCss({
      keepSpecialComments: 0,
      aggressiveMerging: false,
      advanced: false,
      compatibility: {properties: {colors: false}}
    }))
    .pipe(gulp.dest(out.styles))
))

gulp.task('styles:watch', () => {
  $.watch(src.styles, gulp.series('styles:build'))
})

/* -------------------------------- Images ----------------------------------*/

// Resize and copy images
gulp.task('images:normal', () => (
  gulp.src(src.images)
    // Requires `graphicsmagick` or `imagemagick`. Install via Homebrew.
    .pipe($.imageResize({quality: 1}))
    .pipe(gulp.dest(out.images))
))

// Minify and copy images.
gulp.task('images:small', () => (
  gulp.src(src.images)
    .pipe($.imageResize({
      quality: 1,
      width: 640,    // max width
      upscale: false
    }))
    .pipe(gulp.dest(out.images + '/small'))
))

// Crop images to small squares
gulp.task('images:square', () => (
  gulp.src(src.images)
    .pipe($.imageResize({
      quality: 1,
      gravity: 'Center',  // crop relative to center
      crop: true,
      width: 640,
      height: 640,
      upscale: false
    }))
    .pipe(gulp.dest(out.images + '/square'))
))

gulp.task('images:build',
  gulp.parallel('images:normal', 'images:small', 'images:square'))

gulp.task('images:watch', () => {
  $.watch(src.images, gulp.series('images:build'))
})

/* -------------------------------- Server ----------------------------------*/

gulp.task('devserver', () => {
  let proc

  process.on('exit', () => {
    if (proc) proc.kill()
  })

  function restart () {
    if (proc) proc.kill()
    proc = fork('./devserver')
  }

  restart()
  $.watch(['./webpack.config.js', './devserver.js'], restart)
})

/* -------------------------------- Default ---------------------------------*/

gulp.task('buildup', gulp.parallel(
  'static:build',
  'html:build',
  'xml:build',
  'styles:build',
  'images:build'
))

gulp.task('watch', gulp.parallel(
  'html:watch',
  'xml:watch',
  'styles:watch',
  'images:watch',
  'devserver'
))

gulp.task('build', gulp.series('clear', gulp.parallel('buildup', 'scripts:build')))

gulp.task('default', gulp.series('clear', 'buildup', 'watch'))
