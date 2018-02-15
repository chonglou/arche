import React from 'react'
import {Route} from 'react-router'
import Loadable from 'react-loadable';

import nut from './nut'

const plugins = [nut]

const dynamicWrapper = (w) => Loadable({
  loader: () => w,
  loading: () => <div>Loading...</div>
});

export default {
  menus: plugins.reduce((ar, it) => ar.concat(it.menus), []),
  routes: plugins.reduce((ar, it) => ar.concat(it.routes), []).map((it) => {
    return (<Route key={it.path} exact={true} path={it.path} component={dynamicWrapper(it.component)}/>)
  }).concat([<Route key="not-found" component={dynamicWrapper(import ('./NotFound'))}/>])
}
