import React, {Component} from 'react'
import Tabs from 'material-ui/lib/tabs/tabs'
import Tab from 'material-ui/lib/tabs/tab'
import WorkPage from '../components/work-page'
import PlayPage from '../components/play-page'
import CreatePage from '../components/create-page'

export default class TabsComponent extends Component {
  constructor(props) {
    super(props);
  }
  render() {
    return (
      
      <Tabs style={{backgroundColor:'#ffffff'}}>
        <Tab label="Work" style={{backgroundColor:'#283593'}}>
          <WorkPage />
        </Tab>
        <Tab label="Play" style={{backgroundColor:'#283593'}}>
          <PlayPage />
        </Tab>
        <Tab label="Create" style={{backgroundColor:'#283593'}}>
          <CreatePage />
        </Tab>
      </Tabs>
      
    );
  }
}