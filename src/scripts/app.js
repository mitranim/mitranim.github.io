/**
 * Wires the application together.
 */

angular.module('astil.attributes', [
  'astil.attributes.activeSwitch',
  'astil.attributes.contenteditable',
  'astil.attributes.disableAll',
])

angular.module('astil.firebase', [
  'firebase',
  'astil.config',
  'foliant'
])

angular.module('astil', [
  // Templates
  'astil.templates',

  // Configuration
  'astil.config',

  // Components
  'astil.components.appWordsPage'
])
