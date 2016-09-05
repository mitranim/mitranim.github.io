'use strict'

/* ***************************** Dependencies ********************************/

const $ = require('gulp-load-plugins')()
const del = require('del')
const gulp = require('gulp')
const statilOptions = require('./statil')
const webpack = require('webpack')
const webpackConfig = require('./webpack.config')

/* ******************************** Globals **********************************/

const prod = process.env.NODE_ENV === 'production'

const src = {
  html: 'src/html/**/*',
  xml: 'src/xml/**/*',
  robots: 'src/robots.txt',
  stylesCore: 'src/styles/app.scss',
  styles: 'src/styles/**/*.scss',
  images: 'src/images/**/*',
  fonts: 'node_modules/font-awesome/fonts/**/*'
}

const out = {
  html: 'dist',
  xml: 'dist/**/*.xml',
  styles: 'dist/styles',
  images: 'dist/images',
  fonts: 'dist/fonts'
}

function noop () {}

/* ********************************* Tasks ***********************************/

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

/* -------------------------------- Styles ----------------------------------*/

gulp.task('styles:clear', () => (
  del(out.styles).catch(noop)
))

gulp.task('styles:compile', () => (
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

gulp.task('styles:build',
  gulp.series('styles:clear', 'styles:compile'))

gulp.task('styles:watch', () => {
  $.watch(src.styles, gulp.series('styles:build'))
})

/* --------------------------------- HTML -----------------------------------*/

gulp.task('html:clear', () => (
  del(out.html + '/**/*.html').catch(noop)
))

gulp.task('html:compile', () => (
  gulp.src(src.html)
    .pipe($.statil(statilOptions()))
    .pipe($.if(prod, $.minifyHtml({
      empty: true,
      loose: true
    })))
    .pipe(gulp.dest(out.html))
))

// Copy robots.txt.
gulp.task('html:robots', () => (
  gulp.src(src.robots).pipe(gulp.dest(out.html))
))

gulp.task('html:build', gulp.series('html:clear', 'html:compile', 'html:robots'))

gulp.task('html:watch', () => {
  $.watch(src.html, gulp.series('html:build'))
})

/* ---------------------------------- XML -----------------------------------*/

gulp.task('xml:clear', () => (
  del(out.xml).catch(noop)
))

gulp.task('xml:compile', () => (
  gulp.src([src.html, src.xml])
    .pipe($.statil(statilOptions()))
    .pipe($.filter('*feed*'))
    .pipe($.rename('feed.xml'))
    .pipe(gulp.dest(out.html))
))

gulp.task('xml:build', gulp.series('xml:compile'))

gulp.task('xml:watch', () => {
  $.watch(src.html, gulp.series('xml:build'))
  $.watch(src.xml, gulp.series('xml:build'))
})

/* -------------------------------- Images ----------------------------------*/

gulp.task('images:clear', () => (
  del(out.images).catch(noop)
))

// Resize and copy images
gulp.task('images:normal', () => (
  gulp.src(src.images)
    /**
    * Experience so far.
    * {quality: 1} -> reduces size by ≈66% with no resolution change and no visible quality change
    * {quality: 1, width: 1920} -> reduces size by ≈10 times for hi-res images
    */
    .pipe($.imageResize({
      quality: 1,
      width: 1920,    // max width
      upscale: false
    }))
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
  gulp.series('images:clear',
    gulp.parallel('images:normal', 'images:small', 'images:square')))

gulp.task('images:watch', () => {
  $.watch(src.images, gulp.series('images:build'))
})

/* --------------------------------- Fonts ----------------------------------*/

gulp.task('fonts:clear', () => (
  del(out.fonts).catch(noop)
))

gulp.task('fonts:copy', () => (
  gulp.src(src.fonts).pipe(gulp.dest(out.fonts))
))

gulp.task('fonts:build', gulp.series('fonts:copy'))

gulp.task('fonts:watch', () => {
  $.watch(src.fonts, gulp.series('fonts:build'))
})

/* -------------------------------- Server ----------------------------------*/

gulp.task('devserver', () => {
  require('./devserver')
})

/* -------------------------------- Default ---------------------------------*/

gulp.task('build',
  !prod
  ? gulp.parallel(
    'styles:build', 'html:build', 'xml:build', 'fonts:build', 'images:build'
  )
  : gulp.parallel(
    'scripts:build', 'styles:build', 'html:build', 'xml:build', 'fonts:build', 'images:build'
  )
)

gulp.task('watch', gulp.parallel(
  'styles:watch', 'html:watch', 'xml:watch', 'fonts:watch', 'images:watch', 'devserver'
))

gulp.task('default', gulp.series('build', 'watch'))
