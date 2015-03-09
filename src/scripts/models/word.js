angular.module('astil.models.Word', ['Datacore'])
.factory('Word', function(Record) {

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

  }

  /**
   * Export.
   */
  return Word

})
