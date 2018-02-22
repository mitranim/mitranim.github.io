---
papyre: {fn: html, layout: Post}
title: Cheating For Website Performance
description: Frontend tips for speeding up websites
date: 2015-03-11T00:00:00.000Z
---

Have been turning into a bit of a performance nut lately. This is what I've
found useful for speeding up websites. These are mostly frontend optimisations;
I'm not going to delve into server performance here.

## TL:DR

* [Minify Everything](#minify-everything)
* [Concatenate Everything](#concatenate-everything)
* [Use Pjax](#use-pjax)
* [Use Server Rendering](#use-server-rendering)
* [Make Your JavaScript Lazy](#make-your-javascript-lazy)
* [Use Font Icons or Inline SVG](#use-font-icons-or-inline-svg)
* [Serve Static Assets](#serve-static-assets)
* [Reduce Latency](#reduce-latency)
* [Consider a Static Website](#consider-a-static-website)

## Minify Everything

By far the most important thing to optimise is images. There are great free
tools like [graphicsmagick](http://www.graphicsmagick.org) that let you
automatically compress images without visible quality loss, rescale to
different dimensions, crop, etc. They can be a
[part](https://github.com/scalableminds/gulp-image-resize) of your standard
build chain, so there's absolutely no excuse for not using them. See
[example](https://github.com/Mitranim/stylific/blob/master/gulpfile.js) (scroll
down to image processing).

Another important thing to compress is JavaScript. Modern JavaScript libraries
(and hopefully your application's code) tend to be richly commented, bloating
the source size, with the expectation of being minified for production use. With
massive frameworks like Angular, React, or Polymer, the total size easily
rockets past a megabyte. Minification gets it down to manageable size.

Minifying CSS is usually less important, but like everything else, it's a useful
optimisation and there's no excuse for not doing it.

## Concatenate Everything

Network latency is a huge deal. I can't stress this enough. Depending on the
connectivity between your servers and your users, latency could range from 50ms
to as much as a second.

If you serve assets as multiple independent files, the browser has to make
separate network requests for each. Browsers only download a few assets at a
time, stalling other requests, which means any additional, say, stylesheets
delay the _beginning_ of loading for other assets like images or fonts. Even
when everything is cached and elicits a 304 "not modified" response, the browser
still has to wait longer before rendering the entirety of the page.

That's bad. To avoid that, make sure to concatenate assets used on each page,
like stylesheets, scripts, and icons (see below on that).

## Use Pjax

**Update**: see this [in-depth post](/thoughts/cheating-for-performance-pjax) on pjax.

Pjax is a cheap trick that combines `history.pushState` and `ajax` to mimic page
transitions without actually reloading the page.

The basic idea is dead simple and can be implemented in a few lines of code.
Attach a document-level event listener to intercept clicks on `<a>` elements. If
the clicked link leads to an internal page, fetch the page by ajax, replace the
contents of the current page, and replace the URL using `pushState`. For
browsers that don't support this API, you simply fall back to normal page
transitions.

Despite the simplicity, the benefits are stunning. It gives you most of the
advantages enjoyed by SPA (single page applications). The browser gets to keep
the same JavaScript runtime and all downloaded assets, including images, fonts,
stylesheets, etc. This dramatically improves page load times, particularly on
poor connections such as mobile networks. This also lets you maintain a
persistent WebSocket connection while the user navigates your server-rendered
multi-page app!

There are a few [implementations](https://github.com/defunkt/jquery-pjax) in the
wild, but they require clientside _and_ server-side configuration. If you're
like me, this will seem like a waste of time. The biggest benefit of pjax is
keeping the browsing session. Micromanaging partial templates is probably not
worth your time, but everyone's needs are different.

I wrote a [simple pjax library](https://github.com/Mitranim/simple-pjax) that
works with zero config. Check the
[gotchas](https://github.com/Mitranim/simple-pjax#gotchas) to see if it's usable
for your site, then give it a spin or roll your own! The library is also used
on this very site. Inspect the network console to observe the effects.

## Use Server Rendering

There's a trend towards single page applications (SPA) with clientside routing
and rendering. They tend to skip server-side rendering in favor of being
data-driven, usually through a RESTful API. As a result, they tend to have slow
initial page loads. This is bad, particularly on slow connections, which is
typical for mobile.

Practice has shown that for consumer-facing websites, initial load time matters.
On top of that, lack of prerendering costs you SEO. Don't fall into this trap;
server rendering is a sacrifice you don't have to make. Some JavaScript UI
libraries, like React, already support isomorphic routing and rendering, and
other frameworks, like Angular 2 and Ember, are planning to support it. Make
sure to research this feature for your stack of choice.

## Make Your JavaScript Lazy

If your application is JavaScript-heavy, you should use a module system with
lazy loading. This is supported by the ES6 module system, and you can use it
today with [SystemJS](https://github.com/systemjs/systemjs) and, optionally,
[jspm](http://jspm.io). You can also achieve a similar effect with AMD.

The core parts of the application should be bundled into a single file, and big
but optional parts may be imported asynchronously when needed. If your app is
small, you can skip lazy loading and bundle the entire app.

## Use Font Icons or Inline SVG

Most sites need icons. In the past, we had to use raster images. However, in the
days of widespread retina displays, `@font-face`, and SVG, that's a poor option.
Hopefully you have switched to the vector alternatives: icon fonts and SVG
icons. They scale to any display sharpness and are easy to style with CSS.

SVGs can be embedded into the document or base64-encoded directly into your CSS,
eliminating icon flicker on page load. They can also be directly manipulated
with JavaScript for cool visual effects. On the other hand, icon fonts are
easier to set up and use, and cost less bandwidth than embedded SVGs. For most
sites, a mix of both solutions will probably be optimal.

## Serve Static Assets

This goes without saying, but you should double check to make sure your server
is properly configured for static files like images, stylesheets, and scripts.
It should include headers that tell the browser to cache the file, and respond
with 304 for unchanged assets. This eliminates a lot of redownloading, reducing
latency+download time to latency+0.

## Reduce Latency

Network latency is a huge deal. It's a part of each request made by the browser,
even for static assets with 304 responses. The browser blocks page rendering
while downloading the document and anything included in `<head>`, which defines
how snappy or sluggish your site feels. The browser may also wait for the first
few images (Firefox seems to have this tendency), or it may choose to render the
page and later flicker them into view, and latency determines how quickly this
happens.

On many sites, the document is rendered dynamically and involves database
access. This absolutely needs to be fast, but this work is usually done once per
page load. The rest comes from network latency for the document and assets. Make
sure to use a web hosting with low latency times for your target audience. If
your audience is all over the world, pick a server with good average latency and
use a caching proxy / CDN like CloudFlare to reduce latency for static content.

## Consider a Static Website

Simple websites with one maintainer, like a personal page or a blog, don't need
a scripting engine with a database. You can prerender them into HTML files, then
serve with nginx or on a service like GitHub Pages. Dynamic functionality can be
implemented with ajax.

Serving static files is naturally more performant than rendering templates on
each request. They're also automatically subject to caching. When the base
document is cached, some browsers may serve the entire page, including assets,
from the cache, rendering it with zero latency.

Static site generators are [plentiful](https://www.staticgen.com), and if they
don't float your boat, you can write your
[own](https://github.com/Mitranim/statil) in an afternoon.
