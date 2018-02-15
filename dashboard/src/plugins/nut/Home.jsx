import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'

import SignIn from '../../plugins/nut/users/SignIn'
import Logs from '../../plugins/nut/users/Logs'

class Widget extends Component {
  render() {
    const {user} = this.props
    return user.uid
      ? <Logs/>
      : <SignIn/>
  }
}

Widget.propTypes = {
  push: PropTypes.func.isRequired,
  user: PropTypes.object.isRequired
}

export default connect(state => ({user: state.currentUser}), {push})(Widget)
