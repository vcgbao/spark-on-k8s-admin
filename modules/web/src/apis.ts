


const origin = window.location.origin

export async function getSparkApplications() {
	const response = await fetch(
		`${origin}/api/sparkapplication`,
	)
	return await response.json() as Response<SparkApplicationResponse>
}

export async function deleteSparkApplication(appName: string) {
	const response = await fetch(
		`${origin}/api/sparkapplication/${appName}`,
		{ method: 'DELETE' }
	)
	return await response.json() as Response<any>
}

export async function createSparkApplication(spec: any) {
	const response = await fetch(
		`${origin}/api/sparkapplication/create`,
		{
			method: 'POST',
			body: JSON.stringify(spec)
		}
	)

	return await response.json() as Response<any>
}

export async function updateSparkApplication(spec: any) {
	const response = await fetch(
		`${origin}/api/sparkapplication/update`,
		{
			method: 'POST',
			body: JSON.stringify(spec)
		}
	)

	return await response.json() as Response<any>
}


// Scheduled Spark Application

export async function getScheduledSparkApplications() {
	const response = await fetch(
		`${origin}/api/scheduledsparkapplication`,
	)
	return await response.json() as Response<ScheduledSparkApplicationResponse>
}

export async function deleteScheduledSparkApplication(appName: string) {
	const response = await fetch(
		`${origin}/api/scheduledsparkapplication/${appName}`,
		{ method: 'DELETE' }
	)
	return await response.json() as Response<any>
}

export async function createScheduledSparkApplication(spec: any) {
	const response = await fetch(
		`${origin}/api/scheduledsparkapplication/create`,
		{
			method: 'POST',
			body: JSON.stringify(spec)
		}
	)

	return await response.json() as Response<any>
}

export async function updateScheduledSparkApplication(spec: any) {
	const response = await fetch(
		`${origin}/api/scheduledsparkapplication/update`,
		{
			method: 'POST',
			body: JSON.stringify(spec)
		}
	)

	return await response.json() as Response<any>
}



export const SUCCESS = 0
export const ERROR = 1

export interface Response<T> {
	statusCode: number
	message: string
	data: T
}

export interface ScheduledSparkApplicationResponse {
	apiVersion: string
	items: ScheduledSparkApplication[]
	kind: string
}

export interface ScheduledSparkApplication {
	apiVersion: string
	kind: string
	metadata: Metadata
	spec: ScheduledSparkApplicationSpec
	status: ScheduledSparkApplicationStatus
}

export interface ScheduledSparkApplicationSpec {
	schedule: string
	template: Spec
	concurrencyPolicy: string
	successfulRunHistoryLimit: number
	failedRunHistoryLimit: number
}


export interface ScheduledSparkApplicationStatus {
	lastRun: string
	nextRun: string
	lastRunName: string
	scheduleState: string
	pastFailedRunNames?: string[]
	pastSuccessfulRunNames?: string[]
}


export interface SparkApplicationResponse {
	apiVersion: string
	items: SparkApplication[]
	kind: string
}

export interface SparkApplication {
	apiVersion: string
	kind: string
	metadata: Metadata
	spec: Spec
	status: Status
}

export interface Metadata {
	annotations: any
	creationTimestamp: string
	generation: number
	name: string
	namespace: string
	resourceVersion: string
	uid: string
}


export interface Spec {
	driver: Driver
	executor: Executor
	image: string
	imagePullPolicy: string
	mainApplicationFile: string
	mainClass: string
	mode: string
	restartPolicy: RestartPolicy
	sparkVersion: string
	type: string
	volumes: Volume[]
}

export interface Driver {
	coreLimit: string
	cores: number
	labels: Labels
	memory: string
	serviceAccount: string
}

export interface Labels {
	version: string
}

export interface Executor {
	cores: number
	instances: number
	labels: Labels2
	memory: string
}

export interface Labels2 {
	version: string
}

export interface RestartPolicy {
	type: string
}

export interface Volume {
	hostPath: HostPath
	name: string
}

export interface HostPath {
	path: string
	type: string
}

export interface Status {
	applicationState: ApplicationState
	driverInfo: DriverInfo
	lastSubmissionAttemptTime: string
	submissionAttempts: number
	terminationTime: any
}

export interface ApplicationState {
	errorMessage: string
	state: string
}

export interface DriverInfo { }
