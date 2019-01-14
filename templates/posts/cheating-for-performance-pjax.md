## Overview

Optimizing website performance is tricky. There's plenty of articles delving deep into technical detail, like [this great guide](https://developers.google.com/web/fundamentals/performance/critical-rendering-path/analyzing-crp?hl=en) by Google. Naturally, when you make it that hard, most people aren't going to bother.

What if I told you there's a way to dramatically speed up page transitions just by adding a library? With zero or few code changes? And it's overlooked by the contemporary blogosphere?

<p class="font-large">
    <span>Demo time!</span>
    {{externalAnchor "https://mitranim.com/simple-pjax/" "https://mitranim.com/simple-pjax/"}}
</p>

Who benefits from this?

<ul class="list-unstyled">
    <li><span class="fg-blue">✓</span> typical server-rendered sites</li>
    <li><span class="fg-blue">✓</span> statically generated sites</li>
    <li><span class="fg-red">✕</span> but not SPA (they already enjoy clientside routing)</li>
</ul>

As you might have guessed, we're going to exploit clientside routing with `history.pushState`. It's usually considered a domain of client-rendered SPA, but what a mistake that is!

When you think about it, the status quo of content delivery on the web is _insane_. We're forcing visitors to make dozens of network connections and execute massive amounts of JavaScript on _each page load_ on the same site.

<p class="font-large">👎 Typical page transition</p>

<ol>
    <li>Link clicked</li>

    <ul class="list-unstyled">
        <li>✅ download new document
        <li>💀 throw away JS runtime
        <li>💀 throw away websocket connections
        <li>💩 304 requests for stylesheets, scripts, old images, fonts</li>
        <li>✅ download new images if needed</li>
    </ul>

    <li>More work!</li>
    <ul class="list-unstyled">
        <li>💀 create new JS runtime</li>
        <li>💀 rerun all scripts</li>
        <li>🎂 display new document, with images and fonts flickering in</li>
        <li>💀 negotiate new websocket connections</li>
    </ul>
</ol>

With pushstate routing, we can do better.

<p class="font-large">👍 Page transition with pjax</p>

<ol>
    <li>Link clicked</li>

    <ul class="list-unstyled">
        <li>✅ download new document</li>
        <li>✅ download new images if needed</li>
    </ul>

    <li>🎂 display new document 🎉</li>
</ol>

## Implementation

The idea is dead simple. Say a user navigates from page A to page B on your site. Instead of a full page reload, fetch B by ajax, replace A, and update the URL using `history.pushState`. This technique has been termed _`pjax`_.

Here's a super naive example to illustrate the point. (DON'T COPY THIS, SEE BELOW)

```js
document.addEventListener('click', function(event) {
    // Find a clicked <a>, if any.
    const anchor = event.target
    do {
        if (anchor instanceof HTMLAnchorElement) break
    } while (anchor = anchor.parentElement)
    if (!anchor) return

    event.preventDefault()

    const xhr = new XMLHttpRequest()

    xhr.onload = function() {
        if (xhr.status < 200 || xhr.status > 299) return xhr.onerror()
        // Update the URL to match the clicked link.
        history.pushState(null, '', anchor.href)
        // Replace the old document with the new content.
        document.body = xhr.responseXML.body
        window.scrollTo(0, 0)
    }

    xhr.onerror = xhr.onabort = xhr.ontimeout = function() {
        // Ensure a normal page transition.
        history.pushState(null, '', anchor.href)
        location.reload()
    }

    xhr.open('GET', anchor.href)
    // This will automatically parse the response as XML on the fly.
    xhr.responseType = 'document'
    xhr.send(null)
})
```

I have fashioned this into a simple, fully automatic [library](https://github.com/mitranim/simple-pjax). Just drop it into your site and enjoy the benefits. Feedback and contributions are welcome! If you happen to find a better implementation, I'd be happy to hear about it.

## Benefits

Despite the simplicity, the benefits are stunning. This gives your multi-page website most of the advantages enjoyed by SPA. The browser gets to keep the same JavaScript runtime and all downloaded assets, including images, fonts, stylesheets, etc. This dramatically improves page load times, particularly on poor connections such as mobile networks. This also lets you maintain a persistent websocket connection while the user navigates your server-rendered multi-page app!

Also, I can't overstate how wasteful it is to execute all scripts on each new page load, which is typical for most websites. I just checked [wired.com](http://wired.com) and the total execution time of all scripts was **480 ms** _before_ ads kicked in. Each new page reruns all scripts. Using pjax, you can eliminate this waste, keeping your website more responsive and saving the visitors' CPU cycles and battery life.

## Gotchas

You need to watch out for code that modifies the DOM on page load. Most websites have this in the form of analytics and UI widgets. When transitioning to a new page, that code must be re-executed to modify the new document body.

Before a transition, you'll need to perform teardown like unmounting React components or destroying jQuery plugins. Do that in a document-level `simple-pjax-before-transition` event listener.

After a transition, you'll need to run the same setup as on the first page load. Do that in a document-level `simple-pjax-after-transition` event listener.

`simple-pjax` also reruns any inline scripts found in the new document body, which makes it compatible out-of-the-box with common analytics snippets.

You'll also need to take special care of widget libraries with a fragile DOM lifecycle, like Angular or Polymer. They break when document body is replaced. Notably, React is perfectly compatible; just make sure to unmount all components before replacing the body.

## Prior Art

Pjax has been around for a few years. There are a few implementations floating around, like the eponymous jQuery [plugin](https://github.com/defunkt/jquery-pjax). Pjax is baked into Ruby on Rails and YUI. Many sites use it in one form or another.

Why isn't pjax more popular? Maybe because people overengineer it. The libraries I've seen tend to focus on downloading partials (HTML snippets). They require you to micromanage the markup, and some need a special server configuration. I think these people have missed the point. The biggest benefit is keeping the browsing session alive, and this can be achieved with zero configuration or thought. For most sites, this is enough, and additional effort is usually not worth it. Is this wrong? You tell me!

Let's use this technique to improve the web!
