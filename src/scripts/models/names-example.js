angular.module('astil.models.NamesExample', [
  'astil.models.Record', 'astil.models.WordsExample'
])
.factory('NamesExample', function($q, Record, WordsExample) {

  /********************************** Data ***********************************/

  function records() {
    var record = new NamesExample({
      LangId: 'eng',
      Title:  'Names',
    })
    record.$readLS()
    if (!record.Words.length) record.Words = [
      'jasmine', 'katie', 'nariko', 'karen', 'miranda'
    ]
    return [record]
  }

  /********************************** Class **********************************/

  class NamesExample extends WordsExample {

    $path(): string {return Record.prototype.$path() + '/names-example'}

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

  return NamesExample

})
