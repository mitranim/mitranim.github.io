import CMS from 'netlify-cms'
import * as views from '../../src/templates/layouts'

CMS.registerPreviewStyle('/styles/main.css')
CMS.registerPreviewTemplate('index', MdArticlePreview)
CMS.registerPreviewTemplate('works', MdArticlePreview)
CMS.registerPreviewTemplate('demos', MdArticlePreview)
CMS.registerPreviewTemplate('resume', SimpleMdArticlePreview)
CMS.registerPreviewTemplate('posts', PostPreview)
CMS.registerPreviewTemplate('post', PostPreview)

export function MdArticlePreview({entry}) {
  return <views.MdArticle entry={parseEntry(entry)} />
}

export function SimpleMdArticlePreview({entry}) {
  return <views.SimpleMdArticle entry={parseEntry(entry)} />
}

export function PostPreview({entry}) {
  return <views.Post entry={parseEntry(entry)} />
}

function parseEntry(entry) {
  const {path, data} = entry.toJS()
  return {path: path.replace(/src[/]templates[/]/, ''), ...data}
}
