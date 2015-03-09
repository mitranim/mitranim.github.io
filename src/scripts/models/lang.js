angular.module('astil.models.Lang', [
  'Datacore', 'astil.models.Mode'
])
.factory('Lang', function(Record, Mode) {

  /**
   * Class.
   */
  class Lang extends Record {

    /**
     * Type annotations.
     */
    Id:     string;
    Title:  string;
    $modes: Mode[];

    /**
     * Schema.
     */
    get $schema() {return {
      Id:     '',
      Title:  '',
      $modes: [Mode]
    }}

    /**
     * Methods.
     */
    $path(): string {return super.$path() + '/langs'}

  }

  /**
   * Export.
   */
  return Lang

})
