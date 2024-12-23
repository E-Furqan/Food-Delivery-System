package model

type DatabaseEnv struct {
	DATABASE_HOST     string
	DATABASE_USER     string
	DATABASE_PASSWORD string
	DATABASE_NAME     string
	DATABASE_PORT     int
}

type WorkflowEnv struct {
	BASE_URL                  string
	PORT                      string
	DATAPIPELINE_WORKFLOW_URL string
}
