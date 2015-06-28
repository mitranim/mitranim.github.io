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

declare module _ {
  interface LoDashStatic {
    matches(pattern: {}): (object: {}) => boolean;
    sum(object: any[], iterator?: (value: any, index: number) => number): number;
    sum(object: {}, iterator?: (value: any, key: string) => number): number;
  }
}

declare module 'ng-decorate' {
  export var Attribute: typeof ngDecorate.Attribute;
  export var Ambient: typeof ngDecorate.Ambient;
  export var Component: typeof ngDecorate.Component;
  export var Service: typeof ngDecorate.Service;
  export var Controller: typeof ngDecorate.Controller;
  export var autoinject: typeof ngDecorate.autoinject;
  export var bindTwoWay: typeof ngDecorate.bindTwoWay;
  export var bindOneWay: typeof ngDecorate.bindOneWay;
  export var bindString: typeof ngDecorate.bindString;
  export var bindExpression: typeof ngDecorate.bindExpression;
  export var defaults: typeof ngDecorate.defaults;
}

declare module ngDecorate {
  // Class decorators.
  export function Attribute(config: DirectiveConfig): ClassDecorator;
  export function Ambient(config: BaseConfig): ClassDecorator;
  export function Ambient(target: Function): void;
  export function Component(config: DirectiveConfig): ClassDecorator;
  export function Service(config: ServiceConfig): ClassDecorator;
  export function Controller(config: ControllerConfig): ClassDecorator;

  // Property decorators.
  export function autoinject(target: any, key: string);
  export function bindTwoWay(options: BindTwoWayOptions): PropertyDecorator;
  export function bindTwoWay(target: any, key: string): void;
  export function bindOneWay(key: string): PropertyDecorator;
  export function bindOneWay(target: any, key: string): void;
  export function bindString(key: string): PropertyDecorator;
  export function bindString(target: any, key: string): void;
  export function bindExpression(key: string): PropertyDecorator;
  export function bindExpression(target: any, key: string): void;

  // Mutable configuration.
  export const defaults: {
    module?: ng.IModule;
    moduleName?: string;
    controllerAs: string;
    makeTemplateUrl: (selector: string) => string;
  };

  // Abstract interface shared by configuration objects.
  interface BaseConfig {
    // Angular module object. If provided, other module options are ignored, and
    // no new module is declared.
    module?: ng.IModule;

    // Optional name for the new module created for this service or directive.
    // If omitted, the service or directive name is used.
    moduleName?: string;

    // Names of other angular modules this module depends on.
    dependencies?: string[];

    // DEPRECATED in favour of @autoinject.
    // Angular services that will be assigned to the class prototype.
    inject?: string[];

    // DEPRECATED in favour of @autoinject.
    // Angular services that will be assigned to the class as static properties.
    injectStatic?: string[];
  }

  interface DirectiveConfig extends BaseConfig, ng.IDirective {
    // The name of the custom element or attribute. Used to derive module name,
    // directive name, and template url.
    selector: string;
  }

  interface ServiceConfig extends BaseConfig {
    // The name of the service in the angular module system. Mandatory
    // due to minification woes.
    serviceName: string;
  }

  interface ControllerConfig extends BaseConfig {
    // Mandatory controller name.
    controllerName: string;
    // Optional service name. If included, the controller is published to
    // angular's DI as a service under this name.
    serviceName?: string;
  }

  interface ControllerClass extends Function {
    template?: string|Function;
    templateUrl?: string|Function;
    link?: Function;
    compile?: any;
  }

  interface BindTwoWayOptions {
    // Adds `*` to the property descriptor, marking it for `$watchCollection`.
    collection?: boolean;
    // Adds `?` to the property descriptor, marking it optional.
    optional?: boolean;
    // Adds an external property name to the binding.
    key?: string;
  }
}

/******************************** Third Party ********************************/

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
