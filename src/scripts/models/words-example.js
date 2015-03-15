angular.module('astil.models.WordsExample', [
  'astil.models.Record'
])
.factory('WordsExample', function($q, Record) {

  /********************************** Data ***********************************/

  function records() {
    var record = new WordsExample({
      LangId: 'eng',
      Title:  'Words',
    })
    record.$readLS()
    if (!record.Words.length) record.Words = [
      'aurora', 'quasar', 'nanite', 'eridium', 'collapse', 'source'
    ]
    return [record]
  }

  /********************************** Class **********************************/

  class WordsExample extends Record {

    LangId: string;
    Words: string[];

    get $schema() {return {
      LangId: '',
      Title:  '',
      Words:  ['']
    }}

    $path(): string {return super.$path() + '/words-example'}

    // Belongs to Lang.
    $id(): string {return this.LangId}

    /**
     * Fake data.
     */
    static readAll() {
      return $q.when(records())
    }
    static readOne(id) {
      var record = _.find(records(), {Id: id})
      return record ? $q.when(record) : $q.reject()
    }

  }

  return WordsExample

})
