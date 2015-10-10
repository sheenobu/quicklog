

var Header = React.createClass({
	render: function() {
		return (
			<h1>quicklog log search</h1>
		);
	}
});

var SearchRow = React.createClass({
	render: function() {
		var id = this.props.row.id;

		var self = this;

		var fields = JSON.stringify(this.props.row.fields);

		return (
			<tr>
				<td>{this.props.row.fields.timestamp}</td>
				<td>{this.props.row.fields.msg}</td>
				<td><button className="btn btn-primary btn-small" onClick={this.showDetails}>Details</button>
					<div className="modal" role="dialog" ref="modal">
						<div className="modal-dialog" role="document">
							<div className="modal-content">
								<div className="modal-header">
									<button type="button" className="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
									<h4 className="modal-title">Message Details</h4>
								</div>
								<div className="modal-body">
									<label>ID:</label>{id}<br/>
									{fields}
								</div>
							</div>
						</div>
					</div>
				</td>
			</tr>
		);
	},
	showDetails: function() {
		$(this.refs.modal).modal('show');
	}
})

var SearchResults = React.createClass({
	render: function() {
		var rows = this.props.data.map(function(entry) {
			return (
				<SearchRow row={entry}/>
			);
		});

		return (
			<table className="table table-striped table-bordered">
				<thead>
					<tr>
						<th>Timestamp</th>
						<th>Message</th>
						<th></th>
					</tr>
				</thead>
				<tbody>
					{rows}
				</tbody>
			</table>
		);
	}
})

var SearchForm = React.createClass({
	handleSubmit: function(e) {
		e.preventDefault();
		var query = this.refs.query.value.trim();
		if (!query || query == "") {
			return;
		}

		this.props.onSearch(query);
		return;
	},
	render: function() {
		return (
			<form className="form-inline" onSubmit={this.handleSubmit}>
				<div className="form-group">
					<div className="input-group">
						<input placeholder="query" className="form-control" ref="query" type="text"/>
					</div>
				</div>
				<button type="submit" className="btn btn-primary">Search</button>
			</form>
		);
	}
});

var SearchApplication = React.createClass({
	getInitialState: function() {
		return {data: []};
	},
	search: function(query) {
		$.ajax({
			url: "/search",
			dataType: "json",
			cache: false,
			method: "POST",
			contentType : 'application/json',
			data: JSON.stringify({
				"size": 10,
				"fields": [ "*" ],
				"query": {
					"query": query,
				},
				"explain": false,
				"hightlight": {
					"fields": [ "msg", "timestamp" ]
				}
			}),
			success: function(data) {
				this.setState({data: data.hits});
			}.bind(this)
		})
	},
	render: function() {
		return (
			<div>
				<div className="row">
					<div className="col-md-12">
						<Header/>
					</div>
				</div>

				<div className="row">
					<div className="col-md-12">
						<SearchForm onSearch={this.search}/>
					</div>
				</div>
				
				<div className="row">
					<div className="col-md-12">
						<SearchResults data={this.state.data} />
					</div>
				</div>
			</div>
		);
	}
});

ReactDOM.render(
	<SearchApplication />,
	document.getElementById('app')
);

