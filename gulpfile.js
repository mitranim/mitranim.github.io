'use strict';

/**
 * Requires gulp 4.0:
 *   "gulp": "git://github.com/gulpjs/gulp#4.0"
 */

/******************************* Dependencies ********************************/

var $       = require('gulp-load-plugins')();
var bsync   = require('browser-sync').create();
var cheerio = require('cheerio');
var gulp    = require('gulp');
var hjs     = require('highlight.js')
var marked  = require('gulp-marked/node_modules/marked');
var fs      = require('fs');
var flags   = require('yargs').argv;
var pt      = require('path');

/********************************** Globals **********************************/

var src = {
  html: [
    'src/html/**/*',
    'bower_components/font-awesome-svg-png/black/**/*.svg'
  ],
  robots: 'src/robots.txt',
  scripts: [
    'src/app/**/*.js',
    'node_modules/stylific/lib/stylific.js'
  ],
  stylesCore: 'src/styles/app.scss',
  styles: 'src/styles/**/*.scss',
  images: 'src/images/**/*',
  system: './system.config.js'
};

var dest = {
  html: 'mitranim-master',
  scripts: [
    'mitranim-master/app/**/*.js',
    '!mitranim-master/app/views.js'
  ],
  styles: 'mitranim-master/styles',
  images: 'mitranim-master/images',
  app: 'mitranim-master/app'
};

function prod() {
  return flags.prod === true || flags.prod === 'true';
}

function reload(done) {
  bsync.reload();
  done();
}

/***************************** Template Imports ******************************/

/**
 * Utility methods for templates.
 */
var imports = {
  prod: prod,
  lastId: 0,
  uniqId: function() {return 'uniq-id-' + ++imports.lastId},
  lastUniqId: function() {return 'uniq-id-' + imports.lastId},
  bgImg: function(path) {
    return 'style="background-image: url(/images/' + path + ')"'
  },
  truncate: function(html, num) {
    var part = cheerio(html).text().slice(0, num)
    if (part.length === num) part += ' ...'
    return part
  }
};

/********************************** Config ***********************************/

/**
 * marked rendering enhancements.
 */

// Default link renderer func.
var renderLink = marked.Renderer.prototype.link;

// Custom link renderer func that adds target="_blank" to links to other sites.
// Mostly copied from the marked source.
marked.Renderer.prototype.link = function(href, title, text) {
  if (this.options.sanitize) {
    try {
      var prot = decodeURIComponent(unescape(href))
        .replace(/[^\w:]/g, '')
        .toLowerCase();
    } catch (e) {
      return '';
    }
    if (prot.indexOf('javascript:') === 0 || prot.indexOf('vbscript:') === 0) {
      return '';
    }
  }
  var out = '<a href="' + href + '"';
  if (title) {
    out += ' title="' + title + '"';
  }
  if (/^[a-z]+:\/\//.test(href)) {
    out += ' target="_blank"';
  }
  out += '>' + text + '</a>';
  return out;
}

// Default code renderer.
var renderCode = marked.Renderer.prototype.code;

// Custom code renderer that understands a few custom directives.
marked.Renderer.prototype.code = function(code, lang, escaped) {
  // var regexInclude = /#include (.*)(?:\n|$)/g;
  var regexCollapse = /#collapse (.*)(?:\n|$)/g;

  // if (regexInclude.test(code)) {
  //   code = code.replace(regexInclude, function(match, path) {
  //     return fs.readFileSync(path, 'utf8').trim();
  //   });
  // }

  // Remove collapse directives and remember if there were any.
  var collapse = regexCollapse.exec(code);
  if (collapse) {
    var label = collapse[1];
    code = code.replace(regexCollapse, '').trim();
  }

  // Default render with highlighting.
  code = renderCode.call(this, code, lang, escaped).trim();

  // Optionally wrap in collapse.
  if (label) {
    code =
      '<div class="sf-collapse">\n' +
      '  <label class="theme-primary">' + label + '</label>\n' +
      '  <div class="sf-collapse-body">\n' +
           code + '\n' +
      '  </div>\n' +
      '</div>';
  }

  return code;
}

/*********************************** Tasks ***********************************/

/*--------------------------------- Scripts ---------------------------------*/

gulp.task('scripts:clear', function() {
  return gulp.src(dest.scripts, {read: false, allowEmpty: true})
    .pipe($.plumber())
    .pipe($.rimraf());
});

gulp.task('scripts:compile', function() {
  return gulp.src(src.scripts)
    .pipe($.plumber())
    // .pipe($.sourcemaps.init())
    .pipe($.babel({
      // externalHelpers: true,
      modules: 'system',
      optional: [
        'spec.protoToAssign',
        'es7.classProperties',
        'es7.decorators',
        'es7.functionBind',
        'validation.undeclaredVariableCheck'
      ]
    }))
    // .pipe($.sourcemaps.write())
    .pipe(gulp.dest(dest.app));
});
// });

// gulp.task('scripts:build', gulp.series('scripts:clear', 'scripts:compile', 'scripts:env'));
gulp.task('scripts:build', gulp.series('scripts:clear', 'scripts:compile'));

