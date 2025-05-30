// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package querying demonstrates usages of the BigQuery query interface.
package querying

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/golang-samples/bigquery/snippets/bqtestutil"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestQueries(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	testDatasetID, err := bqtestutil.UniqueBQName("snippet_table_tests")
	if err != nil {
		t.Fatalf("couldn't generate unique resource name: %v", err)
	}

	meta := &bigquery.DatasetMetadata{
		Location: "US", // See https://cloud.google.com/bigquery/docs/locations
	}
	if err := client.Dataset(testDatasetID).Create(ctx, meta); err != nil {
		t.Fatalf("failed to create test dataset: %v", err)
	}
	// Cleanup dataset at end of test.
	defer client.Dataset(testDatasetID).DeleteWithContents(ctx)

	// Run query tests in parallel.
	t.Run("group", func(t *testing.T) {
		t.Run("queryBasic", func(t *testing.T) {
			t.Parallel()
			if err := queryBasic(io.Discard, tc.ProjectID); err != nil {
				t.Errorf("queryBasic: %v", err)
			}
		})
		t.Run("queryBatch", func(t *testing.T) {
			t.Parallel()
			tableID := "bigquery_query_batch"
			if err := queryBatch(io.Discard, tc.ProjectID, testDatasetID, tableID); err != nil {
				t.Errorf("queryBatch(%q): %v", testDatasetID, err)
			}
		})
		t.Run("queryDisableCache", func(t *testing.T) {
			t.Parallel()
			if err := queryDisableCache(io.Discard, tc.ProjectID); err != nil {
				t.Errorf("queryDisableCache: %v", err)
			}
		})
		t.Run("queryDryRun", func(t *testing.T) {
			t.Parallel()
			if err := queryDryRun(io.Discard, tc.ProjectID); err != nil {
				t.Errorf("queryDryRun: %v", err)
			}
		})
		t.Run("queryLegacy", func(t *testing.T) {
			t.Parallel()
			sql := "SELECT 17 as foo"
			if err := queryLegacy(io.Discard, tc.ProjectID, sql); err != nil {
				t.Errorf("queryLegacy: %v", err)
			}
		})
		t.Run("queryLegacyLargeResults", func(t *testing.T) {
			t.Parallel()
			tableID := "bigquery_query_legacy_large_results"
			if err := queryLegacyLargeResults(io.Discard, tc.ProjectID, testDatasetID, tableID); err != nil {
				t.Errorf("queryLegacyLargeResults: %v", err)
			}
		})
		t.Run("createTableAndWidenQuery", func(t *testing.T) {
			t.Parallel()
			tableID := "bigquery_add_column_query_append"
			if err := createTableAndWidenQuery(tc.ProjectID, testDatasetID, tableID); err != nil {
				t.Errorf("createTableAndWidenQuery: %v", err)
			}
		})
		t.Run("queryWithDestination", func(t *testing.T) {
			t.Parallel()
			tableID := "bigquery_query_destination_table"
			if err := queryWithDestination(io.Discard, tc.ProjectID, testDatasetID, tableID); err != nil {
				t.Errorf("queryWithDestination: %v", err)
			}
		})
		t.Run("queryWithDestinationCMEK", func(t *testing.T) {
			if bqtestutil.SkipCMEKTests() {
				t.Skip("skipping CMEK tests")
			}
			t.Parallel()
			tableID := "bigquery_query_destination_table_cmek"
			if err := queryWithDestinationCMEK(io.Discard, tc.ProjectID, testDatasetID, tableID); err != nil {
				t.Errorf("queryWithDestinationCMEK: %v", err)
			}
		})
		t.Run("queryWithArrayParams", func(t *testing.T) {
			t.Parallel()
			if err := queryWithArrayParams(io.Discard, tc.ProjectID); err != nil {
				t.Errorf("queryWithArrayParams: %v", err)
			}
		})
		t.Run("queryWithNamedParams", func(t *testing.T) {
			t.Parallel()
			if err := queryWithNamedParams(io.Discard, tc.ProjectID); err != nil {
				t.Errorf("queryWithNamedParams: %v", err)
			}
		})
		t.Run("queryWithPositionalParams", func(t *testing.T) {
			t.Parallel()
			if err := queryWithPositionalParams(io.Discard, tc.ProjectID); err != nil {
				t.Errorf("queryWithPositionalParams: %v", err)
			}
		})
		t.Run("queryWithStructParam", func(t *testing.T) {
			t.Parallel()
			if err := queryWithStructParam(io.Discard, tc.ProjectID); err != nil {
				t.Errorf("queryWithStructParam: %v", err)
			}
		})
		t.Run("queryWithTimestampParam", func(t *testing.T) {
			t.Parallel()
			if err := queryWithTimestampParam(io.Discard, tc.ProjectID); err != nil {
				t.Errorf("queryWithTimestampParam: %v", err)
			}
		})
		t.Run("queryPartitionedTable", func(t *testing.T) {
			t.Parallel()
			tableID := "bigquery_query_partitioned_table"
			if err := preparePartitionedData(tc.ProjectID, testDatasetID, tableID); err != nil {
				t.Fatalf("couldn't setup clustered table: %v", err)
			}
			if err := queryPartitionedTable(io.Discard, tc.ProjectID, testDatasetID, tableID); err != nil {
				t.Errorf("queryPartitionedTable: %v", err)
			}
		})
		t.Run("queryClusteredTable", func(t *testing.T) {
			t.Parallel()
			tableID := "bigquery_query_clustered_table"
			if err := prepareClusteredData(tc.ProjectID, testDatasetID, tableID); err != nil {
				t.Fatalf("couldn't setup clustered table: %v", err)
			}
			if err := queryClusteredTable(io.Discard, tc.ProjectID, testDatasetID, tableID); err != nil {
				t.Errorf("queryClusteredTable: %v", err)
			}
		})
		t.Run("queryJobOptional", func(t *testing.T) {
			t.Parallel()
			if err := queryJobOptional(io.Discard, tc.ProjectID); err != nil {
				t.Errorf("queryJobOptional: %v", err)
			}
		})
	})

}

