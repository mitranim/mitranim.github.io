import * as l from 'https://cdn.jsdelivr.net/npm/@mitranim/js@0.1.44/lang.mjs'
import * as d from 'https://cdn.jsdelivr.net/npm/@mitranim/js@0.1.44/dom.mjs'
import * as u from 'https://cdn.jsdelivr.net/npm/@mitranim/js@0.1.44/url.mjs'
import * as i from 'https://cdn.jsdelivr.net/npm/@mitranim/js@0.1.44/iter.mjs'
import * as dr from 'https://cdn.jsdelivr.net/npm/@mitranim/js@0.1.44/dom_reg.mjs'

dr.Reg.main.setDefiner(customElements)

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
    loc.push()

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
    for (const val of document.querySelectorAll(this.localName)) val.refresh(url)
  }
}

class Loc extends u.Loc {
  // TODO consider adding to `@mitranim/js/url.mjs`→`Query`.
  queryToggle(key, val) {
    let vals = this.query.getAll(key)
    if (vals?.includes(val)) vals = i.remove(vals, val)
    else vals = i.append(vals, val)
    this.query.setAll(key, vals)
    return this
  }
}

class TimeSink extends TagLike {
  static customName = `btn-time-sink`
  static {dr.reg(this)}
  static queryKey() {return `time_sinks`}
}

class Tag extends TagLike {
  static customName = `btn-tag`
  static {dr.reg(this)}
  static queryKey() {return `tags`}
}

class TagLikes extends d.MixNode(HTMLElement) {
  static {dr.reg(this)}

  items() {return this.descs(TagLike)}
  checked() {return i.filter(this.items(), TagLike.isChecked)}
  vals() {return i.map(this.items(), TagLike.val)}
  checkedVals() {return i.map(this.checked(), TagLike.val)}
}

class FilterList extends d.MixNode(HTMLElement) {
  static {dr.reg(this)}

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
    for (const val of document.querySelectorAll(this.localName)) val.refresh(url)
  }
}

function isVisible(val) {return !d.reqElement(val).hidden}

class FilterItem extends d.MixNode(HTMLElement) {
  static {dr.reg(this)}

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

class FilterPlaceholder extends d.MixNode(HTMLParagraphElement) {
  static {dr.reg(this)}

  refresh(found) {
    l.reqBool(found)

    if (found) {
      this.hidden = true
      return
    }

    this.textContent = `Nothing found. Try changing the filters.`
    this.hidden = false
  }
}

function main() {
  const loc = Loc.current()
  TimeSink.refresh(loc)
  Tag.refresh(loc)
  FilterList.refresh(loc)
}

main()
