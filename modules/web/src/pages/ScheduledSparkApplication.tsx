import React, { useState, useEffect } from "react"
import { ScheduledSparkApplicationResponse, getScheduledSparkApplications, ScheduledSparkApplication, deleteScheduledSparkApplication, SUCCESS, updateScheduledSparkApplication, createScheduledSparkApplication } from '../apis'
import { SIToBytes, byteToGi } from '../utils/converter'
import AceEditor from "react-ace";

import "ace-builds/src-noconflict/mode-yaml";
import "ace-builds/src-noconflict/theme-github";
import "ace-builds/src-noconflict/ext-language_tools";


import {
	EllipsisOutlined,
	PlusOutlined
} from '@ant-design/icons';
import Yaml from 'yaml'
import { Table, Spin, Tag, Dropdown, Space, Modal, message, FloatButton } from 'antd';
import type { TableProps, MenuProps } from 'antd';
import { MessageInstance } from "antd/es/message/interface";

const dropdownItems = (
	setDeleteSparkApp: React.Dispatch<React.SetStateAction<string | null>>,
	sparkApp: ScheduledSparkApplication,
	setJobDetail: React.Dispatch<React.SetStateAction<ScheduledSparkApplication | null>>,
	setShowJobEditor: React.Dispatch<React.SetStateAction<EditorMode>>,
	setEditorValue: React.Dispatch<React.SetStateAction<any>>,
	messageApi: MessageInstance
) => {
	const i: MenuProps['items'] = [
		{
			label: <a onClick={() => setJobDetail(sparkApp)} >Job Details</a>,
			key: '1',
		},
		{
			type: 'divider',
		},
		{
			label: <a onClick={() => {
				if (sparkApp.metadata !== undefined && sparkApp.metadata.annotations != undefined && sparkApp.metadata.annotations["kubectl.kubernetes.io/last-applied-configuration"] !== undefined) {
					setEditorValue(toYaml(sparkApp.metadata.annotations["kubectl.kubernetes.io/last-applied-configuration"]))
					setShowJobEditor(EditorMode.Edit)
				} else {
					messageApi.open({
						type: 'warning',
						content: `Loading last applied configuration ${sparkApp.metadata.name} failed.\n The job is missing the "kubectl.kubernetes.io/last-applied-configuration" annotation, indicating it might not have been created with this tool or kubectl. It's also possible it was created by a scheduled Spark application. Content will be cleared for pasting your job definition.`,
						duration: 5
					});
					setEditorValue(null)
					setShowJobEditor(EditorMode.Edit)
				}


			}}>Edit</a>,
			key: '2',
			danger: true,
		},
		{
			label: <a onClick={() => setDeleteSparkApp(sparkApp.metadata.name)}>Delete</a>,
			key: '3',
			danger: true,
		},
	]
	return i

}

const getColumns = (
	setDeleteSparkApp: React.Dispatch<React.SetStateAction<string | null>>,
	setJobDetail: React.Dispatch<React.SetStateAction<ScheduledSparkApplication | null>>,
	setShowJobEditor: React.Dispatch<React.SetStateAction<EditorMode>>,
	setEditorValue: React.Dispatch<React.SetStateAction<any>>,
	messageApi: MessageInstance
) => {
	const cols: TableProps<ScheduledSparkApplication>['columns'] = [
		{
			title: 'Name',
			dataIndex: 'name',
			key: 'name',
			render: (_: any, app: ScheduledSparkApplication, index: number) => <a>{app.metadata.name}</a>,
		},
		{
			title: 'Schedule',
			dataIndex: 'schedule',
			key: 'schedule',
			render: (_: any, app: ScheduledSparkApplication, index: number) => <a>{app.spec.schedule}</a>,
		},
		{
			title: 'LastRun',
			key: 'last_run',
			render: (_: any, app: ScheduledSparkApplication, index: number) => (
				<Tag
					color={
						(app.status.pastSuccessfulRunNames || []).includes(app.status.lastRunName) ? 'green' :
							((app.status.pastFailedRunNames || []).includes(app.status.lastRunName) ? 'red' : 'grey')
					}
				>{app.status.lastRun}</Tag>
			)
		},
		{
			title: 'NextRun',
			key: 'next_run',
			render: (_: any, app: ScheduledSparkApplication, index: number) => app.status.nextRun,
		},
		{
			title: 'Diver Resources',
			dataIndex: 'driver_resources',
			key: 'driver_resources',
			render: (_: any, app: ScheduledSparkApplication, index: number) => app.spec.template.driver.cores + 'cores, ' + byteToGi(SIToBytes(app.spec.template.driver.memory)) + ' Gi'
		},
		{
			title: 'Num Executor',
			key: 'num_executor',
			dataIndex: 'num_executor',
			render: (_: any, app: ScheduledSparkApplication, index: number) => app.spec.template.executor.instances
		},
		{
			title: 'Executor Resources',
			dataIndex: 'executor_resources',
			key: 'executor_resources',
			render: (_: any, app: ScheduledSparkApplication, index: number) => app.spec.template.executor.cores + 'cores ' + byteToGi(SIToBytes(app.spec.template.executor.memory)) + ' Gi'
		},
		{
			title: 'Total Resources',
			dataIndex: 'total_resources',
			key: 'total_resources',
			render: (_: any, app: ScheduledSparkApplication, index: number) => app.spec.template.driver.cores + app.spec.template.executor.instances * app.spec.template.executor.cores + ' cores, ' + byteToGi(SIToBytes(app.spec.template.driver.memory) + (app.spec.template.executor.instances * SIToBytes(app.spec.template.executor.memory))) + ' Gi'
		},
		{
			title: '',
			render: (_: any, record: ScheduledSparkApplication, index: number) => {
				const items = dropdownItems(setDeleteSparkApp, record, setJobDetail, setShowJobEditor, setEditorValue, messageApi)
				return (
					<Dropdown menu={{ items }} trigger={['click']}>
						<Space style={{ paddingLeft: '8px', paddingRight: '8px', paddingTop: '2px', paddingBottom: '2px' }}>
							<EllipsisOutlined style={{ fontSize: 18 }} />
						</Space>
					</Dropdown>
				)
			}
		}
	]
	return cols
}


