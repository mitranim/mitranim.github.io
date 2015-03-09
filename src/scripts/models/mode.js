angular.module('astil.models.Mode', [
  'Datacore', 'astil.models.Word'
])
.factory('Mode', function(Record, Word) {

  /**
   * Class.
   */
  class Mode extends Record {

    /**
     * Type annotations.
     */

    // Strict.
    Id:       string;
    Title:    string;
    Soundset: string;
    LangId:   string;

    // Extended.
    $source:    Word[];
    $generated: Word[];
    $textMode:  string;
    $word:      string;
    $error:     string;

    /**
     * Schema.
     */
    get $schema() {return {
      // Strict.
      Id:       '',
      Title:    '',
      Soundset: '',
      LangId:   '',

      // Extended.
      $source:    [Word],
      $generated: [Word],
      $textMode: function(): string {
        return this.Title === 'Names' ? 'text-capitalise' : 'text-lowercase'
      },
      $word:  '',
      $error: ''
    }}

    /**
     * Methods.
     */

    $path(): string {return super.$path() + '/modes'}

    // Maps own source words to an array of lowercase strings.
    words(): string[] {
      return _.invoke(_.map(this.$source, 'Value'), 'toLowerCase')
    }

  }

  /**
   * Export.
   */
  return Mode

})
