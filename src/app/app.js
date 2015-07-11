import 'stylific';
import React from 'react';
import {Words} from 'words';

if (document.getElementById('foliantComponent')) {
  React.render(<Words />, document.getElementById('foliantComponent'));
}