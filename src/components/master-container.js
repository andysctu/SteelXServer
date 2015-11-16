import React, {Component} from 'react'
import Header from '../components/header'
import BodyContainer from '../components/body-container'
import AppBar from 'material-ui/lib/app-bar'
import TabsComponent from '../components/tabs'
//import ReactCSSTransitionGroup from 'react-addons/lib/ReactCSSTransitionGroup'

var ReactCSSTransitionGroup = require('react-addons-css-transition-group');

export default class MasterContainer extends Component {
	constructor(props) {
		super(props);
		this.state = {data: {}};
	}

	componentDidMount() {
		this.loadInfoFromServer();
	}

	loadInfoFromServer() {
		$.ajax({
			url: this.props.url,
			dataType: 'json',
			cache: false,
			success: function(data) {
				this.setState({data: data});
				console.log("data: " + data["name"]);
			}.bind(this),
			error: function(xhr, status, err) {
				console.error(this.props.url, status, err.toString());
			}.bind(this)
		});
	}

	render() {
		return (
			// <div>
			<ReactCSSTransitionGroup transitionName="example" transitionAppear={false} transitionAppearTimeout={1000} transitionLeaveTimeout={300} transitionEnterTimeout={1000}>
          	<div>
          	<Header data={this.state.data}/>
			<TabsComponent />
			<audio controls autoPlay hidden>
			 <source src='../../assets/Yundi-Li-Beethoven-Pathetique-Sonata-2nd-Movement(cut).mp3' type="audio/mp3"/>
			 Your browser does not support the audio tag.
			</audio> 
			</div>
        	
			
			</ReactCSSTransitionGroup>
			

			
			
		);
	}
};

//<AppBar title='Hello' style={{backgroundColor:'#283593'}}/>