angular.module('astil.stores.Word', [
  'astil.models.Word'
])
.factory('Words', function(Word) {

  /**
   * Class.
   */
  class WordStore extends Word {
    constructor(attrs?) {super(attrs)}
    records: Word[];
  }

  /**
   * Prototype.
   */
  WordStore.prototype.$schema = {
    records: [Word]
  }

  /**
   * Read from localStorage.
   */
  var wordStore = new WordStore()
  wordStore.$readLS()

  /**
   * Default populate.
   */
  if (!wordStore.records.length) {
    var first = [
      'nebula', 'aurora', 'quasar', 'nanite',
      'eridium', 'collapse', 'source'
    ].map(value => new Word({Value: value, ModeId: 'Mode1'}))

    var second = [
      'jasmine', 'katie', 'nariko'
    ].map(value => new Word({Value: value, ModeId: 'Mode2'}))

    wordStore.records = [...first, ...second]
  }

  return wordStore

})
