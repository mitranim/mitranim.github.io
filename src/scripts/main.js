require('simple-pjax')

document.addEventListener('click', onClick)

function onClick({target}) {
  const collapse = findParent(isCollapse, findParent(isCollapseHead, target))
  if (collapse) collapse.classList.toggle('active')
}

function isCollapse(elem) {
  return elem.classList.contains('collapse')
}

function isCollapseHead(elem) {
  return elem.classList.contains('collapse--head')
}

function findParent(test, elem) {
  return !(elem instanceof window.HTMLElement)
    ? null
    : test(elem)
    ? elem
    : findParent(test, elem.parentElement)
}

if (module.hot) {
  module.hot.accept(err => {
    if (err) console.warn(err)
  })
  module.hot.dispose(() => {
    console.clear()
    document.removeEventListener('click', onClick)
  })
}
