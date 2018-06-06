import * as _ from 'lodash'
import * as pt from 'path'
import * as f from 'fpx'
import {compileTemplate} from 'statil/lib/template'
import {icon as faIcon} from '@fortawesome/fontawesome'
import marked from 'marked'
import hljs from 'highlight.js'

const fab = _.keyBy(require('@fortawesome/fontawesome-free-brands').default, 'iconName')
const far = _.keyBy(require('@fortawesome/fontawesome-free-regular').default, 'iconName')
const fas = _.keyBy(require('@fortawesome/fontawesome-free-solid').default, 'iconName')
const iconGroups = {fab, far, fas}

export function addKey(val, key) {
  return React.isValidElement(val) ? React.cloneElement(val, {key}) : null
}

export function innerHtmlProps(text) {
  return f.isString(text) ? {dangerouslySetInnerHTML: {__html: text}} : null
}

// Unfuck FA's disgusting API, and do sensible error reporting
export function faSvg(name, props) {
  f.validate(name, f.isString)

  if (props && props.prefix) {
    const group = iconGroups[props.prefix]
    if (!group) {
      const prefixes = Object.keys(iconGroups)
      throw Error(`Unknown prefix: ${props.prefix}. ` +
                  `Known prefixes are: \`${prefixes.join(', ')}\`.`)
    }
    if (!group[name]) {
      throw Error(`Icon \`${name}\` not found in group \`${props.prefix}\``)
    }
    return faIconToHtml(group[name], props)
  }

  const icons = {}
  for (const key in iconGroups) {
    if (iconGroups[key][name]) icons[key] = iconGroups[key][name]
  }

  const prefixes = Object.keys(icons)

  if (!prefixes.length) {
    throw Error(`Icon not found: \`${name}\``)
  }
  if (prefixes.length > 1) {
    throw Error(`Found icon \`${name}\` under more than one prefix: ` +
                `\`${prefixes.join(', ')}\`. ` +
                `Please disambiguate by passing a prefix in props.`)
  }

  return faIconToHtml(icons[prefixes[0]], props)
}

function faIconToHtml(icon, props) {
  const {html} = faIcon(icon, props)
  if (!html || !html[0]) {
    throw Error(`Unexpected FA failure on icon \`${icon.iconName || JSON.stringify(icon)}\``)
  }
  return html[0]
}

export function iconProps() {
  return innerHtmlProps(faSvg(...arguments))
}

export function renderEntryTemplate(entry) {
  return entry.body
    ? compileTemplate(entry.body, {context: {faSvg}})(entry)
    : ''
}

export function mdProps(text) {
  // Be aware: marked doesn't sanitize HTML.
  return f.isString(text) ? innerHtmlProps(md(text)) : null
}

export function current(path, subpath) {
  return path === subpath ? {'aria-current': ''} : undefined
}

// export function current() {
//   return isCurrent(...arguments) ? {'aria-current': ''} : undefined
// }

// function isCurrent(path, subpath, opts) {
//   return opts && opts.exact ? samePath(subpath, path) : pathStartsWith(subpath, path)
// }

// function toSegments(path) {
//   return f.isString(path) ? path.split('/').filter(Boolean) : undefined
// }

// function samePath(one, other) {
//   return listEqual(toSegments(one), toSegments(other))
// }

// function pathStartsWith(full, start) {
//   if (!f.isString(full) || !f.isString(start)) return false
//   start = toSegments(start)
//   return listEqual(start, toSegments(full).slice(0, start.length))
// }

// function listEqual(one, other) {
//   if (!Array.isArray(one) || !Array.isArray(other) || one.length !== other.length) {
//     return false
//   }
//   for (let i = -1; ++i < one.length;) {
//     if (!Object.is(one[i], other[i])) return false
//   }
//   return true
// }

const linkIcon = faSvg('link')

class MarkedRenderer extends marked.Renderer {
  // Adds ID anchors to headings
  heading(text, level, raw) {
    const id = this.options.headerPrefix + raw.toLowerCase().replace(/[^\w]+/g, '-')
    return (
  `<h${level}>
    <span>${text}</span>
    <a class="heading-anchor undecorate" href="#${id}" id="${id}">${linkIcon}</a>
  </h${level}>\n`
    )
  }

  // Adds target="_blank" to external links. Mostly copied from marked's source.
  link(href, title, text) {
    if (this.options.sanitize) {
      let prot = ''
      try {
        prot = decodeURIComponent(unescape(href))
          .replace(/[^\w:]/g, '')
          .toLowerCase()
      }
      catch (__) {
        return ''
      }
      if (/^(javascript|vbscript):/.test(prot)) return ''
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

  // Understands a few custom directives in code
  code(code, lang, escaped) {
    const reCollapse = /#collapse (.*)(?:\n|$)/g

    // Remove collapse directives and remember if there were any
    const collapse = reCollapse.exec(code)
    let head = ''
    if (collapse) {
      head = collapse[1]
      code = code.replace(reCollapse, '').trim()
    }

    // Default render with highlighting
    code = super.code(code, lang, escaped).trim()

    // Wrap in collapse if detected
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
}

function highlight(code, lang) {
  return lang ? hljs.highlight(lang, code).value : code
}

const markedOptions = {
  renderer: new MarkedRenderer(),
  smartypants: true,
  highlight,
}

export function md(content) {
  return marked(content, markedOptions)
    .replace(/<pre><code/g, '<pre class="padding-1"><code')
    .replace(/<!--\s*:((?:[^:]|:(?!\s*-->))*):\s*-->/g, '$1')
}

// Webpack's polyfill for the 'path' module is incomplete and doesn't have
// `parse` or any other function that allows to get a file's name
// without the extension.
export function fileName(path) {
  const ext = pt.extname(path).replace(/^[.]/, '\\.')
  return pt.basename(path).replace(new RegExp(`${ext}$`), '')
}
