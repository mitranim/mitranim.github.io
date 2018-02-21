import * as _ from 'lodash'
import * as f from 'fpx'
import {renderToStaticMarkup} from 'react-dom/server'
import * as m from './misc'

function HtmlHead({children, entry: {title, description, type, image}}) {
  return (
    <head>
      <base href='/' />
      <meta charSet='utf-8' />
      <meta httpEquiv='X-UA-Compatible' content='IE=edge,chrome=1' />
      <meta name='viewport' content='width=device-width, minimum-scale=1, maximum-scale=2, initial-scale=1, user-scalable=yes' />
      <link rel='icon' href='data:;base64,=' />
      <link rel='stylesheet' type='text/css' href='styles/main.css' />
      {!process.env.PROD ? null :
      <link href='https://fonts.googleapis.com/css?family=Open+Sans:300,400,500,600,300italic,400italic,500italic,600italic&subset=latin' rel='stylesheet' />}
      <title>{title || 'about:mitranim'}</title>
      <meta name='author' content='Nelo Mitranim' />
      {!description ? null :
      <meta name='description' content={description} />}
      {!type ? null :
      <meta property='og:site_name' content='about:mitranim' />}
      {!type ? null :
      <meta property='og:type' content={type} />}
      {!title ? null :
      <meta property='og:title' content={title} />}
      {!description ? null :
      <meta property='og:description' content={description} />}
      {!image ? null :
      <meta property='og:image' content={`/images/${image}`} />}
      <DevNocache />
      {children}
    </head>
  )
}

const focusScript = (
  <script {...m.innerHtmlProps(`document.body.classList.remove('enable-focus-indicators')`)} />
)

function PageBody({className, style, children, entry}) {
  return (
    <body className='col-between-stretch enable-focus-indicators' style={{minHeight: '100vh'}}>
      {focusScript}
      <PageHeader entry={entry} />
      <div className={`flex-1 ${className || ''}`} style={style}>
        {children}
      </div>
      <PageFooter entry={entry} />
      <script src='scripts/main.js' />
      <DevReload />
    </body>
  )
}

function DocpageBody({className, style, children}) {
  return (
    <body className='col-between-stretch enable-focus-indicators'>
      {focusScript}
      <div className={className} style={style}>{children}</div>
      <DevReload />
    </body>
  )
}

function PageHeader({entry: {path}}) {
  return (
    <div className='row-start-stretch sm-flex-wrap margin-2-b text-ellipsis gaps-1-h'>
      <nav className='flex-3 flex-shrink-none row-start-stretch sm-margin-1-b'>
        <a
          href='/'
          className='font-large padding-1 busy-nav'
          {...m.current(path, 'index.mdx')}><span>home</span></a>
        <a
          href='/works'
          className='font-large padding-1 busy-nav'
          {...m.current(path, 'works.mdx')}><span>works</span></a>
        <a
          href='/posts'
          className='font-large padding-1 busy-nav'
          {...m.current(path, 'posts.mdx')}><span>posts</span></a>
        <a
          href='/demos'
          className='font-large padding-1 busy-nav'
          {...m.current(path, 'demos.mdx')}><span>demos</span></a>
      </nav>
      <span className='flex-1 text-right text-blue row-end-center'>
        <span>Updated: {new Date().toDateString()}</span>
      </span>
    </div>
  )
}

const scrollButton = `
  <button class="fg-faded padding-1" onclick="window.scrollTo(0, 0)">
    ${m.faSvg('arrow-up')}
  </button>
`.trim()

function PageFooter() {
  return (
    <footer className='row-between-center margin-4-t margin-2-b'>
      <span className='flex-1 text-left'>
        {((new Date()).getFullYear() > 2014 ? '2014—' + (new Date()).getFullYear() : '2014')}
      </span>
      <span className='flex-1 text-center'>
        <a href='/#contacts' className='decorate-link'>touch me</a>
      </span>
      <span className='flex-1 text-right' {...m.innerHtmlProps(scrollButton)} />
    </footer>
  )
}

export function MdArticle({entry}) {
  return (
    <html>
      <HtmlHead entry={entry} />
      <PageBody entry={entry}>
        <article className='fancy-typography' {...m.mdProps(m.renderEntryTemplate(entry))} />
      </PageBody>
    </html>
  )
}

export function DocpageMdArticle({entry}) {
  return (
    <html>
      <HtmlHead entry={entry} />
      <DocpageBody>
        <article
          className='fancy-typography padding-2-v'
          {...m.mdProps(m.renderEntryTemplate(entry))} />
      </DocpageBody>
    </html>
  )
}

export function Posts({tree: {posts}, entry}) {
  return (
    <html>
      <HtmlHead entry={entry} />
      <PageBody entry={entry}>
        <header className='font-large gaps-2-h margin-2-b'>
          <a href='/feed.xml' target='_blank' className='gaps-0x25-h'>
            <span {...m.iconProps('rss')} />
            <span>RSS</span>
          </a>
          <a
            href='http://feedly.com/i/subscription/feed/https://mitranim.com/feed.xml'
            target='_blank'
            className='gaps-0x25-h'>
            <span className='fg-feedly' {...m.iconProps('rss-square')} />
            <span>Feedly</span>
          </a>
          <a href='https://twitter.com/Mitranim' target='_blank' className='gaps-0x25-h'>
            <span {...m.iconProps('twitter')} />
            <span>Twitter</span>
          </a>
        </header>

        <div className='fancy-typography gaps-2-v'>
          {publicPostList(posts).map(post => {
            const link = `/posts/${post.slug}`
            return (
              <div className='gaps-1-v' key={post.slug}>
                <h2>
                  <a href={link}>{post.title}</a>
                </h2>
                {!post.description ? null :
                <p>
                  <span>{post.description}</span>
                  <a href={link} className='undecorate fg-link'> →</a>
                </p>}
                <p className='fg-faded font-small'>
                  {post.date instanceof Date ? post.date.toDateString() : post.date}
                </p>
              </div>
            )
          })}
        </div>
      </PageBody>
    </html>
  )
}

