import React, {Component} from 'react'
import Tabs from 'material-ui/lib/tabs/tabs'
import Tab from 'material-ui/lib/tabs/tab'
import WorkPage from '../components/work-page'
import PlayPage from '../components/play-page'
import CreatePage from '../components/create-page'
import SwipeableViews from 'react-swipeable-views'

export default class TabsComponent extends Component {
  constructor(props) {
    super(props);
    this._handleTabActive = this._handleTabActive.bind(this);
    this._handleChangeTabs = this._handleChangeTabs.bind(this);
    this.state = {
      slideIndex: 0,
    }
  }
  _handleChangeIndex(index) {
    this.setState({
      slideIndex: index,
    });
  }

  _handleChangeTabs(value) {
    this.setState({
      slideIndex: parseInt(value, 10),
    });
  }

  _handleButtonClick() {
    this.setState({tabsValue: 'c'});
  }

  _handleTabActive(tab){
    this.props.history.pushState(null, tab.props.route);
  }

  render() {
    return (
      <div>
      <Tabs onChange={this._handleChangeTabs.bind(this)} value={this.state.slideIndex + ''} style={{backgroundColor:'#ffffff'}}>
        <Tab label="Work" value="0" style={{backgroundColor:'#283593'}}/>
        <Tab label="Play" value="1" style={{backgroundColor:'#283593'}}/>
        <Tab label="Create" value="2" style={{backgroundColor:'#283593'}}/>
      </Tabs>
      <SwipeableViews index={this.state.slideIndex} onChangeIndex={this._handleChangeIndex.bind(this)} style={{height:1000}}>
        
          <WorkPage />
        
        
          <PlayPage />
        
        
          <CreatePage />
        
      </SwipeableViews>
      </div>
      // <Tabs style={{backgroundColor:'#ffffff'}}>
      //   <Tab label="Work" style={{backgroundColor:'#283593'}}>
      //     <WorkPage />
      //   </Tab>
      //   <Tab label="Play" style={{backgroundColor:'#283593'}}>
      //     <PlayPage />
      //   </Tab>
      //   <Tab label="Create" style={{backgroundColor:'#283593'}}>
      //     <CreatePage />
      //   </Tab>
      // </Tabs>
      
    );
  }
}