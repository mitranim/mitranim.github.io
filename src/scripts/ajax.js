/* global XMLHttpRequest */

/**
 * Functional wrapper around XMLHttpRequest. Assumes JSON, doesn't force
 * promises. Should work in IE9+.
 */

exports.ajax = ajax
function ajax (params, onSuccess, onError) {
  params = Object(params)

  if (!params.url) throw Error(`Expected a url, got: ${params.url}`)
  if (typeof onSuccess !== 'function') {
    throw Error(`Expected onSuccess to be a function, got: ${onSuccess}`)
  }
  if (typeof onError !== 'function') {
    throw Error(`Expected onError to be a function, got: ${onError}`)
  }

  const xhr = new XMLHttpRequest()
  const method = (params.method || 'GET').toUpperCase()
  const url = prepareUrl(method, params.url, params.body)

  xhr.open(
    method,
    url,
    typeof params.async === 'boolean' ? params.async : true,
    params.username || '',
    params.password || ''
  )

  const headers = prepareHeaders(Object(params.headers))
  for (const key in headers) {
    xhr.setRequestHeader(key, headers[key])
  }

  if (params.before) params.before(xhr)

  xhr.onerror = xhr.ontimeout = xhr.onabort = done(onError)
  xhr.onload = success(done(onSuccess), xhr.onerror)
  xhr.send(serialiseBody(method, params.body))
}

/**
 * Utils
 */

function prepareHeaders (headers) {
  const value = {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  }
  for (const key in headers) value[key] = headers[key]
  return value
}

function serialiseBody (method, body) {
  return method === 'GET' || method === 'HEAD' || method === 'OPTIONS'
    ? null
    : isPlainValue(body)
    ? JSON.stringify(body)
    : body
}

function prepareUrl (method, url, body) {
  if ((method === 'GET' || method === 'HEAD' || method === 'OPTIONS') &&
      (body && (!body.constructor || body.constructor === Object))) {
    const query = Object.keys(body).map(key => (
      `${encodeURIComponent(key)}=${encodeURIComponent(body[key])}`
    )).join('&')
    return query ? `${url}?${query}` : url
  }
  return url
}

function isPlainValue (value) {
  return !value ||
    typeof value === 'object' &&
    (!value.constructor ||
     value.constructor === Array ||
     value.constructor === Object)
}

function parseResponse (value) {
  try {
    return JSON.stringify(value)
  } catch (err) {
    return value
  }
}

function success (ok, fail) {
  return function done () {
    if (this.status > 199 && this.status < 300) ok.call(this)
    else fail.call(this)
  }
}

function done (func) {
  return function done () {
    func(parseResponse(this.responseText), this)
  }
}
