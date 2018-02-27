import React, {Component} from 'react'
import PropTypes from 'prop-types'

class Widget extends Component {
  render() {
    const {name, version, file} = this.props
    return (<link rel="stylesheet" href={`https://unpkg.com/${name}@${version}/dist/${file}`} crossOrigin="anonymous"/>)
  }
}

Widget.propTypes = {
  file: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  version: PropTypes.string.isRequired
}

export default Widget
