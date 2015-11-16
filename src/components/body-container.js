import React, {Component} from 'react'
import InfoBlock from './infoBlock'

export default class BodyContainer extends Component {
	constructor(props) {
		super(props);
		this.state = {data: {}};
	}

	componentDidMount() {
	
	}

	render() {
		return (
			<div className="body-container">
				<InfoBlock />

			</div>
			
		);
	}
};