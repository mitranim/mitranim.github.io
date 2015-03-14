angular.module('astil.models.Word', [
  'foliant', 'astil.models.Record'
])
.factory('Word', function($q, Traits, Record) {

  /**
   * Class.
   */
  class Word extends Record {

    /**
     * Type annotations.
     */
    Value:  string;
    ModeId: string;

    /**
     * Schema.
     */
    get $schema() {return {
      Value:  '',
      ModeId: ''
    }}

    /**
     * Methods.
     */

    $id(): string {return this.Value}

    $path(): string {return super.$path() + '/words'}

    $valid(): boolean {
      return typeof this.Value === 'string' && /^[a-zа-я]{2,}$/.test(this.Value)
    }

    static readAll(options: ?{}): Promise {
      if (options == null || typeof options !== 'object') return $q.reject()
      var params = options.params
      if (params == null || typeof params !== 'object') return $q.reject()
      if (!(params.words instanceof Array)) return $q.reject()
      if (!params.words.length) return $q.reject()

      var traits = new Traits(params.words)
      var gen = traits.generator()
      var words = _.times(12, gen).map(word => ({Value: word}))
      return $q.when(words)
    }

  }

  /**
   * Export.
   */
  return Word

})