export function Post({entry, entry: {title, description, date}}) {
  return (
    <html>
      <HtmlHead entry={entry} type='article' />
      <PageBody entry={entry}>
        <article className='fancy-typography gaps-2-v'>
          <header>
            <h1>{title}</h1>
            {!description ? null :
            <h3>{description}</h3>}
            {!date ? null :
            <p className='fg-faded'>
              {date instanceof Date ? date.toDateString() : date}
            </p>}
          </header>
          <div className='fancy-typography' {...m.mdProps(m.renderEntryTemplate(entry))} />
        </article>
        <Disqus />
      </PageBody>
    </html>
  )
}

// Unsure if makes sense
function Disqus() {
  if (process.env.PROD !== 'true') return null

  return null && (
    <div>
      <div id='disqus_thread' />
      <script type='text/javascript'>{
`if (window.DISQUS) {
  window.DISQUS.reset({reload: true})
}
else {
  var disqus_shortname = 'mitranim'
  var script = document.createElement('script')
  script.type = 'text/javascript'
  script.async = true
  script.src = 'https://' + disqus_shortname + '.disqus.com/embed.js'
  document.body.appendChild(script)
}`
        }</script>
      <noscript>
        Please enable JavaScript to view the <a href='https://disqus.com/?ref_noscript' rel='nofollow'>comments powered by Disqus.</a>
      </noscript>
    </div>
  )
}

export function Page404({entry}) {
  return (
    <html>
      <HtmlHead entry={entry} />
      <PageBody entry={entry}>
        <article className='fancy-typography'>
          <h2>{entry.title}</h2>
          <p>Sorry, this page is not found.</p>
          <p><a href='/'>Return to homepage.</a></p>
        </article>
      </PageBody>
    </html>
  )
}

export function Admin() {
  return (
    <html>
      <head>
        <meta charSet='utf-8' />
        <meta httpEquiv='X-UA-Compatible' content='IE=edge,chrome=1' />
        <meta name='viewport' content='width=device-width,initial-scale=1' />
        <link rel='icon' href='data:;base64,=' />
        <link rel='stylesheet' href='/styles/cms.css' />
        <title>Content Manager</title>
        <DevNocache />
      </head>
      <body>
        <DevReload />
        <script src='/scripts/cms.js'></script>
      </body>
    </html>
  )
}

export function DevNocache() {
  if (process.env.PROD) return null
  return [
    <meta httpEquiv='cache-control' content='max-age=0' />,
    <meta httpEquiv='cache-control' content='no-cache' />,
    <meta httpEquiv='expires' content='0' />,
    <meta httpEquiv='expires' content='Tue, 01 Jan 1980 1:00:00 GMT' />,
    <meta httpEquiv='pragma' content='no-cache' />,
  ].map(m.addKey)
}

const lrScript = 'document.write(`<script src="http://${location.hostname}:35729/livereload.js"></${`script`}>`)'  // eslint-disable-line

export function DevReload() {
  if (process.env.PROD) return null
  return (
    <script dangerouslySetInnerHTML={{__html: lrScript}} />
  )
}

export function html(props) {
  return `<!doctype html>${renderWithReact(props)}`
}

export function rssFeed({tree: {posts}}) {
  return `
<?xml version="1.0" encoding="utf-8"?>
<rss version="2.0">
  <channel>
    <title>Nelo Mitranim"s Blog</title>
    <description>Occasional notes, mostly about programming</description>
    <language>en-us</language>
    <docs>http://www.rssboard.org/rss-specification</docs>
    <link type="text/html" href="https://mitranim.com" />
    ${publicPostList(posts).map(post => {
      const link = `https://mitranim.com/posts/${post.slug}`
      return `
        <item>
          <title>${post.title}</title>
          <link type="text/html" href="${link}" />
          <guid isPermaLink="true">${link}</guid>
          ${!post.date ? '' :
          `<pubDate>${post.date instanceof Date ? post.date.toISOString() : post.date}</pubDate>`}
          <author>Nelo Mitranim</author>
          ${!post.description ? '' :
          `<description>${post.description}</description>`}
        </item>
      `
    }).join('\n')}
  </channel>
</rss>
`.trim()
}

function renderWithReact(props) {
  const {layout} = Object(props.entry.papyre)
  const Layout = exports[layout]
  if (!f.isFunction(Layout)) {
    throw Error(`Expected to find layout function ${layout}, got ${Layout}`)
  }
  return renderToStaticMarkup(<Layout {...props} />)
}

function publicPostList(posts) {
  posts = _.filter(posts, post => post.showInList !== false)
  return _.sortBy(posts, post => (
    post.date instanceof Date ? post.date : -Infinity
  )).reverse()
}
