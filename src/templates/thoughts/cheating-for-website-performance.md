Have been turning into a bit of a performance nut lately. This is what I've
found useful for speeding up websites.

## TL:DR

* [Minify Everything](#minify-everything)
* [Concatenate Everything](#concatenate-everything)
* [Serve Static Assets](#serve-static-assets)
* [Use a Lazy Module System](#use-a-lazy-module-system)
* [Avoid JS Where Possible](#avoid-js-where-possible)
* [Use Inline Svg](#use-inline-svg)
* [Reduce Latency](#reduce-latency)
* [Consider a Static Website](#consider-a-static-website)

## Minify Everything

By far the most important thing to optimise is images. There are great free
tools like [graphicsmagick](http://www.graphicsmagick.org) that let you
automatically compress images without visible quality loss, rescale to
different dimensions, crop, etc. They can be a
[part](https://github.com/scalableminds/gulp-image-resize) of your standard
build chain, so there's absolutely no excuse for not using them. See
[example](https://github.com/Mitranim/stylific/blob/master/gulpfile.js).

Another important thing to compress is JavaScript. Modern JavaScript libraries
(and hopefully your application's code) tend to be richly commented, bloating
the source size, with the expectation of being minified for production use. With
massive frameworks like Angular, React, or Polymer, the total size easily
rockets past a megabyte. Minification gets it down to manageable size.

Minifying CSS is a very minor optimisation because there isn't much to compress.
Preprocessors like [LESS](http://lesscss.org) ignore comments, and CSS rules
can't be shortened like JS variable names. Still, it might shave off a few tens
of kilobytes.

## Concatenate Everything

Network latency is a huge deal. I can't stress this enough. Depending on the
connectivity between your servers and your users, latency could range from 50ms
to as much as a second (yes, there are areas with networks _that_ bad).

If you serve assets as multiple independent files, the browser has to make
separate network requests for each. Browsers only download a few assets at a
time, stalling other requests, which means any additional, say, stylesheets
delay the _beginning_ of loading for other assets like images or fonts. Even
when everything is cached and elicits a 304 "not modified" response, the browser
still has to wait longer before rendering the entirety of the page.

That's bad. To avoid that, make sure to concatenate assets used on each page,
like stylesheets, scripts, and icons (see below on that).

## Serve Static Assets

Double check to make sure your server is properly configured for static files
like images, stylesheets, and scripts. It should include headers that tell the
browser to cache the file, and respond with 304 for unchanged assets. This
eliminates a lot of redownloading, reducing latency+download time to latency+0.

## Use a Lazy Module System

For your JavaScript, you want a module system that avoids executing modules
unless absolutely necessary. Not lazy loading — lazy _execution_. You want to
concatenate all scripts, load them upfront (one request, then cached), and
execute only when necessary.

I'm aware of only one such system — the module system in AngularJS (1.x, not
2.x). With all other systems, including CommonJS, RequireJS, and AMD, _all_
loaded modules must run immediately. You don't get to pick and choose, each
module must execute in order to produce its export value. In Angular's system,
most code is wrapped into service functions that remain un-executed until
requested for dependency injection by another component.

What this means is that with a properly structured module dependency tree and a
component-based approach, it doesn't matter how many new JavaScript components
you add to your app, because most of it won't run. It's a huge boon for non-SPA
sites, and appears (to my knowledge) to be impossible with other widespread
module systems. The ES6 module specification also supports this approach, but
it's NYI at the time of writing.

## Avoid JS Where Possible

What's even faster than smart JavaScript? No JavaScript!

Don't rely on JavaScript for UI components that [don't
need](http://mitranim.com/stylific/components/) it. If a page doesn't require
JavaScript for its functionality, exclude it from the page completely. Running
an extremely complicated program on each page load is ridiculous. By skipping
it, you can significantly speed up page rendering. This doesn't preclude you
from using analytics scripts, but you shouldn't need everything at all times.

## Use Inline Svg

There are basically three options for an icon system: sprites, icon fonts, and
SVG icons. Raster images concatenated into one file used to be an option, but
they require pixel-perfect positioning, don't scale, look bad on retina
displays, and so on. Hopefully you've moved on to one of the other options. Icon
fonts like the aptly named [Font Awesome](http://fontawesome.io) offer a very
fine system, but still have limitations. For instance, they can't be inlined
into your stylesheet or document and have to be downloaded on page load,
flickering in. They may also be overridden by user agent rules, if the user
likes to enforce the font of their choosing.

SVG doesn't have these limitations. SVG icons can be inlined directly into your
HTML or, even better, into CSS. The LESS preprocessor even has this feature
[built-in](http://lesscss.org/functions/#misc-functions-data-uri), so all you
need is point it to an SVG file. The result is icons that are built into your
page, display instantly, don't delay other assets, and scale to any display
sharpness. (And there's FA for that,
[too](https://github.com/encharm/Font-Awesome-SVG-PNG).)

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
your audience is all over the world, pick a server with good average latency.

## Consider a Static Website

One of the biggest cheats in the book is making the _document itself_ a static
file. This is usually not possible for dynamically rendered websites that
incorporate user-created content, like wikis or forums. However, it's perfectly
feasible for a one-sided site, like a company presentation site, a personal
blog, or even a [company blog](http://facebook.github.io/react/blog/). Even
websites with backend-reliant features can cheat by being mostly static and
[offloading](/foliant/) backend interaction to an XHR API.

When even the base document is cached, it's no longer subject to latency. Some
browsers (at the time of writing, Chrome and Safari) may serve the entire page
from cache without waiting for network. Zero latency, instant load.

To top it off, static content may be duplicated to several servers all over the
world to route your visitors to the nearest. And with no backend engine
dependency, you can deploy to _any_ server.

Static site generators are [plentiful](https://www.staticgen.com), and if they
don't rock your boat, you can actually write your own in an afternoon. Like
[`statil`](https://github.com/Mitranim/statil) that generates this very website.
