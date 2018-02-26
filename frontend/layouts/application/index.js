import Link from 'next/link'
import Head from 'next/head'
import bootstrap from 'bootstrap/dist/css/bootstrap.css'

export default({
  children,
  title = 'application'
}) => (<div>
  <Head>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no"/>
    <title>{title}</title>
  </Head>
  <style jsx="jsx" global="global">
    {
      bootstrap
    }</style>
  <header>
    head
  </header>

  {children}

  <footer>
    foot
  </footer>
</div>)