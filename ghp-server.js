'use strict'

// Serves files using the same HTML resolution algorithm as GitHub Pages and Netlify.
// Needs a better name, then be published.

const fs = require('fs')
const pt = require('path')
const url = require('url')
const mime = require('mime')
const {promisify} = require('util')
const fsStat = promisify(fs.stat)

exports.serve = serve
async function serve(req, res, settings) {
  const {rootDir} = Object(settings)
  if (typeof rootDir !== 'string') {
    throw Error(`Expected settings with rootDir, got ${rootDir}`)
  }
  try {
    const {status, headers, body} = await createResponse(rootDir, req.url)
    res.writeHead(status, headers)
    if (body && body.pipe) body.pipe(res)
    else res.end(body)
  }
  catch (err) {
    console.error(err)
    res.writeHead(500)
    res.end('500 Internal Server Error')
  }
}

exports.createResponse = createResponse
async function createResponse(rootDir, path) {
  if (typeof rootDir !== 'string') {
    throw Error(`Expected a root directory path, got ${rootDir}`)
  }
  if (typeof path !== 'string') {
    throw Error(`Expected a URL pathname, got ${path}`)
  }

  const pathname = url.parse(path).pathname
    //  one//two//  ->  one/two/
    .replace(/([/]+)/g, '/')
    //  /blah  ->  blah
    .replace(/^[/]/g, '')

  if (pathname === '' || pathname === '/') {
    const fsPath = pt.join(rootDir, pathname, 'index.html')
    if (await isFile(fsPath)) return fileResponse(200, fsPath)
  }
  else if (/[/]$/.test(pathname)) {
    {
      const fsPath = pt.join(rootDir, pathname, 'index.html')
      if (await isFile(fsPath)) return fileResponse(200, fsPath)
    }
    {
      const subpath = pathname.replace(/[/]$/, '')
      const fsPath = pt.join(rootDir, `${subpath}.html`)
      if (await isFile(fsPath)) return {status: 301, headers: {location: `/${subpath}`}}
    }
  }
  else {
    {
      const fsPath = pt.join(rootDir, pathname)
      if (await isFile(fsPath)) return fileResponse(200, fsPath)
    }
    {
      const fsPath = `${pt.join(rootDir, pathname)}.html`
      if (await isFile(fsPath)) return fileResponse(200, fsPath)
    }
    {
      const fsPath = pt.join(rootDir, pathname, 'index.html')
      if (await isFile(fsPath)) return {status: 301, headers: {location: `/${pathname}/`}}
    }
  }

  const fsPath = pt.join(rootDir, '404.html')
  if (await isFile(fsPath)) return fileResponse(404, fsPath)

  return {status: 404, body: '404 Not Found'}
}

async function isFile(path) {
  try {
    return (await fsStat(path)).isFile()
  }
  catch (err) {
    if (err.code === 'ENOENT') return false
    throw err
  }
}

function fileResponse(status, path) {
  const contentType = mime.getType(pt.parse(path).ext)
  return {
    status,
    headers: contentType ? {'content-type': contentType} : undefined,
    body: fs.createReadStream(path),
  }
}
