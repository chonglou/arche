import Link from 'next/link'
import Head from 'next/head'

export default({
  children,
  title = 'application'
}) => (<div>
  <Head>
    <meta charSet="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no"/>
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