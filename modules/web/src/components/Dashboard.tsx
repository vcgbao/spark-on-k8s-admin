

import { SparkApplicationResponse } from '../apis'
import { SIToBytes, byteToGi } from '../utils/converter'

import { Table, Spin, Tag, Card, Col, Row, Statistic } from 'antd';



const Dashboard = ({ sparkApp }: { sparkApp: SparkApplicationResponse }) => {


  return (
    <Row justify='center' style={{ paddingTop: '20px', paddingBottom: '20px' }}>
      <Col span={7} style={{ marginLeft: '10px', marginRight: '10px' }}>
        <Card bordered={true}>
          <Statistic
            title="Running Spark Application"
            value={sparkApp.items.filter((item) => item.status.applicationState.state === "RUNNING").length}
            valueStyle={{ color: '#3f8600' }}
          />
        </Card>
      </Col>
      <Col span={7} style={{ marginLeft: '10px', marginRight: '10px' }}>
        <Card bordered={true}>
          <Statistic
            title="Total Cores"
            value={sparkApp.items
              .filter((item) => item.status.applicationState.state === "RUNNING")
              .reduce((sum, item) => sum + item.spec.driver.cores + item.spec.executor.instances * item.spec.executor.cores, 0)}
            valueStyle={{ color: '#3f8600' }}
          />
        </Card>
      </Col>
      <Col span={7} style={{ marginLeft: '10px', marginRight: '10px' }}>
        <Card bordered={false}>
          <Statistic
            title="Total Memory (Gi)"
            value={sparkApp.items
              .filter((item) => item.status.applicationState.state === "RUNNING")
              .reduce((sum, item) => sum + byteToGi(SIToBytes(item.spec.driver.memory) + item.spec.executor.instances * SIToBytes(item.spec.executor.memory)), 0)}
            valueStyle={{ color: '#3f8600' }}
          />
        </Card>
      </Col>
    </Row>
  )
}

export default Dashboard