// Copyright 2024 Google LLC
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

package spanner

// [START spanner_set_max_commit_delay]

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/spanner"
)

func setMaxCommitDelay(w io.Writer, db string) error {
	// db is the fully-qualified database name of the form `projects/<project>/instances/<instance-id>/database/<database-id>`
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return fmt.Errorf("setMaxCommitDelay.NewClient: %w", err)
	}
	defer client.Close()

	commitDelay := 100 * time.Millisecond
	resp, err := client.ReadWriteTransactionWithOptions(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT Singers (SingerId, FirstName, LastName)
					VALUES (111, 'Virginia', 'Watson')`,
		}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d record(s) inserted.\n", rowCount)
		return nil
	}, spanner.TransactionOptions{CommitOptions: spanner.CommitOptions{MaxCommitDelay: &commitDelay, ReturnCommitStats: true}})
	if err != nil {
		return fmt.Errorf("setMaxCommitDelay.ReadWriteTransactionWithOptions: %w", err)
	}
	fmt.Fprintf(w, "%d mutations in transaction\n", resp.CommitStats.MutationCount)
	return nil
}

// [END spanner_set_max_commit_delay]
