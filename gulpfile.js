'use strict'

/**
 * Requires gulp 4.0:
 *   "gulp": "gulpjs/gulp#4.0"
 *
 * Requires Node.js 4.0+
 */

/* ***************************** Dependencies ********************************/

const $ = require('gulp-load-plugins')()
const bsync = require('browser-sync').create()
const del = require('del')
const flags = require('yargs').boolean('prod').argv
const gulp = require('gulp')
const pt = require('path')
const webpack = require('webpack')
const statilOptions = require('./statil')

/* ******************************** Globals **********************************/

const src = {
  html: 'src/html/**/*',
  xml: 'src/xml/**/*',
  robots: 'src/robots.txt',
  scripts: 'src/scripts/**/*.js',
  scriptsCore: 'src/scripts/app.js',
  stylesCore: 'src/styles/app.scss',
  styles: 'src/styles/**/*.scss',
  images: 'src/images/**/*',
  fonts: 'node_modules/font-awesome/fonts/**/*'
}

const dest = {
  html: 'dist',
  xml: 'dist/**/*.xml',
  scripts: 'dist/scripts',
  styles: 'dist/styles',
  images: 'dist/images',
  fonts: 'dist/fonts'
}

function reload (done) {
  bsync.reload()
  done()
}

/* ********************************* Tasks ***********************************/

/* -------------------------------- Scripts ---------------------------------*/

function scripts (done) {
  webpack({
    entry: './' + src.scriptsCore,
    output: {
      path: pt.join(process.cwd(), dest.scripts),
      filename: 'app.js'
    },
    module: {
      loaders: [
        {
          test: /\.js$/,
          loader: 'babel',
          include: pt.join(process.cwd(), 'src/scripts')
        }
      ]
    },
    plugins: flags.prod ? [new webpack.optimize.UglifyJsPlugin({compress: {warnings: false}})] : [],
    devtool: !flags.prod && typeof done !== 'function' ? 'source-map' : null,
    watch: typeof done !== 'function'
  }, function (err, stats) {
    if (err) {
      throw new Error(err)
    } else {
      const report = stats.toString({
        colors: true,
        chunks: false,
        timings: true,
        version: false,
        hash: false,
        assets: false
      })
      if (report) console.log(report)
    }
    if (typeof done === 'function') done()
    else bsync.reload()
  })
}

gulp.task('scripts:build', scripts)

gulp.task('scripts:build:watch', () => { scripts() })

/* -------------------------------- Styles ----------------------------------*/

gulp.task('styles:clear', function (done) {
  del(dest.styles).then(() => { done() })
})

gulp.task('styles:compile', function () {
  return gulp.src(src.stylesCore)
    .pipe($.plumber())
    .pipe($.sass())
    .pipe($.autoprefixer())
    .pipe($.base64({
      baseDir: '.',
      extensions: ['svg']
    }))
    .pipe($.minifyCss({
      keepSpecialComments: 0,
      aggressiveMerging: false,
      advanced: false
    }))
    .pipe(gulp.dest(dest.styles))
    .pipe(bsync.stream())
})

gulp.task('styles:build',
  gulp.series('styles:clear', 'styles:compile'))

gulp.task('styles:watch', function () {
  $.watch(src.styles, gulp.series('styles:build'))
})

/* --------------------------------- HTML -----------------------------------*/

gulp.task('html:clear', function (done) {
  del(dest.html + '/**/*.html').then(() => { done() })
})

gulp.task('html:compile', function () {
  return gulp.src(src.html)
    .pipe($.statil(statilOptions()))
    .pipe($.if(flags.prod, $.minifyHtml({
      empty: true,
      loose: true
    })))
    .pipe(gulp.dest(dest.html))
})

// Copy robots.txt.
gulp.task('html:robots', function () {
  return gulp.src(src.robots).pipe(gulp.dest(dest.html))
})

gulp.task('html:build', gulp.series('html:clear', 'html:compile', 'html:robots'))

gulp.task('html:watch', function () {
  $.watch(src.html, gulp.series('html:build', reload))
})

/* ---------------------------------- XML -----------------------------------*/

gulp.task('xml:clear', function (done) {
  del(dest.xml).then(() => { done() })
})

gulp.task('xml:compile', function () {
  return gulp.src([src.html, src.xml])
    .pipe($.plumber())
    .pipe($.statil(statilOptions()))
    .pipe($.filter('*feed*'))
    .pipe($.rename('feed.xml'))
    .pipe(gulp.dest(dest.html))
})

gulp.task('xml:build', gulp.series('xml:compile'))

gulp.task('xml:watch', function () {
  $.watch(src.html, gulp.series('xml:build'))
  $.watch(src.xml, gulp.series('xml:build'))
})

/* -------------------------------- Images ----------------------------------*/

gulp.task('images:clear', function (done) {
  del(dest.images).then(() => { done() })
})

// Resize and copy images
gulp.task('images:normal', function () {
  return gulp.src(src.images)
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
    .pipe(gulp.dest(dest.images))
})

// Minify and copy images.
gulp.task('images:small', function () {
  return gulp.src(src.images)
    .pipe($.imageResize({
      quality: 1,
      width: 640,    // max width
      upscale: false
    }))
    .pipe(gulp.dest(dest.images + '/small'))
})

// Crop images to small squares
gulp.task('images:square', function () {
  return gulp.src(src.images)
    .pipe($.imageResize({
      quality: 1,
      gravity: 'Center',  // crop relative to center
      crop: true,
      width: 640,
      height: 640,
      upscale: false
    }))
    .pipe(gulp.dest(dest.images + '/square'))
})

gulp.task('images:build',
  gulp.series('images:clear',
    gulp.parallel('images:normal', 'images:small', 'images:square')))

gulp.task('images:watch', function () {
  $.watch(src.images, gulp.series('images:build', reload))
})

/* --------------------------------- Fonts ----------------------------------*/

gulp.task('fonts:clear', function (done) {
  del(dest.fonts).then(() => { done() })
})

gulp.task('fonts:copy', function () {
  return gulp.src(src.fonts).pipe(gulp.dest(dest.fonts))
})

gulp.task('fonts:build', gulp.series('fonts:copy'))

gulp.task('fonts:watch', function () {
  $.watch(src.fonts, gulp.series('fonts:build', reload))
})

/* -------------------------------- Server ----------------------------------*/

gulp.task('server', function () {
  return bsync.init({
    startPath: '/',
    server: {
      baseDir: dest.html
    },
    port: 11204,
    online: false,
    ui: false,
    files: false,
    ghostMode: false,
    notify: false
  })
})

/* -------------------------------- Default ---------------------------------*/

if (flags.prod) {
  gulp.task('build', gulp.parallel(
    'scripts:build', 'styles:build', 'html:build', 'xml:build', 'fonts:build', 'images:build'
  ))
} else {
  gulp.task('build', gulp.parallel(
    'styles:build', 'html:build', 'xml:build', 'fonts:build', 'images:build'
  ))
}

gulp.task('watch', gulp.parallel(
  'scripts:build:watch', 'styles:watch', 'html:watch', 'xml:watch', 'fonts:watch', 'images:watch'
))

gulp.task('default', gulp.series('build', gulp.parallel('watch', 'server')))