gulp.task('scripts:watch', function() {
  $.watch(src.scripts, gulp.series('scripts:build', reload));
  // $.watch(src.scriptsEnv, gulp.series('scripts:build', reload));
});

/*--------------------------------- Styles ----------------------------------*/

gulp.task('styles:clear', function() {
  return gulp.src(dest.styles, {read: false, allowEmpty: true})
    .pipe($.plumber())
    .pipe($.rimraf());
});

gulp.task('styles:compile', function() {
  return gulp.src(src.stylesCore)
    .pipe($.plumber())
    .pipe($.sass())
    .pipe($.autoprefixer())
    .pipe($.base64({
      baseDir: '.',
      extensions: ['svg']
    }))
    .pipe($.if(prod(), $.minifyCss({
      keepSpecialComments: 0,
      aggressiveMerging: false,
      advanced: false
    })))
    .pipe(gulp.dest(dest.styles))
    .pipe(bsync.reload({stream: true}));
});

gulp.task('styles:build',
  gulp.series('styles:clear', 'styles:compile'));

gulp.task('styles:watch', function() {
  $.watch(src.styles, gulp.series('styles:build'));
  $.watch('./node_modules/stylific/scss/**/*.scss', gulp.series('styles:build'));
});

/*---------------------------------- HTML -----------------------------------*/

gulp.task('html:clear', function() {
  return gulp.src([
      dest.html + '/**/*.html',
      '!' + dest.app + '/**/*'
    ], {read: false, allowEmpty: true})
    .pipe($.plumber())
    .pipe($.rimraf());
});

gulp.task('html:compile', function() {
  var filterMd = $.filter('**/*.md');

  return gulp.src(src.html)
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
        if (lang) return hjs.highlight(lang, code).value;
        return hjs.highlightAuto(code).value;
      }
    }))
    // Add the hljs code class.
    .pipe($.replace(/<pre><code class="(.*)">|<pre><code>/g,
                    '<pre><code class="hljs $1">'))
    // Return the other files.
    .pipe(filterMd.restore())
    // Render all html.
    .pipe($.statil({imports: imports}))
    // Change each `<filename>` into `<filename>/index.html`.
    .pipe($.rename(function(path) {
      switch (path.basename + path.extname) {
        case 'index.html': case '404.html': return;
      }
      path.dirname = pt.join(path.dirname, path.basename);
      path.basename = 'index';
    }))
    .pipe($.if(prod(), $.minifyHtml({
      empty: true
    })))
    // Write to disk.
    .pipe(gulp.dest(dest.html));
});

// Copy robots.txt.
gulp.task('html:robots', function() {
  return gulp.src(src.robots).pipe(gulp.dest(dest.html));
});

gulp.task('html:build', gulp.series('html:clear', 'html:compile', 'html:robots'));

gulp.task('html:watch', function() {
  $.watch(src.html, gulp.series('html:build', reload));
});

/*--------------------------------- Images ----------------------------------*/

gulp.task('images:clear', function() {
  return gulp.src(dest.images, {read: false, allowEmpty: true}).pipe($.rimraf());
});

// Resize and copy images
gulp.task('images:normal', function() {
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
    .pipe(gulp.dest(dest.images));
});

// Minify and copy images.
gulp.task('images:small', function() {
  return gulp.src(src.images)
    .pipe($.imageResize({
      quality: 1,
      width: 640,    // max width
      upscale: false
    }))
    .pipe(gulp.dest(dest.images + '/small'));
});

// Crop images to small squares
gulp.task('images:square', function() {
  return gulp.src(src.images)
    .pipe($.imageResize({
      quality: 1,
      gravity: 'Center',  // crop relative to center
      crop: true,
      width: 640,
      height: 640,
      upscale: false
    }))
    .pipe(gulp.dest(dest.images + '/square'));
});

gulp.task('images:build',
  gulp.series('images:clear',
    gulp.parallel('images:normal', 'images:small', 'images:square')));

gulp.task('images:watch', function() {
  $.watch(src.images, gulp.series('images:build', reload));
});

/*--------------------------------- Server ----------------------------------*/

gulp.task('server', function() {
  return bsync.init({
    startPath: '/',
    server: {
      baseDir: './',
      middleware: function(req, res, next) {
        if (req.url[0] !== '/') req.url = '/'  + req.url;

        if (/node_modules/.test(req.url) || /mitranim-master/.test(req.url) ||
            /system\.config\.js/.test(req.url) ||
            /env\.js/.test(req.url)) {
          next();
          return;
        }

        if (req.url === '/') req.url = '/' + dest.html + '/index.html';
        else req.url = '/' + dest.html + req.url;

        next();
      }
    },
    port: 11204,
    online: false,
    ui: false,
    files: false,
    ghostMode: false,
    notify: true
  });
});

/*--------------------------------- Default ---------------------------------*/

gulp.task('build', gulp.parallel(
  'scripts:build', 'styles:build', 'html:build'
));

gulp.task('watch', gulp.parallel(
  'scripts:watch', 'styles:watch', 'html:watch'
));

gulp.task('default', gulp.series('build', gulp.parallel('watch', 'server')));
