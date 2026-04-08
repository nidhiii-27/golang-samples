# Go Samples: Project Structure and Automation

This document provides a summary of the project structure, testing conventions, and automation workflows for the Go samples, with a focus on the `storage` directory.

## Analysis of `storage` Directory

### 1. Testing Conventions
*   **Framework:** Tests are written using the standard Go `testing` package.
*   **System Testing:** Most tests are **System Tests** that interact with live Google Cloud Storage services. They rely on the `internal/testutil` package for environment setup and authentication.
*   **Resiliency:** Tests utilize `testutil.Retry` to handle eventual consistency in the cloud and to mitigate transient network or API failures.
*   **Resource Cleanup:** Tests use `TestMain` for global setup/teardown and `t.Cleanup` or `defer` to ensure that buckets and objects created during testing are deleted.
*   **Conditional Execution:** Tests are often skipped if required environment variables (e.g., `GOOGLE_SAMPLES_PROJECT`) are not set, preventing failures in non-cloud environments.

### 2. Placement of Tests
*   **Co-location:** Tests are placed **directly alongside the samples** they verify. Every functional subdirectory (e.g., `buckets`, `objects`, `acl`) contains its own `*_test.go` file (e.g., `storage/buckets/buckets_test.go`).

### 3. Key Observations
*   **Modular Design:** The `storage` directory is an independent Go module with its own `go.mod` file.
*   **Snippet Tags:** Files use `// [START ...]` and `// [END ...]` tags to facilitate automatic snippet extraction for official Google Cloud documentation.
*   **Feature Coverage:** Includes samples for modern GCS features:
    *   **`rapid`**: Zonal buckets and appendable objects.
    *   **`control`**: Hierarchical namespaces (Folders/Managed Folders).
    *   **`transfer_manager`**: Parallelized and high-throughput operations.
*   **Interoperability:** The `s3_sdk` directory demonstrates GCS usage via the AWS S3 SDK interoperability mode.

---

## GitHub Actions Workflows

Automation is primarily handled by `.github/workflows/go.yaml`, which triggers on pushes and pull requests to the `main` branch.

### 1. Primary CI Jobs
*   **Build:** Recursively finds all `go.mod` files and runs `go build ./...` to ensure all samples compile.
*   **Lint:**
    *   Enforces formatting using `goimports`.
    *   Ensures dependency health by running `go mod tidy`.
    *   Fails if any changes are detected, ensuring the repo stays clean.
    *   Runs `shellcheck` on all shell (`.sh`) scripts.
*   **Vet:** Runs `go vet ./...` across all modules to catch common code patterns that may cause runtime issues.
*   **Root Tests:** Executes `go test -v` for any tests located in the root directory.

### 2. Additional Automation Tools
*   **Header Checker (`header-checker-lint.yml`):** Validates that all files contain the required Google LLC copyright and Apache-2.0 license headers.
*   **Snippet Bot (`snippet-bot.yml`):** Manages and validates code snippets used in documentation.
*   **Renovate (`renovate.json`):** Automatically updates Go dependencies across all modules in the repository.
