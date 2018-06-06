'use strict'

/**
 * Dependencies
 */

const $ = require('gulp-load-plugins')()
const del = require('del')
const gulp = require('gulp')
const log = require('fancy-log')
const lr = require('livereload')
const papyre = require('papyre')
const pt = require('path')
const webpack = require('webpack')
const {createServer} = require('http')
const {execSync} = require('child_process')
const {PassThrough} = require('stream')
const webpackConfig = require('./webpack.config.js')
const ghp = require('./ghp-server')

/**
 * Globals
 */

const PROD = process.env.NODE_ENV === 'production'
if (PROD) process.env.PROD = 'true'
process.env.VERSION = getCommitHash()

const srcRootDir = 'src'
const srcStaticFiles = 'src/static/**/*'
const srcTemplateEntry = 'src/templates/layouts.js'
const srcStyleFiles = ['src/styles/**/*.scss', 'cms/styles/**/*.scss']
const srcStyleEntryFiles = ['src/styles/main.scss', 'cms/styles/cms.scss']
const srcImagesRaster = 'src/images/**/*.{jpg,png,gif}'
const srcImagesVector = 'src/images/**/*.svg'
const outRootDir = 'public'
const outStyleDir = 'public/styles'
const outImageDir = 'public/images'

const cmsSrcStaticFiles = 'cms/static/**/*'

const GulpErr = err => ({showStack: false, toString: () => err})

function getCommitHash() {
  try {
    return execSync('git rev-parse --short HEAD').toString().trim()
  }
  catch (__) {
    return null
  }
}

/**
 * Tasks
 */

/* Clear */

gulp.task('clear', () => (
  del(`${outRootDir}/*`).catch(console.error.bind(console))
))

/* Static */

gulp.task('app:static:copy', () => (
  gulp.src(srcStaticFiles).pipe(gulp.dest(outRootDir))
))

gulp.task('cms:static:copy', () => (
  gulp.src(cmsSrcStaticFiles).pipe(gulp.dest(outRootDir))
))

gulp.task('static:copy', gulp.parallel('app:static:copy', 'cms:static:copy'))

gulp.task('static:watch', () => {
  $.watch(srcStaticFiles, gulp.series('app:static:copy'))
  $.watch(cmsSrcStaticFiles, gulp.series('cms:static:copy'))
})

/* HTML */

const papyreConfig = Object.assign({}, webpackConfig, {
  entry: pt.resolve(srcTemplateEntry),
})

gulp.task('templates:build', done => {
  papyre.build(papyreConfig, (err, result) => {
    if (err) done(err)
    else {
      log('[papyre]', result.timing)
      papyre.writeEntries(outRootDir, modifyEntries(result.entries)).then(done, done)
    }
  })
})

gulp.task('templates:watch', () => {
  papyre.watch(papyreConfig, (err, result) => {
    if (err) log(err)
    else {
      log('[papyre]', result.timing)
      papyre.writeEntries(outRootDir, modifyEntries(result.entries)).catch(log)
    }
  })
})

function modifyEntries(entries) {
  const out = []
  for (const entry of entries) {
    entry.path = entry.path.replace(/\.mdx?$/, '.html')
    out.push(entry)
    if (/^posts\//.test(entry.path)) {
      const {name} = pt.parse(entry.path)
      out.push({
        path: `thoughts/${name}.html`,
        body: `<meta http-equiv="refresh" content="0;URL='https://mitranim.com/posts/${name}'" />`,
      })
    }
  }
  return out
}

/* Scripts */

gulp.task('scripts:build', done => {
  buildWithWebpack(webpackConfig, done)
})

gulp.task('scripts:watch', () => {
  watchWithWebpack(webpackConfig)
})

function buildWithWebpack(config, done) {
  webpack(config, (err, stats) => {
    if (err) {
      done(err)
    }
    else {
      log('[webpack]', stats.toString(config.stats))
      done(stats.hasErrors() ? GulpErr('webpack error') : null)
    }
  })
}

function watchWithWebpack(config) {
  webpack(config).watch({}, (err, stats) => {
    log('[webpack]', stats.toString(config.stats))
    if (err) log('[webpack]', err.message)
  })
}

/* Styles */

gulp.task('styles:build', () => (
  gulp.src(srcStyleEntryFiles)
    .pipe($.sass({includePaths: [srcRootDir]}))
    .pipe($.autoprefixer({browsers: ['> 1%', 'IE >= 10', 'iOS 7']}))
    .pipe(!PROD ? new PassThrough({objectMode: true}) : $.cleanCss({
      keepSpecialComments: 0,
      aggressiveMerging: false,
      advanced: false,
      // Don't inline `@import: url()`
      processImport: false,
    }))
    .pipe(gulp.dest(outStyleDir))
))

gulp.task('styles:watch', () => {
  $.watch(srcStyleFiles, gulp.series('styles:build'))
})

/* Images */

// Resize and copy images
gulp.task('images:raster:normal', () => (
  gulp.src(srcImagesRaster)
    // Requires `graphicsmagick` or `imagemagick`. Install via Homebrew
    // or the package manager of your Unix distro.
    .pipe($.imageResize({quality: 1}))
    .pipe(gulp.dest(outImageDir))
))

// Minify and copy images
gulp.task('images:raster:small', () => (
  gulp.src(srcImagesRaster)
    .pipe($.imageResize({
      quality: 1,
      width: 640,    // max width
      upscale: false,
    }))
    .pipe(gulp.dest(outImageDir + '/small'))
))

// Crop images to small squares
gulp.task('images:raster:square', () => (
  gulp.src(srcImagesRaster)
    .pipe($.imageResize({
      quality: 1,
      gravity: 'Center',  // crop relative to center
      crop: true,
      width: 640,
      height: 640,
      upscale: false,
    }))
    .pipe(gulp.dest(outImageDir + '/square'))
))

gulp.task('images:raster',
  gulp.parallel('images:raster:normal', 'images:raster:small', 'images:raster:square'))

gulp.task('images:vector', () => (
  gulp.src(srcImagesVector)
    .pipe($.svgo())
    .pipe(gulp.dest(outImageDir))
))

gulp.task('images:build', gulp.parallel('images:raster', 'images:vector'))

gulp.task('images:watch', () => {
  $.watch(srcImagesRaster, gulp.series('images:raster'))
  $.watch(srcImagesVector, gulp.series('images:vector'))
})

/* Devserver */

const PORT = 11204

gulp.task('server', done => {
  lr.createServer({}, err => {if (err) done(err)}).watch(outRootDir)

  createServer((req, res) => {
    ghp.serve(req, res, {rootDir: outRootDir})
  })
  .listen(PORT, err => {
    if (err) done(err)
    else log(`Server listening on http://localhost:${PORT}`)
  })
})

/* Default */

gulp.task('buildup', gulp.parallel(
  'static:copy',
  'styles:build',
  'images:build'
))

gulp.task('watch', gulp.parallel(
  'static:watch',
  'styles:watch',
  'images:watch',
  'templates:watch',
  'scripts:watch',
  'server'
))

gulp.task('build', gulp.series('clear', gulp.parallel(
  'buildup',
  'templates:build',
  'scripts:build'
)))

gulp.task('default', gulp.series('clear', 'buildup', 'watch'))
