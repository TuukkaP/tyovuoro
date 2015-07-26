var Menu = React.createClass({
  render: function () {
    return (
          <ul className="side-nav">
            <li><a href="#">Home</a></li>
            <li><Link to="users" /></li>
          </ul>
    )
  }
});