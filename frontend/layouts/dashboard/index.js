import Link from 'next/link'
import Head from 'next/head'

export default({
  children,
  title = 'dashboard'
}) => (<div>
  <Head>
    <meta charSet='utf-8'/>
    <meta name='viewport' content='initial-scale=1.0, width=device-width'/>
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