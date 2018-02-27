import React, {Component} from 'react'
import withRedux from 'next-redux-wrapper'

import Layout from '../../layouts/dashboard'
import {initStore} from '../../store'

class Widget extends Component {
  static getInitialProps({store}) {
    store.dispatch(serverRenderClock(isServer))
    store.dispatch(addCount())

    return {isServer}
  }
  render() {
    return (<Layout title='sign in'>
      <div>
        <hr/>
        sign in
      </div>
    </Layout>)
  }
}

export default withRedux(initStore, null, {})(Widget)
