/**
 * Shamelessly extracted from Meteor source and translated into ES6
 * https://github.com/meteor/meteor/blob/devel/packages/reactive-dict/reactive-dict.js
 *
 * TODO port EJSON
 */

import {Tracker} from './tracker';
import _ from 'lodash';

// Mock
const EJSON = {
  equals(one, other) {
    one = JSON.parse(JSON.stringify(one));
    other = JSON.parse(JSON.stringify(other));
    return _.isEqual(one, other);
  },
  parse(value) {
    return JSON.parse(value);
  },
  stringify(value) {
    return JSON.stringify(value);
  }
};

// XXX COMPAT WITH 0.9.1 : accept migrationData instead of dictName
export class ReactiveDict {
  constructor(dictName) {
    // this.keys: key -> value
    if (dictName) {
      if (typeof dictName === 'string') {
        // the normal case, argument is a string name.
        // _registerDictForMigrate will throw an error on duplicate name.
        ReactiveDict._registerDictForMigrate(dictName, this);
        this.keys = ReactiveDict._loadMigratedDict(dictName) || {};
      } else if (typeof dictName === 'object') {
        // back-compat case: dictName is actually migrationData
        this.keys = dictName;
      } else {
        throw new Error("Invalid ReactiveDict argument: " + dictName);
      }
    } else {
      // no name given; no migration will be performed
      this.keys = {};
    }

    this.allDeps = new Tracker.Dependency;
    this.keyDeps = {}; // key -> Dependency
    this.keyValueDeps = {}; // key -> Dependency
  }

  // set() began as a key/value method, but we are now overloading it
  // to take an object of key/value pairs, similar to backbone
  // http://backbonejs.org/#Model-set

  set(keyOrObject, value) {
    if ((typeof keyOrObject === 'object') && (value === undefined)) {
      this._setObject(keyOrObject);
      return;
    }
    // the input isn't an object, so it must be a key
    // and we resume with the rest of the function
    let key = keyOrObject;

    value = stringify(value);

    let oldSerializedValue = 'undefined';
    if (_.has(this.keys, key)) oldSerializedValue = this.keys[key];
    if (value === oldSerializedValue) return;
    this.keys[key] = value;

    this.allDeps.changed();
    changed(this.keyDeps[key]);
    if (this.keyValueDeps[key]) {
      changed(this.keyValueDeps[key][oldSerializedValue]);
      changed(this.keyValueDeps[key][value]);
    }
  }

  setDefault(key, value) {
    // for now, explicitly check for undefined, since there is no
    // ReactiveDict.clear().  Later we might have a ReactiveDict.clear(), in which case
    // we should check if it has the key.
    if (this.keys[key] === undefined) {
      this.set(key, value);
    }
  }

  get(key) {
    this._ensureKey(key);
    this.keyDeps[key].depend();
    return parse(this.keys[key]);
  }

  equals(key, value) {
    // Mongo.ObjectID is in the 'mongo' package
    let ObjectID = null;

    // We don't allow objects (or arrays that might include objects) for
    // .equals, because JSON.stringify doesn't canonicalize object key
    // order. (We can make equals have the right return value by parsing the
    // current value and using EJSON.equals, but we won't have a canonical
    // element of keyValueDeps[key] to store the dependency.) You can still use
    // "EJSON.equals(reactiveDict.get(key), value)".
    //
    // XXX we could allow arrays as long as we recursively check that there
    // are no objects
    if (typeof value !== 'string' &&
        typeof value !== 'number' &&
        typeof value !== 'boolean' &&
        typeof value !== 'undefined' &&
        !(value instanceof Date) &&
        !(ObjectID && value instanceof ObjectID) &&
        value !== null)
      throw new Error("ReactiveDict.equals: value must be scalar");
    let serializedValue = stringify(value);

    if (Tracker.active) {
      this._ensureKey(key);

      if (! _.has(this.keyValueDeps[key], serializedValue)) {
        this.keyValueDeps[key][serializedValue] = new Tracker.Dependency;
      }

      let isNew = this.keyValueDeps[key][serializedValue].depend();
      if (isNew) {
        Tracker.onInvalidate(() => {
          // clean up [key][serializedValue] if it's now empty, so we don't
          // use O(n) memory for n = values seen ever
          if (!this.keyValueDeps[key][serializedValue].hasDependents())
            delete this.keyValueDeps[key][serializedValue];
        });
      }
    }

    let oldValue = undefined;
    if (_.has(this.keys, key)) oldValue = parse(this.keys[key]);
    return EJSON.equals(oldValue, value);
  }

  all() {
    this.allDeps.depend();
    let ret = {};
    _.each(this.keys, (value, key) => {
      ret[key] = parse(value);
    });
    return ret;
  }

  clear() {
    let oldKeys = this.keys;
    this.keys = {};

    this.allDeps.changed();

    _.each(oldKeys, (value, key) => {
      changed(this.keyDeps[key]);
      changed(this.keyValueDeps[key][value]);
      changed(this.keyValueDeps[key]['undefined']);
    });
  }

  _setObject(object) {
    _.each(object, (value, key) => {
      this.set(key, value);
    });
  }

  _ensureKey(key) {
    if (!(key in this.keyDeps)) {
      this.keyDeps[key] = new Tracker.Dependency;
      this.keyValueDeps[key] = {};
    }
  }

  // Get a JSON value that can be passed to the constructor to
  // create a new ReactiveDict with the same contents as this one
  _getMigrationData() {
    // XXX sanitize and make sure it's JSONible?
    return this.keys;
  }
}

// XXX come up with a serialization method which canonicalizes object key
// order, which would allow us to use objects as values for equals.
function stringify(value) {
  if (value === undefined) return 'undefined';
  return EJSON.stringify(value);
}

function parse(serialized) {
  if (serialized === undefined || serialized === 'undefined') return undefined;
  return EJSON.parse(serialized);
}

function changed(v) {
  v && v.changed();
}
