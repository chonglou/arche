import Layout from '../layouts/application'

export default() => (<Layout>
  <div>
    <hr/>
    home {process.env.BACKEND}
    <img src='/static/fail.png'/>
  </div>
</Layout>)