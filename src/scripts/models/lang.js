angular.module('astil.models.Lang', [
  'foliant',
  'astil.models.Record'
])
.factory('Lang', function($q, Traits, Record) {

  /********************************** Data ***********************************/

  function records() {
    return [
      {
        Id: 'eng',
        Title: 'English'
      },
      {
        Id: 'gr',
        Title: 'Greek',
        KnownSounds: [
          'α', 'β', 'γ', 'δ', 'ε', 'ζ', 'η', 'θ', 'ι', 'κ', 'λ', 'μ',
          'ν', 'ξ', 'ο', 'π', 'ρ', 'σ', 'ς', 'τ', 'υ', 'φ', 'χ', 'ψ', 'ω'
        ],
        KnownVowels: ['α', 'ε', 'η', 'ι', 'ο', 'υ', 'ω']
      }
    ].slice(0, 1).map(attrs => new Lang(attrs))
  }

  /********************************** Class **********************************/

  class Lang extends Record {

    Id: string;
    Title: string;
    $examples: [];

    get $schema() {return {
      // Strict.
      Id: '',
      Title: '',
      KnownSounds: [''],
      KnownVowels: [''],
      // Extended.
      $names: null,
      $words: null
    }}

    $path(): string {return super.$path() + '/langs'}

    // Produces a customised traits object.
    $traits(): Traits {
      var traits = new Traits()
      if (this.KnownSounds.length) {
        traits.knownSounds = new Traits.StringSet(this.KnownSounds)
      }
      if (this.KnownVowels.length) {
        traits.knownVowels = new Traits.StringSet(this.KnownVowels)
      }
      return traits
    }

    /**
     * Fake data.
     */
    static readAll() {
      return $q.when(records())
    }
    static readOne(id) {
      var record = _.find(records(), {Id: id})
      return record ? $q.when(record) : $q.reject()
    }

  }

  return Lang

})
