/**
 * @class
 */

angular.module('astil.models.Word', ['Datacore'])
.factory('Word', function(Record) {

  /**
   * Class.
   */
  class Word extends Record {

    constructor(attrs?) {super(attrs)}

    /**
     * Attributes.
     */

    Value:  string;
    ModeId: string;

    /**
     * Methods.
     */

    $id(): string {return this.Value}

    $path(): string {return super.$path() + '/words'}

    $valid(): boolean {
      return typeof this.Value === 'string' && /^[a-zа-я]{2,}$/.test(this.Value)
    }

  }

  /**
   * Schema.
   */
  Word.prototype.$schema = {
    Value:  '',
    ModeId: ''
  }

  /**
   * Export.
   */
  return Word

})
