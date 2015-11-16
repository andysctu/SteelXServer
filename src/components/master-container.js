import React, {Component} from 'react'
import Header from '../components/header'

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
			<div>
			<Header data={this.state.data}/>
			</div>
			
		);
	}
};