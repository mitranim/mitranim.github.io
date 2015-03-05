/**
 * Form component mixin.
 */

angular.module('astil.mixins.form', [])
.factory('mixinForm', function($q) {

  return function(self) {

    /**
     * Submits the record associated with the form. If this fails, displays errors.
     */
    self.submit = function(event): Promise {
      // Pre-validate.
      var errors = self.record.$validate()
      if (!_.isEmpty(errors)) {
        self.showErrors(errors)
        console.log("-- errors:", errors)
        return $q.reject()
      }

      return self.record.$save()
      .then(function() {
        console.log("-- successfully saved record:", self.record)
        self.hideErrors()
        self.done = true
      })
      .catch(function(err) {
        console.error("-- err:", err)
        self.showErrors(err.data)
      })
    }

    /**
     * Aligns own record with the given errors and displays them. In the record,
     * each string field is reset to '' to display the error in that field as a
     * placeholder.
     */
    self.showErrors = function(errors) {
      // Reset matching fields to display errors as placeholder text.
      _.each(errors, function(value, key) {
        if (typeof self.record[key] === 'string') self.record[key] = ''
      })
      self.errors = errors
    }

    /**
     * Removes the error messages.
     */
    self.hideErrors = function() {
      delete self.errors
    }

  }

})
