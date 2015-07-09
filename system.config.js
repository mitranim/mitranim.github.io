System.config({
  "baseURL": "/",
  "transpiler": "traceur",
  "paths": {
    "*": "app/*.js",
    "github:*": "jspm_packages/github/*.js",
    "npm:*": "jspm_packages/npm/*.js"
  }
});

System.config({
  "map": {
    "angular": "github:angular/bower-angular@1.4.1",
    "angular-sanitize": "github:angular/bower-angular-sanitize@1.4.1",
    "firebase": "github:firebase/firebase-bower@2.2.7",
    "fireproof": "npm:fireproof@2.5.2",
    "foliant": "npm:foliant@0.0.1",
    "lodash": "npm:lodash@3.9.3",
    "ng-decorate": "npm:ng-decorate@0.0.15",
    "traceur": "github:jmcriffey/bower-traceur@0.0.88",
    "traceur-runtime": "github:jmcriffey/bower-traceur-runtime@0.0.88",
    "github:angular/bower-angular-sanitize@1.4.1": {
      "angular": "github:angular/bower-angular@1.4.1"
    },
    "github:jspm/nodelibs-process@0.1.1": {
      "process": "npm:process@0.10.1"
    },
    "npm:fireproof@2.5.2": {
      "child_process": "github:jspm/nodelibs-child_process@0.1.0",
      "process": "github:jspm/nodelibs-process@0.1.1",
      "systemjs-json": "github:systemjs/plugin-json@0.1.0"
    },
    "npm:foliant@0.0.1": {
      "lodash": "npm:lodash@3.9.3",
      "process": "github:jspm/nodelibs-process@0.1.1"
    },
    "npm:lodash@3.9.3": {
      "process": "github:jspm/nodelibs-process@0.1.1"
    }
  }
});

