System.defaultJSExtensions = true;

System.config({
  'paths': {
    '*': 'mitranim-master/app/*',
    'npm:*': 'node_modules/*'
  }
});

System.config({
  'map': {
    'react': 'npm:react/dist/react',
    'firebase': 'npm:firebase/lib/firebase-web',
    'foliant': 'npm:foliant/dist/index',
    'lodash': 'npm:lodash/index',
    'stylific': 'npm:stylific/lib/stylific',
    'simple-pjax': 'npm:simple-pjax/simple-pjax'
  }
});
