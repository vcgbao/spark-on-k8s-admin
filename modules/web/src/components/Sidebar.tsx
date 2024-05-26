import {
	ScheduleOutlined,
	AppstoreAddOutlined,
} from '@ant-design/icons';
import { Menu, Layout, Image, Space } from 'antd';
import { useNavigate, useLocation } from "react-router-dom";
import sparkLogo from '../assets/sparkicon.svg'
const { Sider } = Layout;


const Sidebar = () => {
	const navigate = useNavigate()
	const location = useLocation()

	return (
		<Sider
			style={{ overflow: 'auto', height: '100vh', position: 'fixed', left: 0, top: 0, bottom: 0 }}
		>
			<Space direction="horizontal" style={{ width: '100%', justifyContent: 'center', paddingTop: '20px', paddingBottom: '40px' }}>
				<Image src={sparkLogo} width={120} preview={false} />
			</Space>
			<Menu
				onClick={() => { }}
				defaultSelectedKeys={['SparkApplication']}
				defaultOpenKeys={['SparkApplication']}
				selectedKeys={location.pathname === '/sparkapplication' ? ['SparkApplication'] : ['ScheduledSparkApplication']}
				mode="inline"
				items={[
					{
						key: 'SparkApplication',
						label: 'Spark App',
						icon: <AppstoreAddOutlined />,
						onClick: () => navigate("/sparkapplication")
					},
					{
						key: 'ScheduledSparkApplication',
						label: 'Scheduled App',
						icon: <ScheduleOutlined />,
						onClick: () => navigate("/scheduledsparkapplication")
					}
				]
				}
				theme='dark'
			/>
		</Sider>
	)
}

export default Sidebar;