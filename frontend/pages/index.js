import Layout from '../layouts/application'

export default() => (<Layout>
  <div>
    <hr/>
    home {process.env.BACKEND}
  </div>
</Layout>)