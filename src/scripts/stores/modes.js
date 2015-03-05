angular.module('astil.stores.Mode', [
  'astil.models.Mode', 'astil.stores.Word'
])
.factory('Modes', function(Mode, Words) {

  /**
   * Class.
   */
  class ModeStore extends Mode {
    constructor(attrs?) {super(attrs)}
    records: Mode[];
  }

  /**
   * Prototype.
   */
  ModeStore.prototype.$schema = {
    records: [Mode]
  }

  /**
   * Read from localStorage.
   */
  var modeStore = new ModeStore()
  modeStore.$readLS()

  /**
   * Default populate.
   */
  if (!modeStore.records.length) {
    modeStore.records = [
      new Mode({
        Id: 'Mode1',
        Title: 'Words',
        LangId: 'Lang1'
      }),
      new Mode({
        Id: 'Mode2',
        Title: 'Names',
        LangId: 'Lang1'
      })
    ]
  }

  /**
   * Assign words.
   */
  _.each(modeStore.records, mode => {
    mode.$source = _.filter(Words.records, {ModeId: mode.Id})
  })

  return modeStore

})
