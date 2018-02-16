import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Layout, Menu} from 'antd';
import {connect} from 'react-redux'
import {Link} from 'react-router-dom'
import {FormattedMessage} from 'react-intl'

const {Header} = Layout;

class Widget extends Component {
  render() {
    return (<Header>
      <div className="logo"/>
      <Menu theme="dark" mode="horizontal" defaultSelectedKeys={[]} style={{
          lineHeight: '64px'
        }}>
        <Menu.Item key="personal">
          <Link to={'/users/sign-in'}>
            <FormattedMessage id={'nut.users.sign-in.title'}/>
          </Link>
        </Menu.Item>
      </Menu>
    </Header>)
  }
}

Widget.propTypes = {
  user: PropTypes.object.isRequired,
  site: PropTypes.object.isRequired
}

export default connect(state => ({user: state.currentUser, site: state.siteInfo}))(Widget)
