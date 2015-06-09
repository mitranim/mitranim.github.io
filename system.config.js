System.config({
  "baseURL": "/",
  "transpiler": "babel",
  "babelOptions": {
    "optional": [
      "runtime"
    ]
  },
  "paths": {
    "*": "app/*.js",
    "github:*": "jspm_packages/github/*.js",
    "npm:*": "jspm_packages/npm/*.js"
  }
});

System.config({
  "map": {
    "angular": "github:angular/bower-angular@1.3.15",
    "angular-sanitize": "github:angular/bower-angular-sanitize@1.3.15",
    "babel": "npm:babel-core@5.4.3",
    "babel-runtime": "npm:babel-runtime@5.4.3",
    "core-js": "npm:core-js@0.9.10",
    "firebase": "github:firebase/firebase-bower@2.2.4",
    "fireproof": "npm:fireproof@2.5.2",
    "foliant": "npm:foliant@0.0.1",
    "lodash": "npm:lodash@3.8.0",
    "ng-decorate": "npm:ng-decorate@0.0.11",
    "github:angular/bower-angular-sanitize@1.3.15": {
      "angular": "github:angular/bower-angular@1.3.15"
    },
    "github:jspm/nodelibs-process@0.1.1": {
      "process": "npm:process@0.10.1"
    },
    "npm:core-js@0.9.10": {
      "process": "github:jspm/nodelibs-process@0.1.1"
    },
    "npm:fireproof@2.5.2": {
      "child_process": "github:jspm/nodelibs-child_process@0.1.0",
      "process": "github:jspm/nodelibs-process@0.1.1",
      "systemjs-json": "github:systemjs/plugin-json@0.1.0"
    },
    "npm:foliant@0.0.1": {
      "lodash": "npm:lodash@3.8.0",
      "process": "github:jspm/nodelibs-process@0.1.1"
    },
    "npm:lodash@3.8.0": {
      "process": "github:jspm/nodelibs-process@0.1.1"
    },
    "npm:ng-decorate@0.0.11": {
      "process": "github:jspm/nodelibs-process@0.1.1"
    }
  }
});

