var Index = React.createClass({
  render: function () {
    return <h2>Home</h2>;
  }
});

// declare our routes and their hierarchy
var routes = (
  <Route handler={App}>
    <DefaultRoute handler={Index}/>
    <Route name="index" path="/" handler={Index}/>
    <Route name="users" path="/users" handler={UserBox}/>
  </Route>
);

Router.run(routes, function (Root) {
  React.render(<Root/>, document.body);
});