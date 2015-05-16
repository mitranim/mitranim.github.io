'use strict'

/**
 * Requires gulp 4.0:
 *   "gulp": "git://github.com/gulpjs/gulp#4.0"
 */

/******************************* Dependencies ********************************/

var $       = require('gulp-load-plugins')()
var bsync   = require('browser-sync').create()
var cheerio = require('cheerio')
var gulp    = require('gulp')
var hjs     = require('highlight.js')
var marked  = require('gulp-marked/node_modules/marked')

/********************************** Globals **********************************/

// Base source directory.
var srcBase = './src/'

// jspm_packages path.
var jspmPath = destBase + '/jspm_packages/'

// Source paths with masks per type
var src = {
  lessCore: srcBase + 'styles/app.less',
  less:     srcBase + 'styles/**/*.less',
  img:      srcBase + 'img/**/*',
  html:     srcBase + 'html/',
  robots:   srcBase + 'robots.txt',
  js:       srcBase + 'app/**/*.ts',
  jsEnv:    srcBase + 'app/env.js',
  views:    srcBase + 'app/**/*.html',
  system:   './system.config.js'
}

// Base destination directory. Expected to be symlinked as another branch's
// directory.
var destBase = './mitranim-master/'

// Destination paths per type
var dest = {
  css:  destBase + 'css/',
  img:  destBase + 'img/',
  html: destBase,
  app:  destBase + 'app/',
}

/********************************* Utilities *********************************/

function prod() {
  return process.env.GULP_BUILD_TYPE === 'production'
}

// Usable in task flows.
function reload(done) {
  bsync.reload()
  done()
}

/***************************** Template Imports ******************************/

/**
 * Utility methods for templates.
 */
var imports = {
  lastId: 0,
  uniqId: function() {return 'static-id-' + ++imports.lastId},
  lastUniqId: function() {return 'static-id-' + imports.lastId},

  bgImg: function(path) {
    return 'style="background-image: url(/img/' + path + ')"'
  },

  truncate: function(html, num) {
    var part = cheerio(html).text().slice(0, num)
    if (part.length === num) part += ' ...'
    return part
  },

  prod: prod
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
  return gulp.src(dest.css, {read: false, allowEmpty: true}).pipe($.rimraf())
})

gulp.task('styles:less', function() {
  return gulp.src(src.lessCore)
    .pipe($.plumber())
    .pipe($.less())
    .pipe($.autoprefixer())
    .pipe($.if(prod(), $.minifyCss({
      keepSpecialComments: 0,
      aggressiveMerging: false,
      advanced: false
    })))
    .pipe(gulp.dest(dest.css))
    .pipe(bsync.reload({stream: true}))
})

gulp.task('styles:watch', function() {
  // Watch our .less files.
  $.watch(src.less, gulp.series('styles'))
  // Watch stylific's .less files.
  $.watch('./bower_components/stylific/**/*.less', gulp.series('styles'))
})

gulp.task('styles', gulp.series('styles:clear', 'styles:less'))

/*--------------------------------- Images ----------------------------------*/

// Clear images
gulp.task('images:clear', function() {
  return gulp.src(dest.img, {read: false, allowEmpty: true}).pipe($.rimraf())
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

// Make short cover images.
gulp.task('images:short', function() {
  return gulp.src(src.img)
    .pipe($.imageResize({
      quality: 1,
      gravity: 'Center',  // crop relative to the center
      crop: true,         // allow to crop to fit
      width: 1920,        // max width
      height: 512,        // max height
      upscale: false
    }))
    .pipe(gulp.dest(dest.img + 'short'))
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
    gulp.parallel(
      'images:normal',
      'images:short',
      'images:small',
      'images:square')))

/*---------------------------------- HTML -----------------------------------*/

// Clear html
gulp.task('html:clear', function() {
  return gulp.src([
    dest.html + '**/*.html',
    '!' + jspmPath
  ], {read: false, allowEmpty: true}).pipe($.rimraf())
})

// Compile html
gulp.task('html:compile', function() {
  var filterMd = $.filter('**/*.md')

  return gulp.src(src.html + '**/*')
    // .pipe($.plumber())
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
    // Render all html.
    .pipe($.statil({
      relativeDir: src.html,
      imports:     imports
    }))
    // Change each `<filename>` into `<filename>/index.html`.
    .pipe($.rename(function(path) {
      switch (path.basename + path.extname) {
        case 'index.html': case '404.html': return
      }
      path.dirname = pt.join(path.dirname, path.basename)
      path.basename = 'index'
    }))
    // Minify when building for production.
    .pipe($.if(prod(), $.minifyHtml({
      // Needed to keep attributes like [contenteditable]
      empty: true
    })))
    // Write to disk.
    .pipe(gulp.dest(dest.html))
    // Reload the browser.
    .pipe(bsync.reload({stream: true}))
})

// Copy robots.txt.
gulp.task('html:robots', function() {
  return gulp.src(src.robots).pipe(gulp.dest(dest.html))
})

gulp.task('html:watch', function() {
  $.watch(src.html + '**/*', gulp.series('html'))
})

// All html tasks
gulp.task('html', gulp.series('html:clear', 'html:compile', 'html:robots'))

/*--------------------------------- Scripts ---------------------------------*/

gulp.task('scripts:clear', function() {
  return gulp.src(dest.app, {read: false, allowEmpty: true}).pipe($.rimraf())
})

gulp.task('scripts:system', function() {
  return gulp.src(src.system)
    .pipe(gulp.dest(dest.html))
})

gulp.task('scripts:app', function() {
  return gulp.src(src.js)
    .pipe($.plumber()) // intentionally dumb error printing
    .pipe($.typescript({
      noExternalResolve: true,
      typescript: require('typescript'),
      target: 'ES5',
      module: 'system'
    }))
    .pipe(gulp.dest(dest.app))
})

gulp.task('scripts:views', function() {
  return gulp.src(src.views)
    .pipe($.plumber())
    .pipe($.if(prod(), $.minifyHtml({empty: true})))
    .pipe($.ngHtml2js({
      moduleName: 'app'
    }))
    .pipe($.concat('views.js'))
    .pipe($.babel({modules: 'system'}))
    .pipe(gulp.dest(dest.app))
})

gulp.task('scripts:env', function() {
  return gulp.src(src.jsEnv)
    .pipe($.if(!prod(), gulp.dest(dest.app)))
})

gulp.task('scripts:watch', function() {
  // Watch scripts.
  $.watch(src.js, gulp.series('scripts', reload))
  $.watch(src.system, gulp.series('scripts', reload))
  $.watch(src.jsEnv, gulp.series('scripts', reload))
  // Watch views.
  $.watch(src.views, gulp.series('scripts', reload))
})

gulp.task('scripts', gulp.series('scripts:clear', gulp.parallel('scripts:system', 'scripts:app', 'scripts:views', 'scripts:env')))

/*--------------------------------- Server ----------------------------------*/

gulp.task('bsync', function() {
  return bsync.init({
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
    // Don't open the window.
    // open: false,
  })
})

/*--------------------------------- Config ----------------------------------*/

// Build
gulp.task('build', gulp.parallel('styles', 'scripts', 'html'))

// Watch
gulp.task('watch', gulp.parallel(
  'styles:watch', 'scripts:watch', 'html:watch'
))

// Default
gulp.task('default', gulp.series('build', 'watch'))

// Serve files
gulp.task('server', gulp.series('build', gulp.parallel('watch', 'bsync')))
