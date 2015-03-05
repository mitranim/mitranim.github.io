/**
 * @class
 */

angular.module('astil.models.Mode', [
  'Datacore', 'astil.models.Word'
])
.factory('Mode', function(Record, Word) {

  /**
   * Class.
   */
  class Mode extends Record {

    constructor(attrs?) {super(attrs)}

    /**
     * Attributes.
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
     * Methods.
     */

    $id(): string {return this.Id}

    $path(): string {return super.$path() + '/modes'}

    // Maps own source words to an array of lowercase strings.
    words(): string[] {
      return _.invoke(_.map(this.$source, 'Value'), 'toLowerCase')
    }

  }

  /**
   * Schema.
   */
  Mode.prototype.$schema = {
    Id: '',

    Title: '',

    Soundset: '',

    LangId: '',

    $source: [Word],

    $generated: [Word],

    $textMode: function(): string {
      if (this.Title === 'Names') return 'text-capitalise'
      return 'text-lowercase'
    },

    $word: ''
  }

  /**
   * Export.
   */
  return Mode

})
