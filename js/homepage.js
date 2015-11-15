var MasterContainer = React.createClass ({
	getInitialState: function() {
		return {data: []};
	},
	componentDidMount: function() {
		this.loadInfoFromServer();
	},
    loadInfoFromServer: function() {
		$.ajax({
			url: this.props.url,
			dataType: 'json',
			cache: false,
			success: function(data) {
				this.setState({data: data});
			}.bind(this),
			error: function(xhr, status, err) {
				console.error(this.props.url, status, err.toString());
			}.bind(this)
		});
	},
	render: function () {
		return (
			<Header name={this.state./>
			<BodyContainer />
		);
	}
});

var Header = React.createClass({
	render: function() {
    	return (
	    	<div>
	    		<h1>Andy Tu</h1>
			</div>
    	);
	}
});

var BodyContainer = React.createClass({
	render: function() {
		return (
			<Container data=/>
		);
	}
});

var Container = React.createClass({

});

ReactDOM.render(
  <MasterContainer url="/info"/>,
  document.getElementById('content')
);

