import CMS from 'netlify-cms'
import * as views from '../../src/templates/layouts'

CMS.registerPreviewStyle('/styles/main.css')
CMS.registerPreviewTemplate('index', MdArticlePreview)

export function MdArticlePreview({entry}) {
  return <views.MdArticle entry={entry.toJS().data} />
}
