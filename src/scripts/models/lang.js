/**
 * @class
 */

angular.module('astil.models.Lang', [
  'Datacore', 'astil.models.Mode'
])
.factory('Lang', function(Record, Mode) {

  /**
   * Class.
   */
  class Lang extends Record {

    constructor(attrs?) {super(attrs)}

    /**
     * Attributes.
     */

    Id:     string;
    Title:  string;
    $modes: Mode[];

    /**
     * Methods.
     */

    $path(): string {return super.$path() + '/langs'}
  }

  /**
   * Schema.
   */
  Lang.prototype.$schema = {
    Id:     '',
    Title:  '',
    $modes: [Mode]
  }

  /**
   * Export.
   */
  return Lang

})
