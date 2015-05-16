/********************************** Modules **********************************/

declare module 'firebase' {
  var x: FirebaseStatic;
  export default x;
}

declare module 'fireproof' {
  var x: FireproofStatic;
  export default x;
}

declare module 'foliant' {
  class StringSet {
    constructor(strings?: string[]);
    add(string: string): void;
    del(string: string): void;
    has(string: string): boolean;
  }
  class Traits {
    constructor(words?: string[]);
    static StringSet: typeof StringSet;
    examine(words: string[]): void;
    generator(): () => string;
    knownSounds: StringSet;
    knownVowels: StringSet;
    // Set of pairs of sounds that occur in the words.
    pairSet: any[];
  }
  export default Traits;
}

declare module 'ng-decorate' {
  export var Attribute;
  export var Component;
  export var Service;
}

declare module _ {
  interface LoDashStatic {
    matches(pattern: {}): (object: {}) => boolean;
    sum(object: any[], iterator?: (value: any, index: number) => number): number;
    sum(object: {}, iterator?: (value: any, key: string) => number): number;
  }
}

/******************************** Third Party ********************************/

interface Promise {
  constructor(handler: (resolve: Function, reject: Function) => void);
  then(callback: Function): Promise;
  catch(callback: Function): Promise;
}

interface FireproofStatic extends FirebaseStatic {
  /**
   * Constructs a new Fireproof reference from a Firebase instance.
   */
  new (firebase: Firebase): Fireproof;
  bless(promiseConstructor: any);
}

interface Fireproof extends ng.IPromise<any> {
  constructor(root: Firebase);
  onAuth(callback: FirebaseAuthCallback): ng.IPromise<any>;
  authWithOAuthRedirect(provider: string, onComplete?: (error: any) => void, options?: Object): ng.IPromise<any>;
  unauth(): void;
  child(path: string): Fireproof;
  on(eventName: string, callback: FirebaseSnapshotCallback): ng.IPromise<any>;
  off(eventName?: string);
  push(value?: any): Fireproof;
  set(value: any, callback?: Function): ng.IPromise<any>;
  update(values: {}): ng.IPromise<any>;
  /**
   * Removes the value at the current ref.
   */
  remove(): ng.IPromise<any>;
  key(): string;
  changePassword(settings: {
    email: string
    oldPassword: string
    newPassword: string
  }): ng.IPromise<any>;
  authWithPassword(settings: {
    email: string
    password: string
  }): ng.IPromise<any>;
  createUser(settings: {
    email: string
    password: string
  }): ng.IPromise<any>;
  resetPassword(settings: {
    email: string
  }): ng.IPromise<any>;
}

interface FirebaseAuthCallback {
  (authData: FirebaseAuthData): void
}

interface FirebaseSnapshotCallback {
  (snapshot: FirebaseDataSnapshot): void
}

/******************************** Extensions *********************************/

interface Window {
  root?: Fireproof;
  log?: (value: any) => void;
  val?: (value: any) => void;
  warn?: (value: any) => void;
  astilEnvironment?: string;
  recordBaseUrl?: string;
}

interface HTMLElement {
  querySelector(selector: string): HTMLElement;
}

interface CustomEvent {
  constructor(eventName: string, options?: {bubbles?: boolean, detail?: any});
}

/********************************** Custom ***********************************/

interface Lang {
  title: string;
  knownSounds?: string[];
  knownVowels?: string[];
}
