import Firebase from 'firebase';
import React from 'react';
import _ from 'lodash';
import {Tracker} from 'tracker';
import {ReactiveVar} from 'reactive-var';
import {ReactiveDict} from 'reactive-dict';

const fbRootUrl = 'https://incandescent-torch-3438.firebaseio.com';

/**
 * References.
 */

export const root = new Firebase(fbRootUrl);

const RefMappers = {
  defaultLang(authData)  {return root.child('foliant/defaults/langs/eng')},
  defaultNames(authData) {return root.child('foliant/defaults/names/eng')},
  defaultWords(authData) {return root.child('foliant/defaults/words/eng')},
  names(authData)        {return authData ? root.child(`foliant/personal/${authData.uid}/names/eng`) : null},
  words(authData)        {return authData ? root.child(`foliant/personal/${authData.uid}/words/eng`) : null}
};

/**
 * Reactive values.
 */

function reactiveFunc(value) {
  let val = new ReactiveVar(value);
  let func = () => val.get();
  func.set = ::val.set;
  return func;
}

export const authData = reactiveFunc(null);
export const Refs = _.mapValues(RefMappers, () => reactiveFunc(null));
export const Values = _.mapValues(RefMappers, () => reactiveFunc(null));

window.authData = authData;

/**
 * Auth handlers.
 */

root.onAuth(newAuthData => {
  // When deauthed, auth anonymously.
  if (!newAuthData) root.authAnonymously(err => {if (err) throw err});

  // Refresh all reactive variables.

  authData.set(newAuthData);

  Object.keys(Refs).forEach(key => {
    // Refresh ref.
    let ref = RefMappers[key](newAuthData);
    Refs[key].set(ref);

    // Refresh value.
    if (ref) {
      let handler = ref.on('value', snap => {
        Values[key].set(snap.val());
      }, () => {
        ref.off('value', handler);
      });
    }
  });
});

// Reactively refresh names and words.
Tracker.autorun(function() {
  let namesRef = Refs.names();
  if (namesRef) {
    namesRef.on('value', snap => {
      if (!snap.val()) {
        let defNamesRef = Refs.defaultNames();
        let handler = defNamesRef.once('value', snap => {
          namesRef.set(snap.val());
        }, () => {
          namesRef.off('value', handler);
        });
      }
    });
  }

  let wordsRef = Refs.words();
  if (wordsRef) {
    wordsRef.on('value', snap => {
      if (!snap.val()) {
        let defWordsRef = Refs.defaultWords();
        let handler = defWordsRef.once('value', snap => {
          wordsRef.set(snap.val());
        }, () => {
          wordsRef.off('value', handler);
        });
      }
    });
  }
});

/**
 * Component extension.
 */

export class Component extends React.Component {
  componentWillMount() {
    if (typeof this.getState === 'function') {
      Tracker.autorun(() => {
        // Assuming this.getState() calls some functions that return
        // reactive data sources
        this.setState(this.getState());
      });
    }
  }

  componentWillUnmount() {
    // ... TODO cleanup?
  }
}
