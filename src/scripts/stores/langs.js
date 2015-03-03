/**
 * Locally saved group of langs and mods.
 */

angular.module('astil.stores.Langs', [
  'astil.models.Lang'
])
.factory('Langs', function($rootScope, Lang) {

  var Langs = new Lang.collection()

  // Try to read from localStorage.
  Langs.$lsGet()

  /**
   * Default values.
   */
  if (!Langs.length) {
    Langs = new Lang.collection([
      {
        title: 'English',
        modes: [
          {
            title:  'Words',
            source: ['nebula', 'aurora', 'quasar', 'graphene', 'nanite', 'orchestra', 'eridium'],
          },
          {
            title:  'Names',
            source: ['jasmine', 'katie', 'nariko'],
          },
        ]
      }
    ])
  }

  $rootScope.$watch(function() {
    return [
      Langs.length,
      _.flatten(_.map(Langs, 'modes')).length,
      flatten(_.map(flatten(_.map(Langs, 'modes')), 'source')).length
    ]
  }, _.after(2, function() {
    Langs.$lsSet()
  }), true)

  return Langs

  /******************************** Utilities ********************************/

  function flatten(value) {
    var buffer = []
    if (!(value instanceof Array)) return buffer
    _.each(value, function(item) {
      if (item instanceof Array) buffer.push.apply(buffer, flatten(item))
      else buffer.push(item)
    })
    return buffer
  }

})