// preparePartitionedData setups up example partitioned/clustered table resources for the query tests.
func preparePartitionedData(projectID, datasetID, tableID string) error {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}

	gcsRef := bigquery.NewGCSReference("gs://cloud-samples-data/bigquery/us-states/us-states-by-date.csv")
	gcsRef.SkipLeadingRows = 1
	gcsRef.Schema = bigquery.Schema{
		{Name: "name", Type: bigquery.StringFieldType},
		{Name: "post_abbr", Type: bigquery.StringFieldType},
		{Name: "date", Type: bigquery.DateFieldType},
	}
	loader := client.Dataset(datasetID).Table(tableID).LoaderFrom(gcsRef)
	loader.TimePartitioning = &bigquery.TimePartitioning{
		Field:      "date",
		Expiration: 90 * 24 * time.Hour,
	}
	loader.WriteDisposition = bigquery.WriteEmpty

	job, err := loader.Run(ctx)
	if err != nil {
		return err
	}
	_, err = job.Wait(ctx)
	return nil
}

// prepareClusteredData setups up example partitioned/clustered table resources for the query tests.
func prepareClusteredData(projectID, datasetID, tableID string) error {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}

	gcsRef := bigquery.NewGCSReference("gs://cloud-samples-data/bigquery/sample-transactions/transactions.csv")
	gcsRef.SkipLeadingRows = 1
	gcsRef.Schema = bigquery.Schema{
		{Name: "timestamp", Type: bigquery.TimestampFieldType},
		{Name: "origin", Type: bigquery.StringFieldType},
		{Name: "destination", Type: bigquery.StringFieldType},
		{Name: "amount", Type: bigquery.NumericFieldType},
	}
	loader := client.Dataset(datasetID).Table(tableID).LoaderFrom(gcsRef)
	loader.TimePartitioning = &bigquery.TimePartitioning{
		Field: "timestamp",
	}
	loader.Clustering = &bigquery.Clustering{
		Fields: []string{"origin", "destination"},
	}
	loader.WriteDisposition = bigquery.WriteEmpty

	job, err := loader.Run(ctx)
	if err != nil {
		return err
	}
	_, err = job.Wait(ctx)
	return nil
}
