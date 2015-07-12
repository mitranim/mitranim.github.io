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

export const authDataVal = new ReactiveVar(null);
export const Refs = _.mapValues(RefMappers, () => new ReactiveVar(null));
export const Values = _.mapValues(RefMappers, () => new ReactiveVar(null));

/**
 * Auth handlers.
 */

root.onAuth(authData => {
  // When deauthed, auth anonymously.
  if (!authData) root.authAnonymously(err => {if (err) throw err});

  /**
   * Refresh all reactive variables.
   */

  authDataVal.set(authData);

  for (let key in Refs) {
    // Refresh ref.
    let ref = RefMappers[key](authData);
    Refs[key].set(ref);

    // Refresh value.
    if (ref) {
      let handler = ref.on('value', snap => {
        Values[key].set(snap.val());
      }, () => {
        ref.off('value', handler);
      });
    }
  }
});

// Reactively refresh names and words.
Tracker.autorun(function() {
  let namesRef = Refs.names.get();
  if (namesRef) {
    namesRef.on('value', snap => {
      if (!snap.val()) {
        let defNamesRef = Refs.defaultNames.get();
        let handler = defNamesRef.once('value', snap => {
          namesRef.set(snap.val());
        }, () => {
          namesRef.off('value', handler);
        });
      }
    });
  }

  let wordsRef = Refs.words.get();
  if (wordsRef) {
    wordsRef.on('value', snap => {
      if (!snap.val()) {
        let defWordsRef = Refs.defaultWords.get();
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
        this.setState(this.getState());
      });
    }
  }

  componentWillUnmount() {
    // ... TODO cleanup
  }
}
