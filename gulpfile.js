'use strict'

/**
 * Requires gulp 4.0:
 *   "gulp": "git://github.com/gulpjs/gulp#4.0"
 */

/******************************* Dependencies ********************************/

var $      = require('gulp-load-plugins')()
var bsync  = require('browser-sync')
var gulp   = require('gulp')
var hjs    = require('highlight.js')
var marked = require('gulp-marked/node_modules/marked')

/********************************** Globals **********************************/

// Base source directory.
var srcBase = './src/'

// Source paths with masks per type
var src = {
  lessCore:  srcBase + 'styles/app.less',
  less:      srcBase + 'styles/**/*.less',
  img:       srcBase + 'img/**/*',
  templates: srcBase + 'templates/'
}

// Base destination directory. Expected to be symlinked as another branch's
// directory.
var destBase = './mitranim-master/'

// Destination paths per type
var dest = {
  css:   destBase + 'css/',
  img:   destBase + 'img/',
  html:  destBase
}

/********************************** Config ***********************************/

/**
 * Change how marked compiles links to add target="_blank" to links to other sites.
 */

// Default link renderer func.
var linkDef = marked.Renderer.prototype.link

// Custom link renderer func that adds target="_blank" to links to other sites.
// Mostly copied from the marked source.
marked.Renderer.prototype.link = function(href, title, text) {
  if (this.options.sanitize) {
    try {
      var prot = decodeURIComponent(unescape(href))
        .replace(/[^\w:]/g, '')
        .toLowerCase()
    } catch (e) {
      return ''
    }
    if (prot.indexOf('javascript:') === 0 || prot.indexOf('vbscript:') === 0) {
      return ''
    }
  }
  var out = '<a href="' + href + '"'
  if (title) {
    out += ' title="' + title + '"'
  }
  if (/^[a-z]+:\/\//.test(href)) {
    out += ' target="_blank"'
  }
  out += '>' + text + '</a>'
  return out
}

/*********************************** Tasks ***********************************/

/*--------------------------------- Styles ----------------------------------*/

gulp.task('styles:clear', function() {
  return gulp.src(dest.css, {read: false}).pipe($.rimraf())
})

gulp.task('styles:less', function() {
  return gulp.src(src.lessCore)
    .pipe($.plumber())
    .pipe($.less())
    .pipe($.autoprefixer())
    .pipe($.minifyCss({
      keepSpecialComments: 0,
      aggressiveMerging: false,
      advanced: false
    }))
    .pipe(gulp.dest(dest.css))
    .pipe(bsync.reload({stream: true}))
})

gulp.task('styles', gulp.series('styles:clear', 'styles:less'))

/*--------------------------------- Images ----------------------------------*/

// Clear images
gulp.task('images:clear', function() {
  return gulp.src(dest.img, {read: false}).pipe($.rimraf())
})

// Resize and copy images
gulp.task('images:normal', function() {
  return gulp.src(src.img)
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
    .pipe(gulp.dest(dest.img))
})

// Minify and copy images.
gulp.task('images:small', function() {
  return gulp.src(src.img)
    .pipe($.imageResize({
      quality: 1,
      width: 640,    // max width
      upscale: false
    }))
    .pipe(gulp.dest(dest.img + 'small'))
})

// Crop images to small squares
gulp.task('images:square', function() {
  return gulp.src(src.img)
    .pipe($.imageResize({
      quality: 1,
      gravity: 'Center',  // crop relative to the center
      crop: true,
      width: 640,
      height: 640,
      upscale: false
    }))
    .pipe(gulp.dest(dest.img + 'square'))
})

// All image tasks.
gulp.task('images',
  gulp.series(
    'images:clear',
    gulp.parallel('images:normal', 'images:small', 'images:square')))

/*-------------------------------- Templates --------------------------------*/

// Clear templates
gulp.task('templates:clear', function() {
  return gulp.src(dest.html + '**/*.html', {read: false}).pipe($.rimraf())
})

// Compile templates
gulp.task('templates:compile', function() {
  var filterMd = $.filter('**/*.md')

  return gulp.src(src.templates + '**/*')
    .pipe($.plumber())
    // Pre-process the markdown files.
    .pipe(filterMd)
    .pipe($.marked({
      gfm:         true,
      tables:      true,
      breaks:      false,
      sanitize:    false,
      smartypants: true,
      pedantic:    false,
      // Code highlighter.
      highlight: function(code, lang) {
        if (lang) return hjs.highlight(lang, code).value
        return hjs.highlightAuto(code).value
      }
    }))
    // Return the other files.
    .pipe(filterMd.restore())
    // Render all the templates.
    .pipe($.statil({
      relativeDir: src.templates
    }))
    // Remove the partials.
    .pipe($.filter(['**/*.html', '!**/partials/**/*', '!**/indices/**/*']))
    // Write to disk.
    .pipe(gulp.dest(dest.html))
    // Reload the browser.
    .pipe(bsync.reload({stream: true}))
})

// All template tasks
gulp.task('templates', gulp.series('templates:clear', 'templates:compile'))

/*--------------------------------- Server ----------------------------------*/

gulp.task('bsync', function() {
  return bsync({
    server: {
      baseDir: destBase
    },
    port: 11204,
    online: false,
    // Don't enable the UI.
    ui: false,
    // Don't watch files (default false, just making sure)
    files: false,
    // Don't sync anything across devices.
    ghostMode: false,
    // Don't show the notification.
    // notify: false
  })
})

/*--------------------------------- Config ----------------------------------*/

// Watch
gulp.task('watch', function() {
  // Watch the .less files
  $.watch(src.less, gulp.series('styles'))
  // Watch the templates
  $.watch(src.templates + '**/*', gulp.series('templates'))

  // Watch stylific's .less files
  $.watch('./bower_components/stylific/**/*.less', gulp.series('styles'))
})

// Build
gulp.task('build', gulp.parallel('styles', 'templates'))

// Default
gulp.task('default', gulp.series('build', 'watch'))

// Serve files
gulp.task('server', gulp.series('build', gulp.parallel('watch', 'bsync')))
