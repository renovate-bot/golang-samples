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

package topics

// [START pubsub_publish_custom_attributes]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub/v2"
)

func publishCustomAttributes(w io.Writer, projectID, topicID string) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	// client.Publisher can be passed a topic ID (e.g. "my-topic") or
	// a fully qualified name (e.g. "projects/my-project/topics/my-topic").
	// If a topic ID is provided, the project ID from the client is used.
	// Reuse this publisher for all publish calls to send messages in batches.
	publisher := client.Publisher(topicID)
	result := publisher.Publish(ctx, &pubsub.Message{
		Data: []byte("Hello world!"),
		Attributes: map[string]string{
			"origin":   "golang",
			"username": "gcp",
		},
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("Get: %w", err)
	}
	fmt.Fprintf(w, "Published message with custom attributes; msg ID: %v\n", id)
	return nil
}

// [END pubsub_publish_custom_attributes]
