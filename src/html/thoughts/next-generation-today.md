Next generation web frameworks are around the corner.
[Angular2](https://angular.io) and [Aurelia](http://aurelia.io) go beta in a few
months. They codify custom elements as the dominant design pattern, and put
spotlight on some future and exotic JavaScript features: ES6 modules, decorated
classes, TypeScript annotations, and other.

This post will guide you through steps to using these features today, in
production, with a tried and tested current generation framework. I'll use
Angular 1.x as an example. By the end of the guide, your production-ready
Angular code may look like this:

```typescript
import {Component} from 'ng-decorate';

@Component({
  selector: 'app-tabset'
})
export class AppTabset {
  constructor() {
    this.activeTab = this.tabs[0];
  }
}
```

This guide is _massive_. I couldn't fit everything I wanted in here. You'll
probably want to read it in chunks, taking breaks. I also plan another post on
custom attributes and the DOM event pattern as the way of inter-component
communication.

## Quicklinks

* [Setup](#setup)
* [Modules](#modules)
* [Components](#components)
* [Angular DI](#dealing-with-angular-di)
* [Services](#services)
* [Demo](#demo)
* [Production Builds](#production-builds)

## Setup

You can start from scratch or grab a complete demo at GitHub:
[https://github.com/Mitranim/ng-next-gen](https://github.com/Mitranim/ng-next-gen).

### Prerequisites For: Everything

We'll start from blank. You'll need several command line utilities for package
management and builds. The guide assumes you have [Node.js](http://nodejs.org)
and [git](http://git-scm.com) installed. Install other tools with `npm`:

```sh
npm install -g jspm gulp tsd

# If you get an EACCESS error, fix permissions or use admin privileges:
sudo npm install -g jspm gulp tsd
```

Create an empty directory and navigate there. Create the following structure:

```sh
./═╦═ dist
   ╚═ src
      ╠═ app
      ╚═ html
```

For simplicity, our sample app won't have style or image compilation, but if it
did, those files would go into `src` under separate folders.

Run `npm init` to create a `package.json` file. We'll need it later.

### Prerequisites For: ES6/7 Code

We'll be writing code from the future of JavaScript. This requires a
_transpiler_ that will convert it into EcmaScript 5 code compatible with current
browsers. At the moment of writing, there are three big transpilers that support
almost all of ES6 and some of ES7: [Babel](https://babeljs.io),
[TypeScript](https://github.com/Microsoft/TypeScript), and
[Traceur](https://github.com/google/traceur-compiler). They're largely
interchangeable; this guide will use TypeScript, but you can pick Babel with
equal results.

Install the transpiler:
```sh
npm i --save-dev Microsoft/TypeScript gulp-typescript
```

We need `typescript` from the repository for its support for the SystemJS module
format. Once version 1.5.0 is released on `npm`, you can install it from there.

All three transpilers support TypeScript type annotations, so even if you're
using e.g. Babel, you can copy and paste the code from this page, and it will
work (Babel will silently remove type annotations). In your gulp pipeline, you'll
need to pass the following options to the babel transpiler:

```javascript
{stage: 0, modules: 'system'}
```

### Prerequisites For: ES6 Modules

By far the biggest ES6 feature is the new, official, module system. It finally
puts an end to the dark age of globals, AMD/CommonJS wars, and the Angular1 DI
monstrosity. This is what the syntax looks like:

```typescript
import _ from 'lodash';                 // default import
import {Attribute} from 'ng-decorate';  // named import
export class MyViewModel {/* ... */};   // named export
```

To use ES6 modules, you need two pieces:

1. A module loader running in the browser that implements the semantics of ES6
modules.

2. A transpiler that converts your ES6 import/export statements into calls to
that loader's API.

That module loader is [SystemJS](https://github.com/systemjs/systemjs). It
implements the complete semantics of ES6 modules, including circular references.
The transpiler converts your import/export statements into calls to the SystemJS
API. You also get lazy asynchronous loading, which is part of the module spec,
for free. SystemJS also consumes the AMD and CommonJS formats, so you can import
any existing libraries.

In addition, we'll use [`jspm`](http://jspm.io). It's the real package manager
for the web that replaces `bower`. It will automatically install SystemJS for
us.

Run `jspm init` to create the configuration. When asked about baseUrl (public
server path), answer `dist`. When asked about the `config.js` file, answer
`system.config.js` instead of `dist/config.js`. Press enter for all other
questions. Don't worry about it installing `traceur` — we'll pretranspile our
files, so it will never be invoked.

Open `system.config.js`, find `"paths"`, and change the import path for
application files:

`"*": "*.js"` → `"*": "app/*.js"`

Install our runtime dependencies:

```sh
jspm install angular npm:ng-decorate npm:foliant npm:stylific
```

### Prerequisites For: TypeScript

If you chose to use Babel, skip this. If you're using TypeScript, you'll want
some setup.

First, create an `src/app/tsconfig.json` with the following:

```json
#include ng-next-gen/src/app/tsconfig.json
```

`tsconfig.json` is a new feature in TypeScript 1.5. It hints your code editor
and compiler at the root of your TypeScript application.

Run `tsd init` to create a `tsd.json` file for our DefinitelyTyped definitions.
Open it and change `typings` to `src/app/typings`. The definitions need to be
inside our `src/app` for the editor and build pipeline to pick them up. Remove
the automatically created `./typings` folder.

Install the definitions (this command requires `tsd` 0.6+):
```sh
tsd install angular -r -s
```

Create an `src/app/lib.d.ts` with the following:

```typescript
#include ng-next-gen/src/app/lib.d.ts
#collapse src/app/lib.d.ts
```

### Build Configuration

Install additional tools:
```sh
npm i --save-dev gulp gulp-load-plugins gulp-rimraf gulp-plumber gulp-ng-html2js gulp-concat gulp-babel gulp-statil gulp-watch browser-sync yargs
```

Create a `gulpfile.js` with the following:

```javascript
#include ng-next-gen/gulpfile.js
#collapse gulpfile.js
```

### HTML

To keep it dead simple, we'll use just one page. Create `src/html/index.html`
with the following:

```html
<!DOCTYPE html>
<html>
  <head>
    <link rel="stylesheet" href="jspm_packages/npm/stylific@0.0.10/css/stylific.css">
  </head>
  <body>
    <sf-article>
      <word-generator></word-generator>
    </sf-article>

    <script src="jspm_packages/es6-module-loader.js"></script>
    <script src="jspm_packages/system.js"></script>
    <script src="system.config.js"></script>
    <script>
      System.import('boot')
    </script>
  </body>
</html>
```

You'll want to bundle your scripts for production; we'll deal with this at the
end of the tutorial.

At this point, we're ready to start coding!

## Modules

Our first step is to take advantage of ES6 modules. We'll disregard angular
"modules" (a more accurate name would be "DI containers"), using just one for
the entire app.

Create `src/app/app.ts`:

```typescript
import 'angular';

// Our one and only angular module.
export var app = angular.module('app', ['ng']);
```

Create `src/app/boot.ts`:

```typescript
import {app} from 'app';

// Bootstrap the app.
angular.element(document).ready(() => {
  angular.bootstrap(document.body, [app.name], {
    strictDi: true
  });
});
```

Why manual bootstrap instead of `ng-app`? This is unavoidable due to the async
nature of ES6 modules. If you include `ng-app` on the page, `angular` will
bootstrap the application before most of your application code runs. At that
point, it will be too late to run services or register directives. Manual
bootstrap solves this problem.

Invoke `gulp` to start up the pipeline and the local server. You should see a
blank page and no console errors. Now it's time to add some content.

## Components

Next generation frameworks use custom elements as building blocks of your
application. This is also the best practice in Angular 1.x, which gives you the
necessary tools in the form of directives. Here's a custom element defined with
the raw Angular 1.x API:

```typescript
import {app} from 'app';

app.directive('wordGenerator', function() {
  return {
    restrict: 'E',
    scope: {},
    templateUrl: 'word-generator/word-generator.html',
    controllerAs: 'self',
    bindToController: true,
    controller: ViewModel
  };
});

class ViewModel {}
```

All of these options are required for a proper custom element definition. This
API is pretty bad. We'll use custom decorators to make it semantic. I'm going
to cheat and import a library designed for this:
[ng-decorate](https://github.com/Mitranim/ng-decorate). We have already
installed it with `jspm`. Create `src/app/words-generator/words-generator.ts`
with:

```typescript
import {Component} from 'ng-decorate';

@Component({
  selector: 'word-generator'
})
class ViewModel {}
```

Much simpler! The decorator takes any directive options and passes them to
Angular, adding some great defaults. `ng-decorate` assumes `templateUrl` to be
`<element-name>/<element-name>.html`, which is exactly how we structure this
app. The decorated class becomes the controller (the viewmodel) of the custom
element.

Why decorators? Because you can put them at the top of a class, and they look
pretty.

You'll notice we didn't tell the decorator which angular module to use. We'll
configure the decorator library to use our main module for everything. Modify
your `src/app/app.ts`:

```diff
import 'angular';
+ import {defaults} from 'ng-decorate';

// Our one and only angular module.
export var app = angular.module('app', ['ng']);

+ // Use this module in all directive and service declarations.
+ defaults.module = app;
```

Let's add a view to this element. This is going to be a heavily simplified
version of the [foliant demo](http://mitranim.com/foliant/) because I'm lazy.

Create a file `src/app/word-generator/word-generator.html` with:

```html
#collapse src/app/word-generator/word-generator.html

<div class="flex pad-ch">
  <!-- Left column: source words -->
  <div class="flex-1 space-out-v">
    <h3 class="info pad">Source Words</h3>
    <form ng-submit="self.add()" class="flex flex-row pad-ch"
          sf-tooltip="{{self.error}}" sf-trigger="{{!!self.error}}">
      <input autofocus class="sf-input flex-11" tabindex="1" ng-model="self.word">
      <button class="sf-btn flex-1 success" tabindex="1">Add</button>
    </form>
    <div ng-repeat="word in self.words" class="flex justify-between pad-ch">
      <span class="flex-11 info pad" style="margin-right: 1rem">{{word}}</span>
      <button class="sf-btn flex-1" ng-click="self.remove(word)">✕</button>
    </div>
  </div>

  <!-- Right column: generated results -->
  <div class="flex-1 space-out-v">
    <h3 class="success pad">Generated Words</h3>
    <form ng-submit="self.generate()" class="flex flex-row">
      <button class="sf-btn pad success flex-1" tabindex="1">Generate</button>
    </form>
    <div ng-repeat="word in self.results" class="flex justify-between pad-ch">
      <button class="sf-btn flex-1" ng-click="self.pick(word)">←</button>
      <span class="flex-11 success" style="margin-left: 1rem">{{word}}</span>
    </div>
    <div ng-if="self.depleted" class="flex justify-between">
      <span class="error pad">(depleted)</span>
    </div>
  </div>
</div>
```

It won't have any functionality yet. We'll need to grab some data over ajax,
which brings us to Angular's dependency injection and services.

## Dealing With Angular DI

If your code runs before the angular application is bootstrapped, how do you get
hold of angular services that are only available through dependency injection?

You could try `injector.get`:

```typescript
var $q = angular.injector(['ng']).get('$q');
// or
var $q = angular.injector(['app', 'ng']).get('$q');
```

But this will give us the wrong instance of the injector. Angular will create
another one during the bootstrap phase, which will produce a different `$q`. Our
old instance of `$q` won't be able to automatically invoke digests in our app.
We also can't synchronously get services from our own application, if we happen
to still have code that is only available through DI.

Bottom line, you can only get hold of angular services during or after the
bootstrap phase by using `module.run`, `module.factory` or other methods that
take advantage of dependency injection. `ng-decorate` abstracts this away by
capturing injected services as static or prototype properties of the decorated
class. Example:

```typescript
import {Ambient} from 'ng-decorate';

@Ambient({
  inject: ['$q'],          // <-- will be assigned to Record.prototype
  injectStatic: ['$http']  // <-- will be assigned to Record
})
export class Record {
  /**
   * Compile-time type information.
   */
  // Prototype property.
  $q: ng.IQService;
  // Static property.
  static $http: ng.IHttpService;

  constructor() {
    console.log(this.$q);
    console.log(Record.$http);
  }
}
```

If you call `new Record()` immediately, it will log `undefined` twice. However,
if you instantiate it in a component, it will already have both services
available.

Finally, to get hold of contextual dependencies like `$scope` or `$element`,
you'll use a stock Angular feature: annotating the controller class with an
`$inject` property.

```typescript
import {Component} from 'ng-decorate';

@Component({
  selector: 'custom-element'
})
class ViewModel {
  // Compile-time type information.
  element: HTMLElement;

  static $inject = ['$element']; // stock Angular feature
  constructor($element) {
    this.element = $element[0];
  }
}
```

Now that we know how to get hold of angular services, let's take advantage of
`$http` and create a model class with ajax capability.

## Services

Create `src/models/words.ts`:

```typescript
#include ng-next-gen/src/app/models/words.ts
```

Whoah what's going on in here? Let's take this slow.

### 1. Service decorator

```
import {Service} from 'ng-decorate';

@Service({
  injectStatic: ['$http'],
  serviceName: 'Words'
})
export class Words {}
```

This is a shortcut to:

```typescript
import {app} from 'app';

app.factory('Words', ['$http', function($http) {
  Words.$http = $http;
  return Words;
}]);
export class Words {}
```

That's basically all it does. You can also include the `inject` option and it'll
assign the injected services to the prototype, same as we saw above in [Dealing
With Angular DI](#dealing-with-angular-di).

This lets you combine ES6 exports with Angular's DI. You can export it the ES6
way and still be able to get hold of injected services. The decorator will also
publish the class to the DI system, which is handy if your app has old parts
that still rely on it.

If you're writing an application from scratch and don't need DI in Karma tests,
replace `Service` with `Ambient`. It doesn't require a service name and doesn't
publish your class to Angular's DI system. Automatic dependency injection will
still work.

```diff
- import {Service} from 'ng-decorate';
+ import {Ambient} from 'ng-decorate';

- @Service({
-   injectStatic: ['$http'],
-   serviceName: 'Words'
- })
+ @Ambient({
+   injectStatic: ['$http']
+ })
export class Words {
```

### 2. Even weirder type annotations... this is not my grandfather's JavaScript!

```typescript
[key: string]: string;

/* ... */

<StringMap>response.data;
```

This is also a part of TypeScript. Simply disregard this if you're using Babel.
The former indicates what kind of data the object can hold, and the latter is
an inline type cast.

### 3. Ajax

```typescript
static readAll() {
  return this.$http({
    url: url,
    method: 'GET'
  })
  .then(response => new Words(<StringMap>response.data));
}
```

We're using the injected `$http` service to grab some example words from the
backend for the [demo](http://mitranim.com/foliant/) on which this component is
based. `this` refers to the class, and the arrow function transforms the
response, converting it into a new instance of this data model. This is a
typical pattern. In a real app, you would have a root model class that
encapsulates ajax, validation and transformation logic.

Another typical pattern is to have aggregator modules that re-export everything
from their folder. Create `src/app/models/all.ts`:

```typescript
export * from './words';
```

This is handy for maintenance reasons.

Now let's wrap this up by adding real functionality to the element.

## Demo

Modify your `src/app/boot.ts`:

```typescript
#collapse src/app/boot.ts

import {app} from 'app';

// Pull the application together.
import 'views';
import 'models/all';
import 'word-generator/word-generator';

// Bootstrap the app.
angular.element(document).ready(() => {
  angular.bootstrap(document.body, [app.name], {
    strictDi: true
  });
});
```

```diff
import {app} from 'app';

+ // Pull the application together.
+ import 'views';
+ import 'models/all';
+ import 'word-generator/word-generator';
```

Replace the contents of `src/app/word-generator/word-generator.ts` with this:

```typescript
#include ng-next-gen/src/app/word-generator/word-generator.ts
#collapse src/app/word-generator/word-generator.ts
```

Return to the page. You should see source words to the left and generated
results to the right. Congratulations! You have written a working Angular 1.x
application that takes advantage of ES6 and ES7 features, types, ES6 modules,
and a truly universal package system. The best part? This is perfectly valid for
production use.

## Production Builds

Until now, we've been importing JavaScript files over XHR. Now we'll take
advantage of `jspm`'s bundling feature to create a single self-executing
JavaScript bundle, and add some HTML templating logic to include only that link
when building for production.

Modify your `src/html/index.html`:

```html
#include ng-next-gen/src/html/index.html
#collapse src/html/index.html
```

```diff
+ <% if (prod()) { %>
+     <script src="build.js"></script>
+ <% } else { %>
    <script src="jspm_packages/es6-module-loader.js"></script>
    <script src="jspm_packages/system.js"></script>
    <script src="config.js"></script>
    <script>
      System.import('boot')
    </script>
+ <% } %>
```

Open `package.json`, find or create `"scripts"`, and add these lines:

```json
"scripts": {
  "start": "gulp",
  "bundle": "jspm bundle-sfx boot --minify",
  "build-prod": "gulp build --prod && npm run bundle",
  "serve-prod": "npm run build-prod && gulp bsync"
}
```

Now run:

```sh
npm run serve-prod
```

You should see exactly the same application, but this time, all scripts are
bundled and minified, with no external imports.

The core magic here is `jspm bundle-sfx boot`, where `boot` is the name of the
application module you're bundling. `jspm` collects this file and its entire
dependency tree into a single file that behaves exactly like our multi-file
setup in development mode.

----

That's it! You can now build modern web applications using future technologies,
with no drawbacks or compromises. Grab the complete demo over at GitHub:
[https://github.com/Mitranim/ng-next-gen](https://github.com/Mitranim/ng-next-gen)
and start playing around.

If you have any questions, grab me over at [Gitter](https://gitter.im/Mitranim)
or [Skype](skype:mitranim.web?chat).
