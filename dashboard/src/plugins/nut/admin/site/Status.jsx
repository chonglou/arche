import React, {Component} from 'react'
import {Row, Col, Collapse, Table, message} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import PropTypes from 'prop-types'
import SyntaxHighlighter from 'react-syntax-highlighter'
import {docco} from 'react-syntax-highlighter/styles/hljs'

import Layout from '../../../../layouts/dashboard'
import {get} from '../../../../ajax'
import {ADMIN} from '../../../../auth'

const Panel = Collapse.Panel

const Hash = ({item}) => (<Table rowKey="key" dataSource={Object.entries(item).map((v) => {
    return {key: v[0], val: v[1]}
  })} columns={[
    {
      title: <FormattedMessage id="attributes.key"/>,
      key: 'key',
      dataIndex: 'key'
    }, {
      title: <FormattedMessage id="attributes.value"/>,
      key: 'val',
      dataIndex: 'val'
    }
  ]}/>)

Hash.propTypes = {
  item: PropTypes.object.isRequired
}

class Widget extends Component {
  state = {
    os: {},
    database: {},
    cache: "",
    jobber: {
      tasks: []
    },
    network: {}
  }
  componentDidMount() {
    get('/api/admin/site/status').then((rst) => {
      this.setState(rst)
    }).catch(message.error);
  }
  render() {
    const {cache, os, network, database, jobber} = this.state

    const title = {
      id: "nut.admin.site.status.title"
    }
    return (<Layout breads={[{
          href: "/admin/site/status",
          label: title
        }
      ]} title={title} roles={[ADMIN]}>
      <Row>
        <Col md={{
            span: 16,
            offset: 2
          }}>
          <Collapse>
            <Panel key="os" header={(<FormattedMessage id="nut.admin.site.status.os"/>)}>
              <Hash item={os}/>
            </Panel>
            <Panel key="network" header={(<FormattedMessage id="nut.admin.site.status.network"/>)}>
              <Hash item={network}/>
            </Panel>
            <Panel key="database" header={(<FormattedMessage id="nut.admin.site.status.database"/>)}>
              <Hash item={database}/>
            </Panel>
            <Panel key="jobber" header={(<FormattedMessage id="nut.admin.site.status.jobber"/>)}>
              <SyntaxHighlighter language="yaml" style={docco}>{jobber.config}</SyntaxHighlighter>
              <ul>
                {jobber.tasks.map((t, i) => (<li key={i}>{t}</li>))}
              </ul>
            </Panel>
            <Panel key="redis" header={(<FormattedMessage id="nut.admin.site.status.cache"/>)}>
              <SyntaxHighlighter style={docco}>{cache}</SyntaxHighlighter>
            </Panel>
          </Collapse>
        </Col>
      </Row>
    </Layout>);
  }
}
Widget.propTypes = {
  intl: intlShape.isRequired
}

const WidgetI = injectIntl(Widget)

export default WidgetI
