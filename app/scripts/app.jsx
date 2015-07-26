// PROJECT VARS

var Router = window.ReactRouter;
var RouteHandler = Router.RouteHandler;
var Link = Router.Link;
var Route = Router.Route;
var DefaultRoute = Router.DefaultRoute;
var Grid = ReactBootstrap.Grid
	, Row = ReactBootstrap.Row
	, Col = ReactBootstrap.Col
	, Nav = ReactBootstrap.Nav
	, NavItem = ReactBootstrap.NavItem
	, Table = ReactBootstrap.Table
	, Button = ReactBootstrap.Button
	, Input = ReactBootstrap.Input
	, ButtonInput = ReactBootstrap.ButtonInput
	, ModalTrigger = ReactBootstrap.ModalTrigger
	, Modal = ReactBootstrap.Modal;
var NavItemLink = ReactRouterBootstrap.NavItemLink

var App = React.createClass({
  render: function () {
    return (
		  <Grid>
		    <Row className='show-grid'>
		      <Col xs={4} md={2}>        
				    <Nav bsStyle='pills' stacked>
		          <NavItemLink to="index" eventKey={1}>Homepage</NavItemLink>
		          <NavItemLink to="users" eventKey={2}>Users!</NavItemLink>
		        </Nav>
        	</Col>
		      <Col xs={12} md={8}>
		        <div className="content">
		          <RouteHandler/>
		        </div>
		      </Col>
		   	</Row>
	    <Footer/>
	    </Grid>
      )}
});