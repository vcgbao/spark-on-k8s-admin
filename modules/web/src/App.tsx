import React from "react"
import Sidebar from './components/Sidebar'
import { Route, Routes, Navigate } from 'react-router-dom';
import SparkApplication from './pages/SparkApplication'
import ScheduledSparkApplication from './pages/ScheduledSparkApplication'

import { Layout } from 'antd';
const { Content } = Layout;

const App = () => {
	return (
		<Layout hasSider>
			<Sidebar />
			<Layout style={{ marginLeft: 200, height: '100vh' }}>
				<Content style={{ height: '100vh' }}>
					<Routes>
						<Route path="/sparkapplication" element={SparkApplication()} />
						<Route path="/scheduledsparkapplication" element={ScheduledSparkApplication()} />
						<Route path="/" element={<Navigate to="/sparkapplication" />} />
					</Routes>
				</Content>
			</Layout>

		</Layout >
	);
}

export default App;
