angular.module('astil.generator', ['foliant'])
.factory('generate', function(Traits) {

  var num = 12

  return function(source: string[]): string[] {

    var traits = new Traits(source)
    var gen = traits.generator()

    var words = [], word

    // Generate the expected number of words, excluding the source words.
    while ((word = gen()) && words.length < num) {
      if (!~source.indexOf(word)) words.push(word)
    }

    return words

  }

})
