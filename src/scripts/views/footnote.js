export function footnote (render) {
  render(
    ['div', {className: 'text-right container'},
      ['p', null,
        'Made with ',
        ['a', {href: 'https://github.com/Mitranim/prax', target: '_blank'}, 'Prax'],
        ' and ',
        ['a', {href: 'https://github.com/Mitranim/alder', target: '_blank'}, 'Alder'],
        '.']]
  )
}
