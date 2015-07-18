System.config({
  'baseURL': '/',
  'paths': {
    '*': 'mitranim-master/app/*.js',
    'npm:*': 'node_modules/*.js'
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
