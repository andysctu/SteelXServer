import React, {Component} from 'react';
import Jumbotron from 'react-bootstrap/lib/Jumbotron';

export default class Header extends Component {
	constructor(props) {
		super(props);
	}
	render() {
		return (
			<div className="bg">
			<Jumbotron style={{padding:25, background:'transparent', textShadow:'black 0.1em 0.1em 0.1em', color:'white'}}>
			    <h1>{this.props.data.name}</h1>
			    <p>{this.props.data.email}<br/>{this.props.data.phone}</p>
		  	</Jumbotron>
		  	</div>
		);
	}
};