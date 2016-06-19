require('simple-pjax')

function onclick ({target}) {
  const collapse = findParent(
    elem => elem.classList.contains('collapse'),
    findParent(elem => elem.classList.contains('collapse--head'), target)
  )
  if (collapse) collapse.classList.toggle('active')
}

function findParent (test, elem) {
  return !(elem instanceof window.HTMLElement)
    ? null
    : test(elem)
    ? elem
    : findParent(test, elem.parentElement)
}

document.addEventListener('click', onclick)

if (module.hot) {
  module.hot.accept()
  module.hot.dispose(() => {
    document.removeEventListener('click', onclick)
  })
}
