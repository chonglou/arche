import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {
  Form,
  Row,
  Col,
  Input,
  Select,
  message
} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'

import Layout from '../../../../layouts/dashboard'
import {post, get} from '../../../../ajax'
import {ADMIN} from '../../../../auth'
import {Submit, formItemLayout} from '../../../../components/form'

const FormItem = Form.Item
const Option = Select.Option

class Widget extends Component {
  state = {
    options: []
  }
  componentDidMount() {
    const {setFieldsValue} = this.props.form
    get('/admin/site/home').then((rst) => {
      setFieldsValue({favicon: rst.favicon, home: rst.home})
      this.setState({options: rst.options})
    }).catch(message.error)
  }
  handleSubmit = (e) => {
    const {formatMessage} = this.props.intl
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (!err) {
        post('/admin/site/home', values).then(() => {
          message.success(formatMessage({id: "flash.success"}))
        }).catch(message.error);
      }
    });
  }
  render() {
    const {formatMessage} = this.props.intl
    const {getFieldDecorator} = this.props.form
    const title = {
      id: "nut.admin.site.home.title"
    }
    return (<Layout breads={[{
          href: "/admin/site/home",
          label: title
        }
      ]} title={title} roles={[ADMIN]}>
      <Row>
        <Col md={{
            span: 16,
            offset: 2
          }}>
          <Form onSubmit={this.handleSubmit}>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "nut.admin.site.home.favicon" />} hasFeedback={true}>
              {
                getFieldDecorator('favicon', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "validator.required"})
                    }
                  ]
                })(<Input/>)
              }
            </FormItem>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.body" />}>
              {
                getFieldDecorator('home')(<Select>
                  {this.state.options.map((it) => (<Option key={it} value={it}>{it}</Option>))}
                </Select>)
              }
            </FormItem>
            <Submit/>
          </Form>
        </Col>
      </Row>
    </Layout>);
  }
}

Widget.propTypes = {
  intl: intlShape.isRequired,
  push: PropTypes.func.isRequired
}

const WidgetF = Form.create()(injectIntl(Widget))

export default connect(state => ({}), {
  push
},)(WidgetF)
