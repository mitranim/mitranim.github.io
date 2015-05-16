/**
 * Static and dynamic references.
 */

import {root} from './root'

export var defaultLang = root.child('foliant/defaults/langs/eng')
export var defaultNames = root.child('foliant/defaults/names/eng')
export var defaultWords = root.child('foliant/defaults/words/eng')

export var makeNames = uid => root.child(`foliant/personal/${uid}/names/eng`)
export var makeWords = uid => root.child(`foliant/personal/${uid}/words/eng`)
