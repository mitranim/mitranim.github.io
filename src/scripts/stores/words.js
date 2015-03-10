angular.module('astil.stores.Word', [
  'astil.models.Word'
])
.factory('Words', function(Word) {

  /**
   * Class.
   */
  class WordStore extends Word {
    /**
     * Type annotations.
     */
    records: Word[];

    /**
     * Schema.
     */
    get $schema() {return {
      records: [Word]
    }}
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
      'aurora', 'quasar', 'nanite', 'eridium', 'collapse', 'source'
    ].map(value => new Word({Value: value, ModeId: 'Mode1'}))

    var second = [
      'jasmine', 'katie', 'nariko', 'karen', 'miranda'
    ].map(value => new Word({Value: value, ModeId: 'Mode2'}))

    wordStore.records = [...first, ...second]
  }

  return wordStore

})
