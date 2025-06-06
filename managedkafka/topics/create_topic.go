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

package topics

// [START managedkafka_create_topic]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/managedkafka/apiv1/managedkafkapb"
	"google.golang.org/api/option"

	managedkafka "cloud.google.com/go/managedkafka/apiv1"
)

func createTopic(w io.Writer, projectID, region, clusterID, topicID string, partitionCount, replicationFactor int32, configs map[string]string, opts ...option.ClientOption) error {
	// projectID := "my-project-id"
	// region := "us-central1"
	// clusterID := "my-cluster"
	// topicID := "my-topic"
	// partitionCount := 10
	// replicationFactor := 3
	// configs := map[string]string{"min.insync.replicas":"1"}
	ctx := context.Background()
	client, err := managedkafka.NewClient(ctx, opts...)
	if err != nil {
		return fmt.Errorf("managedkafka.NewClient got err: %w", err)
	}
	defer client.Close()

	clusterPath := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", projectID, region, clusterID)
	topicPath := fmt.Sprintf("%s/topics/%s", clusterPath, topicID)
	topicConfig := &managedkafkapb.Topic{
		Name:              topicPath,
		PartitionCount:    partitionCount,
		ReplicationFactor: replicationFactor,
		Configs:           configs,
	}

	req := &managedkafkapb.CreateTopicRequest{
		Parent:  clusterPath,
		TopicId: topicID,
		Topic:   topicConfig,
	}
	topic, err := client.CreateTopic(ctx, req)
	if err != nil {
		return fmt.Errorf("client.CreateTopic got err: %w", err)
	}
	fmt.Fprintf(w, "Created topic: %s\n", topic.Name)
	return nil
}

// [END managedkafka_create_topic]
