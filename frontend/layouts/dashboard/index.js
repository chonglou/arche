import Link from 'next/link'
import Head from 'next/head'
import antDesign from 'antd/dist/antd.css'
import antDesignPro from 'ant-design-pro/dist/ant-design-pro.css'

export default({
  children,
  title = 'dashboard'
}) => (<div>
  <Head>
    <meta charSet='utf-8'/>
    <meta name='viewport' content='initial-scale=1.0, width=device-width'/>
    <title>{title}</title>
  </Head>
  <style jsx="jsx" global="global">
    {
      antDesign
    }</style>
  <style jsx="jsx" global="global">
    {
      antDesignPro
    }</style>
  <header>
    head
  </header>

  {children}

  <footer>
    foot
  </footer>
</div>)