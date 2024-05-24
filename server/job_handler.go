/*
* Copyright 2022-2024 Thorsten A. Knieling
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
 */

package server

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
	"github.com/tknie/clu/api"
)

// GetJobExecutionResult implements getJobExecutionResult operation.
//
// Retrieves a specific job result.
//
// GET /tasks/results
func (Handler) GetJobExecutionResult(ctx context.Context, params api.GetJobExecutionResultParams) (r api.GetJobExecutionResultRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetJobFullInfo implements getJobFullInfo operation.
//
// Retrieves a full job definition.
//
// GET /tasks/{jobName}/full
func (Handler) GetJobFullInfo(ctx context.Context, params api.GetJobFullInfoParams) (r api.GetJobFullInfoRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetJobs implements getJobs operation.
//
// Retrieves a list of jobs known by the Interface.
//
// GET /tasks
func (Handler) GetJobs(ctx context.Context, params api.GetJobsParams) (r api.GetJobsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetJobsConfig implements getJobsConfig operation.
//
// Read job configuration section.
//
// GET /admin/config/jobs
func (Handler) GetJobsConfig(ctx context.Context) (r api.GetJobsConfigRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteJobResult implements deleteJobResult operation.
//
// Delete a specific job result.
//
// DELETE /tasks/{jobName}/result/{jobId}
func (Handler) DeleteJobResult(ctx context.Context, params api.DeleteJobResultParams) (r api.DeleteJobResultRes, _ error) {
	return r, ht.ErrNotImplemented
}

// PostJob implements postJob operation.
//
// Create a new Job database.
//
// POST /tasks
func (Handler) PostJob(ctx context.Context, req api.PostJobReq) (r api.PostJobRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetJobResult implements getJobResult operation.
//
// Delete a specific job result.
//
// GET /rest/tasks/{jobName}/{jobId}
func (Handler) GetJobResult(ctx context.Context, params api.GetJobResultParams) (r api.GetJobResultRes, _ error) {
	return r, ht.ErrNotImplemented
}

// TriggerJob implements triggerJob operation.
//
// Trigger a job.
//
// PUT /rest/tasks/{jobName}
func (Handler) TriggerJob(ctx context.Context, params api.TriggerJobParams) (r api.TriggerJobRes, _ error) {
	return r, ht.ErrNotImplemented
}
