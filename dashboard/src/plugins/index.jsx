import React from 'react'
import {Route} from 'react-router'

import nut from './nut'
import donate from './donate'
import NotFound from './NotFound'

const plugins = [nut, donate]

export default {
  menus: plugins.reduce((ar, it) => ar.concat(it.menus), []),
  routes: plugins.reduce((ar, it) => ar.concat(it.routes), []).map((it) => {
    return (< Route key = {
      it.path
    }
    exact = {
      true
    }
    path = {
      it.path
    }
    component = {
      it.component
    } />)
  }).concat([<Route key="not-found" component={NotFound}/>])
}
