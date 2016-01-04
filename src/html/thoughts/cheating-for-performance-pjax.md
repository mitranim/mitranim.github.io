<!-- {% extend('thoughts/index.html') %} -->

## Overview

Optimising website performance is tricky. There's plenty of articles delving
deep into technical detail, like
[this great guide](https://developers.google.com/web/fundamentals/performance/critical-rendering-path/analyzing-crp?hl=en)
by Google. Naturally, when you make it that
hard, most people aren't going to bother.

What if I told you there's a way to dramatically speed up page transitions
just by adding a library? With zero or few code changes? And it's overlooked by
the contemporary blogosphere?

<p class="pad shadow" style="font-size: 1.2em">
  <span>Demo time!</span>
  <a href="http://mitranim.com/simple-pjax/" target="_blank">http://mitranim.com/simple-pjax/</a>
  <a href="https://github.com/Mitranim/simple-pjax" target="_blank" class="fa fa-github inline"></a>
</p>

What kinds of applications this applies to?

<div style="margin-left: 1rem">
  <p><span class="fa fa-check inline theme-text-primary"></span> typical server-rendered sites</p>
  <p><span class="fa fa-check inline theme-text-primary"></span> statically generated sites</p>
  <p><span class="fa fa-times inline theme-text-warn"></span> SPA (they already enjoy clientside routing)</p>
</div>

As you might have guessed, we're going to exploit clientside routing with
`history.pushState`. It's usually considered a domain of client-rendered SPA,
but what a mistake that is!

When you think about it, the status quo of content delivery on the web is
_insane_. We're forcing visitors to make dozens of network connections and
execute massive amounts of JavaScript on _each page load_ on the same site.

<div class="pad-half clamp theme-primary">
  <span class="fa fa-thumbs-o-down inline"></span>
  <span>Typical page transition</span>
</div>

<pre class="whitebox">
click link =>  ✅ download new document            =>  😫 create new JS runtime
               💀 throw away JS runtime                😫 re-execute all scripts
               💀 throw away websocket connections     🎂 display new document
               💩 304 requests for stylesheets         😫 negotiate new websocket connections
               💩 304 requests for scripts
               💩 304 requests for old images / download new images
               💩 304 requests for fonts
</pre>

With pushstate routing, we can do better.

<div class="pad-half clamp theme-primary">
  <span class="fa fa-thumbs-o-up inline"></span>
  <span>Page transition with pjax</span>
</div>

<pre class="whitebox">
click link =>  ✅ download new document          => 🎂 display new document 🎉
                  download new images if needed
</pre>

## Implementation

The idea is dead simple. Say a user navigates from page A to page B on your site.
Instead of a full page reload, fetch B by ajax, replace A, and update the URL
using `history.pushState`. This technique has been termed _`pjax`_.

Here's a super naive example to illustrate the point.

```javascript
document.addEventListener('click', function(event) {
  // Find a clicked <a>, if any.
  var anchor = event.target;
  do {
    if (anchor instanceof HTMLAnchorElement) break;
  } while (anchor = anchor.parentElement);
  if (!anchor) return;

  event.preventDefault();

  var xhr = new XMLHttpRequest();

  xhr.onload = function() {
    if (xhr.status < 200 || xhr.status > 299) return xhr.onerror();
    // Update the URL to match the clicked link.
    history.pushState(null, '', anchor.href);
    // Replace the old document with the new content.
    document.body = xhr.responseXML.body;
    window.scrollTo(0, 0);
  };

  xhr.onerror = xhr.onabort = xhr.ontimeout = function() {
    // Ensure a normal page transition.
    history.pushState(null, '', anchor.href);
    location.reload();
  };

  xhr.open('GET', anchor.href);
  // This will automatically parse the response as XML on the fly.
  xhr.responseType = 'document';
  xhr.send(null);
});
```

I have fashioned this into a simple, fully automatic
[library](https://github.com/Mitranim/simple-pjax). Just drop it into your site
and enjoy the benefits. Feedback and contributions are welcome! If you happen to
find a better implementation, I'd be happy to hear about it.

## Benefits

Despite the simplicity, the benefits are stunning. This gives your multi-page
website most of the advantages enjoyed by SPA. The browser gets to keep the same
JavaScript runtime and all downloaded assets, including images, fonts,
stylesheets, etc. This dramatically improves page load times, particularly on
poor connections such as mobile networks. This also lets you maintain a
persistent websocket connection while the user navigates your server-rendered
multi-page app!

Also, I can't overstate how wasteful it is to execute all scripts on each new
page load, which is typical for most websites. I just checked
[wired.com](http://wired.com) and the total execution time of all scripts was
**480 ms** _before_ ads kicked in. Each new page reruns all scripts. Using pjax,
you can eliminate this waste, keeping your website more responsive and saving
the visitors' CPU cycles and battery life.

## Gotchas

You need to watch out for code that modifies the DOM on page load. Most websites
have this in the form of analytics and UI widgets. When transitioning to a new
page, that code must be re-executed to modify the new document body.

Existing pjax libraries tend to emit a custom event when transitioning to a new
location, and expect you to register a listener to rerun any DOM-related code.

With `simple-pjax`, I opted for something more automatic: it reruns any inline
scripts found in the new document body, and emits the native `DOMContentLoaded`
event as a signal for the code that needs to be rerun. So far, this worked for
me. You'll need to examine your code and decide which approach fits it best.

You'll also need to take special care of widget libraries with a fragile DOM
lifecycle, like Angular or Polymer. When replacing the document body, you'll
need to run their bootstrap process again. (Notably, React doesn't have this
problem.)

## Prior Art

Pjax has been around for a few years. There are a few implementations floating
around, like the eponymous jQuery [plugin](https://github.com/defunkt/jquery-pjax).
Pjax is baked into Ruby on Rails and YUI. Many sites use it in one form or another.

Why isn't pjax more popular? Maybe because people overengineer it. The libraries
I've seen tend to focus on downloading partials (HTML snippets). They require
you to micromanage the markup, and some need a special server configuration. I
think these people have missed the point. The biggest benefit is keeping the
browsing session alive, and this can be achieved with zero configuration or
thought. For most sites, this is enough, and additional effort is usually
not worth it. Is this wrong? You tell me!

Let's use this technique to improve the web. Start with your next website!
