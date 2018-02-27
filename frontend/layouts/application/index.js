import Link from 'next/link'
import Head from 'next/head'

import Cdn from '../cdn'

export default({
  children,
  title = 'application'
}) => (<div>
  <Head>
    <meta charSet="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no"/>
    <Cdn name="bootstrap" version="4.0.0" file="css/bootstrap.min.css"/>
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