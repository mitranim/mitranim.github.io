import Firebase from 'firebase';
import React from 'react';
import _ from 'lodash';

const fbRootUrl = 'https://incandescent-torch-3438.firebaseio.com';

// Root reference for our firebase.
export const root = new Firebase(fbRootUrl);

// When deauthed, auth anonymously.
root.onAuth(authData => {
  if (!authData) root.authAnonymously(err => {if (err) throw err});
});

export class Component extends React.Component {
  _cache = {};
  _cancelers = {};

  constructor() {
    super();

    if (this.subscriptions) {
      for (let key in this.subscriptions) {
        this.subscribeToOne(key, this.subscriptions[key]);
      }
    }
  }

  subscribeToOne(key, mapper) {
    let handler = snap => {
      this._cache[key] = snap.val();
      if (this.isReady) this.setState(this._cache);
    };

    let emitter = null;
    this._cancelers[key] = () => {};

    root.onAuth(authData => {
      if (emitter) {
        emitter.off('value', handler);
        this._cancelers[key] = () => {};
      }

      let cursor = mapper(authData);
      if (typeof cursor === 'string') emitter = root.child(cursor);
      else if (cursor instanceof Firebase) emitter = cursor;

      if (emitter) {
        emitter.on('value', handler, () => {
          this._cache[key] = null;
        });
        this._cancelers[key] = () => {
          if (emitter) emitter.off('value', handler);
        }
      }
    });
  }

  get isReady() {
    for (let key in this._cancelers) {
      if (!(key in this._cache)) return false;
    }
    return true;
  }

  componentWillUnmount() {
    for (let key in this._cancelers) {
      this._cancelers[key]();
    }
  }
}

let authData = null;

root.onAuth(data => {
  authData = data;
});

export const Refs = {
  defaultLang(data)  {return root.child('foliant/defaults/langs/eng')},
  defaultNames(data) {return root.child('foliant/defaults/names/eng')},
  defaultWords(data) {return root.child('foliant/defaults/words/eng')},
  names(data)        {return (data || authData) ? root.child(`foliant/personal/${(data || authData).uid}/names/eng`) : null},
  words(data)        {return (data || authData) ? root.child(`foliant/personal/${(data || authData).uid}/words/eng`) : null}
};

root.onAuth(authData => {
  if (authData) {
    let names = Refs.names(authData);
    let namesHandler = snap => {
      if (_.isEmpty(snap.val())) {
        Refs.defaultNames().once('value', snap => {
          names.set(snap.val());
        });
      }
    };
    names.on('value', namesHandler, () => {names.off('value', namesHandler)});

    let words = Refs.words(authData);
    let wordsHandler = snap => {
      if (_.isEmpty(snap.val())) {
        Refs.defaultNames().once('value', snap => {
          words.set(snap.val());
        });
      }
    };
    words.on('value', wordsHandler, () => {words.off('value', wordsHandler)});
  }
});
