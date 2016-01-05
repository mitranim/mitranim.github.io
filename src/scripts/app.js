import 'simple-pjax'
import './views/foliant'

document.addEventListener('click', event => {
  let elem = event.target
  do {
    if (elem.classList.contains('collapse--head')) break
  } while ((elem = elem.parentElement))

  if (!elem) return

  if ((elem = elem.parentElement) && elem.classList.contains('collapse')) {
    elem.classList.toggle('active')
  }
})
