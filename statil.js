'use strict'

const hljs = require('highlight.js')
const marked = require('marked')
const _ = require('lodash')
const pt = require('path')
const cheerio = require('cheerio')
const flags = require('yargs').boolean('prod').argv

/*
 * Markdown config
 */

marked.setOptions({
  smartypants: true,
  highlight (code, lang) {
    const result = lang ? hljs.highlight(lang, code) : hljs.highlightAuto(code)
    return result.value
  }
})

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
      '<div class="collapse">\n' +
      '  <div class="collapse--head bg-smoke">' + head + '</div>\n' +
      '  <div class="collapse--body">\n' +
           code + '\n' +
      '  </div>\n' +
      '</div>'
  }

  return code
}

/*
 * Statil config
 */

function statilData () {
  const path = _.find(_.keys(require.cache), path => (
    /html-meta/.test(path)
  ))
  if (path) delete require.cache[path]
  return require('./html-meta')
}

module.exports = function statilOptions () {
  return {
    data: statilData(),
    ignorePaths: path => (
      path === 'thoughts/index.html' ||
      /^partials/.test(path)
    ),
    rename: '$&/index.html',
    renameExcept: ['index.html', '404.html'],
    imports: {
      prod: flags.prod,
      truncate (html, length) {
        return _.trunc(cheerio(html).text(), length)
      },
      sortPosts: posts => _.sortBy(posts, post => (
        post.date instanceof Date ? post.date : -Infinity
      )).reverse()
    },
    pipeline: [
      (content, path) => {
        if (pt.extname(path) === '.md') {
          return marked(content).replace(/<pre><code class="(.*)">|<pre><code>/g, '<pre><code class="hljs $1">')
        }
      }
    ]
  }
}
