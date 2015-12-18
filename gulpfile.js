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
const cheerio = require('cheerio')
const del = require('del')
const flags = require('yargs').boolean('prod').argv
const gulp = require('gulp')
const hjs = require('highlight.js')
const marked = require('gulp-marked/node_modules/marked')
const pt = require('path')
const webpack = require('webpack')

/* ******************************** Globals **********************************/

const src = {
  html: 'src/html/**/*',
  xml: [
    'src/html/**/*.yaml',
    'src/html/thoughts/**/*',
    '!src/html/thoughts/index.html',
    'src/xml/**/*'
  ],
  robots: 'src/robots.txt',
  scripts: [
    'src/scripts/**/*.js',
    'node_modules/stylific/lib/stylific.min.js',
    'node_modules/simple-pjax/lib/simple-pjax.min.js'
  ],
  scriptsCore: 'src/scripts/app.js',
  stylesCore: 'src/styles/app.scss',
  styles: [
    'src/styles/**/*.scss',
    'node_modules/stylific/scss/**/*.scss'
  ],
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

/* *************************** Template Imports ******************************/

/**
 * Utility methods for templates.
 */
const imports = {
  prod: flags.prod,
  bgImg: function (path) {
    return 'style="background-image: url(/images/' + path + ')"'
  },
  truncate: function (html, num) {
    let part = cheerio(html).text().slice(0, num)
    if (part.length === num) part += ' ...'
    return part
  }
}

/* ******************************** Config ***********************************/

/**
 * marked rendering enhancements.
 */

// Custom heading renderer func that adds an anchor.
marked.Renderer.prototype.heading = function (text, level, raw) {
  const id = this.options.headerPrefix + raw.toLowerCase().replace(/[^\w]+/g, '-')
  return (
`<h${level}>
  <span>${text}</span>
  <a class="heading-anchor fa fa-link" href="#${id}" id="${id}"></a>
</h${level}>\n`
  )
}

// Custom link renderer func that adds target="_blank" to links to other sites.
// Mostly copied from the marked source.
marked.Renderer.prototype.link = function (href, title, text) {
  if (this.options.sanitize) {
    let prot = ''
    try {
      prot = decodeURIComponent(unescape(href))
        .replace(/[^\w:]/g, '')
        .toLowerCase()
    } catch (e) {
      return ''
    }
    if (prot.indexOf('javascript:') === 0 || prot.indexOf('vbscript:') === 0) {
      return ''
    }
  }
  let out = '<a href="' + href + '"'
  if (title) {
    out += ' title="' + title + '"'
  }
  if (/^[a-z]+:\/\//.test(href)) {
    out += ' target="_blank"'
  }
  out += '>' + text + '</a>'
  return out
}

// Default code renderer.
const renderCode = marked.Renderer.prototype.code

// Custom code renderer that understands a few custom directives.
marked.Renderer.prototype.code = function (code, lang, escaped) {
  const regexCollapse = /#collapse (.*)(?:\n|$)/g

  // Remove collapse directives and remember if there were any.
  const collapse = regexCollapse.exec(code)
  let head = ''
  if (collapse) {
    head = collapse[1]
    code = code.replace(regexCollapse, '').trim()
  }

  // Default render with highlighting.
  code = renderCode.call(this, code, lang, escaped).trim()

  // Optionally wrap in collapse.
  if (head) {
    code =
      '<div class="sf-collapse">\n' +
      '  <div class="sf-collapse-head theme-primary">' + head + '</div>\n' +
      '  <div class="sf-collapse-body">\n' +
           code + '\n' +
      '  </div>\n' +
      '</div>'
  }

  return code
}

/* ********************************* Tasks ***********************************/

/* -------------------------------- Scripts ---------------------------------*/

function scripts (done) {
  const alias = {
    'simple-pjax': 'simple-pjax/dist/simple-pjax'
  }
  if (flags.prod) {
    alias['react'] = 'react/dist/react.min'
    alias['react-dom'] = 'react-dom/dist/react-dom.min'
  }

  webpack({
    entry: './' + src.scriptsCore,
    output: {
      path: pt.join(process.cwd(), dest.scripts),
      filename: 'app.js'
    },
    resolve: {alias},
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
    // devtool: !flags.prod && typeof done !== 'function' ? 'inline-source-map' : null,
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

gulp.task('scripts:build:watch', (_) => {scripts()})

/* -------------------------------- Styles ----------------------------------*/

gulp.task('styles:clear', function (done) {
  del(dest.styles).then((_) => {done()})
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
    .pipe($.if(flags.prod, $.minifyCss({
      keepSpecialComments: 0,
      aggressiveMerging: false,
      advanced: false
    })))
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
  del(dest.html + '/**/*.html').then((_) => {done()})
})

gulp.task('html:compile', function () {
  const filterMd = $.filter('**/*.md', {restore: true})

  return gulp.src(src.html)
    .pipe($.plumber())
    // Pre-process markdown files.
    .pipe(filterMd)
    .pipe($.marked({
      gfm: true,
      tables: true,
      breaks: false,
      sanitize: false,
      smartypants: true,
      pedantic: false,
      // Code highlighter.
      highlight: function (code, lang) {
        if (lang) return hjs.highlight(lang, code).value
        return hjs.highlightAuto(code).value
      }
    }))
    // Add hljs code class.
    .pipe($.replace(/<pre><code class="(.*)">|<pre><code>/g,
                    '<pre><code class="hljs $1">'))
    // Restore other files.
    .pipe(filterMd.restore)
    // Unpack commented HTML parts.
    .pipe($.replace(/<!--\s*:((?:[^:]|:(?!\s*-->))*):\s*-->/g, '$1'))
    // Render all html.
    .pipe($.statil({imports: imports}))
    // Change each `<filename>` into `<filename>/index.html`.
    .pipe($.rename(function (path) {
      switch (path.basename + path.extname) {
        case 'index.html': case '404.html': return
      }
      path.dirname = pt.join(path.dirname, path.basename)
      path.basename = 'index'
    }))
    .pipe($.if(flags.prod, $.minifyHtml({
      empty: true,
      loose: true
    })))
    // Write to disk.
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
  del(dest.xml).then((_) => {done()})
})

gulp.task('xml:compile', function () {
  const filterMd = $.filter('**/*.md', {restore: true})

  return gulp.src(src.xml)
    .pipe($.plumber())
    // Pre-process markdown files.
    .pipe(filterMd)
    .pipe($.marked({
      gfm: true,
      tables: true,
      breaks: false,
      sanitize: false,
      smartypants: true,
      pedantic: false
    }))
    // Restore other files.
    .pipe(filterMd.restore)
    .pipe($.statil({imports: imports}))
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
  del(dest.images).then((_) => {done()})
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
  del(dest.fonts).then((_) => {done()})
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
    'scripts:build', 'styles:build', 'html:build', 'xml:build', 'fonts:build'
  ))
} else {
  gulp.task('build', gulp.parallel(
    'styles:build', 'html:build', 'xml:build', 'fonts:build'
  ))
}

gulp.task('watch', gulp.parallel(
  'scripts:build:watch', 'styles:watch', 'html:watch', 'xml:watch', 'fonts:watch'
))

gulp.task('default', gulp.series('build', gulp.parallel('watch', 'server')))