enum EditorMode {
	Hide = 0,
	New = 1,
	Edit = 2
}



const ScheduledSparkApplicationPage = () => {
	const [sparkApp, setSparkApp] = useState<ScheduledSparkApplicationResponse | null>();
	useEffect(() => {
		getScheduledSparkApplications().then((a) => {
			setSparkApp(a.data)
		})
	}, []);

	const [deleteSparkApp, setDeleteSparkApp] = useState<string | null>(null)
	const [jobDetail, setJobDetail] = useState<ScheduledSparkApplication | null>(null)
	const [messageApi, contextHolder] = message.useMessage()
	const [showJobEditor, setShowJobEditor] = useState<EditorMode>(EditorMode.Hide)
	const [editorValue, setEditorValue] = useState<any>()
	const [submittingJob, setSubmittingJob] = useState<boolean>(false)

	if (sparkApp == null) {
		return (
			<div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '20vh' }}>
				<Spin />
			</div>
		)
	}

	const columns = getColumns(setDeleteSparkApp, setJobDetail, setShowJobEditor, setEditorValue, messageApi)

	return (
		<>
			<Table columns={columns} dataSource={sparkApp.items} style={{ paddingLeft: '20px', paddingRight: '20px', marginTop: '40px' }} />

			<FloatButton
				icon={<PlusOutlined />}
				shape="circle"
				style={{ backgroundColor: '#3399ff', width: '50px', height: '50px' }}
				onClick={() => {
					setEditorValue(null)
					setShowJobEditor(EditorMode.New)
				}}
			/>

			<Modal
				title={`Do you want to delete ${deleteSparkApp}`}
				open={deleteSparkApp != null}
				onOk={async () => {
					const r = await deleteScheduledSparkApplication(deleteSparkApp as string)

					if (r.statusCode == SUCCESS) {
						messageApi.open({
							type: 'success',
							content: `Delete ${deleteSparkApp} successfully`,
							duration: 2
						});
						const newSparkApp = sparkApp
						newSparkApp.items = sparkApp.items.filter((a) => a.metadata.name != deleteSparkApp)
						setSparkApp(newSparkApp)
					} else {
						messageApi.open({
							type: 'error',
							content: `Error while delete ${deleteSparkApp}: ${r.message}`,
							duration: 5
						});
					}
					setDeleteSparkApp(null)
				}}
				onCancel={() => { setDeleteSparkApp(null) }}>
			</Modal>
			<Modal
				title={jobDetail?.metadata.name + ' details'}
				open={jobDetail != null}
				okButtonProps={{ hidden: true }}
				cancelButtonProps={{ hidden: true }}
				onCancel={() => setJobDetail(null)}
				width='50%'
			>
				<AceEditor
					width="100%"
					setOptions={{ useWorker: false }}
					mode="yaml"
					onChange={() => { }}
					name="job-detail-editor"
					editorProps={{ $blockScrolling: false }}
					value={toYaml(JSON.stringify(jobDetail))}
					readOnly={true}
				/>
			</Modal>


			<Modal
				title={'Scheduled Spark Application Editor'}
				open={showJobEditor != EditorMode.Hide}
				onCancel={() => {
					setShowJobEditor(EditorMode.Hide)
					setEditorValue(null)
				}}
				onOk={async () => {
					setSubmittingJob(true)
					const app = Yaml.parse(editorValue)
					const response = showJobEditor === EditorMode.New ? await createScheduledSparkApplication(app) : await updateScheduledSparkApplication(app)
					if (response.statusCode == SUCCESS) {
						messageApi.open({
							type: 'success',
							content: 'Job submit successfully',
							duration: 2
						});
					} else {
						messageApi.open({
							type: 'error',
							content: 'Error while submit job: ' + response.message,
							duration: 5
						})
					}
					setSubmittingJob(false)
					setShowJobEditor(EditorMode.Hide)
					setEditorValue(null)
				}}
				width='50%'
				okButtonProps={{
					loading: submittingJob,
					disabled: submittingJob
				}}
				cancelButtonProps={{
					disabled: submittingJob
				}}
			>
				<AceEditor
					width="100%"
					setOptions={{ useWorker: false }}
					mode="yaml"
					onChange={(v) => { setEditorValue(v) }}
					name="job-editor"
					editorProps={{ $blockScrolling: false }}
					value={editorValue || ""}
				/>
			</Modal>
			{contextHolder}
		</>
	)
}


function toYaml(input: string) {
	const doc = new Yaml.Document()
	doc.contents = JSON.parse(input)
	return doc.toString()
}

export default ScheduledSparkApplicationPage