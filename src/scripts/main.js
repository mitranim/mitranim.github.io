/**
 * Global
 */

require('simple-pjax')

window.addEventListener('click', onClick)

window.addEventListener('keydown', () => {
  document.body.classList.add('enable-focus-indicators')
})

window.addEventListener('mousemove', () => {
  document.body.classList.remove('enable-focus-indicators')
})

/**
 * Utils
 */

function onClick({target}) {
  const collapse = findAncestor(target, isCollapse)
  if (collapse) collapse.classList.toggle('active')
}

function isCollapse(elem) {
  return isElement(elem) && elem.classList.contains('collapse')
}

function isElement(value) {
  return value instanceof window.Element
}

function findAncestor(elem, test) {
  return !elem
    ? undefined
    : test(elem)
    ? elem
    : findAncestor(elem.parentNode, test)
}
