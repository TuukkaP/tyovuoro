var UserBox = React.createClass({
    url: "/api/users",
    fetchUsersFromBackend: function() {
      console.log("Load users")
      $.ajax({
        url: this.url,
        dataType: 'json',
        cache: false,
        success: function(users) {
          this.setState({users: users});
        }.bind(this),
        error: function(xhr, status, err) {
          console.error(url, status, err.toString());
        }.bind(this)
      });
    },
    getInitialState: function() {
      return {
        users: []
      };
    },
    componentDidMount: function() {
      this.fetchUsersFromBackend();
    },
    render: function() {
      return (
        <div className="userBox">
          <h1>Users</h1>
          <UserList users={this.state.users}/>
        </div>
      );
    }
});

var UserList = React.createClass({
  sortByKey: function (array, key) {
      return array.sort(function(a, b) {
        var x = a[key].toLowerCase(); var y = b[key].toLowerCase();
        return ((x < y) ? -1 : ((x > y) ? 1 : 0));
      });
  },
  componentWillReceiveProps: function(nextProps) {
    this.setState({users: this.sortByKey(nextProps.users, "username")});
  },
  childChanged: function (index, updatedUser) {
    console.log(this.state, index, updatedUser);
    var newUsers = this.state.users;
    newUsers[index] = updatedUser;
    this.setState({users: this.sortByKey(newUsers, "username")});
  },
	render: function() {
    // Is this necessary to pass callback to children??
    var childChanged = this.childChanged;
    return (
	  <Table striped bordered condensed hover responsive>
	    <thead>
	      <tr>
	        <th>Username</th>
	        <th>First Name</th>
	        <th>Last Name</th>
          <th></th>
	      </tr>
	    </thead>
	    <tbody>
        	{this.props.users.map(function (user, index) {
            return <User user={user} key={user.id} index={index} updateUserList = {childChanged} />
          })}
          <UserInput />
      </tbody>
      <tfooter>
      </tfooter>
    </Table>
    );
  }
});

var User = React.createClass({
  render: function() {
    return (
      <tr>
        <td>{this.props.user.username}</td>
        <td>{this.props.user.firstname}</td>
        <td>{this.props.user.lastname}</td>
        <td>
          <ModalTrigger modal={<UserEditModal user={this.props.user} key={this.props.user.id} index={this.props.index} updateUserList={this.props.updateUserList}/>}>
          <Button bsStyle="primary" bsSize="small">Edit</Button>
          </ModalTrigger>
        </td>
      </tr>
    );
  }
});

var UserEditModal = React.createClass({
    getInitialState: function () {
      return {
        username: this.props.user.username,
        firstname: this.props.user.firstname,
        lastname: this.props.user.lastname
      };
    },

  validateUsername: function () {
    "use strict";
    let length = this.state.username.length;
      if (length > 4) {
        return 'success'; 
      } else if (length > 0) {
        return 'error';
      }
  },

  handleChange: function () {
    this.setState({
      key: this.refs.idInput.getValue(),
      username: this.refs.usernameInput.getValue(),
      firstname: this.refs.firstnameInput.getValue(),
      lastname: this.refs.lastnameInput.getValue()
    });
  },
  
  handleUserEditSubmit: function (e) {
    e.preventDefault();
    var id = this.refs.idInput.getValue();
    var username = this.refs.usernameInput.getValue();
    var firstname = this.refs.firstnameInput.getValue();
    var lastname = this.refs.lastnameInput.getValue();
    if (!username || !id) {
      console.error(id, username);
      return;
    }
    this.sendData({id: parseInt(id, 10), username: username, firstname: firstname, lastname: lastname});
    console.log("HandleUserEditSubmit done!");
    this.props.onRequestHide();
    return;
  },

  sendData: function (user) {
    $.ajax({
      url: "/api/users/" + user.id,
      dataType: "json",
      method: "PUT",
      data: JSON.stringify(user),
      contentType: "application/json",
      beforeSend: function(xhr) {
        console.log(this);
      },
      success: function() {
        console.log(user, "success", this.props);
        this.props.updateUserList(this.props.index, user);
/*        this.setState({user: user});*/
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(status, err.toString());
      }.bind(this)
    });
  },

  render: function () {
    console.log(this.props);
      return (
        <Modal {...this.props} title={"Edit user:" + this.state.username} animation={true}>
          <form onSubmit={this.handleUserEditSubmit}>
            <div className='modal-body'>
                  <Input type="text"
                      value={this.state.username}
                      placeholder="Username" 
                      onChange={this.handleChange} 
                      bsStyle={this.validateUsername()}
                      ref="usernameInput"
                      hasFeedback />
              <Input type="hidden" value={this.props.user.id} ref="idInput" />
              <Input type="text" value={this.state.firstname} ref="firstnameInput" onChange={this.handleChange} placeholder="First name" />
              <Input type="text" value={this.state.lastname} ref="lastnameInput" onChange={this.handleChange} placeholder="Last name" />
            </div>
            <div className='modal-footer'>
              <ButtonInput type="submit" bsStyle="primary" bsSize="small" value="Edit" ></ButtonInput>
              <Button onClick={this.props.onRequestHide}>Close</Button>
            </div>
          </form>
        </Modal>
      );
    }
});


var UserInput = React.createClass({
    getInitialState: function () {
      return {
        username: ""
      };
    },

  validateUsername: function () {
    "use strict";
    let length = this.state.username.length;
    if (length > 4) {
      return 'success'; 
    } else if (length > 0) {
      return 'error';
    }
  },

  handleChange: function () {
    this.setState({
      username: this.refs.usernameInput.getValue()
    });
  },

  render: function () {
      return (
        <tr>
          <td>
            <Input type="text" 
                value={this.state.username}
                placeholder="Username" 
                onChange={this.handleChange} 
                bsStyle={this.validateUsername()}
                ref="usernameInput"
                hasFeedback />
          </td>
          <td><Input type="text" placeholder="First name" /></td>
          <td><Input type="text" placeholder="Last name" /></td>
          <td><Button bsStyle="primary" bsSize="small">Create</Button></td>
        </tr>
      );
    }
});

/*React.render(
  <UserBox url = "/api/users" />,
  document.getElementById('main')
);

var data = 
[{
"id": 14,
"username": "Data",
"password": "006525ca1a86547f21f7c196ba0b6fafa0747ba161839739c46a7c439d341a67",
"role": "user",
"address": "USS Enterprise",
"firstname": "Mr.",
"lastname": "Data",
"email": "noonien@forever"
},
{
"id": 15,
"username": "Tom",
"password": "5c34f5d261d092fd0216643301cf12545b930a4cd688b7aec215d08c82401b75",
"role": "user",
"address": "Azkaban c/o ankeuttajat",
"firstname": "Tom Lomen",
"lastname": "Valedro",
"email": "avadakedavra.hgw"
},];*/