import Link from 'next/link'
import Head from 'next/head'

import Cdn from '../cdn'

export default({
  children,
  title = 'dashboard'
}) => (<div>
  <Head>
    <meta charSet='utf-8'/>
    <meta name='viewport' content='initial-scale=1.0, width=device-width'/>
    <Cdn name="antd" version="3.2.2" file="antd.min.css"/>
    <Cdn name="ant-design-pro" version="1.1.0" file="ant-design-pro.min.css"/>
    <Cdn name="quill" version="1.2.6" file="quill.snow.css"/>
    <title>{title}</title>
  </Head>
  <header>
    head
  </header>

  {children}

  <footer>
    foot
  </footer>
</div>)