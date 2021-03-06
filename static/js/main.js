
var Header = React.createClass({
	render: function() {
		return (
			<nav className="navbar navbar-default navbar-fixed-top">
				<div className="container-fluid">
				    <div className="navbar-header">
					  <a className="navbar-brand" href="#">quicklog</a>
					</div>
				</div>
			</nav>
		);
	}
});

var Time = React.createClass({
	render: function() {
		var val = moment(this.props.val).format("YYYY-MM-DD hh:mm A");
		return (<span>{val}</span>);
	}
});

var SearchRow = React.createClass({
	render: function() {
		var id = this.props.row.id;

		var self = this;

		var fields = JSON.stringify(this.props.row.fields);

		return (
			<tr>
				<td className="col-xs-4 col-md-3"><Time val={this.props.row.fields.timestamp} /></td>
				<td className="col-xs-6 col-md-7">{this.props.row.fields.message}</td>
				<td className="col-xs-2 col-md-2"><button className="btn btn-primary btn-small pull-right" onClick={this.showDetails}>Details</button>
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
		var rows = this.props.data.hits.map(function(entry) {
			return (
				<SearchRow key={entry.id} row={entry}/>
			);
		});

		var pages = [];

		var total = this.props.data.total_hits;
		if(total != 0){

			var pageSize = this.props.data.request.size;

			for(var i = 0; i < total/pageSize; i++){
				var activeClass="";
				if (this.props.page == i) {
					activeClass = "active";
				}
				pages.push(
					<li className={activeClass}><a href="#" onClick={this.props.setPage.bind(this, i, pageSize)}>{i}</a></li>
				);
			}
		}

		return (
			<div>
				<div className="row">
					<div className="col-md-6">
						<SearchForm onSearch={this.props.onSearch} />
					</div>
					<div className="col-md-6">
						<ul className="pagination pull-right">
							{pages}
						</ul>
					</div>
				</div>
				<table className="center table table-striped table-condensed table-hover">
					<thead>
						<tr>
							<th className="col-xs-2">Timestamp</th>
							<th>Message</th>
							<th></th>
						</tr>
					</thead>
					<tbody>
						{rows}
					</tbody>
				</table>
			</div>
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
				<div className="input-group">
					<input placeholder="Search Query" className="form-control" ref="query" type="text"/>
					<span className="input-group-btn">
						<button type="submit" className="btn btn-primary">Search</button>
					</span>
				</div>
			</form>
		);
	}
});

var SearchApplication = React.createClass({
	getInitialState: function() {
		return {query: "", from: 0, data: { hits: [], total_hits: 0, request: { from: 0, size: 0 } } };
	},
	search: function(query) {

		if (query == "" || query == "*") {
			$.ajax({
				url: "/search",
				dataType: "json",
				cache: false,
				method: "POST",
				contentType : 'application/json',
				data: JSON.stringify({
					"size": 10,
					"from": this.state.from
				}),
				success: function(data) {
					this.setState({query: query, data: data, from: this.state.from});
				}.bind(this)
			})
		} else {
			$.ajax({
				url: "/search",
				dataType: "json",
				cache: false,
				method: "POST",
				contentType : 'application/json',
				data: JSON.stringify({
					"size": 10,
					"query": query,
					"from": this.state.from
				}),
				success: function(data) {
					this.setState({query: query, data: data, from: this.state.from});
				}.bind(this)
			})
		}
	},
	setPage: function(page, size) {
		var st = this.state;
		st.from = page * size;
		st.page = page;
		this.search(st.query);
	},
	render: function() {
		return (
			<div>
			<Header/>

			<div id="wrapper">
				<div className="sidebar" id="sidebar-wrapper">
					<ul className="sidebar-nav">
						<li>
							<a href="#">Live</a>
						</li>
						<li>
							<a href="#">Search</a>
						</li>
						<li>
							<a href="#">Settings</a>
						</li>
					</ul>
				</div>

				<div id="page-content-wrapper">
					<div class="container-fluid">
						<SearchResults page={this.state.page} data={this.state.data} setPage={this.setPage} onSearch={this.search} />
					</div>
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

