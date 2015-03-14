/**
 * Base record class.
 */
angular.module('astil.models.Record', [
  'Datacore', 'astil.config'
])
.factory('Record', function(Datacore, config) {

  class Record extends Datacore {
    $id() {return this.Id || ''}
    $path() {return config.baseUrl}
  }

  return Record

})
