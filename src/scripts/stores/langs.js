angular.module('astil.stores.Lang', [
  'astil.models.Lang', 'astil.stores.Mode'
])
.factory('Langs', function(Lang, Modes) {

  /**
   * Class.
   */
  class LangStore extends Lang {
    constructor(attrs?) {super(attrs)}
    records: Lang[];
  }

  /**
   * Prototype.
   */
  LangStore.prototype.$schema = {
    records: [Lang]
  }

  /**
   * Read from localStorage.
   */
  var langStore = new LangStore()
  langStore.$readLS()

  /**
   * Default populate.
   */
  if (!langStore.records.length) {
    langStore.records = [
      new Lang({
        Id: 'Lang1',
        Title: 'English'
      })
    ]
  }

  /**
   * Assign modes.
   */
  _.each(langStore.records, lang => {
    lang.$modes = _.filter(Modes.records, {LangId: lang.Id})
  })

  return langStore

})
