/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2e

import (
	"fmt"
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"testing"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/gitlab/impl"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/gitlab/tasks"
)

func TestGitlabMrNoteDataFlow(t *testing.T) {

	var gitlab impl.Gitlab
	dataflowTester := e2ehelper.NewDataFlowTester(t, "gitlab", gitlab)

	taskData := &tasks.GitlabTaskData{
		Options: &tasks.GitlabOptions{
			ConnectionId: 1,
			ProjectId:    12345678,
		},
	}
	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_gitlab_api_merge_requests_for_mr_notes_test.csv",
		"_raw_gitlab_api_merge_requests")

	// verify extraction
	dataflowTester.FlushTabler(&models.GitlabMergeRequest{})
	dataflowTester.FlushTabler(&models.GitlabMrLabel{})
	dataflowTester.Subtask(tasks.ExtractApiMergeRequestsMeta, taskData)
	dataflowTester.VerifyTable(
		models.GitlabMergeRequest{},
		fmt.Sprintf("./snapshot_tables/%s_for_mr_notes_test.csv", models.GitlabMergeRequest{}.TableName()),
		[]string{
			"connection_id",
			"gitlab_id",
			"iid",
			"project_id",
			"source_project_id",
			"target_project_id",
			"state",
			"title",
			"web_url",
			"user_notes_count",
			"work_in_progress",
			"source_branch",
			"target_branch",
			"merge_commit_sha",
			"merged_at",
			"gitlab_created_at",
			"closed_at",
			"merged_by_username",
			"description",
			"author_username",
			"author_user_id",
			"component",
			"first_comment_time",
			"review_rounds",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_gitlab_api_merge_request_notes.csv",
		"_raw_gitlab_api_merge_request_notes")

	// verify extraction
	dataflowTester.FlushTabler(&models.GitlabMrNote{})
	dataflowTester.FlushTabler(&models.GitlabMrComment{})
	dataflowTester.Subtask(tasks.ExtractApiMrNotesMeta, taskData)
	dataflowTester.VerifyTable(
		models.GitlabMrNote{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.GitlabMrNote{}.TableName()),
		[]string{
			"connection_id",
			"gitlab_id",
			"merge_request_id",
			"merge_request_iid",
			"noteable_type",
			"author_username",
			"body",
			"gitlab_created_at",
			"confidential",
			"resolvable",
			"is_system",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
	dataflowTester.VerifyTable(
		models.GitlabMrComment{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.GitlabMrComment{}.TableName()),
		[]string{
			"connection_id",
			"gitlab_id",
			"merge_request_id",
			"merge_request_iid",
			"body",
			"author_username",
			"author_user_id",
			"gitlab_created_at",
			"resolvable",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	// verify conversion
	dataflowTester.FlushTabler(&code.Note{})
	dataflowTester.Subtask(tasks.ConvertApiNotesMeta, taskData)
	dataflowTester.VerifyTable(
		code.Note{},
		fmt.Sprintf("./snapshot_tables/%s.csv", code.Note{}.TableName()),
		[]string{
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
			"id",
			"pr_id",
			"type",
			"author",
			"body",
			"resolvable",
			"is_system",
			"created_date",
		},
	)

	// verify conversion
	dataflowTester.FlushTabler(&code.PullRequestComment{})
	dataflowTester.Subtask(tasks.ConvertMrCommentMeta, taskData)
	dataflowTester.VerifyTable(
		code.PullRequestComment{},
		fmt.Sprintf("./snapshot_tables/%s.csv", code.PullRequestComment{}.TableName()),
		[]string{
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
			"id",
			"pull_request_id",
			"body",
			"user_id",
			"created_date",
			"commit_sha",
			"position",
		},
	)
}
