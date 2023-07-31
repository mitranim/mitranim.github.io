import 'https://cdn.jsdelivr.net/npm/@ungap/custom-elements@1.0.0/es.js'
import * as l from 'https://cdn.jsdelivr.net/npm/@mitranim/js@0.1.44/lang.mjs'
import * as d from 'https://cdn.jsdelivr.net/npm/@mitranim/js@0.1.44/dom.mjs'
import * as u from 'https://cdn.jsdelivr.net/npm/@mitranim/js@0.1.44/url.mjs'
import * as i from 'https://cdn.jsdelivr.net/npm/@mitranim/js@0.1.44/iter.mjs'

class TagLike extends d.MixNode(HTMLButtonElement) {
  connectedCallback() {
    this.classList.add(`--busy`)
    this.onclick = this.onClick
  }

  onClick(eve) {
    d.eventKill(eve)
    this.toggle()

    const loc = Loc.current()
    this.mutUrl(loc)

    // Push would be better for desktop, but replace seems nicer for mobile,
    // particularly when this page is opened in a webview from another app.
    // When using push, after modifying filters, it may take multiple
    // slide-left attempts to back out of the webview.
    loc.replace()

    this.constructor.refresh(loc)
    FilterList.refresh(loc)
  }

  queryKey() {return this.constructor.queryKey()}
  static queryKey() {throw Error(`implement in subclass`)}
  isChecked() {return this.hasAttribute(`aria-checked`)}
  check() {this.setChecked(true)}
  uncheck() {this.setChecked()}
  toggle() {this.setChecked(!this.isChecked())}
  mutUrl(url) {url.queryToggle(this.queryKey(), this.val())}
  val() {return this.textContent}

  setChecked(val) {
    if (l.optBool(val)) this.setAttribute(`aria-checked`, `true`)
    else this.removeAttribute(`aria-checked`)
  }

  eq(val) {
    return (
      val?.constructor === this.constructor &&
      val.val() === this.val()
    )
  }

  static isChecked(val) {return val.isChecked()}

  static val(val) {return val.val()}

  refresh(url) {
    this.setChecked(url.query.getAll(this.queryKey())?.includes(this.val()))
  }

  static refresh(url) {
    for (const val of findAll(this)) val.refresh(url)
  }
}

function find(cls) {return d.descendant(document.body, cls)}
function findAll(cls) {return d.descendants(document.body, cls)}

class Loc extends u.Loc {
  // TODO consider adding to `@mitranim/js/url.mjs`→`Query`.
  queryToggle(key, val) {
    this.query.setAll(key, toggle(this.query.getAll(key), val))
    return this
  }
}

// TODO consider adding to `@mitranim/js/iter.mjs`.
function toggle(tar, val) {
  tar = l.laxArr(tar)
  return tar.includes(val) ? i.remove(tar, val) : i.append(tar, val)
}

class TimeSink extends TagLike {
  static queryKey() {return `time_sinks`}
}
customElements.define(`time-sink`, TimeSink, {extends: `button`})

class Tag extends TagLike {
  static queryKey() {return `tags`}
}
customElements.define(`a-tag`, Tag, {extends: `button`})

class TagLikes extends d.MixNode(HTMLElement) {
  items() {return this.descs(TagLike)}
  checked() {return i.filter(this.items(), TagLike.isChecked)}
  vals() {return i.map(this.items(), TagLike.val)}
  checkedVals() {return i.map(this.checked(), TagLike.val)}
}
customElements.define(`tag-likes`, TagLikes)

class FilterList extends d.MixNode(HTMLElement) {
  items() {return this.descs(FilterItem)}

  placeholder() {return this.desc(FilterPlaceholder)}

  refresh(url) {
    l.reqInst(url, u.Url)

    const items = i.arr(this.items())
    if (l.isEmpty(items)) return

    for (const val of items) val.refresh(url)
    this.placeholder()?.refresh(i.some(items, isVisible))
  }

  static refresh(url) {
    for (const val of findAll(this)) val.refresh(url)
  }
}
customElements.define(`filter-list`, FilterList)

function isVisible(val) {return !d.reqElement(val).hidden}

class FilterItem extends d.MixNode(HTMLElement) {
  timeSinks() {return this.descs(TimeSink)}
  tags() {return this.descs(Tag)}

  refresh(url) {
    l.reqInst(url, u.Url)
    this.hidden = false

    this.refreshTimeSinks(url)
    this.refreshTags(url)
  }

  refreshTimeSinks(url) {
    this.refreshWith(i.some, this.timeSinks(), url.query.getAll(TimeSink.queryKey()))
  }

  refreshTags(url) {
    this.refreshWith(i.every, this.tags(), url.query.getAll(Tag.queryKey()))
  }

  // TODO simplify.
  refreshWith(fun, elems, vals) {
    l.reqFun(fun)
    elems = i.arr(elems)
    vals = i.arr(vals)

    if (l.isEmpty(vals)) {
      for (const elem of elems) elem.uncheck()
      return
    }

    const elemVals = new Set()

    for (const elem of elems) {
      const val = elem.val()
      elemVals.add(val)
      elem.setChecked(vals.includes(val))
    }

    if (i.hasLen(vals) && !fun(vals, elemVals.has.bind(elemVals))) {
      this.hidden = true
    }
  }
}
customElements.define(`filter-item`, FilterItem)

class FilterPlaceholder extends d.MixNode(HTMLParagraphElement) {
  refresh(found) {this.hidden = l.reqBool(found)}
}
customElements.define(`filter-placeholder`, FilterPlaceholder, {extends: `p`})

function main() {
  const loc = Loc.current()
  TimeSink.refresh(loc)
  Tag.refresh(loc)
  FilterList.refresh(loc)
}

/*
This detection and delay are hacks for Safari, where custom element classes that
extend built-in classes other than `HTMLElement` seem to be registered
asynchronously by the polyfill we're using. In browsers with full custom
element v1 support, such as Chrome, we can run this synchronously.
*/
if (find(Tag)) main()
else globalThis.requestAnimationFrame(main)
