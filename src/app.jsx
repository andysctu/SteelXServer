import React from 'react';
var ReactDOM = require('react-DOM')
let injectTapEventPlugin = require('react-tap-event-plugin');

import MasterContainer from './components/master-container'
 
//Needed for onTouchTap 
//Can go away when react 1.0 release 
//Check this repo: 
//https://github.com/zilverline/react-tap-event-plugin 
injectTapEventPlugin();

ReactDOM.render(
  <MasterContainer url="/contactInfo"/>,
  document.getElementById('app')
);

