'use strict'

module.exports = {
  posts: {
    'thoughts/cheating-for-performance-pjax.md': {
      slug: 'cheating-for-performance-pjax',
      title: 'Cheating for Performance: Pjax',
      description: 'Faster page transitions, for free',
      date: new Date('2015-07-25T00:00:00.000Z'),
    },
    'thoughts/next-generation-today.md': {
      slug: 'next-generation-today',
      title: 'Next Generation Today',
      description: 'EcmaScript 2015/2016 workflow with current web frameworks',
      date: new Date('2015-05-18T00:00:00.000Z'),
    },
    'thoughts/cheating-for-website-performance.md': {
      slug: 'cheating-for-website-performance',
      title: 'Cheating For Website Performance',
      description: 'Frontend tips for speeding up websites',
      date: new Date('2015-03-11T00:00:00.000Z'),
    },
    'thoughts/keeping-things-simple.md': {
      slug: 'keeping-things-simple',
      title: 'Keeping Things Simple',
      date: new Date('2015-03-10T00:00:00.000Z'),
    },
    'thoughts/old-posts.md': {
      slug: 'old-posts',
      title: 'Old Posts',
      description: 'some old stuff from around the net',
    },
  },

  misc: {
    ngTemplate: `
<div layout="space-out">
  <!-- Left column: source words -->
  <div flex="1" class="space-out">
    <h3 theme="text-primary" layout="space-between">
      <span>Source Words</span>
      <span id="indicator"></span>
    </h3>
    <form ng-submit="self.add()" layout="space-out"
          sf-tooltip="{{self.error}}" sf-trigger="{{!!self.error}}">
      <input flex="11" tabindex="1" ng-model="self.word">
      <button flex="1" theme="primary" tabindex="1">Add</button>
    </form>
    <div ng-repeat="word in self.words" layout="space-between space-out">
      <span flex="11" layout="cross-center" class="pad" style="margin-right: 1rem">{{word}}</span>
      <button flex="1" ng-click="self.remove(word)">✕</button>
    </div>
  </div>

  <!-- Right column: generated results -->
  <div flex="1" class="space-out">
    <h3 theme="text-accent">Generated Words</h3>
    <form ng-submit="self.generate()" layout>
      <button flex="1" theme="accent" tabindex="1">Generate</button>
    </form>
    <div ng-repeat="word in self.results" layout="space-between">
      <button flex="1" ng-click="self.pick(word)">←</button>
      <span flex="11" layout="cross-center" class="pad" style="margin-left: 1rem">{{word}}</span>
    </div>
    <div ng-if="self.depleted" layout="cross-center">
      <span theme="text-warn" class="pad">(depleted)</span>
    </div>
  </div>
</div>`
  }
}
