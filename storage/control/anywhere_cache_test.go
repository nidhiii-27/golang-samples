// Copyright 2026 Google LLC
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

package control

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestAnywhereCache(t *testing.T) {
	// Anywhere Cache tests often require specific project permissions and long runtimes.
	// Skipping by default similar to folders_test.go or if environment is not set.
	tc := testutil.SystemTest(t)
	zone := os.Getenv("GOOGLE_CLOUD_CPP_TEST_ZONE")
	if zone == "" {
		t.Skip("GOOGLE_CLOUD_CPP_TEST_ZONE not set")
	}

	ctx := context.Background()

	// Initialize local storage client.
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	t.Cleanup(func() { storageClient.Close() })

	// Create bucket with UBLA enabled.
	// Use local prefix to avoid dependency on other test files.
	acTestPrefix := "storage-control-ac-test"
	bucketName := testutil.UniqueBucketName(acTestPrefix)
	b := storageClient.Bucket(bucketName)
	attrs := &storage.BucketAttrs{
		UniformBucketLevelAccess: storage.UniformBucketLevelAccess{
			Enabled: true,
		},
	}
	if err := b.Create(ctx, tc.ProjectID, attrs); err != nil {
		t.Fatalf("Bucket.Create(%q): %v", bucketName, err)
	}
	t.Cleanup(func() {
		if err := testutil.DeleteBucketIfExists(ctx, storageClient, bucketName); err != nil {
			log.Printf("Bucket.Delete(%q): %v", bucketName, err)
		}
	})

	anywhereCacheID := zone
	// Use partial matches for resource names to handle project ID/number/placeholder inconsistencies.
	wantPartialName := fmt.Sprintf("buckets/%v/anywhereCaches/%v", bucketName, anywhereCacheID)

	// Create Anywhere Cache.
	buf := &bytes.Buffer{}
	if err := createAnywhereCache(buf, bucketName, zone); err != nil {
		t.Fatalf("createAnywhereCache: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, wantPartialName) {
		t.Errorf("createAnywhereCache: got %q, want to contain %q", got, wantPartialName)
	}

	// Get Anywhere Cache.
	buf.Reset()
	if err := getAnywhereCache(buf, bucketName, anywhereCacheID); err != nil {
		t.Errorf("getAnywhereCache: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, wantPartialName) {
		t.Errorf("getAnywhereCache: got %q, want to contain %q", got, wantPartialName)
	}

	// List Anywhere Caches.
	buf.Reset()
	if err := listAnywhereCaches(buf, bucketName); err != nil {
		t.Errorf("listAnywhereCaches: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, wantPartialName) {
		t.Errorf("listAnywhereCaches: got %q, want to contain %q", got, wantPartialName)
	}

	// Update Anywhere Cache.
	buf.Reset()
	admissionPolicy := "admit-on-second-miss"
	if err := updateAnywhereCache(buf, bucketName, anywhereCacheID, admissionPolicy); err != nil {
		t.Errorf("updateAnywhereCache: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, wantPartialName) {
		t.Errorf("updateAnywhereCache: got %q, want to contain %q", got, wantPartialName)
	}

	// Pause Anywhere Cache.
	buf.Reset()
	if err := pauseAnywhereCache(buf, bucketName, anywhereCacheID); err != nil {
		t.Errorf("pauseAnywhereCache: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, wantPartialName) {
		t.Errorf("pauseAnywhereCache: got %q, want to contain %q", got, wantPartialName)
	}

	// Resume Anywhere Cache.
	buf.Reset()
	if err := resumeAnywhereCache(buf, bucketName, anywhereCacheID); err != nil {
		t.Errorf("resumeAnywhereCache: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, wantPartialName) {
		t.Errorf("resumeAnywhereCache: got %q, want to contain %q", got, wantPartialName)
	}

	// Disable Anywhere Cache.
	buf.Reset()
	if err := disableAnywhereCache(buf, bucketName, anywhereCacheID); err != nil {
		t.Errorf("disableAnywhereCache: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, wantPartialName) {
		t.Errorf("disableAnywhereCache: got %q, want to contain %q", got, wantPartialName)
	}
}
