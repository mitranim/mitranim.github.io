In the [previous post](/thoughts/next-generation-today), we explored new
JavaScript technologies championed by next-gen web frameworks like [Angular
2](http://angular.io) and [Aurelia](http://aurelia.io), and learned how to use
them with a current-generation framework, using Angular 1.x as example.

This is Part 2. We'll explore _custom elements_ — the most distinct feature of
newer frameworks. In particular, we'll examine how they interoperate, and
implement this behaviour in Angular 1.x.

## Quicklinks

* Overview
* Creating a Custom Element
* One-way Binding
* Communicating Through DOM Events

## Overview

Most high-profile web frameworks and libraries, including Angular 1.x, Ember,
and React, let you build reusable components. Here's what a component is:

```sh
component ─┬─ viewmodel class     ⟶    viewmodel instance
           │   │                          │
           │   ├─ methods                 ├─ methods
           │   │                          │
           │   └─ property declarations   ├─ live property bindings
           │                              │
           │                              └─ transient data (state)
           │
           └─ template            ⟶    live view
```

### Conceptual Differences

```sh
component ─┬─

 ⥣ what identifies a component?
```

* Angular: a selector string, used as an HTML tag.
* Ember: a string, used as a Handlebars identifier.
* React: an object, used as an HTML tag.
* _Angular 2, Aurelia, Polymer: a selector string, used as an HTML tag._

```sh
├─ live property bindings

    ⥣ one-way or two-way?
    ⥣ properties or attributes?
    ⥣ what is the binding syntax?
```

* Angular: two-way for exposed properties, one-way for attributes.
* Ember: two-way (optionally one-way) for exposed properties, one-way for attributes.
* React: one-way.
* _Angular 2, Aurelia, Polymer: one-way for properties_ and _attributes, optionally two-way for properties._

Everyone's binding syntax uses or resembles HTML attributes.

```sh
├─ live property bindings

    ⥣ how to propagate data changes or actions back to parent?
```

* Angular:
  * Through two-way binding.
  * Through bound expressions.
  * Through bound callbacks.
* Ember:
  * Through two-way binding.
  * Through bound callbacks.
* React:
  * Through bound callbacks.
  * (Flux) By calling methods of application-level data stores.
* _Angular 2, Aurelia, Polymer:_
  * _Through DOM events._
  * _Optionally through two-way binding._

### Convergence

Newer frameworks have unanimously converged on the following:
* Components are HTML elements, used like other tags.
* Data is bound primarily one-way.
* Changes are propagated primarily through DOM events.

This makes it easier to interoperate with native elements, which don't
understand two-way binding and propagate actions through DOM events. By
extension, this results in automatic compatibility with web components, like the
ones created with Polymer. How cool is that!

## Creating a Custom Element

This article builds on the application from Part 1, using ES6 modules and
TypeScript. If you haven't completed it, clone the
[complete demo](https://github.com/Mitranim/ng-next-gen) from the repository
to get up to speed.

Run the commands listed in the readme, and start up the local server with `npm
start`. You should see the demo page. Let's enhance it by making the generated
words draggable.

In Angular 1.x, you create a custom element with a directive:

```typescript
import {app} from 'app';

app.directive('myElement', function() {
  return {
    restrict: 'E',            // only as tagname
    templateUrl: 'my-element/my-element.html',
    scope: {},                // no scope inheritance
    controller: ViewModel,
    controllerAs: 'self',     // controller as the viewmodel
    bindToController: true    // bind properties to controller, not scope
  };
});

class ViewModel {}
```

This is noisy, so we'll cheat by using a special decorator:

```typescript
import {Component} from 'ng-decorate';

@Component({
  selector: 'my-element'
})
class ViewModel {}
```
