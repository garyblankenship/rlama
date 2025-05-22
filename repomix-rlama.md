This file is a merged representation of the entire codebase, combined into a single document by Repomix.
The content has been processed where content has been compressed (code blocks are separated by ⋮---- delimiter).

# File Summary

## Purpose
This file contains a packed representation of the entire repository's contents.
It is designed to be easily consumable by AI systems for analysis, code review,
or other automated processes.

## File Format
The content is organized as follows:
1. This summary section
2. Repository information
3. Directory structure
4. Repository files (if enabled)
5. Multiple file entries, each consisting of:
  a. A header with the file path (## File: path/to/file)
  b. The full contents of the file in a code block

## Usage Guidelines
- This file should be treated as read-only. Any changes should be made to the
  original repository files, not this packed version.
- When processing this file, use the file path to distinguish
  between different files in the repository.
- Be aware that this file may contain sensitive information. Handle it with
  the same level of security as you would the original repository.

## Notes
- Some files may have been excluded based on .gitignore rules and Repomix's configuration
- Binary files are not included in this packed representation. Please refer to the Repository Structure section for a complete list of file paths, including binary files
- Files matching patterns in .gitignore are excluded
- Files matching default ignore patterns are excluded
- Content has been compressed - code blocks are separated by ⋮---- delimiter
- Files are sorted by Git change count (files with more changes are at the bottom)

# Directory Structure
```
.claude/
  settings.local.json
.github/
  workflows/
    release.yml
cmd/
  add_docs.go
  add_reranker.go
  api.go
  chunk_eval.go
  crawl_add_docs.go
  crawl_rag.go
  delete.go
  hf_browse.go
  install_dependencies.go
  list_chunks.go
  list_docs.go
  list.go
  migrate_batch.go
  migrate_to_internal.go
  migrate_to_qdrant.go
  migration_flags_test.go
  migration_flags.go
  profile.go
  rag.go
  remove_doc.go
  root.go
  run_hf.go
  run.go
  uninstall.go
  update_model.go
  update_reranker.go
  update.go
  watch.go
  web_watch.go
  wizard.go
docs/
  bge_onnx_reranker.md
  chunking_guidelines.md
  reranking_guidelines.md
internal/
  client/
    bge_onnx_reranker_client_test.go
    bge_onnx_reranker_client.go
    bge_reranker_benchmark_test.go
    bge_reranker_client.go
    client_test.go
    llm_client.go
    ollama_client.go
    openai_client.go
    pure_go_onnx_test.go
  config/
    config_test.go
    config.go
  crawler/
    crawl4ai_style_test.go
    crawl4ai_style.go
    crawler_test.go
    crawler.go
  domain/
    document_chunk.go
    document.go
    profile.go
    rag_test.go
    rag.go
  repository/
    profile_repository.go
    rag_repository.go
    repository_test.go
  server/
    server_test.go
    server.go
  service/
    chunker_evaluation.go
    chunker_service.go
    composite_rag_service.go
    config.go
    document_loader.go
    document_service.go
    document_temp_test.go
    embedding_cache.go
    embedding_service.go
    file_watcher.go
    integration_test.go
    migration_service.go
    provider.go
    query_service.go
    rag_service_test.go
    rag_service.go
    reranker_service_test.go
    reranker_service.go
    service_test.go
    watch_service.go
    web_watcher.go
  util/
    format_test.go
    format.go
pkg/
  vector/
    bruteforce_vector_store.go
    hybrid_store.go
    internal_vector_store.go
    qdrant_store.go
    store.go
    vector_test.go
scripts/
  build.sh
  install_deps.sh
test-small-docs/
  clayborn.txt
.gitattributes
.gitignore
.repomixignore
install.ps1
install.sh
main.go
README.md
repomix.config.json
```

# Files

## File: .claude/settings.local.json
````json
{
  "permissions": {
    "allow": [
      "Bash(rg:*)",
      "Bash(sed:*)",
      "Bash(go build:*)",
      "Bash(./rlama profile add:*)",
      "Bash(./rlama profile list)",
      "Bash(go test:*)",
      "Bash(./rlama profile delete:*)",
      "Bash(git checkout:*)",
      "Bash(git add:*)",
      "Bash(grep:*)",
      "Bash(./rlama rag:*)",
      "Bash(go vet:*)",
      "Bash(./rlama:*)",
      "Bash(curl:*)",
      "Bash(git push:*)",
      "Bash(ls:*)",
      "Bash(rlama run:*)",
      "Bash(rlama list:*)",
      "Bash(python -m pip:*)",
      "Bash(python:*)",
      "Bash(go get:*)",
      "Bash(go mod:*)",
      "Bash(cp:*)",
      "Bash(rm:*)",
      "Bash(find:*)",
      "Bash(mv:*)",
      "Bash(go run:*)",
      "Bash(go list:*)",
      "Bash(mkdir:*)",
      "Bash(git clone:*)"
    ],
    "deny": []
  }
}
````

## File: .github/workflows/release.yml
````yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    name: Build and Release
    runs-on: ubuntu-latest
    permissions:
      contents: write
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Build binaries
        run: |
          mkdir -p dist
          GOOS=linux GOARCH=amd64 go build -o dist/rlama_linux_amd64
          GOOS=linux GOARCH=arm64 go build -o dist/rlama_linux_arm64
          GOOS=darwin GOARCH=amd64 go build -o dist/rlama_darwin_amd64
          GOOS=darwin GOARCH=arm64 go build -o dist/rlama_darwin_arm64
          GOOS=windows GOARCH=amd64 go build -o dist/rlama_windows_amd64.exe
          chmod +x dist/*

      # Using GitHub's official actions instead
      - name: Create Release
        id: create_release
        uses: actions/create-release@v1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          tag_name: ${{ github.ref }}
          release_name: Release ${{ github.ref }}
          draft: false
          prerelease: false

      # Upload each asset separately
      - name: Upload Linux AMD64 Binary
        uses: actions/upload-release-asset@v1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          upload_url: ${{ steps.create_release.outputs.upload_url }}
          asset_path: ./dist/rlama_linux_amd64
          asset_name: rlama_linux_amd64
          asset_content_type: application/octet-stream

      - name: Upload Linux ARM64 Binary
        uses: actions/upload-release-asset@v1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          upload_url: ${{ steps.create_release.outputs.upload_url }}
          asset_path: ./dist/rlama_linux_arm64
          asset_name: rlama_linux_arm64
          asset_content_type: application/octet-stream

      - name: Upload macOS AMD64 Binary
        uses: actions/upload-release-asset@v1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          upload_url: ${{ steps.create_release.outputs.upload_url }}
          asset_path: ./dist/rlama_darwin_amd64
          asset_name: rlama_darwin_amd64
          asset_content_type: application/octet-stream

      - name: Upload macOS ARM64 Binary
        uses: actions/upload-release-asset@v1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          upload_url: ${{ steps.create_release.outputs.upload_url }}
          asset_path: ./dist/rlama_darwin_arm64
          asset_name: rlama_darwin_arm64
          asset_content_type: application/octet-stream

      - name: Upload Windows AMD64 Binary
        uses: actions/upload-release-asset@v1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          upload_url: ${{ steps.create_release.outputs.upload_url }}
          asset_path: ./dist/rlama_windows_amd64.exe
          asset_name: rlama_windows_amd64.exe
          asset_content_type: application/octet-stream
````

## File: cmd/delete.go
````go
package cmd
⋮----
import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/repository"
)
⋮----
"fmt"
"strings"
⋮----
"github.com/spf13/cobra"
"github.com/dontizi/rlama/internal/repository"
⋮----
var forceDelete bool
⋮----
var deleteCmd = &cobra.Command{
	Use:   "delete [rag-name]",
	Short: "Delete a RAG system",
	Long:  `Permanently delete a RAG system and all its indexed documents.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]
		repo := repository.NewRagRepository()

		// Check if the RAG exists
		if !repo.Exists(ragName) {
			return fmt.Errorf("the RAG system '%s' does not exist", ragName)
		}

		// Ask for confirmation unless --force is specified
		if !forceDelete {
			fmt.Printf("Are you sure you want to permanently delete the RAG system '%s'? (y/n): ", ragName)
			var response string
			fmt.Scanln(&response)
			
			response = strings.ToLower(strings.TrimSpace(response))
			if response != "y" && response != "yes" {
				fmt.Println("Deletion cancelled.")
				return nil
			}
		}

		// Delete the RAG
		err := repo.Delete(ragName)
		if err != nil {
			return err
		}

		fmt.Printf("The RAG system '%s' has been successfully deleted.\n", ragName)
		return nil
	},
}
⋮----
// Check if the RAG exists
⋮----
// Ask for confirmation unless --force is specified
⋮----
var response string
⋮----
// Delete the RAG
⋮----
func init()
````

## File: cmd/list.go
````go
package cmd
⋮----
import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/repository"
	"github.com/dontizi/rlama/internal/util"
)
⋮----
"fmt"
"os"
"text/tabwriter"
⋮----
"github.com/spf13/cobra"
"github.com/dontizi/rlama/internal/repository"
"github.com/dontizi/rlama/internal/util"
⋮----
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available RAG systems",
	Long:  `Display a list of all RAG systems that have been created.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := repository.NewRagRepository()
		ragNames, err := repo.ListAll()
		if err != nil {
			return err
		}

		if len(ragNames) == 0 {
			fmt.Println("No RAG systems found.")
			return nil
		}

		fmt.Printf("Available RAG systems (%d found):\n\n", len(ragNames))
		
		// Use tabwriter for aligned display
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tMODEL\tCREATED ON\tDOCUMENTS\tSIZE")
		
		for _, name := range ragNames {
			rag, err := repo.Load(name)
			if err != nil {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", name, "error", "error", "error", "error")
				continue
			}
			
			// Format the date
			createdAt := rag.CreatedAt.Format("2006-01-02 15:04:05")
			
			// Calculate total size
			var totalSize int64
			for _, doc := range rag.Documents {
				totalSize += doc.Size
			}
			
			// Format the size
			sizeStr := util.FormatSize(totalSize)
			
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n", 
				rag.Name, rag.ModelName, createdAt, len(rag.Documents), sizeStr)
		}
		w.Flush()
		
		return nil
	},
}
⋮----
// Use tabwriter for aligned display
⋮----
// Format the date
⋮----
// Calculate total size
var totalSize int64
⋮----
// Format the size
⋮----
func init()
````

## File: cmd/migrate_batch.go
````go
package cmd
⋮----
import (
	"fmt"
	"strings"

	"github.com/dontizi/rlama/internal/repository"
	"github.com/dontizi/rlama/internal/service"
	"github.com/spf13/cobra"
)
⋮----
"fmt"
"strings"
⋮----
"github.com/dontizi/rlama/internal/repository"
"github.com/dontizi/rlama/internal/service"
"github.com/spf13/cobra"
⋮----
var (
	// Batch migration specific flags
	batchFromVectorStore         string
	batchToVectorStore           string
	batchRagNames                []string
	batchQdrantCollectionPrefix  string
	batchContinueOnError         bool

	// Common flag sets for batch migration
	batchQdrantFlags    QdrantFlags
	batchMigrationFlags MigrationFlags
)
⋮----
// Batch migration specific flags
⋮----
// Common flag sets for batch migration
⋮----
var migrateBatchCmd = &cobra.Command{
	Use:   "migrate-batch",
	Short: "Migrate multiple RAG systems between vector stores in batch",
	Long: `Migrate multiple RAG systems between different vector store types in a single operation.

This command allows you to migrate multiple RAGs at once, which is useful for:
- Migrating entire environments (dev to prod, local to cloud)
- Bulk upgrades from internal to Qdrant storage
- Infrastructure changes affecting multiple RAGs

The command will process each RAG individually and provide detailed progress and error reporting.

Examples:
  # Migrate all RAGs from internal to Qdrant
  rlama migrate-batch --from=internal --to=qdrant \
    --qdrant-host=production.qdrant.com

  # Migrate specific RAGs to Qdrant Cloud
  rlama migrate-batch --from=internal --to=qdrant \
    --rags=docs,wiki,knowledge-base \
    --qdrant-host=xyz.qdrant.cloud \
    --qdrant-apikey=your-api-key \
    --qdrant-grpc

  # Migrate with backup and cleanup
  rlama migrate-batch --from=internal --to=qdrant \
    --backup \
    --delete-old \
    --continue-on-error

  # Migrate from Qdrant back to internal (all RAGs)
  rlama migrate-batch --from=qdrant --to=internal \
    --backup`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate parameters
		if batchFromVectorStore == "" || batchToVectorStore == "" {
			return fmt.Errorf("both --from and --to parameters are required")
		}

		if batchFromVectorStore == batchToVectorStore {
			return fmt.Errorf("source and target vector stores cannot be the same")
		}

		validStores := map[string]bool{"internal": true, "qdrant": true}
		if !validStores[batchFromVectorStore] || !validStores[batchToVectorStore] {
			return fmt.Errorf("vector store must be 'internal' or 'qdrant'")
		}

		// Get list of RAGs to migrate
		var ragNames []string
		if len(batchRagNames) > 0 {
			ragNames = batchRagNames
		} else {
			// Get all RAGs if none specified
			ragRepo := repository.NewRagRepository()
			allRags, err := ragRepo.ListAll()
			if err != nil {
				return fmt.Errorf("failed to list RAGs: %w", err)
			}
			ragNames = allRags
		}

		if len(ragNames) == 0 {
			fmt.Println("📄 No RAGs found to migrate")
			return nil
		}

		// Filter RAGs by source vector store type
		filteredRags, err := filterRagsByVectorStore(ragNames, batchFromVectorStore)
		if err != nil {
			return fmt.Errorf("failed to filter RAGs: %w", err)
		}

		if len(filteredRags) == 0 {
			fmt.Printf("📄 No RAGs found using %s vector store\n", batchFromVectorStore)
			return nil
		}

		fmt.Printf("🚀 Starting batch migration of %d RAGs from %s to %s\n", 
			len(filteredRags), batchFromVectorStore, batchToVectorStore)

		// Create migration service
		migrationService := service.NewMigrationService()

		// Track results
		var successful, failed int
		var failedRags []string

		// Migrate each RAG
		for i, ragName := range filteredRags {
			fmt.Printf("\n📋 [%d/%d] Migrating RAG '%s'...\n", i+1, len(filteredRags), ragName)

			// Get flag values
			qdrantHost, qdrantPort, qdrantAPIKey, _, qdrantUseGRPC := GetQdrantFlagValues(&batchQdrantFlags)
			createBackup, backupPath, verify, deleteOld := GetMigrationFlagValues(&batchMigrationFlags)

			// Create migration options
			migrationOpts := service.MigrationOptions{
				TargetVectorStore:    batchToVectorStore,
				QdrantHost:           qdrantHost,
				QdrantPort:           qdrantPort,
				QdrantAPIKey:         qdrantAPIKey,
				QdrantCollectionName: getCollectionName(ragName, batchQdrantCollectionPrefix),
				QdrantGRPC:           qdrantUseGRPC,
				CreateBackup:         createBackup,
				BackupPath:           backupPath,
				VerifyAfterMigration: verify,
				DeleteOldData:        deleteOld,
			}

			// Perform migration
			var err error
			if batchToVectorStore == "qdrant" {
				err = migrationService.MigrateToQdrant(ragName, migrationOpts)
			} else {
				err = migrationService.MigrateToInternal(ragName, migrationOpts)
			}

			if err != nil {
				failed++
				failedRags = append(failedRags, ragName)
				fmt.Printf("❌ Failed to migrate RAG '%s': %v\n", ragName, err)
				
				if !batchContinueOnError {
					return fmt.Errorf("migration failed for RAG '%s', stopping batch operation", ragName)
				}
			} else {
				successful++
				fmt.Printf("✅ Successfully migrated RAG '%s'\n", ragName)
			}
		}

		// Summary
		fmt.Printf("\n📊 Batch Migration Summary:\n")
		fmt.Printf("   ✅ Successful: %d\n", successful)
		fmt.Printf("   ❌ Failed: %d\n", failed)

		if len(failedRags) > 0 {
			fmt.Printf("   🔴 Failed RAGs: %s\n", strings.Join(failedRags, ", "))
		}

		if failed > 0 {
			return fmt.Errorf("batch migration completed with %d failures", failed)
		}

		fmt.Println("🎉 Batch migration completed successfully!")
		return nil
	},
}
⋮----
// Validate parameters
⋮----
// Get list of RAGs to migrate
var ragNames []string
⋮----
// Get all RAGs if none specified
⋮----
// Filter RAGs by source vector store type
⋮----
// Create migration service
⋮----
// Track results
var successful, failed int
var failedRags []string
⋮----
// Migrate each RAG
⋮----
// Get flag values
⋮----
// Create migration options
⋮----
// Perform migration
var err error
⋮----
// Summary
⋮----
func init()
⋮----
// Required flags
⋮----
// Add common flag sets
⋮----
// Batch-specific Qdrant flag
⋮----
// filterRagsByVectorStore returns only RAGs that use the specified vector store type
func filterRagsByVectorStore(ragNames []string, vectorStoreType string) ([]string, error)
⋮----
var filtered []string
⋮----
// Check vector store type (default to internal if not set)
⋮----
// getCollectionName generates the Qdrant collection name
func getCollectionName(ragName, prefix string) string
````

## File: cmd/migrate_to_internal.go
````go
package cmd
⋮----
import (
	"fmt"

	"github.com/dontizi/rlama/internal/service"
	"github.com/spf13/cobra"
)
⋮----
"fmt"
⋮----
"github.com/dontizi/rlama/internal/service"
"github.com/spf13/cobra"
⋮----
var (
	// Flag holders for migrate-to-internal command
	internalMigrationFlags MigrationFlags
)
⋮----
// Flag holders for migrate-to-internal command
⋮----
var migrateToInternalCmd = &cobra.Command{
	Use:   "migrate-to-internal [rag-name]",
	Short: "Migrate a RAG system from Qdrant back to internal vector storage",
	Long: `Migrate an existing RAG system that uses Qdrant vector database back to internal vector storage.

This command will:
1. Load the existing RAG with Qdrant vector storage
2. Connect to the Qdrant instance and retrieve all vectors
3. Create internal vector storage with appropriate dimensions
4. Transfer all vectors and metadata to internal storage
5. Update the RAG configuration to use internal storage
6. Verify the migration was successful

This is useful for:
- Moving to offline environments
- Reducing infrastructure dependencies
- Testing or development scenarios
- Backup/disaster recovery

Examples:
  # Basic migration to internal storage
  rlama rag migrate-to-internal my-rag

  # Migration with backup
  rlama rag migrate-to-internal my-rag \
    --backup \
    --backup-path=/safe/backup/location

  # Migration and cleanup
  rlama rag migrate-to-internal my-rag \
    --backup \
    --delete-old`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]

		// Get flag values
		createBackup, backupPath, verify, deleteOld := GetMigrationFlagValues(&internalMigrationFlags)

		// Create migration options
		migrationOpts := service.MigrationOptions{
			TargetVectorStore:    "internal",
			CreateBackup:        createBackup,
			BackupPath:          backupPath,
			VerifyAfterMigration: verify,
			DeleteOldData:       deleteOld,
		}

		// Create migration service
		migrationService := service.NewMigrationService()

		fmt.Printf("🔄 Starting migration of RAG '%s' to internal storage...\n", ragName)

		if migrationOpts.CreateBackup {
			fmt.Printf("💾 Backup will be created at: %s\n", migrationOpts.BackupPath)
		}

		// Perform migration
		err := migrationService.MigrateToInternal(ragName, migrationOpts)
		if err != nil {
			return fmt.Errorf("❌ Migration failed: %w", err)
		}

		fmt.Printf("✅ Successfully migrated RAG '%s' to internal storage!\n", ragName)
		return nil
	},
}
⋮----
// Get flag values
⋮----
// Create migration options
⋮----
// Create migration service
⋮----
// Perform migration
⋮----
func init()
⋮----
// Add common migration control flags
````

## File: cmd/migrate_to_qdrant.go
````go
package cmd
⋮----
import (
	"fmt"

	"github.com/dontizi/rlama/internal/service"
	"github.com/spf13/cobra"
)
⋮----
"fmt"
⋮----
"github.com/dontizi/rlama/internal/service"
"github.com/spf13/cobra"
⋮----
var (
	// Flag holders for migrate-to-qdrant command
	qdrantFlags    QdrantFlags
	migrationFlags MigrationFlags
)
⋮----
// Flag holders for migrate-to-qdrant command
⋮----
var migrateToQdrantCmd = &cobra.Command{
	Use:   "migrate-to-qdrant [rag-name]",
	Short: "Migrate a RAG system from internal storage to Qdrant vector database",
	Long: `Migrate an existing RAG system that uses internal vector storage to use Qdrant vector database instead.

This command will:
1. Load the existing RAG with internal vector storage
2. Connect to the specified Qdrant instance
3. Create a collection with appropriate dimensions
4. Transfer all vectors and metadata to Qdrant
5. Update the RAG configuration to use Qdrant
6. Verify the migration was successful

Examples:
  # Basic migration to local Qdrant
  rlama rag migrate-to-qdrant my-rag

  # Migration to remote Qdrant with custom settings
  rlama rag migrate-to-qdrant my-rag \
    --qdrant-host=production.qdrant.com \
    --qdrant-port=6333 \
    --qdrant-collection=migrated-docs \
    --qdrant-grpc

  # Migration with backup
  rlama rag migrate-to-qdrant my-rag \
    --backup \
    --backup-path=/safe/backup/location

  # Migration to Qdrant Cloud
  rlama rag migrate-to-qdrant my-rag \
    --qdrant-host=xyz.qdrant.cloud \
    --qdrant-apikey=your-api-key \
    --qdrant-grpc`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]

		// Get flag values
		qdrantHost, qdrantPort, qdrantAPIKey, qdrantCollection, qdrantUseGRPC := GetQdrantFlagValues(&qdrantFlags)
		createBackup, backupPath, verify, deleteOld := GetMigrationFlagValues(&migrationFlags)

		// Create migration options
		migrationOpts := service.MigrationOptions{
			TargetVectorStore:    "qdrant",
			QdrantHost:          qdrantHost,
			QdrantPort:          qdrantPort,
			QdrantAPIKey:        qdrantAPIKey,
			QdrantCollectionName: qdrantCollection,
			QdrantGRPC:          qdrantUseGRPC,
			CreateBackup:        createBackup,
			BackupPath:          backupPath,
			VerifyAfterMigration: verify,
			DeleteOldData:       deleteOld,
		}

		// Set default collection name if not provided
		if migrationOpts.QdrantCollectionName == "" {
			migrationOpts.QdrantCollectionName = ragName
		}

		// Create migration service
		migrationService := service.NewMigrationService()

		fmt.Printf("🔄 Starting migration of RAG '%s' to Qdrant...\n", ragName)
		fmt.Printf("📋 Target: %s:%d, Collection: %s\n", 
			migrationOpts.QdrantHost, migrationOpts.QdrantPort, migrationOpts.QdrantCollectionName)

		if migrationOpts.CreateBackup {
			fmt.Printf("💾 Backup will be created at: %s\n", migrationOpts.BackupPath)
		}

		// Perform migration
		err := migrationService.MigrateToQdrant(ragName, migrationOpts)
		if err != nil {
			return fmt.Errorf("❌ Migration failed: %w", err)
		}

		fmt.Printf("✅ Successfully migrated RAG '%s' to Qdrant!\n", ragName)
		return nil
	},
}
⋮----
// Get flag values
⋮----
// Create migration options
⋮----
// Set default collection name if not provided
⋮----
// Create migration service
⋮----
// Perform migration
⋮----
func init()
⋮----
// Add common flag sets
````

## File: cmd/migration_flags_test.go
````go
package cmd
⋮----
import (
	"testing"

	"github.com/spf13/cobra"
)
⋮----
"testing"
⋮----
"github.com/spf13/cobra"
⋮----
func TestQdrantFlags(t *testing.T)
⋮----
var flags QdrantFlags
⋮----
// Add Qdrant flags
⋮----
// Test flag registration
⋮----
// Check that all expected flags are registered
⋮----
// Test default values by parsing empty args
⋮----
// Check default values
⋮----
func TestMigrationFlags(t *testing.T)
⋮----
var flags MigrationFlags
⋮----
// Add migration flags
⋮----
// Test default values
⋮----
func TestFlagValueParsing(t *testing.T)
⋮----
var qdrantFlags QdrantFlags
var migrationFlags MigrationFlags
⋮----
// Test parsing custom values
⋮----
// Verify Qdrant flag values
⋮----
// Verify migration flag values
````

## File: cmd/migration_flags.go
````go
package cmd
⋮----
import "github.com/spf13/cobra"
⋮----
// QdrantFlags holds common Qdrant connection configuration flags
type QdrantFlags struct {
	Host       *string
	Port       *int
	APIKey     *string
	Collection *string
	UseGRPC    *bool
}
⋮----
// MigrationFlags holds common migration control flags
type MigrationFlags struct {
	CreateBackup         *bool
	BackupPath           *string
	VerifyAfterMigration *bool
	DeleteOldData        *bool
}
⋮----
// AddQdrantFlags adds standard Qdrant connection flags to a cobra command
func AddQdrantFlags(cmd *cobra.Command, flags *QdrantFlags, collectionUsage string)
⋮----
// AddMigrationControlFlags adds standard migration control flags to a cobra command
func AddMigrationControlFlags(cmd *cobra.Command, flags *MigrationFlags)
⋮----
// AddBatchMigrationFlags adds flags specific to batch migration operations
func AddBatchMigrationFlags(cmd *cobra.Command, ragNames *[]string, continueOnError *bool)
⋮----
// GetQdrantFlagValues returns the actual values from QdrantFlags pointers
func GetQdrantFlagValues(flags *QdrantFlags) (host string, port int, apiKey string, collection string, useGRPC bool)
⋮----
// GetMigrationFlagValues returns the actual values from MigrationFlags pointers
func GetMigrationFlagValues(flags *MigrationFlags) (createBackup bool, backupPath string, verify bool, deleteOld bool)
````

## File: cmd/remove_doc.go
````go
package cmd
⋮----
import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/service"
)
⋮----
"fmt"
"strings"
⋮----
"github.com/spf13/cobra"
"github.com/dontizi/rlama/internal/service"
⋮----
var forceRemoveDoc bool
⋮----
var removeDocCmd = &cobra.Command{
	Use:   "remove-doc [rag-name] [doc-id]",
	Short: "Remove a document from a RAG system",
	Long: `Remove a specific document from a RAG system by its ID.
Example: rlama remove-doc my-docs document.pdf
	
The document ID is typically the filename. You can see document IDs by using the
"rlama list-docs [rag-name]" command.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]
		docID := args[1]

		// Get Ollama client from root command
		ollamaClient := GetOllamaClient()

		// Create necessary services
		ragService := service.NewRagService(ollamaClient)

		// Load the RAG
		rag, err := ragService.LoadRag(ragName)
		if err != nil {
			return err
		}

		// Find the document
		doc := rag.GetDocumentByID(docID)
		if doc == nil {
			return fmt.Errorf("document with ID '%s' not found in RAG '%s'", docID, ragName)
		}

		// Ask for confirmation unless --force is specified
		if !forceRemoveDoc {
			fmt.Printf("Are you sure you want to remove document '%s' from RAG '%s'? (y/n): ", 
				doc.Name, ragName)
			var response string
			fmt.Scanln(&response)
			
			response = strings.ToLower(strings.TrimSpace(response))
			if response != "y" && response != "yes" {
				fmt.Println("Document removal cancelled.")
				return nil
			}
		}

		// Remove the document
		removed := rag.RemoveDocument(docID)
		if !removed {
			return fmt.Errorf("failed to remove document '%s'", docID)
		}

		// Save the RAG
		err = ragService.UpdateRag(rag)
		if err != nil {
			return fmt.Errorf("error saving the RAG: %w", err)
		}

		fmt.Printf("Successfully removed document '%s' from RAG '%s'.\n", doc.Name, ragName)
		return nil
	},
}
⋮----
// Get Ollama client from root command
⋮----
// Create necessary services
⋮----
// Load the RAG
⋮----
// Find the document
⋮----
// Ask for confirmation unless --force is specified
⋮----
var response string
⋮----
// Remove the document
⋮----
// Save the RAG
⋮----
func init()
````

## File: docs/bge_onnx_reranker.md
````markdown
# BGE ONNX Reranker Implementation

This document describes the Go-native BGE reranker implementation using ONNX runtime.

## Overview

The BGE ONNX reranker provides a faster alternative to the original Python subprocess-based implementation by using:

1. **Pre-exported ONNX models** - No need to export models yourself
2. **Python ONNX microservice** - Runs ONNX inference in a dedicated HTTP server
3. **Go HTTP client** - Communicates with the microservice for reranking

## Architecture

```
┌─────────────────┐    HTTP     ┌──────────────────────┐
│ Go Application  │ ──────────► │ Python ONNX Server   │
│                 │             │                      │
│ BGEONNXReranker │             │ - onnxruntime        │
│ Client          │             │ - transformers       │
│                 │             │ - model.onnx         │
└─────────────────┘             └──────────────────────┘
```

## Performance Benefits

The ONNX implementation provides significant performance improvements:

- **8-15 seconds** vs 20-30 seconds for the original PyTorch models
- **Persistent server** - No subprocess startup overhead
- **Optimized inference** - ONNX runtime optimizations
- **Batch processing** - Multiple pairs in single request

## Setup Requirements

### 1. Download Pre-exported ONNX Model

```bash
mkdir -p ./models
cd ./models
git clone https://huggingface.co/corto-ai/bge-reranker-large-onnx
```

### 2. Install Python Dependencies

```bash
pip install onnxruntime transformers numpy
```

### 3. Verify Installation

```bash
go test ./internal/client -v -run TestBGEONNXRerankerClient
```

## Usage

### Basic Usage

```go
import "github.com/dontizi/rlama/internal/client"

// Create ONNX reranker client
modelDir := "./models/bge-reranker-large-onnx"
client, err := client.NewBGEONNXRerankerClient(modelDir)
if err != nil {
    log.Fatal(err)
}
defer client.Cleanup() // Important: stops the Python server

// Rerank query-passage pairs
pairs := [][]string{
    {"What is a cat?", "A cat is a small domesticated carnivorous mammal."},
    {"What is a cat?", "The weather is nice today."},
}

scores, err := client.ComputeScores(pairs, true) // true = normalize scores
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Scores: %v\n", scores) // [0.95, 0.05] (first pair more relevant)
```

### Integration with Existing Reranker Service

The ONNX client implements the same interface as the original BGE client:

```go
type RerankerClient interface {
    ComputeScores(pairs [][]string, normalize bool) ([]float64, error)
    GetModelName() string
}
```

To integrate, modify the reranker service to choose between implementations:

```go
func NewRerankerClient(modelName string, useONNX bool) RerankerClient {
    if useONNX {
        modelDir := "./models/bge-reranker-large-onnx"
        return client.NewBGEONNXRerankerClient(modelDir)
    }
    return client.NewBGERerankerClient(modelName)
}
```

## Available ONNX Models

Several pre-exported ONNX models are available on Hugging Face:

- `corto-ai/bge-reranker-large-onnx` - Standard ONNX version (recommended)
- `swulling/bge-reranker-large-onnx-o4` - O4 optimized version
- `swulling/bge-reranker-base-onnx-o4` - Base model, O4 optimized  
- `EmbeddedLLM/bge-reranker-base-onnx-o4-o2-gpu` - GPU optimized

## Implementation Details

### Microservice Approach

The implementation uses a Python HTTP server that:

1. **Loads ONNX model** once at startup
2. **Tokenizes input** using HuggingFace transformers
3. **Runs ONNX inference** with optimized runtime
4. **Returns scores** via JSON API

### Input Format

The BGE reranker expects input in the format:
```
query + " </s> " + passage
```

### Output Format

- **Normalized scores**: Sigmoid applied to logits (0.0 to 1.0)
- **Raw scores**: Direct logits output (any real number)

### Error Handling

The client handles common errors:
- Invalid pair format (not exactly 2 elements)
- Server connection failures
- ONNX runtime errors
- Tokenization errors

## Testing

Run the test suite to verify functionality:

```bash
# Basic functionality tests
go test ./internal/client -v -run TestBGEONNXRerankerClient

# Performance tests  
go test ./internal/client -v -run TestBGEONNXRerankerClient_Performance

# Benchmark against original implementation
go test ./internal/client -bench=BenchmarkBGEReranker
```

## Troubleshooting

### Common Issues

1. **"Model directory not found"**
   - Ensure ONNX model is downloaded to correct path
   - Check file permissions

2. **"Failed to start Python server"**
   - Verify Python dependencies are installed
   - Check port 8765 is available
   - Ensure Python is in PATH

3. **"Invalid input name: token_type_ids"**
   - This indicates ONNX model doesn't expect token_type_ids
   - Fixed in current implementation

### Performance Tuning

1. **Batch Size**: Process multiple pairs in single request
2. **Server Persistence**: Keep server running between requests
3. **Model Selection**: Use base model for faster inference if acceptable

## Future Improvements

1. **Pure Go Implementation**: Direct ONNX runtime without Python
2. **GPU Acceleration**: Use CUDA-enabled ONNX models
3. **Model Caching**: Cache tokenizer and model in memory
4. **Connection Pooling**: Reuse HTTP connections

## References

- [ONNX Runtime Go bindings](https://github.com/yalue/onnxruntime_go)
- [BGE Reranker Paper](https://arxiv.org/abs/2309.07597)
- [ONNX Model Hub](https://huggingface.co/models?library=onnx)
````

## File: internal/client/bge_onnx_reranker_client_test.go
````go
package client
⋮----
import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)
⋮----
"os"
"path/filepath"
"testing"
"time"
⋮----
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"
⋮----
func TestBGEONNXRerankerClient(t *testing.T)
⋮----
// Skip test if model directory doesn't exist
⋮----
// Cleanup
⋮----
// Score should be between 0 and 1 when normalized
⋮----
// First pair should have higher score than second pair (more relevant)
⋮----
// All scores should be normalized between 0 and 1
⋮----
// Without normalization, scores can be any real number (logits)
⋮----
func TestBGEONNXRerankerClient_Performance(t *testing.T)
⋮----
// Test performance with multiple pairs
⋮----
// Should be faster than the original Python subprocess approach
⋮----
// Relevant pairs should score higher than irrelevant ones
````

## File: internal/client/bge_onnx_reranker_client.go
````go
package client
⋮----
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)
⋮----
"bytes"
"encoding/json"
"fmt"
"io"
"math/rand"
"net/http"
"os"
"os/exec"
"path/filepath"
"time"
⋮----
// BGEONNXRerankerClient handles BGE reranking using ONNX runtime via HTTP microservice
type BGEONNXRerankerClient struct {
	serverURL   string
	httpClient  *http.Client
	modelDir    string
	serverProc  *exec.Cmd
}
⋮----
// NewBGEONNXRerankerClient creates a new ONNX-based BGE reranker client
func NewBGEONNXRerankerClient(modelDir string) (*BGEONNXRerankerClient, error)
⋮----
// Find an available port
⋮----
// Start the Python ONNX server
⋮----
// findAvailablePort finds an available port for the ONNX server
func findAvailablePort() int
⋮----
// Start from a base port and add a random offset to avoid conflicts
⋮----
// startONNXServer starts a Python HTTP server that runs ONNX inference
func (c *BGEONNXRerankerClient) startONNXServer(port int) error
⋮----
// Create the Python server script
⋮----
// Write server script to temporary file
⋮----
// Start the server process
⋮----
// Wait for server to be ready
⋮----
// isServerReady checks if the server is responding
func (c *BGEONNXRerankerClient) isServerReady() bool
⋮----
// ComputeScores calculates relevance scores between queries and passages using ONNX
func (c *BGEONNXRerankerClient) ComputeScores(pairs [][]string, normalize bool) ([]float64, error)
⋮----
var response map[string]interface{}
⋮----
// GetModelName returns the model identifier
func (c *BGEONNXRerankerClient) GetModelName() string
⋮----
// Cleanup properly stops the server and frees resources
func (c *BGEONNXRerankerClient) Cleanup() error
````

## File: internal/client/bge_reranker_benchmark_test.go
````go
package client
⋮----
import (
	"os"
	"path/filepath"
	"testing"
)
⋮----
"os"
"path/filepath"
"testing"
⋮----
func BenchmarkBGEReranker(b *testing.B)
⋮----
// Test data
⋮----
func BenchmarkBGEReranker_SinglePair(b *testing.B)
````

## File: internal/client/pure_go_onnx_test.go
````go
package client
⋮----
import (
	"os"
	"path/filepath"
	"testing"

	ort "github.com/yalue/onnxruntime_go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)
⋮----
"os"
"path/filepath"
"testing"
⋮----
ort "github.com/yalue/onnxruntime_go"
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"
⋮----
func TestPureGoONNXRuntime(t *testing.T)
⋮----
// Skip test if model directory doesn't exist
⋮----
// Create dummy input tensors for inspection
⋮----
// Output tensor
⋮----
// Create session
⋮----
// Verify tensor shapes
⋮----
// Fill input tensors with dummy data
⋮----
// Simple dummy tokenization:
// token_ids: [0, 1, 2, 3, ..., 10] + padding with 1 (pad token)
// attention_mask: [1, 1, 1, ...] for real tokens, [0, 0, ...] for padding
⋮----
inputIdsData[i] = int64(i)  // Some dummy token IDs
attentionMaskData[i] = 1    // Attention for real tokens
⋮----
inputIdsData[i] = 1         // Pad token ID
attentionMaskData[i] = 0    // No attention for padding
⋮----
// Run inference
⋮----
// Check output
⋮----
// Convert to probability using sigmoid
⋮----
// Score should be a reasonable probability (0-1)
⋮----
func TestONNXRuntimeCapabilities(t *testing.T)
⋮----
// Test tensor creation and manipulation
⋮----
// Test data access
⋮----
// Test data modification
⋮----
// Verify data was set
⋮----
// Test int64 tensors (for input_ids, attention_mask)
⋮----
// Test float32 tensors (for outputs)
````

## File: internal/service/composite_rag_service.go
````go
package service
⋮----
import (
	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/domain"
)
⋮----
"github.com/dontizi/rlama/internal/client"
"github.com/dontizi/rlama/internal/domain"
⋮----
// CompositeRagService implements RagService by composing focused services
// This replaces the monolithic RagServiceImpl with a cleaner architecture
type CompositeRagService struct {
	documentService DocumentService
	queryService    QueryService
	watchService    WatchService
	ollamaClient    *client.OllamaClient
}
⋮----
// NewCompositeRagService creates a new composite RAG service
func NewCompositeRagService(llmClient client.LLMClient, ollamaClient *client.OllamaClient) RagService
⋮----
// Create focused services
⋮----
// Create the composite service first, then create watch service with it
⋮----
// Now create watch service with the composite service
⋮----
// NewCompositeRagServiceWithConfig creates a new composite RAG service with configuration options
func NewCompositeRagServiceWithConfig(llmClient client.LLMClient, ollamaClient *client.OllamaClient, config *ServiceConfig) RagService
⋮----
// Create focused services with configuration
⋮----
// CreateRagWithOptions implements RagService.CreateRagWithOptions
func (crs *CompositeRagService) CreateRagWithOptions(modelName, ragName, folderPath string, options DocumentLoaderOptions) error
⋮----
// GetRagChunks implements RagService.GetRagChunks
func (crs *CompositeRagService) GetRagChunks(ragName string, filter ChunkFilter) ([]*domain.DocumentChunk, error)
⋮----
// LoadRag implements RagService.LoadRag
func (crs *CompositeRagService) LoadRag(ragName string) (*domain.RagSystem, error)
⋮----
// Query implements RagService.Query
func (crs *CompositeRagService) Query(rag *domain.RagSystem, query string, contextSize int) (string, error)
⋮----
// AddDocsWithOptions implements RagService.AddDocsWithOptions
func (crs *CompositeRagService) AddDocsWithOptions(ragName string, folderPath string, options DocumentLoaderOptions) error
⋮----
// UpdateModel implements RagService.UpdateModel
func (crs *CompositeRagService) UpdateModel(ragName string, newModel string) error
⋮----
// Load the RAG
⋮----
// Update the model
⋮----
// Save the updated RAG
⋮----
// UpdateRag implements RagService.UpdateRag
func (crs *CompositeRagService) UpdateRag(rag *domain.RagSystem) error
⋮----
// UpdateRerankerModel implements RagService.UpdateRerankerModel
func (crs *CompositeRagService) UpdateRerankerModel(ragName string, model string) error
⋮----
// ListAllRags implements RagService.ListAllRags
func (crs *CompositeRagService) ListAllRags() ([]string, error)
⋮----
// GetOllamaClient implements RagService.GetOllamaClient
func (crs *CompositeRagService) GetOllamaClient() *client.OllamaClient
⋮----
// SetPreferredEmbeddingModel implements RagService.SetPreferredEmbeddingModel
func (crs *CompositeRagService) SetPreferredEmbeddingModel(model string)
⋮----
// Directory watching methods - delegate to WatchService
func (crs *CompositeRagService) SetupDirectoryWatching(ragName string, dirPath string, watchInterval int, options DocumentLoaderOptions) error
⋮----
func (crs *CompositeRagService) DisableDirectoryWatching(ragName string) error
⋮----
func (crs *CompositeRagService) CheckWatchedDirectory(ragName string) (int, error)
⋮----
// Web watching methods - delegate to WatchService
func (crs *CompositeRagService) SetupWebWatching(ragName string, websiteURL string, watchInterval int, options domain.WebWatchOptions) error
⋮----
func (crs *CompositeRagService) DisableWebWatching(ragName string) error
⋮----
func (crs *CompositeRagService) CheckWatchedWebsite(ragName string) (int, error)
````

## File: internal/service/config.go
````go
package service
⋮----
import (
	"fmt"
	"os"
	"strconv"
)
⋮----
"fmt"
"os"
"strconv"
⋮----
// ServiceConfig holds all configuration needed for service creation
type ServiceConfig struct {
	// Client Configuration
	OllamaHost     string
	OllamaPort     string
	OpenAIAPIKey   string
	DataDirectory  string
	
	// Profile Configuration
	APIProfileName string
	
	// Document Processing Configuration
	ChunkSize        int
	ChunkOverlap     int
	ChunkingStrategy string
	
	// Embedding Configuration
	EmbeddingModel string
	
	// Reranking Configuration
	RerankerEnabled   bool
	RerankerModel     string
	RerankerWeight    float64
	RerankerThreshold float64
	UseONNXReranker   bool
	ONNXModelDir      string
	
	// Vector Store Configuration
	VectorStoreType      string
	QdrantHost           string
	QdrantPort           int
	QdrantAPIKey         string
	QdrantCollectionName string
	QdrantGRPC           bool
	
	// Debugging and Logging
	Verbose bool
}
⋮----
// Client Configuration
⋮----
// Profile Configuration
⋮----
// Document Processing Configuration
⋮----
// Embedding Configuration
⋮----
// Reranking Configuration
⋮----
// Vector Store Configuration
⋮----
// Debugging and Logging
⋮----
// NewServiceConfig creates a new service configuration with defaults
func NewServiceConfig() *ServiceConfig
⋮----
// Client defaults
⋮----
// Document processing defaults
⋮----
// Reranking defaults
⋮----
// Vector store defaults
⋮----
// Validate checks that the configuration is valid
func (sc *ServiceConfig) Validate() error
⋮----
// Validate Ollama configuration
⋮----
// Validate port is numeric
⋮----
// Validate chunk configuration
⋮----
// Validate reranker weight
⋮----
// Validate vector store configuration
⋮----
// GetOllamaURL returns the full Ollama URL
func (sc *ServiceConfig) GetOllamaURL() string
⋮----
// ToDocumentLoaderOptions converts the config to DocumentLoaderOptions
func (sc *ServiceConfig) ToDocumentLoaderOptions() DocumentLoaderOptions
⋮----
// Clone creates a copy of the configuration
func (sc *ServiceConfig) Clone() *ServiceConfig
⋮----
// WithProfile returns a copy of the config with the specified profile
func (sc *ServiceConfig) WithProfile(profileName string) *ServiceConfig
⋮----
// WithVectorStore returns a copy of the config with the specified vector store settings
func (sc *ServiceConfig) WithVectorStore(storeType string, host string, port int, apiKey string) *ServiceConfig
⋮----
// getEnvWithDefault gets an environment variable with a default value
func getEnvWithDefault(key, defaultValue string) string
````

## File: internal/service/document_service.go
````go
package service
⋮----
import (
	"strings"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/repository"
)
⋮----
"strings"
⋮----
"github.com/dontizi/rlama/internal/client"
"github.com/dontizi/rlama/internal/domain"
"github.com/dontizi/rlama/internal/repository"
⋮----
// DocumentService handles document management operations for RAG systems
type DocumentService interface {
	// CreateRAG creates a new RAG system with documents from the specified folder
	CreateRAG(modelName, ragName, folderPath string, options DocumentLoaderOptions) error
	
	// AddDocuments adds documents from a folder to an existing RAG system
	AddDocuments(ragName string, folderPath string, options DocumentLoaderOptions) error
	
	// GetChunks retrieves document chunks from a RAG system with optional filtering
	GetChunks(ragName string, filter ChunkFilter) ([]*domain.DocumentChunk, error)
	
	// LoadRAG loads a RAG system from the repository
	LoadRAG(ragName string) (*domain.RagSystem, error)
	
	// UpdateRAG saves changes to a RAG system
	UpdateRAG(rag *domain.RagSystem) error
	
	// ListRAGs returns all available RAG system names
	ListRAGs() ([]string, error)
}
⋮----
// CreateRAG creates a new RAG system with documents from the specified folder
⋮----
// AddDocuments adds documents from a folder to an existing RAG system
⋮----
// GetChunks retrieves document chunks from a RAG system with optional filtering
⋮----
// LoadRAG loads a RAG system from the repository
⋮----
// UpdateRAG saves changes to a RAG system
⋮----
// ListRAGs returns all available RAG system names
⋮----
// DocumentServiceImpl implements the DocumentService interface
type DocumentServiceImpl struct {
	documentLoader   *DocumentLoader
	embeddingService *EmbeddingService
	ragRepository    *repository.RagRepository
}
⋮----
// NewDocumentService creates a new DocumentService instance
func NewDocumentService(llmClient client.LLMClient) DocumentService
⋮----
// CreateRAG implements DocumentService.CreateRAG
func (ds *DocumentServiceImpl) CreateRAG(modelName, ragName, folderPath string, options DocumentLoaderOptions) error
⋮----
// Load documents from folder
⋮----
// Create new RAG system
⋮----
// Generate embeddings for all documents
⋮----
// Save the RAG system
⋮----
// AddDocuments implements DocumentService.AddDocuments
func (ds *DocumentServiceImpl) AddDocuments(ragName string, folderPath string, options DocumentLoaderOptions) error
⋮----
// Load existing RAG
⋮----
// Load new documents
⋮----
// Add documents to RAG
⋮----
// Generate embeddings for new documents
⋮----
// Save updated RAG
⋮----
// GetChunks implements DocumentService.GetChunks
func (ds *DocumentServiceImpl) GetChunks(ragName string, filter ChunkFilter) ([]*domain.DocumentChunk, error)
⋮----
var filteredChunks []*domain.DocumentChunk
⋮----
// LoadRAG implements DocumentService.LoadRAG
func (ds *DocumentServiceImpl) LoadRAG(ragName string) (*domain.RagSystem, error)
⋮----
// UpdateRAG implements DocumentService.UpdateRAG
func (ds *DocumentServiceImpl) UpdateRAG(rag *domain.RagSystem) error
⋮----
// ListRAGs implements DocumentService.ListRAGs
func (ds *DocumentServiceImpl) ListRAGs() ([]string, error)
⋮----
// generateEmbeddings generates embeddings for all chunks in the RAG system
func (ds *DocumentServiceImpl) generateEmbeddings(rag *domain.RagSystem, modelName string) error
⋮----
// Create chunker service with default options since RAG doesn't store these directly
⋮----
ChunkSize:    1000, // Use sensible defaults
⋮----
// Generate chunks for all documents
var allChunks []*domain.DocumentChunk
⋮----
// Generate embeddings for all chunks
⋮----
// Add chunks to RAG
⋮----
// matchesFilter checks if a chunk matches the given filter criteria
func (ds *DocumentServiceImpl) matchesFilter(chunk *domain.DocumentChunk, filter ChunkFilter) bool
````

## File: internal/service/document_temp_test.go
````go
package service
⋮----
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dontizi/rlama/internal/domain"
)
⋮----
"os"
"path/filepath"
"testing"
⋮----
"github.com/dontizi/rlama/internal/domain"
⋮----
func TestCreateTempDirForDocuments(t *testing.T)
⋮----
// Create test documents
⋮----
// Create temporary directory
⋮----
// Verify directory exists
⋮----
// Verify files were created
⋮----
// Verify file contents
⋮----
// Clean up
⋮----
// Verify directory was cleaned up
⋮----
func TestCleanupTempDirWithEmpty(t *testing.T)
⋮----
// Test cleanup with empty path (should not panic)
````

## File: internal/service/integration_test.go
````go
package service
⋮----
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dontizi/rlama/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)
⋮----
"os"
"path/filepath"
"testing"
⋮----
"github.com/dontizi/rlama/internal/client"
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"
⋮----
func TestONNXRerankerIntegration(t *testing.T)
⋮----
// Skip test if model directory doesn't exist
⋮----
// Validate configuration
⋮----
// Check that DocumentLoaderOptions includes ONNX settings
⋮----
// Verify that the query service was created successfully
⋮----
// Check if it's using ONNX (this is indirect since we can't easily check the internal state)
⋮----
// Verify service creation was successful
⋮----
// Test creating a RAG service for a model
⋮----
func TestRerankerServiceInterface(t *testing.T)
⋮----
// Test that both implementations satisfy the RerankerClient interface
⋮----
var _ RerankerClient = pythonClient
⋮----
var _ RerankerClient = onnxClient
var _ CleanupableRerankerClient = onnxClient
⋮----
// Test with ONNX reranker (which needs cleanup)
⋮----
// Cleanup should not error
⋮----
// Test with Python reranker (which doesn't need cleanup)
⋮----
// Cleanup should not error even when not needed
````

## File: internal/service/migration_service.go
````go
package service
⋮----
import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/repository"
	"github.com/dontizi/rlama/pkg/vector"
	"github.com/dontizi/rlama/internal/config"
)
⋮----
"fmt"
"os"
"path/filepath"
"time"
⋮----
"github.com/dontizi/rlama/internal/domain"
"github.com/dontizi/rlama/internal/repository"
"github.com/dontizi/rlama/pkg/vector"
"github.com/dontizi/rlama/internal/config"
⋮----
// MigrationOptions contains configuration for RAG migrations
type MigrationOptions struct {
	// Target configuration
	TargetVectorStore    string
	QdrantHost           string
	QdrantPort           int
	QdrantAPIKey         string
	QdrantCollectionName string
	QdrantGRPC           bool

	// Migration control
	CreateBackup         bool
	BackupPath           string
	VerifyAfterMigration bool
	DeleteOldData        bool
}
⋮----
// Target configuration
⋮----
// Migration control
⋮----
// MigrationService handles RAG system migrations between vector stores
type MigrationService struct {
	ragRepository *repository.RagRepository
}
⋮----
// NewMigrationService creates a new migration service
func NewMigrationService() *MigrationService
⋮----
// MigrateToQdrant migrates a RAG from internal storage to Qdrant
func (ms *MigrationService) MigrateToQdrant(ragName string, opts MigrationOptions) error
⋮----
// Step 1: Load existing RAG
⋮----
// Check if already using Qdrant
⋮----
// Step 2: Create backup if requested
⋮----
// Step 3: Extract vectors and metadata from current store
⋮----
// Step 4: Create new Qdrant-based hybrid store
⋮----
// Step 5: Transfer all vectors to Qdrant
⋮----
// Step 6: Update RAG configuration
⋮----
// Step 7: Verify migration if requested
⋮----
// Step 8: Save updated RAG
⋮----
// Step 9: Clean up old data if requested
⋮----
// MigrateToInternal migrates a RAG from Qdrant to internal storage
func (ms *MigrationService) MigrateToInternal(ragName string, opts MigrationOptions) error
⋮----
// Check if already using internal storage
⋮----
// Step 3: Extract vectors from Qdrant
⋮----
// Step 4: Create new internal hybrid store
⋮----
// Step 5: Transfer all vectors to internal storage
⋮----
// VectorData represents a vector with its associated metadata
type VectorData struct {
	ID       string
	Vector   []float32
	Content  string
	Metadata string
}
⋮----
// extractVectorsFromRAG extracts all vectors and metadata from a RAG system
func (ms *MigrationService) extractVectorsFromRAG(rag *domain.RagSystem) ([]VectorData, error)
⋮----
var vectors []VectorData
⋮----
// Extract vectors from all chunks
⋮----
// transferVectorsToStore transfers vectors to the target hybrid store
func (ms *MigrationService) transferVectorsToStore(vectors []VectorData, store *vector.EnhancedHybridStore) error
⋮----
// createBackup creates a backup of the RAG before migration
func (ms *MigrationService) createBackup(ragName, backupPath string) error
⋮----
// Create backup directory
⋮----
// For now, this is a placeholder - in a full implementation, we would:
// 1. Copy the RAG's info.json file
// 2. Copy vector storage files (if internal)
// 3. Export Qdrant data (if Qdrant)
// 4. Create a manifest file with backup info
⋮----
// verifyMigration verifies that the migration was successful
func (ms *MigrationService) verifyMigration(rag *domain.RagSystem, expectedVectorCount int) error
⋮----
// Perform a simple search to verify the vector store is working
⋮----
// Use the first chunk's embedding as a test query
⋮----
// Perform a search
⋮----
// Check that we can find the exact vector we searched for
⋮----
// deleteOldInternalData removes old internal vector files
func (ms *MigrationService) deleteOldInternalData(ragName string) error
````

## File: internal/service/provider.go
````go
package service
⋮----
import (
	"fmt"
	"sync"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/repository"
)
⋮----
"fmt"
"sync"
⋮----
"github.com/dontizi/rlama/internal/client"
"github.com/dontizi/rlama/internal/repository"
⋮----
// ServiceProvider manages the creation and lifecycle of all services
// It implements dependency injection pattern for centralized service management
type ServiceProvider struct {
	config *ServiceConfig
	
	// Cached clients (lazy initialization)
	ollamaClient     *client.OllamaClient
	llmClient        client.LLMClient
	clientMutex      sync.RWMutex
	
	// Cached services (lazy initialization)  
	documentService  DocumentService
	queryService     QueryService
	watchService     WatchService
	ragService       RagService
	serviceMutex     sync.RWMutex
	
	// Repositories
	ragRepository    *repository.RagRepository
	
	// Service factories (for testing/mocking)
	documentServiceFactory func(client.LLMClient) DocumentService
	queryServiceFactory    func(client.LLMClient, *client.OllamaClient, DocumentService) QueryService
	watchServiceFactory    func(DocumentService, RagService) WatchService
}
⋮----
// Cached clients (lazy initialization)
⋮----
// Cached services (lazy initialization)
⋮----
// Repositories
⋮----
// Service factories (for testing/mocking)
⋮----
// NewServiceProvider creates a new service provider with the given configuration
func NewServiceProvider(config *ServiceConfig) (*ServiceProvider, error)
⋮----
// Default factories
⋮----
// GetConfig returns a copy of the current configuration
func (sp *ServiceProvider) GetConfig() *ServiceConfig
⋮----
// UpdateConfig updates the service provider configuration and clears cached services
func (sp *ServiceProvider) UpdateConfig(config *ServiceConfig) error
⋮----
// Clear cached clients and services to force recreation with new config
⋮----
// GetOllamaClient returns the Ollama client (cached after first creation)
func (sp *ServiceProvider) GetOllamaClient() *client.OllamaClient
⋮----
// Double-check after acquiring write lock
⋮----
// GetLLMClient returns the appropriate LLM client based on configuration and model
func (sp *ServiceProvider) GetLLMClient(modelName string) (client.LLMClient, error)
⋮----
// For profile-based clients, create fresh instances
⋮----
// For direct OpenAI models
⋮----
// Default to Ollama client
⋮----
// GetDocumentService returns the document service (cached after first creation)
func (sp *ServiceProvider) GetDocumentService() DocumentService
⋮----
// Create LLM client for embeddings
⋮----
// Fallback to Ollama client
⋮----
// GetQueryService returns the query service (cached after first creation)
func (sp *ServiceProvider) GetQueryService() QueryService
⋮----
// Create dependencies
⋮----
// GetWatchService returns the watch service (cached after first creation)
func (sp *ServiceProvider) GetWatchService() WatchService
⋮----
// GetEmbeddingService returns the embedding service
func (sp *ServiceProvider) GetEmbeddingService() *EmbeddingService
⋮----
// GetRagService returns the composite RAG service (cached after first creation)
func (sp *ServiceProvider) GetRagService() RagService
⋮----
// CreateRagServiceForModel creates a RAG service configured for a specific model
func (sp *ServiceProvider) CreateRagServiceForModel(modelName string) (RagService, error)
⋮----
// Use configuration-aware service if ONNX reranker is enabled
⋮----
// SetDocumentServiceFactory allows injecting a custom document service factory (for testing)
func (sp *ServiceProvider) SetDocumentServiceFactory(factory func(client.LLMClient) DocumentService)
⋮----
sp.documentService = nil // Clear cached service
⋮----
// SetQueryServiceFactory allows injecting a custom query service factory (for testing)
func (sp *ServiceProvider) SetQueryServiceFactory(factory func(client.LLMClient, *client.OllamaClient, DocumentService) QueryService)
⋮----
sp.queryService = nil // Clear cached service
⋮----
// SetWatchServiceFactory allows injecting a custom watch service factory (for testing)
func (sp *ServiceProvider) SetWatchServiceFactory(factory func(DocumentService, RagService) WatchService)
⋮----
sp.watchService = nil // Clear cached service
⋮----
// Reset clears all cached services and clients (useful for testing)
func (sp *ServiceProvider) Reset()
````

## File: internal/service/query_service.go
````go
package service
⋮----
import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/domain"
)
⋮----
"fmt"
"math"
"sort"
"strings"
⋮----
"github.com/dontizi/rlama/internal/client"
"github.com/dontizi/rlama/internal/domain"
⋮----
// ChunkSearchResult wraps a document chunk with its similarity score
type ChunkSearchResult struct {
	Chunk *domain.DocumentChunk
	Score float64
}
⋮----
// QueryService handles search and retrieval operations for RAG systems
type QueryService interface {
	// Query processes a search query against a RAG system and returns a response
	Query(rag *domain.RagSystem, query string, contextSize int) (string, error)
	
	// SetPreferredEmbeddingModel sets the preferred model for generating embeddings
	SetPreferredEmbeddingModel(model string)
	
	// UpdateRerankerModel updates the reranker model for a specific RAG system
	UpdateRerankerModel(ragName string, model string) error
}
⋮----
// Query processes a search query against a RAG system and returns a response
⋮----
// SetPreferredEmbeddingModel sets the preferred model for generating embeddings
⋮----
// UpdateRerankerModel updates the reranker model for a specific RAG system
⋮----
// QueryServiceImpl implements the QueryService interface
type QueryServiceImpl struct {
	embeddingService *EmbeddingService
	rerankerService  *RerankerService
	documentService  DocumentService
	llmClient        client.LLMClient
	ollamaClient     *client.OllamaClient
}
⋮----
// NewQueryService creates a new QueryService instance
func NewQueryService(llmClient client.LLMClient, ollamaClient *client.OllamaClient, documentService DocumentService) QueryService
⋮----
// NewQueryServiceWithConfig creates a new QueryService instance with configuration options
func NewQueryServiceWithConfig(llmClient client.LLMClient, ollamaClient *client.OllamaClient, documentService DocumentService, config *ServiceConfig) QueryService
⋮----
var rerankerService *RerankerService
⋮----
// Query implements QueryService.Query
func (qs *QueryServiceImpl) Query(rag *domain.RagSystem, query string, contextSize int) (string, error)
⋮----
// Generate embedding for the query
⋮----
// Search for similar chunks
⋮----
// Apply reranking if enabled
⋮----
// Log the error but continue without reranking
⋮----
// Build context from chunks
⋮----
// Generate response using LLM
⋮----
// SetPreferredEmbeddingModel implements QueryService.SetPreferredEmbeddingModel
func (qs *QueryServiceImpl) SetPreferredEmbeddingModel(model string)
⋮----
// UpdateRerankerModel implements QueryService.UpdateRerankerModel
func (qs *QueryServiceImpl) UpdateRerankerModel(ragName string, model string) error
⋮----
// searchSimilarChunks finds chunks similar to the query embedding
func (qs *QueryServiceImpl) searchSimilarChunks(rag *domain.RagSystem, queryEmbedding []float32, limit int) ([]ChunkSearchResult, error)
⋮----
var scoredChunks []ChunkSearchResult
⋮----
// Calculate similarity scores for all chunks
⋮----
continue // Skip chunks without embeddings
⋮----
// Sort by similarity score (highest first)
⋮----
// Return top chunks
⋮----
// applyReranking applies reranking to the similar chunks
func (qs *QueryServiceImpl) applyReranking(rag *domain.RagSystem, query string, chunks []ChunkSearchResult) ([]ChunkSearchResult, error)
⋮----
// Prepare documents for reranking
var documents []string
⋮----
// Get reranker model
⋮----
rerankerModel = rag.ModelName // Fall back to main model
⋮----
// For now, simulate reranking since the actual reranking API is complex
// In a real implementation, this would use the Rerank method with proper SearchResults
var filteredResults []struct{ Index int; Score float64 }
⋮----
// Simulate reranker scores (in practice this would come from the reranker service)
⋮----
Score: 0.8, // Placeholder score
⋮----
// Combine vector and reranker scores
var rerankedChunks []ChunkSearchResult
⋮----
// Combine scores using the configured weight
⋮----
// Sort by combined score
⋮----
// buildContext constructs the context string from the selected chunks
func (qs *QueryServiceImpl) buildContext(chunks []ChunkSearchResult, maxLength int) string
⋮----
var contextParts []string
⋮----
// Truncate if necessary
⋮----
// generateResponse generates the final response using the LLM
func (qs *QueryServiceImpl) generateResponse(modelName, query, context string) (string, error)
⋮----
// calculateCosineSimilarity calculates cosine similarity between two vectors
func (qs *QueryServiceImpl) calculateCosineSimilarity(a, b []float32) float64
⋮----
var dotProduct, normA, normB float64
````

## File: internal/service/watch_service.go
````go
package service
⋮----
import (
	"fmt"

	"github.com/dontizi/rlama/internal/domain"
)
⋮----
"fmt"
⋮----
"github.com/dontizi/rlama/internal/domain"
⋮----
// WatchService handles file and web monitoring for RAG systems
type WatchService interface {
	// SetupDirectoryWatching enables monitoring of a directory for changes
	SetupDirectoryWatching(ragName string, dirPath string, watchInterval int, options DocumentLoaderOptions) error
	
	// DisableDirectoryWatching disables directory monitoring for a RAG system
	DisableDirectoryWatching(ragName string) error
	
	// CheckWatchedDirectory checks for changes in the watched directory and returns count of new documents
	CheckWatchedDirectory(ragName string) (int, error)
	
	// SetupWebWatching enables monitoring of a website for changes
	SetupWebWatching(ragName string, websiteURL string, watchInterval int, options domain.WebWatchOptions) error
	
	// DisableWebWatching disables web monitoring for a RAG system
	DisableWebWatching(ragName string) error
	
	// CheckWatchedWebsite checks for changes on the watched website and returns count of new documents
	CheckWatchedWebsite(ragName string) (int, error)
}
⋮----
// SetupDirectoryWatching enables monitoring of a directory for changes
⋮----
// DisableDirectoryWatching disables directory monitoring for a RAG system
⋮----
// CheckWatchedDirectory checks for changes in the watched directory and returns count of new documents
⋮----
// SetupWebWatching enables monitoring of a website for changes
⋮----
// DisableWebWatching disables web monitoring for a RAG system
⋮----
// CheckWatchedWebsite checks for changes on the watched website and returns count of new documents
⋮----
// WatchServiceImpl implements the WatchService interface
type WatchServiceImpl struct {
	documentService DocumentService
	ragService      RagService
	fileWatcher     *FileWatcher
	webWatcher      *WebWatcher
}
⋮----
// NewWatchService creates a new WatchService instance
func NewWatchService(documentService DocumentService, ragService RagService) WatchService
⋮----
// SetupDirectoryWatching implements WatchService.SetupDirectoryWatching
func (ws *WatchServiceImpl) SetupDirectoryWatching(ragName string, dirPath string, watchInterval int, options DocumentLoaderOptions) error
⋮----
// Load the RAG system
⋮----
// Configure directory watching
⋮----
// Save the updated RAG
⋮----
// DisableDirectoryWatching implements WatchService.DisableDirectoryWatching
func (ws *WatchServiceImpl) DisableDirectoryWatching(ragName string) error
⋮----
// Disable directory watching
⋮----
// CheckWatchedDirectory implements WatchService.CheckWatchedDirectory
func (ws *WatchServiceImpl) CheckWatchedDirectory(ragName string) (int, error)
⋮----
// Check if directory watching is enabled
⋮----
// Use the file watcher to check for changes
⋮----
// SetupWebWatching implements WatchService.SetupWebWatching
func (ws *WatchServiceImpl) SetupWebWatching(ragName string, websiteURL string, watchInterval int, options domain.WebWatchOptions) error
⋮----
// Configure web watching
⋮----
// DisableWebWatching implements WatchService.DisableWebWatching
func (ws *WatchServiceImpl) DisableWebWatching(ragName string) error
⋮----
// Disable web watching
⋮----
// CheckWatchedWebsite implements WatchService.CheckWatchedWebsite
func (ws *WatchServiceImpl) CheckWatchedWebsite(ragName string) (int, error)
⋮----
// Check if web watching is enabled
⋮----
// Use the web watcher to check for changes
````

## File: internal/util/format_test.go
````go
package util
⋮----
import "testing"
⋮----
func TestFormatSize(t *testing.T)
````

## File: internal/util/format.go
````go
package util
⋮----
import "fmt"
⋮----
// FormatSize returns a human-readable string for the size
func FormatSize(size int64) string
⋮----
const (
		_  = iota
		KB = 1 << (10 * iota)
````

## File: pkg/vector/bruteforce_vector_store.go
````go
package vector
⋮----
import (
	"encoding/gob"
	"fmt"
	"math"
	"os"
	"sort"
)
⋮----
"encoding/gob"
"fmt"
"math"
"os"
"sort"
⋮----
// Define key type for vector IDs
type vectorID string
⋮----
// BruteForceVectorStore implements a vector store using brute-force linear search
// This provides a simple, straightforward implementation without any indexing optimizations
type BruteForceVectorStore struct {
	items map[string][]float32
	dims  int
}
⋮----
// Ensure BruteForceVectorStore implements VectorStoreInterface
var _ VectorStoreInterface = (*BruteForceVectorStore)(nil)
⋮----
// NewBruteForceVectorStore creates a new vector store
func NewBruteForceVectorStore(dimensions int) *BruteForceVectorStore
⋮----
// Add adds a vector to the store
func (s *BruteForceVectorStore) Add(id string, vector []float32)
⋮----
// Store vector in items map
⋮----
// Remove removes a vector from the store
func (s *BruteForceVectorStore) Remove(id string)
⋮----
// Remove from items map
⋮----
// computeCosineSimilarity calculates cosine similarity between two vectors
func computeCosineSimilarity(a, b []float32) float64
⋮----
// Check for empty vectors to prevent index out of range errors
⋮----
// Check for length mismatch
⋮----
// Log the error but return a default value instead of panicking
⋮----
var dotProduct float64
var normA float64
var normB float64
⋮----
// Handle the case where one of the norms is zero
⋮----
// Search searches for similar vectors using cosine similarity
func (s *BruteForceVectorStore) Search(query []float32, limit int) []SearchResult
⋮----
// Compute similarity for all vectors
⋮----
// Sort by similarity score in descending order
⋮----
// Limit results
⋮----
// Save saves the vector store to disk
func (s *BruteForceVectorStore) Save(path string) error
⋮----
// Load loads the vector store from disk
func (s *BruteForceVectorStore) Load(path string) error
⋮----
// File doesn't exist, use empty storage
````

## File: pkg/vector/internal_vector_store.go
````go
package vector
⋮----
import (
	"encoding/gob"
	"fmt"
	"math"
	"os"
	"sort"
	"sync"
)
⋮----
"encoding/gob"
"fmt"
"math"
"os"
"sort"
"sync"
⋮----
// InternalVectorStore implements VectorStoreInterface using a thread-safe brute-force search.
// This provides an internal, memory-based vector store with optimized linear search.
type InternalVectorStore struct {
	vectors map[string][]float32
	dims    int
	mutex   sync.RWMutex
}
⋮----
// Ensure InternalVectorStore implements VectorStoreInterface
var _ VectorStoreInterface = (*InternalVectorStore)(nil)
⋮----
// NewInternalVectorStore creates a new internal vector store.
func NewInternalVectorStore(dimensions int) *InternalVectorStore
⋮----
// Add inserts or updates a vector in the store.
func (s *InternalVectorStore) Add(id string, vector []float32)
⋮----
// Store the vector
⋮----
// computeCosineSimilarityOptimized calculates cosine similarity between two vectors
func computeCosineSimilarityOptimized(a, b []float32) float64
⋮----
var dotProduct, normA, normB float64
⋮----
// Search finds the k-nearest neighbors for a query vector.
func (s *InternalVectorStore) Search(query []float32, limit int) []SearchResult
⋮----
// Compute similarity for all vectors
⋮----
// Sort by similarity score in descending order
⋮----
// Limit results
⋮----
// Remove removes a vector from the store.
func (s *InternalVectorStore) Remove(id string)
⋮----
// Save persists the vector store to disk.
func (s *InternalVectorStore) Save(path string) error
⋮----
// Load reconstructs the vector store from disk.
func (s *InternalVectorStore) Load(path string) error
⋮----
// Check if file exists
⋮----
// File doesn't exist, initialize empty store
⋮----
var saveData struct {
		Vectors map[string][]float32
		Dims    int
	}
````

## File: pkg/vector/qdrant_store.go
````go
package vector
⋮----
import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)
⋮----
"context"
"crypto/tls"
"fmt"
"time"
⋮----
"github.com/qdrant/go-client/qdrant"
"google.golang.org/grpc"
"google.golang.org/grpc/credentials"
"google.golang.org/grpc/credentials/insecure"
"google.golang.org/protobuf/proto"
⋮----
// QdrantStore implements VectorStoreInterface using Qdrant vector database
type QdrantStore struct {
	client         qdrant.PointsClient
	collections    qdrant.CollectionsClient
	conn           *grpc.ClientConn
	collectionName string
	dims           uint64
	timeout        time.Duration
}
⋮----
// Ensure QdrantStore implements VectorStoreInterface
var _ VectorStoreInterface = (*QdrantStore)(nil)
⋮----
// NewQdrantStore creates and configures a new Qdrant client and store
func NewQdrantStore(host string, port int, collectionName string, dims int, apiKey string, useGRPC bool, createCollectionIfNotExists bool) (*QdrantStore, error)
⋮----
var conn *grpc.ClientConn
var err error
⋮----
// Setup gRPC connection options
var dialOpts []grpc.DialOption
⋮----
// For Qdrant Cloud or secured instances, typically use TLS
⋮----
// ensureCollectionExists creates the collection if it doesn't exist
func (s *QdrantStore) ensureCollectionExists(ctx context.Context) error
⋮----
// Add implements VectorStoreInterface - adds a vector without payload
func (s *QdrantStore) Add(id string, vector []float32)
⋮----
// UpsertPointWithPayload adds or updates a point with its vector and payload
func (s *QdrantStore) UpsertPointWithPayload(id string, vector []float32, payload map[string]interface
⋮----
// Search performs a vector search in Qdrant
func (s *QdrantStore) Search(query []float32, limit int) []SearchResult
⋮----
// SearchWithFilter performs a vector search with optional payload filtering
func (s *QdrantStore) SearchWithFilter(query []float32, limit int, filter *qdrant.Filter) []SearchResult
⋮----
var originalID string
⋮----
// Remove deletes a point from Qdrant
func (s *QdrantStore) Remove(id string)
⋮----
// Save is a no-op for QdrantStore as Qdrant server handles persistence
func (s *QdrantStore) Save(path string) error
⋮----
// Load is a no-op for QdrantStore as connection is established at construction
func (s *QdrantStore) Load(path string) error
⋮----
// Close closes the gRPC connection to Qdrant
func (s *QdrantStore) Close() error
````

## File: scripts/build.sh
````bash
#!/bin/bash
# Build script for RLAMA

VERSION=$(grep "Version = " cmd/root.go | cut -d'"' -f2)
PLATFORMS=("windows/amd64" "windows/386" "darwin/amd64" "darwin/arm64" "linux/amd64" "linux/arm64")
BINARY_NAME="rlama"

echo "Building RLAMA v${VERSION}..."

rm -rf ./dist
mkdir -p ./dist

for platform in "${PLATFORMS[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$BINARY_NAME'_'$GOOS'_'$GOARCH
    
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    echo "Building for $GOOS/$GOARCH..."
    env GOOS=$GOOS GOARCH=$GOARCH go build -o ./dist/$output_name
    
    if [ $? -ne 0 ]; then
        echo "Error building for $GOOS/$GOARCH"
    else
        echo "Successfully built for $GOOS/$GOARCH"
    fi
done

echo "Creating archives..."
cd ./dist
for file in rlama_*
do
    zip "${file}.zip" "$file"
done

echo "Build completed!"
````

## File: test-small-docs/clayborn.txt
````
Clayborn Blankenship was born in 1850 in Virginia. He was known for his farming and family connections.
````

## File: .gitattributes
````
# Auto detect text files and perform LF normalization
* text=auto
````

## File: .repomixignore
````
# Add patterns to ignore here, one per line
# Example:
# *.log
# tmp/
models/
LICENSE
````

## File: install.sh
````bash
#!/bin/bash

set -e

# Colors for messages
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

echo "
█▀█ █   ▄▀█ █▀▄▀█ ▄▀█
█▀▄ █▄▄ █▀█ █░▀░█ █▀█

Retrieval-Augmented Language Model Adapter
"

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Determine OS and architecture
get_os_arch() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    # Convert architecture to standard format
    case "$arch" in
        x86_64)
            arch="amd64"
            ;;
        aarch64|arm64)
            arch="arm64"
            ;;
        *)
            echo -e "${RED}Unsupported architecture: $arch${NC}"
            exit 1
            ;;
    esac
    
    # Handle macOS naming
    if [ "$os" = "darwin" ]; then
        os="darwin"
    elif [ "$os" = "linux" ]; then
        os="linux"
    else
        echo -e "${RED}Unsupported operating system: $os${NC}"
        exit 1
    fi
    
    echo "${os}_${arch}"
}

# Check if Ollama is installed
if ! command_exists ollama; then
    echo -e "${YELLOW}⚠️ Ollama is not installed.${NC}"
    echo "RLAMA requires Ollama to function."
    echo "You can install Ollama with:"
    echo "curl -fsSL https://ollama.com/install.sh | sh"
    
    read -p "Do you want to install Ollama now? (y/n): " install_ollama
    if [[ "$install_ollama" =~ ^[Yy]$ ]]; then
        echo "Installing Ollama..."
        curl -fsSL https://ollama.com/install.sh | sh
    else
        echo "Please install Ollama before using RLAMA."
    fi
fi

# Check if Ollama is running
if ! curl -s http://localhost:11434/api/version &>/dev/null; then
    echo -e "${YELLOW}⚠️ The Ollama service doesn't seem to be running.${NC}"
    echo "Please start Ollama before using RLAMA."
fi

# Check if the llama3 model is available
if command_exists ollama; then
    if ! ollama list 2>/dev/null | grep -q "llama3"; then
        echo -e "${YELLOW}⚠️ The llama3 model is not available in Ollama.${NC}"
        echo "For a better experience, you should install it with:"
        echo "ollama pull llama3"
    fi
fi

# Create installation directory
INSTALL_DIR="/usr/local/bin"
DATA_DIR="$HOME/.rlama"

echo "Installing RLAMA..."

# Determine OS and architecture for downloading the correct binary
OS_ARCH=$(get_os_arch)
BINARY_NAME="rlama_${OS_ARCH}"
DOWNLOAD_URL="https://github.com/dontizi/rlama/releases/latest/download/${BINARY_NAME}"

echo "Downloading RLAMA for $OS_ARCH..."

# Create a temporary directory for downloading
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

# Download the binary
if command_exists curl; then
    curl -L -o rlama $DOWNLOAD_URL
elif command_exists wget; then
    wget -O rlama $DOWNLOAD_URL
else
    echo -e "${RED}Error: Neither curl nor wget is installed.${NC}"
    exit 1
fi

# Make it executable
chmod +x rlama

# Install
echo "Installing executable..."
mkdir -p "$DATA_DIR"

# Try to install to /usr/local/bin, fall back to ~/.local/bin if permission denied
if [ -w "$INSTALL_DIR" ]; then
    mv rlama "$INSTALL_DIR/"
else
    echo "Cannot write to $INSTALL_DIR, trying alternative location..."
    LOCAL_BIN="$HOME/.local/bin"
    mkdir -p "$LOCAL_BIN"
    mv rlama "$LOCAL_BIN/"
    
    # Add to PATH if not already there
    if [[ ":$PATH:" != *":$LOCAL_BIN:"* ]]; then
        echo "Adding $LOCAL_BIN to your PATH..."
        echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.bashrc"
        echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.zshrc" 2>/dev/null || true
        export PATH="$HOME/.local/bin:$PATH"
    fi
    
    INSTALL_DIR="$LOCAL_BIN"
fi

# Clean up
cd - > /dev/null
rm -rf "$TEMP_DIR"

echo -e "${GREEN}RLAMA has been successfully installed to $INSTALL_DIR/rlama!${NC}"
echo "You can now use RLAMA by running the 'rlama' command."
echo "Run 'rlama --help' to see available commands."
echo ""
echo "You can also use RLAMA with the following commands:"
echo "- rlama rag [model] [rag-name] [folder-path] : Create a new RAG system"
echo "- rlama run [rag-name] : Run a RAG system"
echo "- rlama list : List all available RAG systems"
echo "- rlama delete [rag-name] : Delete a RAG system"
echo ""
echo "Example: rlama rag llama3 myrag ./documents"
````

## File: repomix.config.json
````json
{
  "$schema": "https://repomix.com/schemas/latest/schema.json",
  "input": {
    "maxFileSize": 52428800
  },
  "output": {
    "filePath": "repomix-rlama.md",
    "style": "markdown",
    "parsableStyle": false,
    "fileSummary": true,
    "directoryStructure": true,
    "files": true,
    "removeComments": false,
    "removeEmptyLines": false,
    "compress": false,
    "topFilesLength": 5,
    "showLineNumbers": false,
    "copyToClipboard": false,
    "git": {
      "sortByChanges": true,
      "sortByChangesMaxCommits": 100,
      "includeDiffs": false
    }
  },
  "include": [],
  "ignore": {
    "useGitignore": true,
    "useDefaultPatterns": true,
    "customPatterns": []
  },
  "security": {
    "enableSecurityCheck": true
  },
  "tokenCount": {
    "encoding": "o200k_base"
  }
}
````

## File: cmd/add_reranker.go
````go
package cmd
⋮----
import (
	"fmt"

	"github.com/dontizi/rlama/internal/service"
	"github.com/spf13/cobra"
)
⋮----
"fmt"
⋮----
"github.com/dontizi/rlama/internal/service"
"github.com/spf13/cobra"
⋮----
var (
	rerankerModel     string
	rerankerWeight    float64
	rerankerThreshold float64
	rerankerTopK      int
	disableReranker   bool
	rerankerSilent    bool
)
⋮----
var addRerankerCmd = &cobra.Command{
	Use:   "add-reranker [rag-name]",
	Short: "Configure reranking for a RAG system",
	Long: `Configure reranking settings for a RAG system (note: reranking is enabled by default).
Example: rlama add-reranker my-rag --model reranker-model

Reranking improves retrieval accuracy by applying a second-stage ranking to initial search results.
This uses a cross-encoder approach to evaluate the relevance of each document to the query.

Use --disable flag to turn off reranking if needed.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]

		// Get Ollama client from root command
		ollamaClient := GetOllamaClient()

		// Create RAG service
		ragService := service.NewRagService(ollamaClient)

		// Load the RAG
		rag, err := ragService.LoadRag(ragName)
		if err != nil {
			return fmt.Errorf("error loading RAG: %w", err)
		}

		if disableReranker {
			// Disable reranking
			rag.RerankerEnabled = false
			fmt.Printf("Reranking disabled for RAG '%s'\n", ragName)
		} else {
			// Enable reranking with specified settings
			rag.RerankerEnabled = true

			if rerankerModel != "" {
				rag.RerankerModel = rerankerModel
			} else if rag.RerankerModel == "" {
				// If not set, use the same model as the RAG
				rag.RerankerModel = rag.ModelName
			}

			// Set weight if specified
			if cmd.Flags().Changed("weight") {
				rag.RerankerWeight = rerankerWeight
			} else if rag.RerankerWeight == 0 {
				// Set default if not already set
				rag.RerankerWeight = 0.7
			}

			// Set threshold if specified
			if cmd.Flags().Changed("threshold") {
				rag.RerankerThreshold = rerankerThreshold
			}

			// Set max results to return if specified
			if cmd.Flags().Changed("topk") {
				rag.RerankerTopK = rerankerTopK
			} else if rag.RerankerTopK == 0 {
				// Set default if not already set
				rag.RerankerTopK = 5
			}
			
			// Set silent mode if specified
			if cmd.Flags().Changed("silent") {
				rag.RerankerSilent = rerankerSilent
			}

			fmt.Printf("Reranking enabled for RAG '%s'\n", ragName)
			fmt.Printf("  Model: %s\n", rag.RerankerModel)
			fmt.Printf("  Weight: %.2f\n", rag.RerankerWeight)
			fmt.Printf("  Threshold: %.2f\n", rag.RerankerThreshold)
			fmt.Printf("  Max results: %d\n", rag.RerankerTopK)
			if rag.RerankerSilent {
				fmt.Printf("  Silent mode: enabled (warnings and info messages suppressed)\n")
			}
		}

		// Update the RAG
		err = ragService.UpdateRag(rag)
		if err != nil {
			return fmt.Errorf("error updating RAG: %w", err)
		}

		return nil
	},
}
⋮----
// Get Ollama client from root command
⋮----
// Create RAG service
⋮----
// Load the RAG
⋮----
// Disable reranking
⋮----
// Enable reranking with specified settings
⋮----
// If not set, use the same model as the RAG
⋮----
// Set weight if specified
⋮----
// Set default if not already set
⋮----
// Set threshold if specified
⋮----
// Set max results to return if specified
⋮----
// Set silent mode if specified
⋮----
// Update the RAG
⋮----
func init()
⋮----
// Configure reranking options
````

## File: cmd/api.go
````go
package cmd
⋮----
import (
	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/server"
)
⋮----
"github.com/spf13/cobra"
"github.com/dontizi/rlama/internal/server"
⋮----
var (
	apiPort string
)
⋮----
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the RLAMA API server",
	Long: `Start an HTTP API server for RLAMA, allowing RAG operations via RESTful endpoints.
	
Example: rlama api --port 11249

Available endpoints:
- POST /rag: Query a RAG system
  Request body: { "rag_name": "my-docs", "prompt": "How many documents do you have?", "context_size": 20 }
  
- GET /health: Check server health`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get Ollama client with configured host and port
		ollamaClient := GetOllamaClient()
		
		// Create and start the server
		srv := server.NewServer(apiPort, ollamaClient)
		return srv.Start()
	},
}
⋮----
// Get Ollama client with configured host and port
⋮----
// Create and start the server
⋮----
func init()
⋮----
// Add port flag
````

## File: cmd/chunk_eval.go
````go
package cmd
⋮----
import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/service"
	"github.com/spf13/cobra"
)
⋮----
"fmt"
"os"
"path/filepath"
"time"
⋮----
"github.com/dontizi/rlama/internal/domain"
"github.com/dontizi/rlama/internal/service"
"github.com/spf13/cobra"
⋮----
var (
	targetFile      string
	outputDetailed  bool
	compareAll      bool
	customChunkSize int
	customOverlap   int
	customStrategy  string
)
⋮----
// chunkEvalCmd represents the command to evaluate chunking strategies
var chunkEvalCmd = &cobra.Command{
	Use:   "chunk-eval",
	Short: "Evaluate and optimize chunking strategies for different documents",
	Long: `Evaluate and compare different chunking strategies for a given document.
This command allows you to:
- Test a specific chunking configuration on a document
- Automatically compare multiple strategies to find the best one
- Get detailed metrics on chunking quality

Examples:
  rlama chunk-eval --file=document.md
  rlama chunk-eval --file=code.go --strategy=semantic --size=1000 --overlap=100
  rlama chunk-eval --file=document.txt --compare-all --detailed`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if file exists
		if targetFile == "" {
			return fmt.Errorf("please specify a file with --file")
		}

		fileInfo, err := os.Stat(targetFile)
		if err != nil {
			return fmt.Errorf("error accessing file: %w", err)
		}

		if fileInfo.IsDir() {
			return fmt.Errorf("the specified path is a directory, not a file")
		}

		// Load file content
		content, err := os.ReadFile(targetFile)
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}

		// Create document from file
		doc := &domain.Document{
			ID:      filepath.Base(targetFile),
			Name:    filepath.Base(targetFile),
			Path:    targetFile,
			Content: string(content),
		}

		// Create evaluator
		chunkerService := service.NewChunkerService(service.DefaultChunkingConfig())
		evaluator := service.NewChunkingEvaluator(chunkerService)

		fmt.Printf("Analyzing document: %s (%d characters)\n", doc.Name, len(doc.Content))

		// If --compare-all is specified, compare all strategies
		if compareAll {
			fmt.Println("\nComparing all available chunking strategies...")
			startTime := time.Now()

			results := evaluator.CompareChunkingStrategies(doc)

			fmt.Printf("\nAnalysis completed in %.2f seconds\n", time.Since(startTime).Seconds())

			// Display top 5 strategies
			fmt.Println("\n=== Top 5 strategies for this document ===")
			fmt.Println("Rank | Strategy        | Size   | Overlap | Score  | Chunks | Time (ms)")
			fmt.Println("-----|----------------|--------|---------|--------|--------|----------")

			limit := 5
			if len(results) < limit {
				limit = len(results)
			}

			for i := 0; i < limit; i++ {
				fmt.Printf("%4d | %-15s | %6d | %7d | %.4f | %6d | %6d\n",
					i+1,
					results[i].Strategy,
					results[i].ChunkSize,
					results[i].ChunkOverlap,
					results[i].SemanticCoherenceScore,
					results[i].TotalChunks,
					results[i].ProcessingTimeMs)
			}

			// Show details of the best strategy
			if len(results) > 0 && outputDetailed {
				fmt.Println("\nDetails of the best strategy:")
				evaluator.PrintEvaluationResults(results[0])
			}

			// Recommended configuration
			if len(results) > 0 {
				best := results[0]
				fmt.Printf("\nRecommended configuration for this document:\n")
				fmt.Printf("  --strategy=%s --size=%d --overlap=%d\n",
					best.Strategy, best.ChunkSize, best.ChunkOverlap)
			}

			return nil
		}

		// Otherwise, evaluate a specific configuration
		config := service.DefaultChunkingConfig()

		// Use custom parameters if specified
		if customStrategy != "" {
			config.ChunkingStrategy = customStrategy
		}

		if customChunkSize > 0 {
			config.ChunkSize = customChunkSize
		}

		if customOverlap >= 0 {
			config.ChunkOverlap = customOverlap
		}

		fmt.Printf("\nEvaluating strategy: %s (size: %d, overlap: %d)\n",
			config.ChunkingStrategy, config.ChunkSize, config.ChunkOverlap)

		// Evaluate the strategy
		metrics := evaluator.EvaluateChunkingStrategy(doc, config)

		// Display results
		if outputDetailed {
			evaluator.PrintEvaluationResults(metrics)
		} else {
			// Simplified display
			fmt.Println("\n=== Evaluation Results ===")
			fmt.Printf("Coherence score: %.4f\n", metrics.SemanticCoherenceScore)
			fmt.Printf("Number of chunks: %d\n", metrics.TotalChunks)
			fmt.Printf("Average chunk size: %.2f characters\n", metrics.AverageChunkSize)
			fmt.Printf("Split sentences: %d (%.1f%%)\n",
				metrics.ChunksWithSplitSentences,
				float64(metrics.ChunksWithSplitSentences)/float64(metrics.TotalChunks)*100)
			fmt.Printf("Split paragraphs: %d (%.1f%%)\n",
				metrics.ChunksWithSplitParagraphs,
				float64(metrics.ChunksWithSplitParagraphs)/float64(metrics.TotalChunks)*100)
			fmt.Printf("Redundancy rate: %.1f%%\n", metrics.RedundancyRate*100)
		}

		// Suggest other strategies if score is low
		if metrics.SemanticCoherenceScore < 0.7 {
			fmt.Println("\nThe coherence score is relatively low. Try comparing other strategies with:")
			fmt.Printf("  rlama chunk-eval --file=%s --compare-all\n", targetFile)
		}

		return nil
	},
}
⋮----
// Check if file exists
⋮----
// Load file content
⋮----
// Create document from file
⋮----
// Create evaluator
⋮----
// If --compare-all is specified, compare all strategies
⋮----
// Display top 5 strategies
⋮----
// Show details of the best strategy
⋮----
// Recommended configuration
⋮----
// Otherwise, evaluate a specific configuration
⋮----
// Use custom parameters if specified
⋮----
// Evaluate the strategy
⋮----
// Display results
⋮----
// Simplified display
⋮----
// Suggest other strategies if score is low
⋮----
func init()
⋮----
// Flags
````

## File: cmd/hf_browse.go
````go
package cmd
⋮----
import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)
⋮----
"fmt"
"os/exec"
"runtime"
"strings"
⋮----
"github.com/spf13/cobra"
⋮----
var (
	browseQuantization string
	browseLimit        int
	browseOpen         bool
)
⋮----
var hfBrowseCmd = &cobra.Command{
	Use:   "hf-browse [search-term]",
	Short: "Browse GGUF models on Hugging Face",
	Long: `Search and browse GGUF models available on Hugging Face.
Optionally open the browser to view the model details.

Examples:
  rlama hf-browse llama3       # Search for Llama 3 models
  rlama hf-browse mistral --open # Search and open browser
  rlama hf-browse "open instruct" --limit 5  # Limit results`,
	RunE: func(cmd *cobra.Command, args []string) error {
		searchTerm := ""
		if len(args) > 0 {
			searchTerm = strings.Join(args, " ")
		}

		// Build the search URL
		searchURL := "https://huggingface.co/models?search=gguf"
		if searchTerm != "" {
			searchURL += "+" + strings.ReplaceAll(searchTerm, " ", "+")
		}

		if browseOpen {
			var err error
			switch runtime.GOOS {
			case "linux":
				err = exec.Command("xdg-open", searchURL).Start()
			case "windows":
				err = exec.Command("rundll32", "url.dll,FileProtocolHandler", searchURL).Start()
			case "darwin":
				err = exec.Command("open", searchURL).Start()
			default:
				err = fmt.Errorf("unsupported platform")
			}
			if err != nil {
				return fmt.Errorf("error opening browser: %w", err)
			}
			fmt.Printf("Opened browser to search for GGUF models: %s\n", searchURL)
		}

		fmt.Println("To use a Hugging Face model with RLAMA, use the format:")
		fmt.Println("  rlama rag hf.co/username/repository my-rag ./docs")
		
		if browseQuantization != "" {
			fmt.Println("\nSpecify quantization:")
			fmt.Println("  rlama rag hf.co/username/repository:" + browseQuantization + " my-rag ./docs")
		}
		
		fmt.Println("\nOr try out a model directly with:")
		fmt.Println("  rlama run-hf username/repository")
		
		return nil
	},
}
⋮----
// Build the search URL
⋮----
var err error
⋮----
func init()
````

## File: cmd/list_chunks.go
````go
package cmd
⋮----
import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/service"
)
⋮----
"fmt"
"strings"
⋮----
"github.com/spf13/cobra"
"github.com/dontizi/rlama/internal/service"
⋮----
var (
	showChunkContent bool
	documentFilter   string
)
⋮----
var listChunksCmd = &cobra.Command{
	Use:   "list-chunks [rag-name]",
	Short: "Inspect document chunks in a RAG system",
	Long: `List and filter document chunks with various options.
	
Examples:
  # Basic chunk listing
  rlama list-chunks my-docs
  
  # Show chunk contents
  rlama list-chunks my-docs --show-content
  
  # Filter chunks from API documentation
  rlama list-chunks my-docs --document=api
  
  # Combine filters
  rlama list-chunks my-docs --document=readme --show-content`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]

		ollamaClient := GetOllamaClient()
		ragService := service.NewRagService(ollamaClient)

		// Get chunks with filters
		chunks, err := ragService.GetRagChunks(ragName, service.ChunkFilter{
			DocumentSubstring: documentFilter,
			ShowContent:       showChunkContent,
		})
		if err != nil {
			return err
		}

		// Display results
		fmt.Printf("Found %d chunks in RAG '%s'\n", len(chunks), ragName)
		for _, chunk := range chunks {
			fmt.Printf("\nChunk ID: %s\n", chunk.ID)
			fmt.Printf("Document: %s\n", chunk.DocumentID)
			fmt.Printf("Position: %d/%d\n", chunk.ChunkNumber+1, chunk.TotalChunks)
			
			if showChunkContent {
				fmt.Printf("Content:\n%s\n", strings.TrimSpace(chunk.Content))
			}
		}
		return nil
	},
}
⋮----
// Get chunks with filters
⋮----
// Display results
⋮----
func init()
````

## File: cmd/list_docs.go
````go
package cmd
⋮----
import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/service"
	"github.com/dontizi/rlama/internal/util"
)
⋮----
"fmt"
"os"
"text/tabwriter"
⋮----
"github.com/spf13/cobra"
"github.com/dontizi/rlama/internal/service"
"github.com/dontizi/rlama/internal/util"
⋮----
var listDocsCmd = &cobra.Command{
	Use:   "list-docs [rag-name]",
	Short: "List all documents in a RAG system",
	Long: `Display a list of all documents in a specified RAG system.
Example: rlama list-docs my-docs`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]

		// Get Ollama client from root command
		ollamaClient := GetOllamaClient()

		// Create necessary services
		ragService := service.NewRagService(ollamaClient)

		// Load the RAG
		rag, err := ragService.LoadRag(ragName)
		if err != nil {
			return err
		}

		if len(rag.Documents) == 0 {
			fmt.Printf("No documents found in RAG '%s'.\n", ragName)
			return nil
		}

		fmt.Printf("Documents in RAG '%s' (%d found):\n\n", ragName, len(rag.Documents))
		
		// Use tabwriter for aligned display
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tPATH\tSIZE\tCONTENT TYPE")
		
		for _, doc := range rag.Documents {
			sizeStr := util.FormatSize(doc.Size)
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", doc.ID, doc.Name, sizeStr, doc.ContentType)
		}
		w.Flush()
		
		return nil
	},
}
⋮----
// Get Ollama client from root command
⋮----
// Create necessary services
⋮----
// Load the RAG
⋮----
// Use tabwriter for aligned display
⋮----
func init()
````

## File: cmd/run_hf.go
````go
package cmd
⋮----
import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)
⋮----
"fmt"
"strings"
⋮----
"github.com/spf13/cobra"
⋮----
var (
	runHfQuantization string
)
⋮----
var runHfCmd = &cobra.Command{
	Use:   "run-hf [huggingface-model]",
	Short: "Run a Hugging Face GGUF model with Ollama",
	Long: `Run a Hugging Face GGUF model directly using Ollama.
This is convenient for testing models before creating a RAG system with them.

Examples:
  rlama run-hf bartowski/Llama-3.2-1B-Instruct-GGUF
  rlama run-hf mlabonne/Meta-Llama-3.1-8B-Instruct-abliterated-GGUF --quant Q5_K_M`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		modelPath := args[0]
		
		// Prepare the model reference
		if !strings.Contains(modelPath, "/") {
			return fmt.Errorf("invalid model format. Use 'username/repository' format")
		}
		
		// Get the Ollama client
		ollamaClient := GetOllamaClient()
		
		fmt.Printf("Running Hugging Face model: %s\n", modelPath)
		if runHfQuantization != "" {
			fmt.Printf("Using quantization: %s\n", runHfQuantization)
		}
		
		return ollamaClient.RunHuggingFaceModel(modelPath, runHfQuantization)
	},
}
⋮----
// Prepare the model reference
⋮----
// Get the Ollama client
⋮----
func init()
````

## File: cmd/watch.go
````go
package cmd
⋮----
import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/service"
)
⋮----
"fmt"
"strconv"
⋮----
"github.com/spf13/cobra"
"github.com/dontizi/rlama/internal/service"
⋮----
var (
	watchExcludeDirs  []string
	watchExcludeExts  []string
	watchProcessExts  []string
	watchChunkSize    int
	watchChunkOverlap int
)
⋮----
var watchCmd = &cobra.Command{
	Use:   "watch [rag-name] [directory-path] [interval]",
	Short: "Set up directory watching for a RAG system",
	Long: `Configure a RAG system to automatically watch a directory for new files and add them to the RAG.
The interval is specified in minutes. Use 0 to only check when the RAG is used.

Example: rlama watch my-docs ./documents 60
This will check the ./documents directory every 60 minutes for new files.

Use rlama watch-off [rag-name] to disable watching.`,
	Args: cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]
		dirPath := args[1]
		
		// Default interval is 0 (only check when RAG is used)
		interval := 0
		
		// If an interval is provided, parse it
		if len(args) > 2 {
			var err error
			interval, err = strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid interval: %s", args[2])
			}
			
			if interval < 0 {
				return fmt.Errorf("interval must be non-negative")
			}
		}
		
		// Get service provider
		provider := GetServiceProvider()
		ragService := provider.GetRagService()
		
		// Set up loader options based on flags
		loaderOptions := service.DocumentLoaderOptions{
			ExcludeDirs:  watchExcludeDirs,
			ExcludeExts:  watchExcludeExts,
			ProcessExts:  watchProcessExts,
			ChunkSize:    watchChunkSize,
			ChunkOverlap: watchChunkOverlap,
		}
		
		// Set up directory watching
		err := ragService.SetupDirectoryWatching(ragName, dirPath, interval, loaderOptions)
		if err != nil {
			return err
		}
		
		// Provide feedback based on the interval
		if interval == 0 {
			fmt.Printf("Directory watching set up for RAG '%s'. Directory '%s' will be checked each time the RAG is used.\n", 
				ragName, dirPath)
		} else {
			fmt.Printf("Directory watching set up for RAG '%s'. Directory '%s' will be checked every %d minutes.\n", 
				ragName, dirPath, interval)
		}
		
		// Perform an initial check
		docsAdded, err := ragService.CheckWatchedDirectory(ragName)
		if err != nil {
			return fmt.Errorf("error during initial directory check: %w", err)
		}
		
		if docsAdded > 0 {
			fmt.Printf("Added %d new documents from '%s'.\n", docsAdded, dirPath)
		} else {
			fmt.Printf("No new documents found in '%s'.\n", dirPath)
		}
		
		return nil
	},
}
⋮----
// Default interval is 0 (only check when RAG is used)
⋮----
// If an interval is provided, parse it
⋮----
var err error
⋮----
// Get service provider
⋮----
// Set up loader options based on flags
⋮----
// Set up directory watching
⋮----
// Provide feedback based on the interval
⋮----
// Perform an initial check
⋮----
var watchOffCmd = &cobra.Command{
	Use:   "watch-off [rag-name]",
	Short: "Disable directory watching for a RAG system",
	Long:  `Disable automatic directory watching for a RAG system.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]
		
		// Get service provider
		provider := GetServiceProvider()
		ragService := provider.GetRagService()
		
		// Disable directory watching
		err := ragService.DisableDirectoryWatching(ragName)
		if err != nil {
			return err
		}
		
		fmt.Printf("Directory watching disabled for RAG '%s'.\n", ragName)
		return nil
	},
}
⋮----
// Disable directory watching
⋮----
var checkWatchedCmd = &cobra.Command{
	Use:   "check-watched [rag-name]",
	Short: "Check a RAG's watched directory for new files",
	Long:  `Manually check a RAG's watched directory for new files and add them to the RAG.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]
		
		// Get service provider
		provider := GetServiceProvider()
		ragService := provider.GetRagService()
		
		// Check the watched directory
		docsAdded, err := ragService.CheckWatchedDirectory(ragName)
		if err != nil {
			return err
		}
		
		if docsAdded > 0 {
			fmt.Printf("Added %d new documents to RAG '%s'.\n", docsAdded, ragName)
		} else {
			fmt.Printf("No new documents found for RAG '%s'.\n", ragName)
		}
		
		return nil
	},
}
⋮----
// Check the watched directory
⋮----
func init()
⋮----
// Add exclusion and processing flags
````

## File: cmd/web_watch.go
````go
package cmd
⋮----
import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/service"
)
⋮----
"fmt"
"strconv"
⋮----
"github.com/spf13/cobra"
"github.com/dontizi/rlama/internal/domain"
"github.com/dontizi/rlama/internal/service"
⋮----
var (
	webWatchMaxDepth     int
	webWatchConcurrency  int
	webWatchExcludePaths []string
	webWatchChunkSize    int
	webWatchChunkOverlap int
)
⋮----
var webWatchCmd = &cobra.Command{
	Use:   "web-watch [rag-name] [website-url] [interval]",
	Short: "Set up website watching for a RAG system",
	Long: `Configure a RAG system to automatically watch a website for new content and add it to the RAG.
The interval is specified in minutes. Use 0 to only check when the RAG is used.

Example: rlama web-watch my-docs https://docs.example.com 60
This will check the website every 60 minutes for new content.

Use rlama web-watch-off [rag-name] to disable watching.`,
	Args: cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]
		websiteURL := args[1]
		
		// Default interval is 0 (only check when RAG is used)
		interval := 0
		
		// If an interval is provided, parse it
		if len(args) > 2 {
			var err error
			interval, err = strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid interval: %s", args[2])
			}
			
			if interval < 0 {
				return fmt.Errorf("interval must be non-negative")
			}
		}
		
		// Get service provider
		provider := GetServiceProvider()
		ragService := provider.GetRagService()
		
		// Set up web watch options
		webWatchOptions := domain.WebWatchOptions{
			MaxDepth:     webWatchMaxDepth,
			Concurrency:  webWatchConcurrency,
			ExcludePaths: webWatchExcludePaths,
			ChunkSize:    webWatchChunkSize,
			ChunkOverlap: webWatchChunkOverlap,
		}
		
		// Set up website watching
		err := ragService.SetupWebWatching(ragName, websiteURL, interval, webWatchOptions)
		if err != nil {
			return err
		}
		
		// Provide feedback based on the interval
		if interval == 0 {
			fmt.Printf("Website watching set up for RAG '%s'. Website '%s' will be checked each time the RAG is used.\n", 
				ragName, websiteURL)
		} else {
			fmt.Printf("Website watching set up for RAG '%s'. Website '%s' will be checked every %d minutes.\n", 
				ragName, websiteURL, interval)
		}
		
		// Perform an initial check
		docsAdded, err := ragService.CheckWatchedWebsite(ragName)
		if err != nil {
			return fmt.Errorf("error during initial website check: %w", err)
		}
		
		if docsAdded > 0 {
			fmt.Printf("Added %d new pages from '%s'.\n", docsAdded, websiteURL)
		} else {
			fmt.Printf("No new content found at '%s'.\n", websiteURL)
		}
		
		return nil
	},
}
⋮----
// Default interval is 0 (only check when RAG is used)
⋮----
// If an interval is provided, parse it
⋮----
var err error
⋮----
// Get service provider
⋮----
// Set up web watch options
⋮----
// Set up website watching
⋮----
// Provide feedback based on the interval
⋮----
// Perform an initial check
⋮----
var webWatchOffCmd = &cobra.Command{
	Use:   "web-watch-off [rag-name]",
	Short: "Disable website watching for a RAG system",
	Long:  `Disable automatic website watching for a RAG system.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]
		
		// Get Ollama client from root command
		ollamaClient := GetOllamaClient()
		
		// Create RAG service
		ragService := service.NewRagService(ollamaClient)
		
		// Disable website watching
		err := ragService.DisableWebWatching(ragName)
		if err != nil {
			return err
		}
		
		fmt.Printf("Website watching disabled for RAG '%s'.\n", ragName)
		return nil
	},
}
⋮----
// Get Ollama client from root command
⋮----
// Create RAG service
⋮----
// Disable website watching
⋮----
var checkWebWatchedCmd = &cobra.Command{
	Use:   "check-web-watched [rag-name]",
	Short: "Check a RAG's watched website for new content",
	Long:  `Manually check a RAG's watched website for new content and add it to the RAG.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]
		
		// Get Ollama client from root command
		ollamaClient := GetOllamaClient()
		
		// Create RAG service
		ragService := service.NewRagService(ollamaClient)
		
		// Check the watched website
		pagesAdded, err := ragService.CheckWatchedWebsite(ragName)
		if err != nil {
			return err
		}
		
		if pagesAdded > 0 {
			fmt.Printf("Added %d new pages to RAG '%s'.\n", pagesAdded, ragName)
		} else {
			fmt.Printf("No new content found for RAG '%s'.\n", ragName)
		}
		
		return nil
	},
}
⋮----
// Check the watched website
⋮----
func init()
⋮----
// Add web watching specific flags
````

## File: docs/chunking_guidelines.md
````markdown
# Optimization Guide for Chunking Strategies

This document provides guidelines to optimize chunking strategies based on different types of documents. Chunking is a crucial step in a RAG (Retrieval-Augmented Generation) system that directly impacts the quality of the responses.

## Table of Contents

1. [Introduction to Chunking](#introduction-to-chunking)  
2. [Available Chunking Strategies](#available-chunking-strategies)  
3. [Recommendations by Document Type](#recommendations-by-document-type)  
4. [Evaluation and Optimization](#evaluation-and-optimization)  
5. [Best Practices](#best-practices)

## Introduction to Chunking

Chunking is the process of dividing a document into smaller units (chunks) that are indexed and retrieved independently. The goal is to create chunks that:

- Contain enough context to be useful  
- Are small enough to be specific  
- Preserve semantic units (sentences, paragraphs)  
- Minimize redundancy while ensuring complete coverage  

Optimal chunking improves both retrieval accuracy and the quality of generated responses.

## Available Chunking Strategies

RLAMA offers several chunking strategies, each optimized for different types of content:

### 1. Fixed Chunking (`fixed`)

- **Description**: Splits the text into fixed-size chunks, trying not to cut words.  
- **Advantages**: Simple, predictable, works on all types of content.  
- **Disadvantages**: Does not respect semantic structure, may split sentences and paragraphs.  
- **Recommended for**: Unstructured documents, heterogeneous content.

### 2. Semantic Chunking (`semantic`)

- **Description**: Divides content by respecting natural boundaries like paragraphs and sections.  
- **Advantages**: Preserves semantic context, improves response quality.  
- **Disadvantages**: May produce chunks of highly variable size.  
- **Recommended for**: Articles, structured documents, user manuals.

### 3. Hybrid Chunking (`hybrid`)

- **Description**: Adapts the strategy based on the detected document type.  
- **Advantages**: Combines the strengths of other strategies.  
- **Disadvantages**: Increased complexity, may be less predictable.  
- **Recommended for**: Mixed corpora with various document types.

### 4. Hierarchical Chunking (`hierarchical`)

- **Description**: Creates a two-level structure with larger parent chunks and smaller child chunks.  
- **Advantages**: Captures both the big picture and finer details.  
- **Disadvantages**: More complex indexing, uses more storage.  
- **Recommended for**: Very long documents, books, full technical documentation.

## Recommendations by Document Type

### Markdown/Documentation Files

- **Recommended strategy**: `semantic` or `hybrid`  
- **Chunk size**: 1000–1500 characters (~250–375 tokens)  
- **Overlap**: 10% of chunk size  
- **Why**: Markdown documents generally have a clear structure (sections, subsections) that semantic chunking can leverage.

### Source Code

- **Recommended strategy**: `hybrid` (which uses code-aware chunking)  
- **Chunk size**: 500–1000 characters (~125–250 tokens)  
- **Overlap**: 5–10% of chunk size  
- **Why**: Code has defined structure (functions, classes) and code-aware chunking preserves these logical units.

### Long Texts/Articles

- **Recommended strategy**: `semantic` or `hierarchical`  
- **Chunk size**: 1500–2000 characters (~375–500 tokens)  
- **Overlap**: 15–20% of chunk size  
- **Why**: Long texts benefit from strategies that respect paragraphs and sections, with higher overlap to maintain context.

### HTML/Web Pages

- **Recommended strategy**: `hybrid` (which uses HTML-aware chunking)  
- **Chunk size**: 1000–1500 characters (~250–375 tokens)  
- **Overlap**: 10–15% of chunk size  
- **Why**: HTML content has structure defined by tags that specialized chunking can exploit.

### Unstructured Texts

- **Recommended strategy**: `fixed` or parameterized `semantic`  
- **Chunk size**: 800–1200 characters (~200–300 tokens)  
- **Overlap**: 20% of chunk size  
- **Why**: Without a clear structure, higher overlap helps preserve context.

## Evaluation and Optimization

RLAMA provides tools to evaluate and optimize your chunking strategies. The `chunk-eval` tool lets you test different configurations on your specific documents.

### Using the Evaluation Tool

```bash
# Evaluate a specific configuration
rlama chunk-eval --file=your_document.txt --strategy=semantic --size=1500 --overlap=150

# Compare all available strategies
rlama chunk-eval --file=your_document.txt --compare-all --detailed
```

### Evaluation Metrics

- **Semantic Coherence Score**: Overall quality of chunking (0–1, higher = better)  
- **Cut Sentences/Paragraphs**: Number of chunks that break semantic units (fewer = better)  
- **Redundancy Rate**: Percentage of duplicated content due to overlap  
- **Content Coverage**: Percentage of the original document covered by the chunks

### Recommended Optimization Process

1. Start by comparing all strategies on your corpus.  
2. Identify the top 2–3 strategies based on the metrics.  
3. Fine-tune parameters (size, overlap) for those strategies.  
4. Test optimized configurations on representative queries.  
5. Measure impact on final responses, not just the metrics.

## Best Practices

### General Tips

- **Adapt the strategy to the content**: There’s no universal configuration; tailor it to your documents.  
- **Favor semantic coherence**: Natural boundaries (paragraphs, sections) usually make better chunk breakpoints.  
- **Avoid overly small chunks**: Chunks under 100 tokens often lack context.  
- **Limit overly large chunks**: Chunks over 500 tokens may be too generic and hurt precision.  
- **Test with real queries**: Final impact is measured by response quality, not just metrics.

### What to Avoid

- **Splitting sentences mid-way**: This fragments information and leads to incoherent chunks.  
- **Ignoring document structure**: Using existing structure (headers, sections) generally improves results.  
- **Too much overlap**: Beyond 25%, redundancy can become more harmful than helpful.  
- **Highly variable chunk sizes**: Wide variation in size can bias retrieval.

---

## Conclusion

Optimizing chunking is an iterative process that requires testing and tuning. The recommendations provided here serve as a starting point, and the metrics generated by the evaluation tool will help you refine your approach for your specific use case.

For questions or suggestions, feel free to open an issue on the project’s GitHub repository.
````

## File: docs/reranking_guidelines.md
````markdown
# RLAMA Reranking Documentation

## Overview
Reranking in RLAMA is a feature that improves retrieval accuracy by applying a second-stage ranking to initial search results using a cross-encoder approach. This helps ensure more relevant documents are prioritized in responses to queries.

## Features
- Enabled by default for all RAG systems
- Configurable weights between vector similarity and reranking scores
- Adjustable result limits and thresholds
- Custom model support for reranking

## Default Configuration
- TopK: 5 results (maximum number of results after reranking)
- Initial retrieval: 20 documents
- Reranker weight: 0.7 (70% reranker score, 30% vector similarity)
- Score threshold: 0.0 (no minimum score requirement)
- Model: Uses the same model as the RAG system by default

## Usage

### Command Line Interface

1. **Configure Reranking for a RAG System**
```bash
rlama add-reranker my-rag [options]
```

Available options:
- `--model`: Specify a custom model for reranking (defaults to RAG model)
- `--weight`: Set the weight for reranker scores (0-1)
- `--threshold`: Set minimum score threshold for results
- `--topk`: Set maximum number of results to return
- `--disable`: Disable reranking for this RAG

Examples:
```bash
# Configure with custom model
rlama add-reranker my-rag --model reranker-model

# Adjust weights and limits
rlama add-reranker my-rag --weight 0.8 --topk 10

# Disable reranking
rlama add-reranker my-rag --disable
```

### Programmatic Usage

1. **Creating a RAG with Reranking**
```go
err := ragService.CreateRagWithOptions("llama3.2", "my-rag", documentPath, service.DocumentLoaderOptions{
    ChunkSize: 200,
    ChunkOverlap: 50,
    EnableReranker: true,  // Reranking is enabled by default
})
```

2. **Customizing Reranking Options**
```go
options := service.RerankerOptions{
    TopK: 10,                // Return top 10 results
    InitialK: 30,           // Retrieve 30 initial results
    RerankerModel: "custom-model",  // Use custom model
    ScoreThreshold: 0.5,    // Minimum relevance score
    RerankerWeight: 0.8,    // 80% reranker, 20% vector similarity
}
```

## How It Works

1. **Initial Retrieval**: The system first retrieves an initial set of documents using vector similarity search (default: top 20 documents).

2. **Reranking Process**:
   - Each retrieved document is evaluated using a cross-encoder model
   - The model scores document relevance on a scale of 0 to 1
   - Final scores combine vector similarity and reranking scores based on weights
   - Results are sorted by final score and limited to TopK

3. **Scoring Formula**:
finalScore = (rerankerWeight × rerankerScore) + ((1 - rerankerWeight) × vectorScore) 

## Performance Considerations

- Reranking adds additional processing time as each document needs to be evaluated
- The InitialK parameter affects both accuracy and performance
- Larger TopK values increase processing time
- Consider disabling reranking for applications requiring minimal latency

## Best Practices

1. **Model Selection**
   - Use the same model as your RAG system for consistency
   - Choose models that excel at cross-encoding for better results

2. **Parameter Tuning**
   - Start with default weights (0.7) and adjust based on results
   - Increase InitialK for better recall at the cost of performance
   - Set appropriate thresholds based on your use case

3. **Performance Optimization**
   - Limit TopK to necessary minimum
   - Consider chunk size impact on reranking performance
   - Monitor and adjust InitialK based on result quality

## Troubleshooting

Common issues and solutions:

1. **Slow Response Times**
   - Reduce InitialK or TopK values
   - Consider using a lighter reranking model
   - Check if chunk sizes are appropriate

2. **Poor Result Quality**
   - Increase reranker weight
   - Adjust score threshold
   - Increase InitialK for more candidate documents

3. **Model Compatibility**
   - Ensure the reranking model supports the required operations
   - Check model availability in your Ollama installation

## Advanced Configuration

For specific use cases, you can fine-tune the reranking system by:

1. **Custom Scoring**
   - Adjust weights based on document types
   - Implement custom thresholds for different queries

2. **Model Chaining**
   - Use different models for initial retrieval and reranking
   - Combine multiple reranking passes with different models

## Examples

1. **Basic Usage with Default Settings**
```bash
rlama add-reranker my-documents
```

2. **High-Precision Configuration**
```bash
rlama add-reranker research-papers --weight 0.9 --threshold 0.7 --topk 3
```

3. **Large-Scale Configuration**
```bash
rlama add-reranker large-corpus --topk 20 --weight 0.6
```
```

This README provides a comprehensive guide to understanding and using RLAMA's reranking functionality, based on the implementation shown in the provided code files.
````

## File: internal/client/client_test.go
````go
package client
⋮----
import (
	"testing"
)
⋮----
"testing"
⋮----
func TestNewOllamaClient(t *testing.T)
````

## File: internal/config/config_test.go
````go
package config
⋮----
import "testing"
⋮----
func TestGetDataDir(t *testing.T)
````

## File: internal/config/config.go
````go
package config
⋮----
import (
	"os"
	"path/filepath"
)
⋮----
"os"
"path/filepath"
⋮----
var (
	// DataDir is the directory where RLAMA stores all its data
	DataDir string
)
⋮----
// DataDir is the directory where RLAMA stores all its data
⋮----
func init()
⋮----
// Check for data directory from environment variable first
⋮----
// If not set in environment, use default location in user's home directory
⋮----
// Ensure the data directory exists
⋮----
// GetDataDir returns the data directory path
// Priority: environment variable > default (~/.rlama)
func GetDataDir() string
⋮----
// Check if environment variable is set
⋮----
// Use default location
⋮----
return ".rlama" // Fallback to current directory
````

## File: internal/crawler/crawl4ai_style_test.go
````go
package crawler
⋮----
import (
	"net/url"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)
⋮----
"net/url"
"strings"
"testing"
⋮----
"github.com/PuerkitoBio/goquery"
⋮----
func TestCrawl4AIStyleConverter(t *testing.T)
⋮----
// Sample HTML content
⋮----
// Set up the test
⋮----
// Create converter and run conversion
⋮----
// Check result
⋮----
// Verify content was processed correctly
⋮----
// Verify unwanted elements were removed
````

## File: internal/crawler/crawl4ai_style.go
````go
package crawler
⋮----
import (
	"net/url"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/PuerkitoBio/goquery"
)
⋮----
"net/url"
"strings"
⋮----
htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
"github.com/PuerkitoBio/goquery"
⋮----
// Crawl4AIStyleConverter provides enhanced HTML to Markdown conversion
// inspired by Crawl4AI's approach to create LLM-friendly markdown content
type Crawl4AIStyleConverter struct{}
⋮----
// NewCrawl4AIStyleConverter creates a new Markdown converter with enhancements
func NewCrawl4AIStyleConverter() *Crawl4AIStyleConverter
⋮----
// ConvertHTMLToMarkdown converts HTML content to Markdown with optimizations
func (c *Crawl4AIStyleConverter) ConvertHTMLToMarkdown(doc *goquery.Document, baseURL *url.URL) (string, error)
⋮----
// Pre-process the document
⋮----
// Extract main content
⋮----
// Get HTML content from the main content
⋮----
// Convert to markdown
⋮----
// Post-process markdown to clean it up
⋮----
// cleanDocument removes unwanted elements from the HTML document
func (c *Crawl4AIStyleConverter) cleanDocument(doc *goquery.Document)
⋮----
// Remove unwanted elements that typically don't contain useful content
⋮----
// extractMainContent finds the main content node of the document
func (c *Crawl4AIStyleConverter) extractMainContent(doc *goquery.Document) *goquery.Selection
⋮----
// Try to find the main content area using common selectors
⋮----
// Verify the selection has substantive content
⋮----
if len(strings.TrimSpace(text)) > 200 { // If it has more than 200 chars of text
⋮----
// Fallback: If no main content area could be determined, use body
⋮----
// postProcessMarkdown cleans up the generated markdown
func (c *Crawl4AIStyleConverter) postProcessMarkdown(markdown string) string
⋮----
// Replace multiple blank lines with a single blank line
⋮----
// Remove trailing whitespace from each line
⋮----
// Trim leading and trailing whitespace from the entire string
````

## File: internal/crawler/crawler_test.go
````go
package crawler
⋮----
import "testing"
⋮----
func TestNewWebCrawler(t *testing.T)
````

## File: internal/domain/document_chunk.go
````go
package domain
⋮----
import (
	"fmt"
	"time"
)
⋮----
"fmt"
"time"
⋮----
// DocumentChunk represents a portion of a document with metadata
type DocumentChunk struct {
	ID          string    `json:"id"`
	DocumentID  string    `json:"documentId"`
	Content     string    `json:"content"`
	StartPos    int       `json:"start_pos"`
	EndPos      int       `json:"end_pos"`
	ChunkIndex  int       `json:"chunk_index"`
	Embedding   []float32 `json:"-"` // Not serialized to JSON
	CreatedAt   time.Time `json:"created_at"`
	Metadata    map[string]string `json:"metadata"`
	ChunkNumber int       `json:"chunkNumber"`
	TotalChunks int       `json:"totalChunks"`
}
⋮----
Embedding   []float32 `json:"-"` // Not serialized to JSON
⋮----
// NewDocumentChunk creates a new chunk from a document
func NewDocumentChunk(doc *Document, content string, startPos, endPos, chunkIndex int) *DocumentChunk
⋮----
// Generate a unique ID for the chunk
⋮----
// Create metadata for the chunk
⋮----
"chunk_position": fmt.Sprintf("%d of %d", chunkIndex+1, 0), // Total will be updated later
⋮----
// GetMetadataString returns a formatted string of the chunk's metadata
func (c *DocumentChunk) GetMetadataString() string
⋮----
// UpdateTotalChunks updates the chunk position metadata with the total chunk count
func (c *DocumentChunk) UpdateTotalChunks(total int)
````

## File: internal/domain/rag_test.go
````go
package domain
⋮----
import (
	"testing"
)
⋮----
"testing"
⋮----
func TestNewRagSystem(t *testing.T)
````

## File: internal/repository/repository_test.go
````go
package repository
⋮----
import (
	"testing"
)
⋮----
"testing"
⋮----
func TestNewRagRepository(t *testing.T)
````

## File: internal/server/server_test.go
````go
package server
⋮----
import "testing"
⋮----
func TestNewServer(t *testing.T)
````

## File: internal/service/file_watcher.go
````go
package service
⋮----
import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dontizi/rlama/internal/domain"
)
⋮----
"fmt"
"os"
"path/filepath"
"time"
⋮----
"github.com/dontizi/rlama/internal/domain"
⋮----
// FileWatcher is responsible for watching directories for file changes
type FileWatcher struct {
	ragService RagService
}
⋮----
// NewFileWatcher creates a new file watcher service
func NewFileWatcher(ragService RagService) *FileWatcher
⋮----
// CheckAndUpdateRag checks for new files in the watched directory and updates the RAG
func (fw *FileWatcher) CheckAndUpdateRag(rag *domain.RagSystem) (int, error)
⋮----
return 0, nil // Watching not enabled
⋮----
// Check if the directory exists
⋮----
// Get the last modified time of the directory
⋮----
// If the directory hasn't been modified since last check, no need to proceed
⋮----
// Convert watch options to document loader options
⋮----
// Get existing document paths to avoid re-processing
⋮----
// Create a document loader
⋮----
// Load all documents from the directory
⋮----
// Filter out existing documents
var newDocs []*domain.Document
⋮----
// Update last watched time even if no new documents
⋮----
// Create chunker service with options from the RAG
⋮----
// Process each new document - chunk and prepare for embeddings
var allChunks []*domain.DocumentChunk
⋮----
// Chunk the document
⋮----
// Update total chunks in metadata
⋮----
// Generate embeddings for all chunks
⋮----
// Add documents and chunks to the RAG
⋮----
// Update last watched time
⋮----
// Save the updated RAG
⋮----
// getLastModifiedTime gets the latest modification time in a directory
func getLastModifiedTime(dirPath string) time.Time
⋮----
var lastModTime time.Time
⋮----
// Walk through the directory and find the latest modification time
⋮----
return nil // Skip errors
⋮----
// StartWatcherDaemon starts a background daemon to watch directories
// for all RAGs that have watching enabled with intervals
func (fw *FileWatcher) StartWatcherDaemon(interval time.Duration)
⋮----
// checkAllRags checks all RAGs with watching enabled
func (fw *FileWatcher) checkAllRags()
⋮----
// Get all RAGs
⋮----
// Check if watching is enabled and if interval has passed
````

## File: internal/service/service_test.go
````go
package service
⋮----
import "testing"
⋮----
func TestNewRagService(t *testing.T)
````

## File: pkg/vector/hybrid_store.go
````go
package vector
⋮----
import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/blevesearch/bleve/v2"
)
⋮----
"fmt"
"os"
"path/filepath"
"sort"
⋮----
"github.com/blevesearch/bleve/v2"
⋮----
// DocumentData represents the data structure for Bleve indexing
type DocumentData struct {
	ID       string `json:"id"`
	Content  string `json:"content"`
	Metadata string `json:"metadata"`
}
⋮----
// EnhancedHybridStore combines vector search and BM25 text search
type EnhancedHybridStore struct {
	VectorStore VectorStoreInterface `json:"-"`
	TextIndex   bleve.Index          `json:"-"`
	WeightBM25  float64              `json:"weight_bm25"`
	// Maps for quick access to content and metadata
	contentCache  map[string]string `json:"-"`
	metadataCache map[string]string `json:"-"`
}
⋮----
// Maps for quick access to content and metadata
⋮----
// Ensure EnhancedHybridStore implements VectorStoreInterface
var _ VectorStoreInterface = (*EnhancedHybridStore)(nil)
⋮----
// HybridStoreConfig holds configuration for creating an EnhancedHybridStore
type HybridStoreConfig struct {
	IndexPath            string
	Dimensions           int
	VectorStoreType      string // "internal", "qdrant"
	QdrantHost           string
	QdrantPort           int
	QdrantAPIKey         string
	QdrantCollectionName string
	QdrantGRPC           bool
}
⋮----
VectorStoreType      string // "internal", "qdrant"
⋮----
// NewEnhancedHybridStore creates a new enhanced hybrid store
func NewEnhancedHybridStore(indexPath string, dimensions int) (*EnhancedHybridStore, error)
⋮----
// NewEnhancedHybridStoreWithConfig creates a new enhanced hybrid store with full configuration
func NewEnhancedHybridStoreWithConfig(config HybridStoreConfig) (*EnhancedHybridStore, error)
⋮----
// Create index directory if needed
⋮----
// Create or open Bleve index
var textIndex bleve.Index
var err error
⋮----
// In-memory index
⋮----
// Check if index already exists
⋮----
// Create new index
⋮----
// Open existing index
⋮----
// Create vector store based on configuration
var vectorStore VectorStoreInterface
⋮----
true, // createCollectionIfNotExists
⋮----
// Default to internal vector store
⋮----
WeightBM25:    0.3, // 30% BM25, 70% vector by default
⋮----
// AddDocument adds a document to both the vector and text indexes
func (hs *EnhancedHybridStore) AddDocument(id string, content string, metadata string, vector []float32) error
⋮----
// Add to vector store with payload if it's a QdrantStore
⋮----
// For internal stores, use the standard Add method
⋮----
// Add to cache
⋮----
// Add to text index
⋮----
// Add implements the VectorStoreInterface
func (hs *EnhancedHybridStore) Add(id string, vector []float32)
⋮----
// Remove removes a document from both indexes
func (hs *EnhancedHybridStore) Remove(id string)
⋮----
// Remove from vector store
⋮----
// Remove from caches
⋮----
// Remove from text index (ignore errors for interface compatibility)
⋮----
// GetContent returns a document's content
func (hs *EnhancedHybridStore) GetContent(id string) string
⋮----
// GetMetadata returns a document's metadata
func (hs *EnhancedHybridStore) GetMetadata(id string) string
⋮----
// HybridSearchResult représente un résultat de recherche hybride
type HybridSearchResult struct {
	ID             string  `json:"id"`
	VectorScore    float64 `json:"vector_score"`
	TextScore      float64 `json:"text_score"`
	CombinedScore  float64 `json:"combined_score"`
}
⋮----
// HybridSearch performs a combined vector and text search
func (hs *EnhancedHybridStore) HybridSearch(queryVector []float32, queryText string, limit int) ([]HybridSearchResult, error)
⋮----
// Execute vector search
vectorResults := hs.VectorStore.Search(queryVector, limit*2) // Get more results for fusion
⋮----
// Execute BM25 text search
⋮----
// Store normalized scores in maps
⋮----
// Normalize vector scores
⋮----
// Normalize text scores
⋮----
// Combine scores with weighting
var hybridResults []HybridSearchResult
⋮----
// If a document is only in one result set, give it a minimum score in the other
⋮----
vectorScore = 0.01 // Minimum score to not completely eliminate
⋮----
textScore = 0.01 // Minimum score to not completely eliminate
⋮----
// Weighted combined score
⋮----
// Sort by combined score in descending order
⋮----
// Limit results
⋮----
// Search implements the basic vector search interface
func (hs *EnhancedHybridStore) Search(query []float32, limit int) []SearchResult
⋮----
// Save saves both indexes
func (hs *EnhancedHybridStore) Save(vectorPath string) error
⋮----
// Save vector store
⋮----
// Bleve index is saved automatically if on disk
⋮----
// Load loads the store from a file
func (hs *EnhancedHybridStore) Load(path string) error
⋮----
// Load vector store
⋮----
// Bleve index is managed separately
⋮----
// Close properly closes the indexes
func (hs *EnhancedHybridStore) Close() error
⋮----
// SortHybridResults sorts results by combined score in descending order
func SortHybridResults(results []HybridSearchResult)
````

## File: pkg/vector/store.go
````go
package vector
⋮----
import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
)
⋮----
"encoding/json"
"fmt"
"math"
"os"
"sort"
⋮----
// VectorStoreInterface defines the common interface for vector stores
type VectorStoreInterface interface {
	Add(id string, vector []float32)
	Search(query []float32, limit int) []SearchResult
	Remove(id string)
	Save(path string) error
	Load(path string) error
}
⋮----
// VectorItem represents an item in the vector storage
type VectorItem struct {
	ID      string    `json:"id"`
	Vector  []float32 `json:"vector"`
}
⋮----
// SearchResult represents a search result
type SearchResult struct {
	ID       string  `json:"id"`
	Score    float64 `json:"score"`
}
⋮----
// Store is a simple vector storage with cosine similarity search
type Store struct {
	Items []VectorItem `json:"items"`
}
⋮----
// Ensure Store implements VectorStoreInterface
var _ VectorStoreInterface = (*Store)(nil)
⋮----
// NewStore creates a new vector storage
func NewStore() *Store
⋮----
// Add adds a vector to the storage
func (s *Store) Add(id string, vector []float32)
⋮----
// Check if the ID already exists
⋮----
// Replace the existing vector
⋮----
// Add a new vector
⋮----
// Search searches for the most similar vectors
func (s *Store) Search(query []float32, limit int) []SearchResult
⋮----
var results []SearchResult
⋮----
// Calculate cosine similarity for each vector
⋮----
// Sort by descending score
⋮----
// Limit the number of results
⋮----
// cosineSimilarity calculates the cosine similarity between two vectors
func cosineSimilarity(a, b []float32) float64
⋮----
var dotProduct float64
var normA, normB float64
⋮----
// Save saves the vector storage to a file
func (s *Store) Save(path string) error
⋮----
// Load loads the vector storage from a file
func (s *Store) Load(path string) error
⋮----
// Check if the file exists
⋮----
// File doesn't exist, use empty storage
⋮----
// Remove removes a vector from the storage by its ID
func (s *Store) Remove(id string)
````

## File: pkg/vector/vector_test.go
````go
package vector
⋮----
import "testing"
⋮----
func TestNewEnhancedHybridStore(t *testing.T)
````

## File: install.ps1
````powershell
# Windows installation script for RLAMA
Write-Host "
 ____  _       _    __  __    _    
|  _ \| |     / \  |  \/  |  / \   
| |_) | |    / _ \ | |\/| | / _ \  
|  _ <| |___/ ___ \| |  | |/ ___ \ 
|_| \_\_____/_/   \_\_|  |_/_/   \_\
                                  
Retrieval-Augmented Language Model Adapter for Windows
"

# Determine architecture
$arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
$binaryName = "rlama_windows_$arch.exe"

# Create installation directories
$dataDir = "$env:USERPROFILE\.rlama"
$installDir = "$env:LOCALAPPDATA\RLAMA"

Write-Host "Installing RLAMA..."
Write-Host "Downloading RLAMA for Windows $arch..."

# Create directories if they don't exist
New-Item -ItemType Directory -Force -Path $dataDir | Out-Null
New-Item -ItemType Directory -Force -Path $installDir | Out-Null

# Download the binary
$downloadUrl = "https://github.com/dontizi/rlama/releases/latest/download/$binaryName"
$outputPath = "$installDir\rlama.exe"

try {
    Invoke-WebRequest -Uri $downloadUrl -OutFile $outputPath
    
    # Add to PATH if not already there
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($currentPath -notlike "*$installDir*") {
        [Environment]::SetEnvironmentVariable("Path", "$currentPath;$installDir", "User")
        Write-Host "Added RLAMA to your PATH. You may need to restart your terminal."
    }
    
    Write-Host "RLAMA has been successfully installed to $outputPath!"
    Write-Host "You can now use RLAMA by running the 'rlama' command."
} catch {
    Write-Host "Error downloading RLAMA: $_"
    exit 1
}
````

## File: main.go
````go
package main
⋮----
import (
	"github.com/dontizi/rlama/cmd"
)
⋮----
"github.com/dontizi/rlama/cmd"
⋮----
func main()
⋮----
// Execute the root command
````

## File: cmd/install_dependencies.go
````go
package cmd
⋮----
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)
⋮----
"fmt"
"os"
"os/exec"
"path/filepath"
"runtime"
⋮----
"github.com/spf13/cobra"
⋮----
var installDependenciesCmd = &cobra.Command{
	Use:   "install-dependencies",
	Short: "Install necessary dependencies for RLAMA",
	Long:  `Install system and Python dependencies for optimal RLAMA performance, including the BGE reranker.`,
	Run: func(cmd *cobra.Command, args []string) {
		installDependencies()
	},
}
⋮----
func init()
⋮----
func installDependencies()
⋮----
// Find the installation script path
⋮----
// Use an alternative solution
⋮----
// The scripts directory is presumed to be in the same directory as the executable
// or in the parent directory for development environments
⋮----
// Check if the script exists
⋮----
// Try in the parent directory (for development)
⋮----
// Execute the script
var cmd *exec.Cmd
⋮----
// On Windows, use bash.exe (WSL or Git Bash)
⋮----
// On Unix-like, execute the script directly
⋮----
// This function is used if the install_deps.sh script is not found
func installDependenciesFallback()
⋮----
func installPythonDependencies()
⋮----
// Determine the Python command to use
⋮----
// Install Python dependencies
````

## File: cmd/uninstall.go
````go
package cmd
⋮----
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)
⋮----
"fmt"
"os"
"os/exec"
"path/filepath"
"runtime"
"strings"
⋮----
"github.com/spf13/cobra"
⋮----
var forceUninstall bool
⋮----
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall RLAMA and all its files",
	Long:  `Completely uninstall RLAMA by removing the executable and all associated data files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. Check if the user confirmed the deletion
		if !forceUninstall {
			fmt.Print("This action will remove RLAMA and all your data. Are you sure? (y/n): ")
			var response string
			fmt.Scanln(&response)

			response = strings.ToLower(strings.TrimSpace(response))
			if response != "y" && response != "yes" {
				fmt.Println("Uninstallation cancelled.")
				return nil
			}
		}

		// 2. Delete the data directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("unable to determine user directory: %w", err)
		}

		dataDir := filepath.Join(homeDir, ".rlama")
		fmt.Printf("Removing data directory: %s\n", dataDir)

		if _, err := os.Stat(dataDir); err == nil {
			err = os.RemoveAll(dataDir)
			if err != nil {
				return fmt.Errorf("unable to remove data directory: %w", err)
			}
			fmt.Println("✓ Data directory removed")
		} else {
			fmt.Println("Data directory doesn't exist or has already been removed")
		}

		// 3. Remove the executable
		var executablePath string
		if runtime.GOOS == "windows" {
			// Try different locations where it might be installed
			localAppData := os.Getenv("LOCALAPPDATA")
			possiblePaths := []string{
				filepath.Join(localAppData, "RLAMA", "rlama.exe"),
				filepath.Join(os.Getenv("ProgramFiles"), "RLAMA", "rlama.exe"),
				filepath.Join(os.Getenv("ProgramFiles(x86)"), "RLAMA", "rlama.exe"),
				filepath.Join(homeDir, "AppData", "Local", "RLAMA", "rlama.exe"),
			}

			for _, path := range possiblePaths {
				if _, err := os.Stat(path); err == nil {
					executablePath = path
					break
				}
			}
		} else {
			executablePath = "/usr/local/bin/rlama"
		}

		fmt.Printf("Removing executable: %s\n", executablePath)

		if executablePath == "" && runtime.GOOS == "windows" {
			fmt.Println("Could not find RLAMA executable. If RLAMA is installed elsewhere, you may need to:")
			fmt.Println("1. Run Command Prompt as Administrator")
			fmt.Println("2. Navigate to the installation directory")
			fmt.Println("3. Manually delete the rlama.exe file")
			fmt.Println("\nRLAMA data directory has been removed successfully.")
			return nil
		}

		if _, err := os.Stat(executablePath); err == nil {
			// On macOS and Linux, we probably need sudo
			var err error
			if runtime.GOOS == "windows" {
				// On Windows, try to remove directly
				err = os.Remove(executablePath)
				if err != nil {
					// If direct removal fails, try with elevated privileges using full PowerShell path
					fmt.Println("Need elevated privileges to remove the executable")
					powershellPath := filepath.Join(os.Getenv("SystemRoot"), "System32", "WindowsPowerShell", "v1.0", "powershell.exe")
					err = execCommand(powershellPath, "-Command", fmt.Sprintf("Start-Process -Verb RunAs -FilePath 'cmd.exe' -ArgumentList '/c del \"%s\"'", executablePath))
				}
			} else if isRoot() {
				// If we're already root on Unix systems
				err = os.Remove(executablePath)
			} else {
				fmt.Println("You may need to enter your password to remove the executable")
				err = execCommand("sudo", "rm", executablePath)
			}

			if err != nil {
				if runtime.GOOS == "windows" {
					return fmt.Errorf("unable to remove executable: %w\nTry running the command prompt as administrator and run 'rlama uninstall' again", err)
				}
				return fmt.Errorf("unable to remove executable: %w", err)
			}
			fmt.Println("✓ Executable removed")
		} else {
			fmt.Println("Executable doesn't exist or has already been removed")
		}

		fmt.Println("\nRLAMA has been successfully uninstalled.")
		return nil
	},
}
⋮----
// 1. Check if the user confirmed the deletion
⋮----
var response string
⋮----
// 2. Delete the data directory
⋮----
// 3. Remove the executable
var executablePath string
⋮----
// Try different locations where it might be installed
⋮----
// On macOS and Linux, we probably need sudo
var err error
⋮----
// On Windows, try to remove directly
⋮----
// If direct removal fails, try with elevated privileges using full PowerShell path
⋮----
// If we're already root on Unix systems
⋮----
// execCommand executes a system command
func execCommand(name string, args ...string) error
⋮----
func init()
⋮----
// isRoot returns true if the current process is running with root/admin privileges
// This is a safe wrapper around os.Geteuid() which doesn't exist on Windows
func isRoot() bool
⋮----
// On Windows, check if we have admin privileges using a different method
// However, this is not easily determined, so we'll return false
// and let the code try direct removal first
⋮----
// On Unix systems, check if euid is 0 (root)
````

## File: cmd/update_reranker.go
````go
package cmd
⋮----
import (
	"fmt"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/service"
	"github.com/spf13/cobra"
)
⋮----
"fmt"
⋮----
"github.com/dontizi/rlama/internal/client"
"github.com/dontizi/rlama/internal/service"
"github.com/spf13/cobra"
⋮----
var updateRerankerCmd = &cobra.Command{
	Use:   "update-reranker [rag-name]",
	Short: "Updates the reranking model of an existing RAG",
	Long:  `Configures an existing RAG to use the default BGE Reranker model.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		updateReranker(args[0])
	},
}
⋮----
func init()
⋮----
func updateReranker(ragName string)
⋮----
// Load the RAG service
⋮----
// Update the reranking model
````

## File: internal/client/bge_reranker_client.go
````go
package client
⋮----
import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)
⋮----
"bytes"
"encoding/json"
"fmt"
"os/exec"
"strings"
⋮----
// BGERerankerClient handles interactions with the BGE Reranker model via Python
type BGERerankerClient struct {
	modelName string
	useFP16   bool
	silent    bool
}
⋮----
// NewBGERerankerClient creates a new instance of BGERerankerClient
func NewBGERerankerClient(modelName string) *BGERerankerClient
⋮----
// NewBGERerankerClientWithOptions creates a new instance of BGERerankerClient with additional options
func NewBGERerankerClientWithOptions(modelName string, silent bool) *BGERerankerClient
⋮----
// Check dependencies and model
⋮----
// Only check model if dependencies are available
⋮----
// GetModelName returns the model name used by this client
func (c *BGERerankerClient) GetModelName() string
⋮----
// ComputeScores calculates relevance scores between queries and passages
func (c *BGERerankerClient) ComputeScores(pairs [][]string, normalize bool) ([]float64, error)
⋮----
// Convert Go boolean to Python boolean
⋮----
// Prepare input data
⋮----
// Execute the Python script
⋮----
// Extract just the JSON part from the output
// First try to find the first valid JSON object
⋮----
// Find where the JSON object ends
⋮----
// Extract just the JSON portion
⋮----
// Parse the output
var result map[string]interface{}
⋮----
// If parsing fails, try to extract the JSON using regex as a fallback
⋮----
// Try one more approach - find the first line that looks like valid JSON
⋮----
// Successfully parsed this line as JSON
⋮----
// If we still don't have a result, return the error
⋮----
// Check for error
⋮----
// Extract scores
⋮----
// CheckDependencies checks if FlagEmbedding is installed
func (c *BGERerankerClient) CheckDependencies() error
⋮----
// CheckModelExists verifies that the model exists and is accessible
func (c *BGERerankerClient) CheckModelExists() error
⋮----
// Find the first line that starts with '{' (likely our JSON)
⋮----
var jsonLine []byte
⋮----
// If no JSON line found, try to extract it from the whole output
⋮----
// Find where the JSON object ends
⋮----
// If we still couldn't find JSON data, return error
````

## File: internal/repository/profile_repository.go
````go
package repository
⋮----
import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dontizi/rlama/internal/domain"
)
⋮----
"encoding/json"
"fmt"
"os"
"path/filepath"
"time"
⋮----
"github.com/dontizi/rlama/internal/domain"
⋮----
// ProfileRepository gère le stockage des profils API
type ProfileRepository struct {
	basePath string
}
⋮----
// NewProfileRepository crée une nouvelle instance de ProfileRepository
func NewProfileRepository() *ProfileRepository
⋮----
// use ~/.rlama/profiles as default directory
⋮----
// Create the directory if it doesn't exist
⋮----
// getProfilePath returns the full path for a given profile
func (r *ProfileRepository) getProfilePath(name string) string
⋮----
// Exists checks if a profile exists
func (r *ProfileRepository) Exists(name string) bool
⋮----
// Save saves a profile
func (r *ProfileRepository) Save(profile *domain.APIProfile) error
⋮----
// Load loads a profile
func (r *ProfileRepository) Load(name string) (*domain.APIProfile, error)
⋮----
var profile domain.APIProfile
⋮----
// Delete deletes a profile
func (r *ProfileRepository) Delete(name string) error
⋮----
// ListAll returns a list of all profiles
func (r *ProfileRepository) ListAll() ([]string, error)
⋮----
var profileNames []string
⋮----
// Remove the .json extension
````

## File: internal/repository/rag_repository.go
````go
package repository
⋮----
import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/pkg/vector"
	"github.com/dontizi/rlama/internal/config"
)
⋮----
"encoding/json"
"fmt"
"os"
"path/filepath"
⋮----
"github.com/dontizi/rlama/internal/domain"
"github.com/dontizi/rlama/pkg/vector"
"github.com/dontizi/rlama/internal/config"
⋮----
// RagRepository manages the persistence of RAG systems
type RagRepository struct {
	basePath string
}
⋮----
// NewRagRepository creates a new instance of RagRepository
func NewRagRepository() *RagRepository
⋮----
// Create the folder if it doesn't exist
⋮----
// getRagPath returns the complete path for a given RAG
func (r *RagRepository) getRagPath(ragName string) string
⋮----
// getRagInfoPath returns the path of the RAG information file
func (r *RagRepository) getRagInfoPath(ragName string) string
⋮----
// getRagVectorStorePath returns the path of the vector storage file
func (r *RagRepository) getRagVectorStorePath(ragName string) string
⋮----
// Exists checks if a RAG exists
func (r *RagRepository) Exists(ragName string) bool
⋮----
// Save saves a RAG system
func (r *RagRepository) Save(rag *domain.RagSystem) error
⋮----
// Create the folder for this RAG
⋮----
// Save RAG information
ragInfo := *rag // Copy to avoid modifying the original
⋮----
// Serialize and save the info.json file
⋮----
// Save the Vector Store (only for internal stores, Qdrant handles its own persistence)
⋮----
// Load loads a RAG system
func (r *RagRepository) Load(ragName string) (*domain.RagSystem, error)
⋮----
// Check if the RAG exists
⋮----
// Load RAG information
⋮----
var ragInfo domain.RagSystem
⋮----
// Create a new Vector Store with the correct dimensions and configuration
⋮----
dimensions = 1536 // Default fallback for older RAGs
⋮----
var hybridStore *vector.EnhancedHybridStore
⋮----
// Create hybrid store with Qdrant configuration
⋮----
// Create internal hybrid store and load from file
⋮----
// ListAll returns the list of all available RAG systems
func (r *RagRepository) ListAll() ([]string, error)
⋮----
// Check if the base folder exists
⋮----
return []string{}, nil // No RAGs available
⋮----
// Read the folder contents
⋮----
var ragNames []string
⋮----
// Check if it's a valid RAG folder (contains info.json)
⋮----
// Delete deletes a RAG system
func (r *RagRepository) Delete(ragName string) error
⋮----
// Delete the complete RAG folder
````

## File: internal/server/server.go
````go
package server
⋮----
import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/service"
)
⋮----
"encoding/json"
"fmt"
"io"
"log"
"net/http"
"time"
⋮----
"github.com/dontizi/rlama/internal/client"
"github.com/dontizi/rlama/internal/service"
⋮----
// Server represents the API server
type Server struct {
	port        string
	ragService  service.RagService
	ollamaClient *client.OllamaClient
}
⋮----
// NewServer creates a new API server
func NewServer(port string, ollamaClient *client.OllamaClient) *Server
⋮----
port = "11249" // Default port
⋮----
// Start starts the API server
func (s *Server) Start() error
⋮----
// Register routes
⋮----
// Start the server
⋮----
// RagQueryRequest represents the request body for RAG queries
type RagQueryRequest struct {
	RagName       string `json:"rag_name"`
	Model         string `json:"model,omitempty"`
	Prompt        string `json:"prompt"`
	ContextSize   int    `json:"context_size,omitempty"`
	MaxWorkers    int    `json:"max_workers,omitempty"` // Added for parallel processing
}
⋮----
MaxWorkers    int    `json:"max_workers,omitempty"` // Added for parallel processing
⋮----
// RagQueryResponse represents the response for RAG queries
type RagQueryResponse struct {
	Response string `json:"response"`
	Success  bool   `json:"success"`
	Error    string `json:"error,omitempty"`
}
⋮----
// Handle RAG queries
func (s *Server) handleRagQuery(w http.ResponseWriter, r *http.Request)
⋮----
// Only allow POST requests
⋮----
// Read request body
⋮----
// Parse request
var req RagQueryRequest
⋮----
// Validate request
⋮----
// Set default context size if not provided
⋮----
// Load the RAG system
⋮----
// If model is specified and different from RAG's model, update it temporarily
⋮----
// Check if Ollama model is available
⋮----
// Use the original model of the RAG
⋮----
// Temporarily update the model if needed
⋮----
// Set parallel workers if specified
⋮----
// Update the RAG service with the new embedding service
⋮----
// Query the RAG system
⋮----
// Restore original model
⋮----
// Send the response
⋮----
// Handle health check requests
func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request)
⋮----
// Helper function to send error responses
func sendErrorResponse(w http.ResponseWriter, message string, statusCode int)
````

## File: internal/service/chunker_evaluation.go
````go
package service
⋮----
import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/dontizi/rlama/internal/domain"
)
⋮----
"fmt"
"math"
"strings"
"time"
⋮----
"github.com/dontizi/rlama/internal/domain"
⋮----
// ChunkingEvaluationMetrics contains evaluation metrics for a chunking strategy
type ChunkingEvaluationMetrics struct {
	// Basic metrics
	TotalChunks      int     // Total number of chunks produced
	AverageChunkSize float64 // Average chunk size in characters
	SizeStdDeviation float64 // Standard deviation of chunk sizes
	MaxChunkSize     int     // Size of the largest chunk
	MinChunkSize     int     // Size of the smallest chunk

	// Coherence metrics
	ChunksWithSplitSentences  int     // Number of chunks that split sentences
	ChunksWithSplitParagraphs int     // Number of chunks that split paragraphs
	SemanticCoherenceScore    float64 // Estimated semantic coherence score

	// Performance metrics
	ProcessingTimeMs int64 // Processing time in ms

	// Distribution metrics
	ContentCoverage float64 // % of original content covered by chunks
	RedundancyRate  float64 // Redundancy rate due to overlap

	// Strategy info
	Strategy     string // Strategy used
	ChunkSize    int    // Configured chunk size
	ChunkOverlap int    // Configured overlap
}
⋮----
// Basic metrics
TotalChunks      int     // Total number of chunks produced
AverageChunkSize float64 // Average chunk size in characters
SizeStdDeviation float64 // Standard deviation of chunk sizes
MaxChunkSize     int     // Size of the largest chunk
MinChunkSize     int     // Size of the smallest chunk
⋮----
// Coherence metrics
ChunksWithSplitSentences  int     // Number of chunks that split sentences
ChunksWithSplitParagraphs int     // Number of chunks that split paragraphs
SemanticCoherenceScore    float64 // Estimated semantic coherence score
⋮----
// Performance metrics
ProcessingTimeMs int64 // Processing time in ms
⋮----
// Distribution metrics
ContentCoverage float64 // % of original content covered by chunks
RedundancyRate  float64 // Redundancy rate due to overlap
⋮----
// Strategy info
Strategy     string // Strategy used
ChunkSize    int    // Configured chunk size
ChunkOverlap int    // Configured overlap
⋮----
// ChunkingEvaluator evaluates different chunking strategies
type ChunkingEvaluator struct {
	chunkerService *ChunkerService
	// References for semantic evaluation
	sentenceEndings  []string
	paragraphMarkers []string
}
⋮----
// References for semantic evaluation
⋮----
// NewChunkingEvaluator creates a new chunking evaluator
func NewChunkingEvaluator(chunkerService *ChunkerService) *ChunkingEvaluator
⋮----
// EvaluateChunkingStrategy evaluates a chunking strategy with the given parameters
func (ce *ChunkingEvaluator) EvaluateChunkingStrategy(doc *domain.Document, config ChunkingConfig) ChunkingEvaluationMetrics
⋮----
// Create a temporary chunker service with the configuration to test
⋮----
// Generate chunks with the strategy to evaluate
⋮----
// Calculate basic metrics
⋮----
// Calculate chunk sizes
⋮----
// Mean and standard deviation
⋮----
// Calculate standard deviation
⋮----
// Evaluate coverage and redundancy
⋮----
// Track covered characters
⋮----
// Calculate redundancy rate
⋮----
// Check for split sentences and paragraphs
⋮----
// Calculate an approximate semantic coherence score based on the metrics above
// Higher score = better estimated coherence
⋮----
// countChunksWithSplitSentences counts chunks that split a sentence
func (ce *ChunkingEvaluator) countChunksWithSplitSentences(chunks []*domain.DocumentChunk, originalContent string) int
⋮----
// Check the beginning of the chunk
⋮----
// Check if the previous character is a sentence ending marker
// Make sure that StartPos-1 is within the valid range
⋮----
// Check if we're in the middle of a sentence
⋮----
// Check the end of the chunk
⋮----
// If the last character is not a sentence ending and the next is not a sentence beginning
⋮----
// countChunksWithSplitParagraphs counts chunks that split a paragraph
func (ce *ChunkingEvaluator) countChunksWithSplitParagraphs(chunks []*domain.DocumentChunk, originalContent string) int
⋮----
// For simplicity, we consider a paragraph is split if:
// 1. The chunk doesn't start after a paragraph marker
// 2. The chunk doesn't end before a paragraph marker
⋮----
// Check the beginning
⋮----
// Check if there's a paragraph marker before
⋮----
// Check the end
⋮----
// Check if there's a paragraph marker after
⋮----
// calculateSemanticCoherenceScore calculates an estimated semantic coherence score
func (ce *ChunkingEvaluator) calculateSemanticCoherenceScore(metrics ChunkingEvaluationMetrics, totalChunks int) float64
⋮----
// Factors penalizing semantic coherence
⋮----
// Size factor: penalize highly variable sizes
⋮----
// Calculate score (inverted so that higher is better)
// Lower values = fewer split sentences/paragraphs and more consistency
⋮----
// Ensure the score is between 0 and 1
⋮----
// CompareChunkingStrategies runs a comparative evaluation of different chunking
// configurations and returns the results sorted by relevance for this document
func (ce *ChunkingEvaluator) CompareChunkingStrategies(doc *domain.Document) []ChunkingEvaluationMetrics
⋮----
var results []ChunkingEvaluationMetrics
⋮----
// Define the different strategies and configurations to test
⋮----
overlapRates := []float64{0.05, 0.1, 0.2} // as percentage of chunk size
⋮----
// Calculate overlap in characters
⋮----
// Evaluate this configuration
⋮----
// Sort results by semantic coherence score (from highest to lowest)
// Use a simple bubble sort for readability
⋮----
// GetOptimalChunkingConfig returns the optimal chunking configuration for the given document
func (ce *ChunkingEvaluator) GetOptimalChunkingConfig(doc *domain.Document) ChunkingConfig
⋮----
// If no results, return default configuration
⋮----
// Take the best configuration (first after sorting)
⋮----
// PrintEvaluationResults displays evaluation results in a readable format
func (ce *ChunkingEvaluator) PrintEvaluationResults(metrics ChunkingEvaluationMetrics)
````

## File: internal/service/chunker_service.go
````go
package service
⋮----
import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dontizi/rlama/internal/domain"
)
⋮----
"fmt"
"path/filepath"
"regexp"
"strings"
⋮----
"github.com/dontizi/rlama/internal/domain"
⋮----
// ChunkingConfig holds configuration for the document chunking process
type ChunkingConfig struct {
	ChunkSize        int    // Target size of each chunk in characters
	ChunkOverlap     int    // Number of characters to overlap between chunks
	IncludeMetadata  bool   // Whether to include metadata in chunk content
	ChunkingStrategy string // Strategy to use: "fixed", "semantic", "hybrid", "hierarchical"
}
⋮----
ChunkSize        int    // Target size of each chunk in characters
ChunkOverlap     int    // Number of characters to overlap between chunks
IncludeMetadata  bool   // Whether to include metadata in chunk content
ChunkingStrategy string // Strategy to use: "fixed", "semantic", "hybrid", "hierarchical"
⋮----
// DefaultChunkingConfig returns a default configuration for chunking
func DefaultChunkingConfig() ChunkingConfig
⋮----
ChunkSize:        1500, // Smaller chunks (~375 tokens) for better retrieval
ChunkOverlap:     150,  // 10% overlap
⋮----
ChunkingStrategy: "hybrid", // Default to hybrid strategy
⋮----
// ChunkerService handles splitting documents into manageable chunks
type ChunkerService struct {
	config ChunkingConfig
}
⋮----
// NewChunkerService creates a new chunker service with the specified config
func NewChunkerService(config ChunkingConfig) *ChunkerService
⋮----
// ChunkDocument splits a document into smaller chunks with metadata
// based on the selected chunking strategy
func (cs *ChunkerService) ChunkDocument(doc *domain.Document) []*domain.DocumentChunk
⋮----
// For very small documents, just return a single chunk regardless of strategy
⋮----
// Apply different chunking strategies based on configuration
var chunks []*domain.DocumentChunk
⋮----
// For auto strategy, use the evaluator to determine optimal configuration
⋮----
// Create a temporary chunker with the optimal configuration
⋮----
// Use the optimal chunker to create chunks
⋮----
// Store chunking strategy info in chunk metadata
⋮----
// Evaluate to get the metrics
⋮----
// Print analysis information
⋮----
// Default hybrid approach - choose strategy based on content type
⋮----
// Fallback to hybrid if invalid strategy specified
⋮----
// Safeguard: If we somehow got no chunks, fall back to fixed-size chunking
⋮----
// createHybridChunks selects the best chunking strategy based on document type
func (cs *ChunkerService) createHybridChunks(doc *domain.Document, content string, chunkSize int, overlap int) []*domain.DocumentChunk
⋮----
// Check file extension and content characteristics to determine best strategy
⋮----
// Determine if content is markdown (either by extension or content analysis)
⋮----
// Determine if content is HTML
⋮----
// Determine if content is code
⋮----
// Apply appropriate strategy based on content type
⋮----
} else if len(content) > chunkSize*5 { // Very long document
⋮----
// Default to paragraph-based chunking for general text
⋮----
// createSemanticChunks creates chunks based on semantic boundaries in the text
func (cs *ChunkerService) createSemanticChunks(doc *domain.Document, content string, chunkSize int, overlap int) []*domain.DocumentChunk
⋮----
// For semantic chunking, we prioritize natural text boundaries
// This is similar to paragraph-based but with more attention to headers and sections
⋮----
// Check if the document has headers (markdown-style or HTML-style)
⋮----
// If the document has headers, chunk by sections
⋮----
// Otherwise use paragraph chunking
⋮----
// createSectionBasedChunks splits content based on headers and sections
func (cs *ChunkerService) createSectionBasedChunks(doc *domain.Document, content string, chunkSize int, overlap int) []*domain.DocumentChunk
⋮----
// Find markdown headers or HTML headers
⋮----
// Split content by headers
⋮----
// Add a dummy header for the first section if it doesn't have one
⋮----
// Process each section
⋮----
// Skip empty sections
⋮----
// Combine header with its content
var sectionContent string
⋮----
// If section is too large, split it further
⋮----
// Create sub-chunks for this section
⋮----
// Update positions and indices
⋮----
// Create a single chunk for this section
⋮----
// createHierarchicalChunks creates a two-level chunking structure
func (cs *ChunkerService) createHierarchicalChunks(doc *domain.Document, content string, chunkSize int, overlap int) []*domain.DocumentChunk
⋮----
// For hierarchical chunking, we first split into major sections,
// then we further split each section if needed
⋮----
// Split into major sections (try headers first, fall back to large chunks)
⋮----
// Split by headers for major sections
⋮----
// Add a dummy header for the first section if it doesn't have one
⋮----
// Process each section
⋮----
// Skip empty sections
⋮----
// Combine header with its content
var sectionContent string
⋮----
// For each major section, create sub-chunks
⋮----
// If section is large enough to need sub-chunks
⋮----
// Create sub-chunks with paragraph-based approach
⋮----
// Update positions and indices for sub-chunks
⋮----
// No clear sections, create artificial major chunks
⋮----
// First create large parent chunks with minimal overlap
⋮----
// Try to break at paragraph boundaries
⋮----
// Then create smaller sub-chunks for each major chunk
⋮----
// Update positions and indices for sub-chunks
⋮----
// createMarkdownBasedChunks optimizes chunking for markdown documents
func (cs *ChunkerService) createMarkdownBasedChunks(doc *domain.Document, content string, chunkSize int, overlap int) []*domain.DocumentChunk
⋮----
// For markdown content, respect header structure
⋮----
// createHTMLBasedChunks optimizes chunking for HTML documents
func (cs *ChunkerService) createHTMLBasedChunks(doc *domain.Document, content string, chunkSize int, overlap int) []*domain.DocumentChunk
⋮----
// For HTML content, try to respect tag structure
// This is a simplified implementation - a full HTML parser would be more accurate
⋮----
// Look for major HTML structural elements
⋮----
// If we find structural elements, try to use them for chunking
⋮----
var chunks []*domain.DocumentChunk
⋮----
// Create chunks for content before first section if needed
⋮----
// Process each section
⋮----
// Handle gaps between sections
⋮----
// If section is too large, split it further
⋮----
// Strip HTML tags for better text chunking
⋮----
// Update positions and indices
⋮----
// Use section as a chunk
⋮----
// Handle content after the last section
⋮----
// Fall back to paragraph-based chunking if no clear structure is found
⋮----
// createCodeBasedChunks optimizes chunking for code documents
func (cs *ChunkerService) createCodeBasedChunks(doc *domain.Document, content string, chunkSize int, overlap int) []*domain.DocumentChunk
⋮----
// For code content, try to respect function/class boundaries
⋮----
// Look for function/class definitions in various languages
// This is a simplified approach - a language-specific parser would be more accurate
⋮----
// Split by function/class definitions
⋮----
// Handle content before first function if needed
⋮----
// Process each function match
⋮----
// Set end to the beginning of the next function
⋮----
// Handle content between last function end and current function start
⋮----
// Extract function content
⋮----
// If function is too large, split it further
⋮----
// Split by logical blocks (like try/catch, if/else)
⋮----
// Use function as a chunk
⋮----
// Fall back to line-based chunking for code with no clear structure
⋮----
// createFixedSizeChunks creates chunks of fixed size with overlap
func (cs *ChunkerService) createFixedSizeChunks(doc *domain.Document, content string, chunkSize int, overlap int) []*domain.DocumentChunk
⋮----
// Calculate total chunks needed
⋮----
// Create chunks with overlap
⋮----
// Adjust start position to avoid breaking words
⋮----
// Try to end at a natural break
⋮----
// Skip empty chunks
⋮----
// createParagraphBasedChunks creates chunks based on paragraph boundaries
func (cs *ChunkerService) createParagraphBasedChunks(doc *domain.Document, content string, chunkSize int, overlap int) []*domain.DocumentChunk
⋮----
// Split by paragraphs to maintain semantic units
⋮----
var currentChunk strings.Builder
⋮----
paraSize := len(para) + 2 // +2 for newlines
⋮----
// If this paragraph alone exceeds chunk size, we need to split it
⋮----
// If we have content in the current chunk, finalize it first
⋮----
// Now split the large paragraph into fixed-size chunks
⋮----
// Update startPos for next chunk
⋮----
// If adding this paragraph would exceed chunk size and we have content
⋮----
// Create a chunk from accumulated content
⋮----
// Handle overlap for the next chunk
⋮----
// Calculate where to start the next chunk with overlap
⋮----
// Start the new chunk with the end of the previous one
⋮----
// Add the paragraph to the current chunk
⋮----
// Handle the last paragraph
⋮----
// splitParagraphIntoChunks splits a single large paragraph into multiple chunks
func (cs *ChunkerService) splitParagraphIntoChunks(doc *domain.Document, paragraph string, chunkSize int, overlap int, startOffset int, chunkIndexOffset int) []*domain.DocumentChunk
⋮----
// For very large paragraphs, split by sentences if possible
⋮----
// If paragraph doesn't have clear sentences or has very few, use character chunking
⋮----
// Character-based chunking
⋮----
// Try not to break words
⋮----
// Sentence-based chunking for more semantic coherence
var currentChunk strings.Builder
⋮----
sentenceSize := len(sentence) + 2 // +2 for ". "
⋮----
// If adding this sentence exceeds the chunk size and we have content
⋮----
// Calculate new start position
⋮----
// Add the sentence to the current chunk
⋮----
// Handle the last sentence
````

## File: internal/service/embedding_cache.go
````go
package service
⋮----
import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)
⋮----
"crypto/sha256"
"encoding/hex"
"sync"
"time"
⋮----
// EmbeddingCache provides caching for embeddings to avoid regeneration of the same content
type EmbeddingCache struct {
	cache       map[string]CachedEmbedding
	mutex       sync.RWMutex
	maxSize     int           // Taille maximale du cache
	ttl         time.Duration // Durée de vie des entrées
	lastCleanup time.Time     // Dernière fois que le cache a été nettoyé
}
⋮----
maxSize     int           // Taille maximale du cache
ttl         time.Duration // Durée de vie des entrées
lastCleanup time.Time     // Dernière fois que le cache a été nettoyé
⋮----
// CachedEmbedding represents a cached embedding with metadata
type CachedEmbedding struct {
	Embedding  []float32
	CreatedAt  time.Time
	AccessedAt time.Time
	UseCount   int // Pour garder trace des éléments les plus utilisés
}
⋮----
UseCount   int // Pour garder trace des éléments les plus utilisés
⋮----
// NewEmbeddingCache creates a new embedding cache
func NewEmbeddingCache(maxSize int, ttl time.Duration) *EmbeddingCache
⋮----
// Get retrieves an embedding from the cache
func (c *EmbeddingCache) Get(text string, modelName string) ([]float32, bool)
⋮----
// Vérifier si l'entrée est expirée
⋮----
// Mettre à jour les statistiques d'accès (sans lock d'écriture pour la performance)
⋮----
// Set adds an embedding to the cache
func (c *EmbeddingCache) Set(text string, modelName string, embedding []float32)
⋮----
// Vérifier si le cache doit être nettoyé
⋮----
// Ajouter au cache
⋮----
// generateKey creates a unique key for the cache
func (c *EmbeddingCache) generateKey(text string, modelName string) string
⋮----
// cleanup removes expired or least used entries when cache is full
func (c *EmbeddingCache) cleanup()
⋮----
// Supprimer les entrées expirées
⋮----
// Si le cache est toujours trop grand, supprimer les entrées les moins utilisées
⋮----
type keyScore struct {
			key   string
			score float64 // Combined score (usage and recency)
		}
⋮----
score float64 // Combined score (usage and recency)
⋮----
// Calculate a score for each entry (combined usage and recency)
⋮----
usageScore := 1.0 / float64(1+entry.UseCount) // Higher usage = smaller score
combinedScore := recencyScore * usageScore    // Smaller = better
⋮----
// Sort by score
⋮----
// Remove entries with the highest score (least useful)
````

## File: internal/service/reranker_service_test.go
````go
package service
⋮----
import (
	"fmt"
	"testing"

	"github.com/dontizi/rlama/internal/domain"
	"github.com/stretchr/testify/assert"
)
⋮----
"fmt"
"testing"
⋮----
"github.com/dontizi/rlama/internal/domain"
"github.com/stretchr/testify/assert"
⋮----
// TestRerankerOptionsDefaultValues checks that the default values are correct
func TestRerankerOptionsDefaultValues(t *testing.T)
⋮----
// Get the default options
⋮----
// Check that the default values are correct
⋮----
// TestApplyTopKLimit checks that the TopK limit is applied correctly
func TestApplyTopKLimit(t *testing.T)
⋮----
// Create sorted results to simulate the output before applying TopK
⋮----
expected: 15, // Cannot return more than what exists
⋮----
expected: 10, // Should not limit if TopK=0
⋮----
// Apply the TopK limit manually (reproduce the logic of Rerank)
var limited []RankedResult
⋮----
// Check that the number is correct
⋮----
// createDummyRankedResults creates a set of dummy results for testing
func createDummyRankedResults(count int) []RankedResult
⋮----
// Reproduce the sorting function to test
func TestSortingByScore(t *testing.T)
⋮----
// Create results in a mixed order
⋮----
// Sort the results (same logic as in Rerank)
// Sort by final score (descending)
⋮----
// Check that the results are sorted correctly
⋮----
// Check the exact order
⋮----
// TestRerankerIntegration checks the integration of reranking in the RAG service
func TestRerankerIntegration(t *testing.T)
⋮----
// This test will integrate reranking in a complete RAG service
// As it requires external dependencies, it will be marked as an integration test
⋮----
// TODO: Implement an integration test with a real RAG service
// This can be done later by using the existing structs and functions
````

## File: scripts/install_deps.sh
````bash
#!/bin/bash

# Installation script for RLAMA dependencies
# This script attempts to install the necessary tools for text extraction
# and reranking with BGE

echo "Installing dependencies for RLAMA..."

# Operating system detection
OS=$(uname -s)
echo "Detected operating system: $OS"

# Function to check if a program is installed
is_installed() {
  command -v "$1" >/dev/null 2>&1
}

# macOS
if [ "$OS" = "Darwin" ]; then
  echo "Installing dependencies for macOS..."
  
  # Check if Homebrew is installed
  if ! is_installed brew; then
    echo "Homebrew not found. Installing Homebrew..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
  fi
  
  # Install tools
  echo "Installing text extraction tools..."
  brew install poppler  # For pdftotext
  brew install tesseract  # For OCR
  brew install tesseract-lang  # Additional languages for Tesseract
  
  # Python and tools
  if ! is_installed pip3; then
    brew install python
  fi
  
  pip3 install pdfminer.six docx2txt xlsx2csv
  
# Linux
elif [ "$OS" = "Linux" ]; then
  echo "Installing dependencies for Linux..."
  
  # Try apt-get (Debian/Ubuntu)
  if is_installed apt-get; then
    echo "Package manager apt-get detected"
    sudo apt-get update
    sudo apt-get install -y poppler-utils tesseract-ocr tesseract-ocr-fra python3-pip
    sudo apt-get install -y catdoc unrtf
  
  # Try yum (CentOS/RHEL)
  elif is_installed yum; then
    echo "Package manager yum detected"
    sudo yum update
    sudo yum install -y poppler-utils tesseract tesseract-langpack-fra python3-pip
    sudo yum install -y catdoc
  
  # Try pacman (Arch Linux)
  elif is_installed pacman; then
    echo "Package manager pacman detected"
    sudo pacman -Syu
    sudo pacman -S poppler tesseract tesseract-data-fra python-pip
  
  # Try zypper (openSUSE)
  elif is_installed zypper; then
    echo "Package manager zypper detected"
    sudo zypper refresh
    sudo zypper install poppler-tools tesseract-ocr python3-pip
  
  else
    echo "No known package manager detected. Please install the dependencies manually."
  fi
  
  # Install Python packages
  pip3 install --user pdfminer.six docx2txt xlsx2csv

# Windows (via WSL)
elif [[ "$OS" == MINGW* ]] || [[ "$OS" == MSYS* ]] || [[ "$OS" == CYGWIN* ]]; then
  echo "Windows system detected."
  echo "It is recommended to use WSL (Windows Subsystem for Linux) for better performance."
  echo "You can install the dependencies manually:"
  echo "1. Install Python: https://www.python.org/downloads/windows/"
  echo "2. Install Python packages: pip install pdfminer.six docx2txt xlsx2csv FlagEmbedding torch transformers"
  echo "3. For OCR, install Tesseract: https://github.com/UB-Mannheim/tesseract/wiki"
  
  # Try to install Python packages with pip in Windows
  if is_installed pip; then
    echo "Installing Python dependencies on Windows..."
    pip install --user pdfminer.six docx2txt xlsx2csv
    pip install --user -U FlagEmbedding torch transformers
  elif is_installed pip3; then
    echo "Installing Python dependencies on Windows..."
    pip3 install --user pdfminer.six docx2txt xlsx2csv
    pip3 install --user -U FlagEmbedding torch transformers
  fi
fi

# Install common Python dependencies
echo "Installing common Python dependencies..."
if is_installed pip3; then
  pip3 install --user pdfminer.six docx2txt xlsx2csv
  echo "Installing dependencies for BGE reranker..."
  pip3 install --user -U FlagEmbedding torch transformers
elif is_installed pip; then
  pip install --user pdfminer.six docx2txt xlsx2csv
  echo "Installing dependencies for BGE reranker..."
  pip install --user -U FlagEmbedding torch transformers
else
  echo "⚠️ Pip is not installed. Cannot install Python dependencies."
  echo "Please install pip then run: pip install -U FlagEmbedding pdfminer.six docx2txt xlsx2csv"
fi

echo "Installation completed!"
echo ""
echo "To use the BGE reranker, run: rlama update-reranker [rag-name]"
echo "This will configure your RAG to use the BAAI/bge-reranker-v2-m3 model for reranking."
````

## File: cmd/add_docs.go
````go
package cmd
⋮----
import (
	"fmt"

	"github.com/dontizi/rlama/internal/service"
	"github.com/spf13/cobra"
)
⋮----
"fmt"
⋮----
"github.com/dontizi/rlama/internal/service"
"github.com/spf13/cobra"
⋮----
var (
	addDocsExcludeDirs      []string
	addDocsExcludeExts      []string
	addDocsProcessExts      []string
	addDocsChunkSize        int
	addDocsChunkOverlap     int
	addDocsChunkingStrategy string
	addDocsDisableReranker  bool
	addDocsRerankerModel    string
	addDocsRerankerWeight   float64
)
⋮----
var addDocsCmd = &cobra.Command{
	Use:   "add-docs [rag-name] [folder-path]",
	Short: "Add documents to an existing RAG system",
	Long: `Add documents from a folder to an existing RAG system.
Example: rlama add-docs my-docs ./new-documents
	
This will load documents from the specified folder, generate embeddings,
and add them to the existing RAG system.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]
		folderPath := args[1]

		// Get service provider
		provider := GetServiceProvider()
		ragService := provider.GetRagService()

		// Set up loader options based on flags
		loaderOptions := service.DocumentLoaderOptions{
			ExcludeDirs:      addDocsExcludeDirs,
			ExcludeExts:      addDocsExcludeExts,
			ProcessExts:      addDocsProcessExts,
			ChunkSize:        addDocsChunkSize,
			ChunkOverlap:     addDocsChunkOverlap,
			ChunkingStrategy: addDocsChunkingStrategy,
			EnableReranker:   !addDocsDisableReranker,
			RerankerModel:    addDocsRerankerModel,
			RerankerWeight:   addDocsRerankerWeight,
		}

		// Pass the options to the service
		err := ragService.AddDocsWithOptions(ragName, folderPath, loaderOptions)
		if err != nil {
			return err
		}

		fmt.Printf("Documents from '%s' added to RAG '%s' successfully.\n", folderPath, ragName)
		return nil
	},
}
⋮----
// Get service provider
⋮----
// Set up loader options based on flags
⋮----
// Pass the options to the service
⋮----
func init()
⋮----
// Add exclusion and processing flags
⋮----
// Add chunking options
⋮----
// Add reranking options
````

## File: cmd/profile.go
````go
package cmd
⋮----
import (
	"fmt"
	"os"
	"text/tabwriter"

	// "time"

	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/repository"
	"github.com/spf13/cobra"
)
⋮----
"fmt"
"os"
"text/tabwriter"
⋮----
// "time"
⋮----
"github.com/dontizi/rlama/internal/domain"
"github.com/dontizi/rlama/internal/repository"
"github.com/spf13/cobra"
⋮----
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage API profiles",
	Long:  `Create, list, and manage API profiles for different providers.`,
}
⋮----
var profileAddCmd = &cobra.Command{
	Use:   "add [name] [provider] [api-key]",
	Short: "Add a new API profile",
	Long: `Add a new API profile for a specific provider.
Examples: 
  rlama profile add openai-work openai sk-...your-api-key...
  rlama profile add lmstudio openai-api none --base-url http://localhost:1234/v1`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		provider := args[1]
		apiKey := args[2]

		// Get base URL flag
		baseURL, _ := cmd.Flags().GetString("base-url")

		// Validate the provider
		switch provider {
		case "openai":
			// Official OpenAI API
		case "openai-api":
			// Generic OpenAI-compatible API
			if baseURL == "" {
				return fmt.Errorf("base-url is required for openai-api provider")
			}
		default:
			return fmt.Errorf("unsupported provider: %s. Supported providers: openai, openai-api", provider)
		}

		// Create the repository
		profileRepo := repository.NewProfileRepository()

		// Check if the profile already exists
		if profileRepo.Exists(name) {
			return fmt.Errorf("profile '%s' already exists", name)
		}

		// Create and save the profile
		profile := domain.NewAPIProfile(name, provider, apiKey)
		profile.BaseURL = baseURL
		if err := profileRepo.Save(profile); err != nil {
			return err
		}

		if baseURL != "" {
			fmt.Printf("Profile '%s' for '%s' (base URL: %s) added successfully.\n", name, provider, baseURL)
		} else {
			fmt.Printf("Profile '%s' for '%s' added successfully.\n", name, provider)
		}
		return nil
	},
}
⋮----
// Get base URL flag
⋮----
// Validate the provider
⋮----
// Official OpenAI API
⋮----
// Generic OpenAI-compatible API
⋮----
// Create the repository
⋮----
// Check if the profile already exists
⋮----
// Create and save the profile
⋮----
var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all API profiles",
	Long:  `Display a list of all configured API profiles.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		profileRepo := repository.NewProfileRepository()

		profiles, err := profileRepo.ListAll()
		if err != nil {
			return err
		}

		if len(profiles) == 0 {
			fmt.Println("No API profiles found.")
			return nil
		}

		fmt.Printf("Available API profiles (%d found):\n\n", len(profiles))

		// Use tabwriter to align the display
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tPROVIDER\tBASE URL\tCREATED ON\tLAST USED")

		for _, name := range profiles {
			profile, err := profileRepo.Load(name)
			if err != nil {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", name, "error", "error", "error", "error")
				continue
			}

			// Format dates
			createdAt := profile.CreatedAt.Format("2006-01-02 15:04:05")
			lastUsed := "never"
			if !profile.LastUsedAt.IsZero() {
				lastUsed = profile.LastUsedAt.Format("2006-01-02 15:04:05")
			}

			baseURL := profile.BaseURL
			if baseURL == "" {
				baseURL = "default"
			}

			// Hide the API key
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				profile.Name, profile.Provider, baseURL, createdAt, lastUsed)
		}
		w.Flush()

		return nil
	},
}
⋮----
// Use tabwriter to align the display
⋮----
// Format dates
⋮----
// Hide the API key
⋮----
var profileDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete an API profile",
	Long:  `Delete an API profile by name.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		profileRepo := repository.NewProfileRepository()

		// Check if the profile exists
		if !profileRepo.Exists(name) {
			return fmt.Errorf("profile '%s' does not exist", name)
		}

		// Ask for confirmation
		fmt.Printf("Are you sure you want to delete profile '%s'? (y/n): ", name)
		var response string
		fmt.Scanln(&response)

		if response != "y" && response != "Y" {
			fmt.Println("Deletion cancelled.")
			return nil
		}

		// Delete the profile
		if err := profileRepo.Delete(name); err != nil {
			return err
		}

		fmt.Printf("Profile '%s' deleted successfully.\n", name)
		return nil
	},
}
⋮----
// Check if the profile exists
⋮----
// Ask for confirmation
⋮----
var response string
⋮----
// Delete the profile
⋮----
func init()
⋮----
// Add flags for profile add command
````

## File: cmd/update_model.go
````go
package cmd
⋮----
import (
	"fmt"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/repository"
	"github.com/spf13/cobra"
)
⋮----
"fmt"
⋮----
"github.com/dontizi/rlama/internal/client"
"github.com/dontizi/rlama/internal/repository"
"github.com/spf13/cobra"
⋮----
var updateModelCmd = &cobra.Command{
	Use:   "update-model [rag-name] [new-model]",
	Short: "Update the Ollama model used by a RAG system",
	Long: `Change the Ollama model used by an existing RAG system.
Example: rlama update-model my-docs llama3.2
	
Note: This does not regenerate embeddings. For optimal results, you may want to
recreate the RAG with the new model instead.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]
		newModel := args[1]

		// Get Ollama client from root command
		ollamaClient := GetOllamaClient()

		// Check if this is a Hugging Face model
		if client.IsHuggingFaceModel(newModel) {
			// Extract model name and quantization
			hfModelName := client.GetHuggingFaceModelName(newModel)
			quantization := client.GetQuantizationFromModelRef(newModel)

			fmt.Printf("Detected Hugging Face model. Pulling %s", hfModelName)
			if quantization != "" {
				fmt.Printf(" with quantization %s", quantization)
			}
			fmt.Println("...")

			// Pull the model from Hugging Face
			if err := ollamaClient.PullHuggingFaceModel(hfModelName, quantization); err != nil {
				return fmt.Errorf("error pulling Hugging Face model: %w", err)
			}

			fmt.Println("Successfully pulled Hugging Face model.")
		} else {
			// For regular Ollama models
			if err := ollamaClient.CheckOllamaAndModel(newModel); err != nil {
				return err
			}
		}

		// Load the RAG
		repo := repository.NewRagRepository()
		rag, err := repo.Load(ragName)
		if err != nil {
			return err
		}

		// Update the model
		oldModel := rag.ModelName
		rag.ModelName = newModel

		// Save the RAG
		if err := repo.Save(rag); err != nil {
			return fmt.Errorf("error saving the RAG: %w", err)
		}

		fmt.Printf("Successfully updated RAG '%s' model from '%s' to '%s'.\n",
			ragName, oldModel, newModel)
		fmt.Println("Note: Embeddings have not been regenerated. For optimal results, consider recreating the RAG.")

		// Check if the profile exists if specified
		if updateModelProfileName != "" {
			profileRepo := repository.NewProfileRepository()
			if !profileRepo.Exists(updateModelProfileName) {
				return fmt.Errorf("profile '%s' does not exist", updateModelProfileName)
			}

			// Update the profile in the RAG
			rag.APIProfileName = updateModelProfileName
			fmt.Printf("Using profile '%s' for model '%s'\n", updateModelProfileName, newModel)
		}

		return nil
	},
}
⋮----
// Get Ollama client from root command
⋮----
// Check if this is a Hugging Face model
⋮----
// Extract model name and quantization
⋮----
// Pull the model from Hugging Face
⋮----
// For regular Ollama models
⋮----
// Load the RAG
⋮----
// Update the model
⋮----
// Save the RAG
⋮----
// Check if the profile exists if specified
⋮----
// Update the profile in the RAG
⋮----
var updateModelProfileName string
⋮----
func init()
````

## File: internal/client/llm_client.go
````go
package client
⋮----
// LLMClient is a common interface for language model clients
type LLMClient interface {
	GenerateCompletion(model, prompt string) (string, error)
	GenerateEmbedding(model, text string) ([]float32, error)
	CheckLLMAndModel(modelName string) error
}
⋮----
// Adapt existing methods of OllamaClient to implement LLMClient
func (c *OllamaClient) CheckLLMAndModel(modelName string) error
⋮----
// Adapt OpenAIClient methods to implement LLMClient
⋮----
// IsOpenAIModel checks if a model is an OpenAI model
func IsOpenAIModel(modelName string) bool
⋮----
// OpenAI models generally start with "gpt-" or "text-"
⋮----
// StartsWith checks if a string starts with a prefix
func StartsWith(s, prefix string) bool
⋮----
// GetLLMClient returns the appropriate client based on the model
func GetLLMClient(modelName string, ollamaClient *OllamaClient) LLMClient
⋮----
// GetLLMClientWithProfile returns the appropriate client based on profile or model
func GetLLMClientWithProfile(modelName, profileName string, ollamaClient *OllamaClient) (LLMClient, error)
⋮----
// If a profile is specified, use it
⋮----
// Otherwise fall back to model-based selection
⋮----
// GetLLMClientFromProfile returns a client based on the specified profile
func GetLLMClientFromProfile(profileName string) (LLMClient, error)
````

## File: internal/client/ollama_client.go
````go
package client
⋮----
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)
⋮----
"bytes"
"encoding/json"
"fmt"
"io/ioutil"
"net/http"
"os"
"os/exec"
"strings"
⋮----
// Default connection settings for Ollama
const (
	DefaultOllamaHost = "localhost"
	DefaultOllamaPort = "11434"
)
⋮----
// OllamaClient is a client for the Ollama API
type OllamaClient struct {
	BaseURL string
	Client  *http.Client
}
⋮----
// EmbeddingRequest is the structure of the request for the /api/embeddings API
type EmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}
⋮----
// EmbeddingResponse is the structure of the response for the /api/embeddings API
type EmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}
⋮----
// GenerationRequest is the structure of the request for the /api/generate API
type GenerationRequest struct {
	Model    string  `json:"model"`
	Prompt   string  `json:"prompt"`
	Context  []int   `json:"context,omitempty"`
	Options  Options `json:"options,omitempty"`
	Format   string  `json:"format,omitempty"`
	Template string  `json:"template,omitempty"`
	Stream   bool    `json:"stream"`
}
⋮----
// Options for the /api/generate API
type Options struct {
	Temperature float64 `json:"temperature,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
	TopK        int     `json:"top_k,omitempty"`
	NumPredict  int     `json:"num_predict,omitempty"`
}
⋮----
// GenerationResponse is the structure of the response for the /api/generate API
type GenerationResponse struct {
	Model     string `json:"model"`
	Response  string `json:"response"`
	Context   []int  `json:"context"`
	CreatedAt string `json:"created_at"`
	Done      bool   `json:"done"`
}
⋮----
// NewOllamaClient creates a new Ollama client
// If host or port are empty, the default values are used
// If OLLAMA_HOST is defined, it is used as the default value
func NewOllamaClient(host, port string) *OllamaClient
⋮----
// Check for OLLAMA_HOST environment variable
⋮----
// Default values and protocol
⋮----
// If OLLAMA_HOST is set, parse it
⋮----
// Handle if OLLAMA_HOST includes protocol
⋮----
// Extract protocol and host
⋮----
// No protocol specified, use standard pattern
⋮----
// Command flags override environment variables
⋮----
// Check if host includes protocol
⋮----
// NewDefaultOllamaClient creates a new Ollama client with the default values
// Kept for compatibility with existing code
func NewDefaultOllamaClient() *OllamaClient
⋮----
// GenerateEmbedding generates an embedding for the given text
func (c *OllamaClient) GenerateEmbedding(model, text string) ([]float32, error)
⋮----
var embeddingResp EmbeddingResponse
⋮----
// GenerateCompletion generates a response for the given prompt
func (c *OllamaClient) GenerateCompletion(model, prompt string) (string, error)
⋮----
var genResp GenerationResponse
⋮----
// IsOllamaRunning checks if Ollama is installed and running
func (c *OllamaClient) IsOllamaRunning() (bool, error)
⋮----
// CheckOllamaAndModel verifies if Ollama is running and if the specified model is available
func (c *OllamaClient) CheckOllamaAndModel(modelName string) error
⋮----
// Check if Ollama is running
⋮----
// Check if model is available (optional)
// This check could be added here
⋮----
// RunHuggingFaceModel prepares a Hugging Face model for use with Ollama
func (c *OllamaClient) RunHuggingFaceModel(hfModelPath string, quantization string) error
⋮----
// PullHuggingFaceModel pulls a Hugging Face model into Ollama without running it
func (c *OllamaClient) PullHuggingFaceModel(hfModelPath string, quantization string) error
⋮----
// IsHuggingFaceModel checks if a model name is a Hugging Face model reference
func IsHuggingFaceModel(modelName string) bool
⋮----
// GetHuggingFaceModelName extracts the repository name from a Hugging Face model reference
func GetHuggingFaceModelName(modelRef string) string
⋮----
// Strip any prefix
⋮----
// Strip any quantization suffix
⋮----
// GetQuantizationFromModelRef extracts the quantization suffix from a model reference
func GetQuantizationFromModelRef(modelRef string) string
````

## File: internal/client/openai_client.go
````go
package client
⋮----
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/dontizi/rlama/internal/repository"
)
⋮----
"bytes"
"encoding/json"
"fmt"
"io/ioutil"
"net/http"
"os"
"time"
⋮----
"github.com/dontizi/rlama/internal/repository"
⋮----
// OpenAIClient is a client for the OpenAI API
type OpenAIClient struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}
⋮----
// OpenAICompletionRequest is the structure for completion requests to OpenAI
type OpenAICompletionRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	Temperature float64         `json:"temperature,omitempty"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
}
⋮----
// OpenAIMessage represents a message in the format expected by OpenAI
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
⋮----
// OpenAICompletionResponse is the structure of the OpenAI API response
type OpenAICompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int           `json:"index"`
		Message      OpenAIMessage `json:"message"`
		FinishReason string        `json:"finish_reason"`
	} `json:"choices"`
⋮----
// OpenAIEmbeddingRequest is the structure for embedding requests to OpenAI
type OpenAIEmbeddingRequest struct {
	Input string `json:"input"`
	Model string `json:"model"`
}
⋮----
// OpenAIEmbeddingResponse is the structure of the OpenAI embedding API response
type OpenAIEmbeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
⋮----
// NewOpenAIClient creates a new OpenAI client for the official API
func NewOpenAIClient() *OpenAIClient
⋮----
// Use API key from environment
⋮----
// NewGenericOpenAIClient creates a new OpenAI-compatible client with custom base URL
func NewGenericOpenAIClient(baseURL, apiKey string) *OpenAIClient
⋮----
// NewOpenAIClientWithProfile creates a new OpenAI client with a specific profile
func NewOpenAIClientWithProfile(profileName string) (*OpenAIClient, error)
⋮----
// Create a profile repository
⋮----
// If no profile is specified, use the environment variable
⋮----
// Load the specified profile
⋮----
// Check that it's an OpenAI or OpenAI-compatible profile
⋮----
// Update last used date
⋮----
// Use BaseURL from profile if available, otherwise default
⋮----
// GenerateCompletion generates a response from a prompt with OpenAI
func (c *OpenAIClient) GenerateCompletion(model, prompt string) (string, error)
⋮----
// Note: API key may be empty for local OpenAI-compatible servers
⋮----
// Build the request
⋮----
Temperature: 0.7, // Default value, can be configured
⋮----
// Create the HTTP request
⋮----
// Add necessary headers
⋮----
// Send the request
⋮----
// Check status code
⋮----
// Decode the response
var completionResp OpenAICompletionResponse
⋮----
// Check that there is at least one choice
⋮----
// Return the content of the response
⋮----
// GenerateEmbedding generates an embedding for the given text using OpenAI
func (c *OpenAIClient) GenerateEmbedding(model, text string) ([]float32, error)
⋮----
var embeddingResp OpenAIEmbeddingResponse
⋮----
// Check that there is at least one embedding
⋮----
// Return the embedding
⋮----
// CheckOpenAIAndModel checks if OpenAI is accessible and if the model is available
func (c *OpenAIClient) CheckOpenAIAndModel(modelName string) error
⋮----
// Only require API key for official OpenAI endpoint
⋮----
// We could check the validity of the model here
// but for now, we assume the model is valid if the API key is set
````

## File: internal/domain/document.go
````go
package domain
⋮----
import (
	"path/filepath"
	"regexp"
	"strings"
	"time"
)
⋮----
"path/filepath"
"regexp"
"strings"
"time"
⋮----
// Document represents a document indexed in the RAG system
type Document struct {
	ID          string    `json:"id"`
	Path        string    `json:"path"`
	Name        string    `json:"name"`
	Content     string    `json:"content"`
	Metadata    string    `json:"metadata"`
	Embedding   []float32 `json:"-"` // Do not serialize to JSON
	CreatedAt   time.Time `json:"created_at"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	URL         string    `json:"url,omitempty"` // Source URL for web documents
}
⋮----
Embedding   []float32 `json:"-"` // Do not serialize to JSON
⋮----
URL         string    `json:"url,omitempty"` // Source URL for web documents
⋮----
// NewDocument creates a new instance of Document
func NewDocument(path string, content string) *Document
⋮----
// Clean the extracted content
⋮----
// cleanExtractedText cleans the extracted text to improve its quality
func cleanExtractedText(text string) string
⋮----
// Replace non-printable characters with spaces
⋮----
// Replace sequences of more than 2 newlines with 2 newlines
⋮----
// Replace sequences of more than 2 spaces with 1 space
⋮----
// Remove lines that contain only special characters or numbers
⋮----
var cleanedLines []string
⋮----
// Check if the line contains at least some letters
⋮----
// guessContentType tries to determine the content type based on the file extension
func guessContentType(path string) string
````

## File: internal/domain/profile.go
````go
package domain
⋮----
import (
	"time"
)
⋮----
"time"
⋮----
// APIProfile represents a profile for API keys
type APIProfile struct {
	Name       string    `json:"name"`
	Provider   string    `json:"provider"` // "openai", "openai-api", "anthropic", etc.
	APIKey     string    `json:"api_key"`
	BaseURL    string    `json:"base_url,omitempty"` // For custom OpenAI-compatible endpoints
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	LastUsedAt time.Time `json:"last_used_at,omitempty"`
}
⋮----
Provider   string    `json:"provider"` // "openai", "openai-api", "anthropic", etc.
⋮----
BaseURL    string    `json:"base_url,omitempty"` // For custom OpenAI-compatible endpoints
⋮----
// NewAPIProfile creates a new API profile
func NewAPIProfile(name, provider, apiKey string) *APIProfile
````

## File: internal/service/rag_service_test.go
````go
package service
⋮----
import (
	"testing"

	"github.com/dontizi/rlama/internal/domain"
	"github.com/stretchr/testify/assert"
)
⋮----
"testing"
⋮----
"github.com/dontizi/rlama/internal/domain"
"github.com/stretchr/testify/assert"
⋮----
// TestRagRerankerTopK checks that the reranking is configured correctly and limits the results to 5 by default
func TestRagRerankerTopK(t *testing.T)
⋮----
// Create a RAG with a custom model using the constructor
⋮----
// Check that the default reranking values are correct
⋮----
// Check that the default reranking options are consistent
⋮----
// Test with different TopK values
⋮----
topK:     0, // 0 means use the default value
⋮----
topK:     -1, // Invalid value, should use the default of the RAG
⋮----
// Simulate the logic of Query() to determine the context size
⋮----
// If contextSize is 0 (auto), use:
// - RerankerTopK of the RAG if defined
// - Otherwise the default TopK (5)
// - 20 if reranking is disabled
⋮----
contextSize = options.TopK // 5 by default
⋮----
contextSize = 20 // 20 by default if reranking is disabled
⋮----
// Check that contextSize corresponds to the expected value
⋮----
// Test the case where reranking is disabled
⋮----
// Context size set to 0 should default to 20 because reranking is disabled
⋮----
contextSize = options.TopK // 5 by default
⋮----
contextSize = 20 // 20 by default if reranking is disabled
⋮----
// Restore the state
````

## File: internal/crawler/crawler.go
````go
package crawler
⋮----
import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dontizi/rlama/internal/domain"
	"golang.org/x/net/html/charset"
)
⋮----
"fmt"
"net/http"
"net/url"
"strings"
"sync"
"time"
⋮----
"github.com/PuerkitoBio/goquery"
"github.com/dontizi/rlama/internal/domain"
"golang.org/x/net/html/charset"
⋮----
// WebCrawler manage web crawling operations
type WebCrawler struct {
	client       *http.Client
	baseURL      *url.URL
	maxDepth     int
	concurrency  int
	excludePaths []string
	visited      map[string]bool
	visitedMutex sync.Mutex
	useSitemap   bool     // Option to use sitemap
	singleURL    bool     // Option to crawl only the specified URL
	urlsList     []string // Custom list of URLs to crawl
}
⋮----
useSitemap   bool     // Option to use sitemap
singleURL    bool     // Option to crawl only the specified URL
urlsList     []string // Custom list of URLs to crawl
⋮----
// NewWebCrawler creates a new web crawler
func NewWebCrawler(urlStr string, maxDepth, concurrency int, excludePaths []string) (*WebCrawler, error)
⋮----
useSitemap:   true,  // By default, use sitemap if available
singleURL:    false, // By default, do normal crawling
urlsList:     nil,   // By default, no custom list
⋮----
// isWebContent checks if a URL points to text/HTML content rather than binary files
func isWebContent(urlStr string) bool
⋮----
// Extensions to ignore (binary files, etc.)
⋮----
// CrawlWebsite crawls the website and returns the documents
func (wc *WebCrawler) CrawlWebsite() ([]domain.Document, error)
⋮----
// If single URL mode, only crawl the base URL
⋮----
// If custom list of URLs, use this list
⋮----
// Otherwise, normal behavior with sitemap or standard crawling
// Try to find a sitemap first
⋮----
// If no sitemap or option disabled, continue with standard crawling
⋮----
// crawlSingleURL crawls only the base URL without following any links
func (wc *WebCrawler) crawlSingleURL() ([]domain.Document, error)
⋮----
var documents []domain.Document
⋮----
// Fetch and parse the single URL
⋮----
// crawlURLsList crawls the specific list of URLs provided by the user
func (wc *WebCrawler) crawlURLsList() ([]domain.Document, error)
⋮----
var wg sync.WaitGroup
var mu sync.Mutex
⋮----
// Check if the URL should be excluded
⋮----
// Use the existing URL crawling function
⋮----
// Log any errors but continue with the documents we have
⋮----
// crawlStandard performs the standard crawling
func (wc *WebCrawler) crawlStandard() ([]domain.Document, error)
⋮----
// Don't crawl deeper if we've reached the maximum depth
⋮----
// Find the links on the page
⋮----
// extractLinks gets all valid links from a page
func (wc *WebCrawler) extractLinks(urlStr string) ([]string, error)
⋮----
var links []string
⋮----
// Convert to absolute URL
⋮----
// Check if the URL is on the same domain
⋮----
// Check the exclusions
⋮----
// resolveURL converts a relative URL to absolute
func (wc *WebCrawler) resolveURL(href string) (string, error)
⋮----
// isSameDomain checks if a URL is on the same domain as the base URL
func (wc *WebCrawler) isSameDomain(urlStr string) bool
⋮----
// convertToMarkdown converts HTML content to Markdown
func (wc *WebCrawler) convertToMarkdown(doc *goquery.Document) string
⋮----
// Remove unwanted elements
⋮----
// Get the main content
var content string
⋮----
// Basic cleanup
⋮----
// fetchAndParseURL fetches and parses a single URL
func (wc *WebCrawler) fetchAndParseURL(urlStr string) (*domain.Document, error)
⋮----
// Use convertToMarkdown instead of extractMarkdownFromHTML
⋮----
// getRelativePath returns the relative path of a URL to the base URL
func (wc *WebCrawler) getRelativePath(urlStr string) string
⋮----
// extractContentAsMarkdown extracts main content from an HTML document and converts it to Markdown
func extractContentAsMarkdown(doc *goquery.Document) (string, error)
⋮----
// Create a Crawl4AI style converter
⋮----
// Use the enhanced converter for HTML to Markdown conversion
⋮----
// SetUseSitemap sets whether to use sitemap for crawling
func (wc *WebCrawler) SetUseSitemap(useSitemap bool)
⋮----
// SetSingleURLMode sets whether to crawl only the specified URL without following links
func (wc *WebCrawler) SetSingleURLMode(singleURL bool)
⋮----
// SetURLsList sets a custom list of URLs to crawl
func (wc *WebCrawler) SetURLsList(urlsList []string)
⋮----
// parseSitemap parses a sitemap XML and returns the list of URLs
func (wc *WebCrawler) parseSitemap(sitemapURL string) ([]string, error)
⋮----
// Use goquery to parse the XML
⋮----
var urls []string
⋮----
// Find all <loc> tags in the sitemap
⋮----
// crawlURLsFromSitemap crawls all URLs found in the sitemap
func (wc *WebCrawler) crawlURLsFromSitemap(urls []string) ([]domain.Document, error)
⋮----
// Mark as visited
````

## File: internal/service/embedding_service.go
````go
package service
⋮----
import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/domain"
)
⋮----
"fmt"
"os"
"os/exec"
"sync"
⋮----
"github.com/dontizi/rlama/internal/client"
"github.com/dontizi/rlama/internal/domain"
⋮----
// EmbeddingService manages the generation of embeddings for documents
type EmbeddingService struct {
	llmClient          client.LLMClient
	maxWorkers         int    // Number of parallel workers for embedding generation
	preferredEmbedding string // Preferred embedding model to try first
}
⋮----
maxWorkers         int    // Number of parallel workers for embedding generation
preferredEmbedding string // Preferred embedding model to try first
⋮----
// NewEmbeddingService creates a new instance of EmbeddingService
func NewEmbeddingService(llmClient client.LLMClient) *EmbeddingService
⋮----
maxWorkers: 3, // Default to 3 workers
⋮----
// SetMaxWorkers sets the maximum number of parallel workers for embedding generation
func (es *EmbeddingService) SetMaxWorkers(workers int)
⋮----
// GetLLMClient returns the underlying LLM client
func (es *EmbeddingService) GetLLMClient() client.LLMClient
⋮----
// SetPreferredEmbeddingModel sets the preferred embedding model to try first
func (es *EmbeddingService) SetPreferredEmbeddingModel(model string)
⋮----
// GenerateEmbeddings generates embeddings for a list of documents
func (es *EmbeddingService) GenerateEmbeddings(docs []*domain.Document, modelName string) error
⋮----
// Determine which embedding model to try first
var embeddingModel string
⋮----
// Process all documents
⋮----
// Generate embedding with the preferred/default model first
⋮----
// If snowflake-arctic-embed2 fails, try to pull it automatically (Ollama only)
⋮----
// Attempt to pull the embedding model automatically (only for Ollama clients)
var pullErr error
⋮----
// Try again with the pulled model
⋮----
// If pulling failed or embedding still fails, fallback to the specified model
⋮----
// DetectEmbeddingDimension detects the dimension of embeddings from the model
func (es *EmbeddingService) DetectEmbeddingDimension(modelName string) (int, error)
⋮----
// Generate a test embedding to detect dimension
⋮----
// Try preferred embedding model first
⋮----
// Fallback to the main model
⋮----
// GenerateQueryEmbedding generates an embedding for a query
func (es *EmbeddingService) GenerateQueryEmbedding(query string, modelName string) ([]float32, error)
⋮----
// Generate embedding with the preferred/default model
⋮----
// If snowflake-arctic-embed2 fails, try to pull it (but only if not already tried)
⋮----
// We don't need to show the warning again if already shown in GenerateEmbeddings
// Attempt to pull the model (this is a no-op if we already tried, and only for Ollama)
var pullErr error
⋮----
// Try again with the pulled model
⋮----
// If pulling failed or embedding still fails, fallback to the specified model
⋮----
// GenerateChunkEmbeddings generates embeddings for document chunks in parallel
func (es *EmbeddingService) GenerateChunkEmbeddings(chunks []*domain.DocumentChunk, modelName string) error
⋮----
// Create a wait group to synchronize goroutines
var wg sync.WaitGroup
⋮----
// Create a channel to limit concurrency
⋮----
// Create a channel for errors
⋮----
// Create a mutex for printing progress
var progressMutex sync.Mutex
var completedChunks int
⋮----
// Check if we need to pull the model (attempt only once)
⋮----
var modelCheckMutex sync.Mutex
⋮----
// Process chunks in parallel
⋮----
// Add to wait group before starting goroutine
⋮----
// Start a goroutine to process this chunk
⋮----
// Acquire semaphore slot (this limits concurrency)
⋮----
// Generate embedding
⋮----
// If the model fails and we haven't checked it yet
⋮----
// Only print the warning and attempt to pull once
⋮----
var pullErr error
⋮----
// Try again with the pulled model
⋮----
// Use the specified model instead if the embedding model failed
⋮----
// Update the chunk with the embedding
⋮----
// Update progress
⋮----
// Wait for all goroutines to complete
⋮----
// Check if any errors occurred
⋮----
return err // Return the first error encountered
⋮----
fmt.Println() // Add a newline after progress indicator
⋮----
// Track if we've already tried to pull the model to avoid multiple attempts
var attemptedModelPull = make(map[string]bool)
⋮----
// pullEmbeddingModel attempts to pull the embedding model via Ollama
func (es *EmbeddingService) pullEmbeddingModel(modelName string) error
⋮----
// Check if we've already tried to pull this model
⋮----
// Mark that we've attempted to pull this model
⋮----
// Check if Ollama CLI is available
⋮----
// Run the ollama pull command
⋮----
cmd.Stdout = os.Stdout // Display output to the user
````

## File: internal/service/reranker_service.go
````go
package service
⋮----
import (
	"fmt"
	"sort"
	"strings"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/pkg/vector"
)
⋮----
"fmt"
"sort"
"strings"
⋮----
"github.com/dontizi/rlama/internal/client"
"github.com/dontizi/rlama/internal/domain"
"github.com/dontizi/rlama/pkg/vector"
⋮----
// RerankerOptions defines configuration options for the reranking process
type RerankerOptions struct {
	// TopK is the number of documents to return after reranking
	TopK int

	// InitialK is the number of documents to retrieve from the initial search
	// before reranking (should be >= TopK)
	InitialK int

	// RerankerModel is the model to use for reranking
	// If empty, will default to the same model used for embedding
	RerankerModel string

	// ScoreThreshold is the minimum relevance score (0-1) for a document to be included
	// Documents with scores below this threshold will be filtered out
	ScoreThreshold float64

	// RerankerWeight defines the weight of the reranker score vs vector similarity
	// 0.0 = use only vector similarity, 1.0 = use only reranker scores
	RerankerWeight float64

	// AdaptiveFiltering when true, uses content relevance to select chunks
	// rather than a fixed top-k approach
	AdaptiveFiltering bool
	
	// Silent suppresses warnings and informational output from the reranker
	Silent bool
}
⋮----
// TopK is the number of documents to return after reranking
⋮----
// InitialK is the number of documents to retrieve from the initial search
// before reranking (should be >= TopK)
⋮----
// RerankerModel is the model to use for reranking
// If empty, will default to the same model used for embedding
⋮----
// ScoreThreshold is the minimum relevance score (0-1) for a document to be included
// Documents with scores below this threshold will be filtered out
⋮----
// RerankerWeight defines the weight of the reranker score vs vector similarity
// 0.0 = use only vector similarity, 1.0 = use only reranker scores
⋮----
// AdaptiveFiltering when true, uses content relevance to select chunks
// rather than a fixed top-k approach
⋮----
// Silent suppresses warnings and informational output from the reranker
⋮----
// DefaultRerankerOptions returns the default options for reranking
func DefaultRerankerOptions() RerankerOptions
⋮----
RerankerWeight:    0.7, // 70% reranker score, 30% vector similarity
⋮----
// RerankerClient interface for different reranker implementations
type RerankerClient interface {
	ComputeScores(pairs [][]string, normalize bool) ([]float64, error)
	GetModelName() string
}
⋮----
// CleanupableRerankerClient interface for clients that need cleanup
type CleanupableRerankerClient interface {
	RerankerClient
	Cleanup() error
}
⋮----
// RerankerService handles document reranking using cross-encoder models
type RerankerService struct {
	ollamaClient      *client.OllamaClient
	bgeRerankerClient RerankerClient
	useONNX           bool
}
⋮----
// NewRerankerService creates a new instance of RerankerService
func NewRerankerService(ollamaClient *client.OllamaClient) *RerankerService
⋮----
// Create the BGE reranker client with the default model (Python implementation)
⋮----
// NewRerankerServiceWithOptions creates a new instance of RerankerService with configuration options
func NewRerankerServiceWithOptions(ollamaClient *client.OllamaClient, useONNX bool, onnxModelDir string) *RerankerService
⋮----
var bgeRerankerClient RerankerClient
⋮----
// Create ONNX-based reranker client
⋮----
// Create Python-based reranker client
⋮----
// RankedResult represents a document with its relevance score after reranking
type RankedResult struct {
	Chunk         *domain.DocumentChunk
	VectorScore   float64
	RerankerScore float64
	FinalScore    float64
}
⋮----
// Rerank takes initial retrieval results and reruns them through a cross-encoder for more accurate ranking
func (rs *RerankerService) Rerank(
	query string,
	rag *domain.RagSystem,
	initialResults []vector.SearchResult,
	options RerankerOptions,
) ([]RankedResult, error)
⋮----
// Create an empty result if no documents were found
⋮----
// Always use BGE Reranker if available
⋮----
// Use the BGE model configured in the client
⋮----
// Code to perform reranking with BGE
⋮----
// Prepare the pairs for batch processing
⋮----
// Get scores
⋮----
// In case of failure, return to standard model
⋮----
// Process scores and return results
⋮----
// Calculate final score as weighted combination of vector and reranker scores
⋮----
// Add to results if above threshold
⋮----
// Sort by final score (descending)
⋮----
// Only apply Top-K limit if we're not using adaptive filtering
⋮----
// If BGE is not available or failed, fall back to the standard model
// Use the model specified in options or the one from RAG
⋮----
// Check if the model is a BGE reranker model
⋮----
var rankedResults []RankedResult
⋮----
// Always recreate the BGE client with the current silent setting
⋮----
// Use BGE reranker for BGE models
⋮----
// Get all scores at once using the BGE reranker
scores, err := rs.bgeRerankerClient.ComputeScores(pairs, true) // normalize=true to get 0-1 scores
⋮----
// Process the scores
⋮----
// Calculate final score as weighted combination of vector and reranker scores
⋮----
// Add to results if above threshold
⋮----
// Use the existing Ollama-based reranker for other models
// Prepare the cross-encoder prompt template
var promptTemplate string
⋮----
// Enhanced prompt for adaptive content filtering
⋮----
// Original prompt for standard reranking
⋮----
// Get chunks and score them
⋮----
// Prepare prompt for this chunk
⋮----
// Get reranking score from the model
⋮----
// Parse the response as a float (score)
var score float64
⋮----
// If parsing fails, use a default score (based on vector similarity)
⋮----
// Ensure score is in range [0,1]
⋮----
// Sort by final score (descending)
⋮----
// Only apply Top-K limit if we're not using adaptive filtering
⋮----
// Cleanup properly cleans up resources if the reranker client supports it
func (rs *RerankerService) Cleanup() error
⋮----
// IsUsingONNX returns true if the service is using ONNX implementation
func (rs *RerankerService) IsUsingONNX() bool
````

## File: internal/service/web_watcher.go
````go
package service
⋮----
import (
	"fmt"
	"strings"
	"time"

	"github.com/dontizi/rlama/internal/crawler"
	"github.com/dontizi/rlama/internal/domain"
)
⋮----
"fmt"
"strings"
"time"
⋮----
"github.com/dontizi/rlama/internal/crawler"
"github.com/dontizi/rlama/internal/domain"
⋮----
// WebWatcher is responsible for watching websites for content changes
type WebWatcher struct {
	ragService RagService
}
⋮----
// NewWebWatcher creates a new web watcher service
func NewWebWatcher(ragService RagService) *WebWatcher
⋮----
// CheckAndUpdateRag checks for new content on the watched website and updates the RAG
func (ww *WebWatcher) CheckAndUpdateRag(rag *domain.RagSystem) (int, error)
⋮----
return 0, nil // Watching not enabled
⋮----
// Create a webcrawler to fetch the site content
⋮----
// Start crawling
⋮----
// Update last watched time even if no new documents
⋮----
// Ensure all documents have a valid URL
var validDocuments []*domain.Document // Changed to use pointers
⋮----
doc := &documents[i] // Get the address of the document
⋮----
// Build a URL based on the path or a unique identifier
⋮----
// Get existing document URLs and content hashes
⋮----
// Filter documents to keep only the new ones
var newDocuments []*domain.Document
⋮----
doc := validDocuments[i] // doc is already a pointer
⋮----
// Debug logging
⋮----
// Check both the URL and the content
⋮----
// Add to the list to avoid duplicates in this session
⋮----
// If no new documents after filtering, update the timestamp and terminate
⋮----
// Process the crawled documents directly without going through the file system
// Create chunker service
⋮----
var allChunks []*domain.DocumentChunk
var processedDocs []*domain.Document
⋮----
// Process each new document directly
⋮----
// Create a unique ID based on the URL
⋮----
// Ensure the URL is preserved
⋮----
// Add to the list of processed documents
⋮----
// Chunk the document
⋮----
// Update the chunk metadata
⋮----
// Generate embeddings for all chunks
⋮----
// Add the documents and chunks to the RAG
⋮----
// Update last watched time
⋮----
// Save the updated RAG
⋮----
// Function to normalize URLs (remove trailing slashes, etc.)
func normalizeURL(url string) string
⋮----
// Remove the trailing slash if it exists
⋮----
// Convert to lowercase
⋮----
// Other normalizations if needed...
⋮----
// Function to generate a simple hash of the content
func getContentHash(content string) string
⋮----
// Simplify the content for comparison (remove spaces, etc.)
⋮----
// If the content is very short, use the entire content
⋮----
// For longer content, take the beginning and the end
// for better identification
⋮----
// StartWebWatcherDaemon starts a background daemon to watch websites
func (ww *WebWatcher) StartWebWatcherDaemon(interval time.Duration)
⋮----
// checkAllRags checks all RAGs with web watching enabled
func (ww *WebWatcher) checkAllRags()
⋮----
// Get all RAGs
⋮----
// Check if web watching is enabled and if interval has passed
````

## File: cmd/crawl_add_docs.go
````go
package cmd
⋮----
import (
	"fmt"

	"github.com/dontizi/rlama/internal/crawler"
	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/service"
	"github.com/spf13/cobra"
)
⋮----
"fmt"
⋮----
"github.com/dontizi/rlama/internal/crawler"
"github.com/dontizi/rlama/internal/domain"
"github.com/dontizi/rlama/internal/service"
"github.com/spf13/cobra"
⋮----
var (
	addCrawlMaxDepth         int
	addCrawlConcurrency      int
	addCrawlExcludePaths     []string
	addCrawlUseSitemap       bool
	addCrawlSingleURL        bool
	addCrawlURLsList         []string
	addCrawlChunkSize        int
	addCrawlChunkOverlap     int
	addCrawlChunkingStrategy string
	addCrawlDisableReranker  bool
	addCrawlRerankerModel    string
	addCrawlRerankerWeight   float64
)
⋮----
var crawlAddDocsCmd = &cobra.Command{
	Use:   "crawl-add-docs [rag-name] [website-url]",
	Short: "Add website content to an existing RAG system",
	Long: `Add content from a website to an existing RAG system.
Example: rlama crawl-add-docs my-docs https://blog.example.com
	
This will crawl the website, extract content, generate embeddings,
and add them to the existing RAG system.

Control the crawling behavior with these flags:
  --max-depth=3         Maximum depth of pages to crawl
  --concurrency=10      Number of concurrent crawlers
  --exclude-path=/tag   Skip specific path patterns (comma-separated)
  --use-sitemap         Use sitemap.xml if available for comprehensive coverage
  --single-url          Process only the specified URL without following links
  --urls-list=url1,url2 Provide a comma-separated list of specific URLs to crawl
  --chunk-size=1000     Character count per chunk
  --chunk-overlap=200   Overlap between chunks in characters
  --chunking-strategy=hybrid  Chunking strategy to use (fixed, semantic, hybrid, hierarchical)
  --disable-reranker    Disable reranking for this content
  --reranker-model=model  Model to use for reranking
  --reranker-weight=0.7   Weight for reranker scores vs vector scores (0-1)`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]
		websiteURL := args[1]

		// Get Ollama client from root command
		ollamaClient := GetOllamaClient()

		// Create necessary services
		ragService := service.NewRagService(ollamaClient)

		// Load existing RAG to get model name
		_, err := ragService.LoadRag(ragName)
		if err != nil {
			return err
		}

		// Create new crawler
		webCrawler, err := crawler.NewWebCrawler(websiteURL, addCrawlMaxDepth, addCrawlConcurrency, addCrawlExcludePaths)
		if err != nil {
			return fmt.Errorf("error initializing web crawler: %w", err)
		}

		// Define crawling options
		webCrawler.SetUseSitemap(addCrawlUseSitemap)
		webCrawler.SetSingleURLMode(addCrawlSingleURL)

		// If specific URL list, define it
		if len(addCrawlURLsList) > 0 {
			webCrawler.SetURLsList(addCrawlURLsList)
		}

		// Show the crawling mode
		if len(addCrawlURLsList) > 0 {
			fmt.Printf("URLs list mode: crawling %d specific URLs\n", len(addCrawlURLsList))
		} else if addCrawlSingleURL {
			fmt.Println("Single URL mode: only the specified URL will be crawled (no links will be followed)")
		} else if addCrawlUseSitemap {
			fmt.Println("Sitemap mode enabled: will try to use sitemap.xml for comprehensive coverage")
		} else {
			fmt.Println("Standard crawling mode: will follow links to the specified depth")
		}

		fmt.Printf("Crawling website '%s' to add content to RAG '%s'...\n", websiteURL, ragName)

		// Start crawling
		documents, err := webCrawler.CrawlWebsite()
		if err != nil {
			return fmt.Errorf("error crawling website: %w", err)
		}

		if len(documents) == 0 {
			return fmt.Errorf("no content found when crawling %s", websiteURL)
		}

		fmt.Printf("Retrieved %d pages from website. Processing content...\n", len(documents))

		// Convert []domain.Document to []*domain.Document
		var docPointers []*domain.Document
		for i := range documents {
			docPointers = append(docPointers, &documents[i])
		}

		// Create temporary directory to store crawled content
		tempDir := service.CreateTempDirForDocuments(docPointers)
		if tempDir != "" {
			defer service.CleanupTempDir(tempDir)
		}

		// Set up loader options
		loaderOptions := service.DocumentLoaderOptions{
			ChunkSize:        addCrawlChunkSize,
			ChunkOverlap:     addCrawlChunkOverlap,
			ChunkingStrategy: addCrawlChunkingStrategy,
			EnableReranker:   !addCrawlDisableReranker,
			RerankerModel:    addCrawlRerankerModel,
			RerankerWeight:   addCrawlRerankerWeight,
		}

		// Pass the options to the service
		err = ragService.AddDocsWithOptions(ragName, tempDir, loaderOptions)
		if err != nil {
			return err
		}

		fmt.Printf("Website content from '%s' added to RAG '%s' successfully.\n", websiteURL, ragName)
		return nil
	},
}
⋮----
// Get Ollama client from root command
⋮----
// Create necessary services
⋮----
// Load existing RAG to get model name
⋮----
// Create new crawler
⋮----
// Define crawling options
⋮----
// If specific URL list, define it
⋮----
// Show the crawling mode
⋮----
// Start crawling
⋮----
// Convert []domain.Document to []*domain.Document
var docPointers []*domain.Document
⋮----
// Create temporary directory to store crawled content
⋮----
// Set up loader options
⋮----
// Pass the options to the service
⋮----
func init()
⋮----
// Add crawling specific flags
⋮----
// Add chunking flags
⋮----
// Add reranking options
````

## File: cmd/update.go
````go
package cmd
⋮----
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)
⋮----
"encoding/json"
"fmt"
"io"
"net/http"
"os"
"os/exec"
"path/filepath"
"runtime"
"strings"
"time"
⋮----
"github.com/spf13/cobra"
⋮----
var forceUpdate bool
⋮----
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
⋮----
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check and install RLAMA updates",
	Long: `Check if a new version of RLAMA is available and install it if so.
Example: rlama update

By default, the command asks for confirmation before installing the update.
Use the --force flag to update without confirmation.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Checking for RLAMA updates...")

		// Use the doUpdate function which properly handles Windows updates
		return doUpdate("", forceUpdate)
	},
}
⋮----
// Use the doUpdate function which properly handles Windows updates
⋮----
// checkForUpdates checks if updates are available by querying the GitHub API
func checkForUpdates() (*GitHubRelease, bool, error)
⋮----
// Query the GitHub API to get the latest release
⋮----
// Parse the JSON response
var release GitHubRelease
⋮----
// Check if the version is newer
⋮----
func init()
⋮----
// Fonction modifiée pour gérer le cas spécifique de Windows
func doUpdate(version string, force bool) error
⋮----
// Si aucune version n'est fournie, obtenir la dernière version
var latestVersion string
var err error
⋮----
// Vérifier si une mise à jour est nécessaire
⋮----
// Demander confirmation, sauf si --force est utilisé
⋮----
var response string
⋮----
// Obtenir le chemin de l'exécutable actuel
⋮----
// Créer un répertoire pour les fichiers de mise à jour si nécessaire
⋮----
// Télécharger la nouvelle version
⋮----
// Chemin pour le nouveau binaire
⋮----
// Télécharger le nouveau binaire
⋮----
// Nettoyer en cas d'erreur
⋮----
// Rendre le nouveau binaire exécutable
⋮----
// Sous Windows, nous devons utiliser une approche différente
⋮----
// Sur les autres plateformes, nous pouvons remplacer directement
⋮----
// Nouvelle fonction pour gérer la mise à jour sous Windows
func doWindowsUpdate(originalPath, newPath string) error
⋮----
// Create temporary batch script in a location we know exists
⋮----
os.MkdirAll(tempDir, 0755) // Ensure the directory exists
⋮----
// Simple batch script that waits for the process to end and then replaces the file
⋮----
// Write the batch script
⋮----
// Run the batch script in a new window
⋮----
// getLatestVersion récupère la dernière version disponible depuis GitHub
func getLatestVersion() (string, error)
⋮----
// Return the version without the 'v' prefix
⋮----
// downloadFile downloads a file from a URL to a local path
// with better error handling and retry attempts
func downloadFile(url string, filepath string) error
⋮----
// Create an HTTP client with timeout
⋮----
Timeout: 120 * time.Second, // 2 minute timeout
⋮----
// Create the file
⋮----
// Maximum retry attempts
⋮----
var lastErr error
⋮----
// Increase delay for next retry
⋮----
// Get the data
⋮----
// Add a user agent to avoid some download restrictions
⋮----
// Send the request
⋮----
// Check server response
⋮----
// Reset file position
⋮----
// Create a progress bar if the file is large enough
⋮----
if resp.ContentLength > 1024*1024 { // If larger than 1MB
⋮----
// Write the body to file
⋮----
// Success
⋮----
// If we get here, all retries failed
````

## File: .gitignore
````
documents/
rlama
dist/
.DS_Store
bin/
hook-executed.log
rlama.exe
*.exe
tests/
test-docs/
go.mod
go.sum
````

## File: cmd/run.go
````go
package cmd
⋮----
import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/service"
	"github.com/spf13/cobra"
)
⋮----
"bufio"
"fmt"
"os"
"strings"
⋮----
"github.com/dontizi/rlama/internal/domain"
"github.com/dontizi/rlama/internal/service"
"github.com/spf13/cobra"
⋮----
var (
	contextSize      int
	promptTemplate   string
	printChunks      bool
	streamOutput     bool
	apiProfileName   string
	maxTokens        int
	temperature      float64
	autoRetrievalAPI bool
	useGUI           bool
	showContext      bool
)
⋮----
var runCmd = &cobra.Command{
	Use:   "run [rag-name]",
	Short: "Run a RAG system",
	Long: `Run a previously created RAG system. 
Starts an interactive session to interact with the RAG system.
Example: rlama run rag1`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]

		// Get service provider
		provider := GetServiceProvider()
		ragService := provider.GetRagService()

		// Load the RAG system
		rag, err := ragService.LoadRag(ragName)
		if err != nil {
			return err
		}

		fmt.Printf("RAG '%s' loaded. Model: %s\n", rag.Name, rag.ModelName)
		if showContext {
			fmt.Printf("Debug info: RAG contains %d documents and %d total chunks\n",
				len(rag.Documents), len(rag.Chunks))
			fmt.Printf("Chunking strategy: %s, Size: %d, Overlap: %d\n",
				rag.ChunkingStrategy,
				rag.WatchOptions.ChunkSize,
				rag.WatchOptions.ChunkOverlap)
			if rag.RerankerEnabled {
				fmt.Printf("Reranking: Enabled (model: %s, weight: %.2f)\n",
					rag.RerankerModel, rag.RerankerWeight)
				defaultOpts := service.DefaultRerankerOptions()
				if contextSize <= 0 {
					fmt.Printf("Using default TopK: %d\n", defaultOpts.TopK)
				} else {
					fmt.Printf("Using custom TopK: %d\n", contextSize)
				}
			} else {
				fmt.Println("Reranking: Disabled")
			}
		}
		fmt.Println("Type your question (or 'exit' to quit):")

		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("> ")
			if !scanner.Scan() {
				break
			}

			question := scanner.Text()
			if question == "exit" {
				break
			}

			if strings.TrimSpace(question) == "" {
				continue
			}

			fmt.Println("\nSearching documents for relevant information...")

			checkWatchedResources(rag, ragService)

			// If debug mode is enabled, get the chunks manually first
			if showContext {
				// Use embedding service from provider
				embeddingService := provider.GetEmbeddingService()
				queryEmbedding, err := embeddingService.GenerateQueryEmbedding(question, rag.ModelName)
				if err != nil {
					fmt.Printf("Error generating embedding: %s\n", err)
				} else {
					results := rag.HybridStore.Search(queryEmbedding, contextSize)

					// Show detailed results
					fmt.Printf("\n--- Debug: Retrieved %d chunks ---\n", len(results))
					for i, result := range results {
						chunk := rag.GetChunkByID(result.ID)
						if chunk != nil {
							fmt.Printf("%d. [Score: %.4f] %s\n", i+1, result.Score, chunk.GetMetadataString())
							if i < 3 { // Show content for top 3 chunks only to avoid overload
								fmt.Printf("   Preview: %s\n", truncateString(chunk.Content, 100))
							}
						}
					}
					fmt.Println("--- End Debug ---")
				}
			}

			answer, err := ragService.Query(rag, question, contextSize)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				continue
			}

			fmt.Println("\n--- Answer ---")
			fmt.Println(answer)
			fmt.Println()
		}

		return nil
	},
}
⋮----
// Get service provider
⋮----
// Load the RAG system
⋮----
// If debug mode is enabled, get the chunks manually first
⋮----
// Use embedding service from provider
⋮----
// Show detailed results
⋮----
if i < 3 { // Show content for top 3 chunks only to avoid overload
⋮----
// Helper function to truncate string for preview
func truncateString(s string, maxLen int) string
⋮----
func init()
⋮----
// Add flags
⋮----
func checkWatchedResources(rag *domain.RagSystem, ragService service.RagService)
⋮----
// Check watched directory if enabled with on-use check
⋮----
// Check watched website if enabled with on-use check
````

## File: cmd/wizard.go
````go
package cmd
⋮----
import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/dontizi/rlama/internal/crawler"
	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/service"
	"github.com/spf13/cobra"
)
⋮----
"bufio"
"bytes"
"fmt"
"io"
"os"
"os/exec"
"strings"
⋮----
"github.com/AlecAivazis/survey/v2"
"github.com/dontizi/rlama/internal/crawler"
"github.com/dontizi/rlama/internal/domain"
"github.com/dontizi/rlama/internal/service"
"github.com/spf13/cobra"
⋮----
// Structure to parse the JSON output of Ollama list
type OllamaModel struct {
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	ModifiedAt string `json:"modified_at"`
	Digest     string `json:"digest"`
}
⋮----
var (
	// Variables for the local wizard
	localWizardModel        string
	localWizardName         string
	localWizardPath         string
	localWizardChunkSize    int
	localWizardChunkOverlap int
	localWizardExcludeDirs  []string
	localWizardExcludeExts  []string
	localWizardProcessExts  []string
)
⋮----
// Variables for the local wizard
⋮----
// Renamed to avoid conflict with snowflake_wizard.go
⋮----
var localWizardCmd = &cobra.Command{
	Use:   "wizard",
	Short: "Interactive wizard to create a local RAG",
	Long: `Start an interactive wizard that guides you through creating a RAG system.
This makes it easy to set up a new RAG without remembering all command options.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print("\n🧙 Welcome to the RLAMA Local RAG Wizard! 🧙\n\n")

		reader := bufio.NewReader(os.Stdin)

		// Étape 1: Nom du RAG
		fmt.Print("Enter a name for your RAG: ")
		ragName, _ := reader.ReadString('\n')
		ragName = strings.TrimSpace(ragName)
		if ragName == "" {
			return fmt.Errorf("RAG name cannot be empty")
		}

		// Declare modelName at the function level so it's available everywhere
		var modelName string

		// Step 2: Model selection
		fmt.Println("\nStep 2: Select a model")

		// Get the list of available Ollama models via the ollama list command
		fmt.Println("Retrieving available Ollama models...")

		// First try with ollama list without --json for better compatibility
		// and capture stderr for debugging
		var stdout, stderr bytes.Buffer
		ollamaCmd := exec.Command("ollama", "list")
		ollamaCmd.Stdout = &stdout
		ollamaCmd.Stderr = &stderr

		// Configuration for the command execution
		ollamaHost := os.Getenv("OLLAMA_HOST")
		if cmd.Flag("host").Changed {
			ollamaHost = cmd.Flag("host").Value.String()
		}

		if ollamaHost != "" {
			// Set the OLLAMA_HOST environment variable for the command
			ollamaCmd.Env = append(os.Environ(), fmt.Sprintf("OLLAMA_HOST=%s", ollamaHost))
		}

		// Execute the command
		err := ollamaCmd.Run()
		if err != nil {
			fmt.Println("❌ Failed to list Ollama models.")
			if stderr.Len() > 0 {
				fmt.Printf("Error details: %s\n", stderr.String())
			}
			fmt.Println("Make sure Ollama is installed and running.")
			fmt.Println("Continuing without model list. You'll need to enter a model name manually.")
		}

		// Parse the output of ollama list (text format)
		modelsOutput := stdout.String()
		var modelNames []string

		if modelsOutput != "" {
			// Typical format:
			// NAME             ID            SIZE    MODIFIED
			// llama3           xxx...xxx     4.7 GB  X days ago

			// Skip the first line (headers)
			lines := strings.Split(modelsOutput, "\n")
			for i, line := range lines {
				if i == 0 || strings.TrimSpace(line) == "" {
					continue
				}

				fields := strings.Fields(line)
				if len(fields) >= 1 {
					modelNames = append(modelNames, fields[0])
				}
			}

			// Display models in our format
			if len(modelNames) > 0 {
				fmt.Println("\nAvailable models:")
				for i, name := range modelNames {
					fmt.Printf("  %d. %s\n", i+1, name)
				}

				// Allow the user to choose a model
				fmt.Print("\nChoose a model (number) or enter model name: ")
				modelChoice, _ := reader.ReadString('\n')
				modelChoice = strings.TrimSpace(modelChoice)

				// Check if the user entered a number
				var modelNumber int
				modelName = "" // Initialize here too

				if _, err := fmt.Sscanf(modelChoice, "%d", &modelNumber); err == nil {
					// The user entered a number
					if modelNumber > 0 && modelNumber <= len(modelNames) {
						modelName = modelNames[modelNumber-1]
					} else {
						fmt.Println("Invalid selection. Please enter a valid model name manually.")
					}
				} else {
					// The user entered a name directly
					modelName = modelChoice
				}
			}
		}

		// If no model was selected, ask manually
		if modelName == "" {
			fmt.Print("Enter model name [llama3]: ")
			inputName, _ := reader.ReadString('\n')
			inputName = strings.TrimSpace(inputName)
			if inputName == "" {
				modelName = "llama3"
			} else {
				modelName = inputName
			}
		}

		// New Step 3: Choose between local documents or website
		fmt.Println("\nStep 3: Choose document source")
		fmt.Println("1. Local document folder")
		fmt.Println("2. Crawl a website")
		fmt.Print("\nSelect option (1/2): ")
		sourceChoice, _ := reader.ReadString('\n')
		sourceChoice = strings.TrimSpace(sourceChoice)

		var folderPath string
		var websiteURL string
		var maxDepth, concurrency int
		var excludePaths []string
		var useWebCrawler bool
		var useSitemap bool

		if sourceChoice == "2" {
			// Website crawler option
			useWebCrawler = true

			// Ask for the website URL
			fmt.Print("\nEnter website URL to crawl: ")
			websiteURL, _ = reader.ReadString('\n')
			websiteURL = strings.TrimSpace(websiteURL)
			if websiteURL == "" {
				return fmt.Errorf("website URL cannot be empty")
			}

			// Maximum crawl depth
			fmt.Print("Maximum crawl depth [2]: ")
			depthStr, _ := reader.ReadString('\n')
			depthStr = strings.TrimSpace(depthStr)
			maxDepth = 2 // default value
			if depthStr != "" {
				fmt.Sscanf(depthStr, "%d", &maxDepth)
			}

			// Concurrency
			fmt.Print("Number of concurrent crawlers [5]: ")
			concurrencyStr, _ := reader.ReadString('\n')
			concurrencyStr = strings.TrimSpace(concurrencyStr)
			concurrency = 5 // default value
			if concurrencyStr != "" {
				fmt.Sscanf(concurrencyStr, "%d", &concurrency)
			}

			// Paths to exclude
			fmt.Print("Paths to exclude (comma-separated): ")
			excludePathsStr, _ := reader.ReadString('\n')
			excludePathsStr = strings.TrimSpace(excludePathsStr)
			if excludePathsStr != "" {
				excludePaths = strings.Split(excludePathsStr, ",")
				for i := range excludePaths {
					excludePaths[i] = strings.TrimSpace(excludePaths[i])
				}
			}

			// Ask if the user wants to use the sitemap
			useSitemapPrompt := &survey.Confirm{
				Message: "Use sitemap.xml if available (recommended for better coverage)?",
				Default: true,
			}
			err = survey.AskOne(useSitemapPrompt, &useSitemap)
			if err != nil {
				return err
			}
		} else {
			// Option local folder (existing code)
			useWebCrawler = false

			fmt.Print("\nEnter path to document folder: ")
			folderPath, _ = reader.ReadString('\n')
			folderPath = strings.TrimSpace(folderPath)
			if folderPath == "" {
				return fmt.Errorf("folder path cannot be empty")
			}
		}

		// Step 4: Chunking options
		fmt.Println("\nStep 4: Chunking options")

		// Add chunking strategy selection
		fmt.Println("\nChunking strategies:")
		fmt.Println("  auto     - Automatically selects the best strategy for each document")
		fmt.Println("  fixed    - Splits text into fixed-size chunks")
		fmt.Println("  semantic - Respects natural boundaries like paragraphs")
		fmt.Println("  hybrid   - Adapts strategy based on document type")
		fmt.Println("  hierarchical - Creates two-level structure for long documents")

		fmt.Print("Chunking strategy [auto]: ")
		chunkingStrategyStr, _ := reader.ReadString('\n')
		chunkingStrategyStr = strings.TrimSpace(chunkingStrategyStr)
		chunkingStrategy := "auto"
		if chunkingStrategyStr != "" {
			chunkingStrategy = chunkingStrategyStr
		}

		fmt.Print("Chunk size [1000]: ")
		chunkSizeStr, _ := reader.ReadString('\n')
		chunkSizeStr = strings.TrimSpace(chunkSizeStr)
		chunkSize := 1000
		if chunkSizeStr != "" {
			fmt.Sscanf(chunkSizeStr, "%d", &chunkSize)
		}

		fmt.Print("Chunk overlap [200]: ")
		overlapStr, _ := reader.ReadString('\n')
		overlapStr = strings.TrimSpace(overlapStr)
		overlap := 200
		if overlapStr != "" {
			fmt.Sscanf(overlapStr, "%d", &overlap)
		}

		// Step 5: File filtering (optional)
		fmt.Println("\nStep 5: File filtering (optional)")

		fmt.Print("Exclude directories (comma-separated): ")
		excludeDirsStr, _ := reader.ReadString('\n')
		excludeDirsStr = strings.TrimSpace(excludeDirsStr)
		var excludeDirs []string
		if excludeDirsStr != "" {
			excludeDirs = strings.Split(excludeDirsStr, ",")
			for i := range excludeDirs {
				excludeDirs[i] = strings.TrimSpace(excludeDirs[i])
			}
		}

		fmt.Print("Exclude extensions (comma-separated): ")
		excludeExtsStr, _ := reader.ReadString('\n')
		excludeExtsStr = strings.TrimSpace(excludeExtsStr)
		var excludeExts []string
		if excludeExtsStr != "" {
			excludeExts = strings.Split(excludeExtsStr, ",")
			for i := range excludeExts {
				excludeExts[i] = strings.TrimSpace(excludeExts[i])
			}
		}

		fmt.Print("Process only these extensions (comma-separated): ")
		processExtsStr, _ := reader.ReadString('\n')
		processExtsStr = strings.TrimSpace(processExtsStr)
		var processExts []string
		if processExtsStr != "" {
			processExts = strings.Split(processExtsStr, ",")
			for i := range processExts {
				processExts[i] = strings.TrimSpace(processExts[i])
			}
		}

		// Step 6: Confirmation and creation
		fmt.Println("\nStep 6: Review and create")
		fmt.Println("RAG configuration:")
		fmt.Printf("- Name: %s\n", ragName)
		fmt.Printf("- Model: %s\n", modelName)

		if useWebCrawler {
			fmt.Printf("- Source: Website - %s\n", websiteURL)
			fmt.Printf("- Crawl depth: %d\n", maxDepth)
			fmt.Printf("- Concurrency: %d\n", concurrency)
			if len(excludePaths) > 0 {
				fmt.Printf("- Exclude paths: %s\n", strings.Join(excludePaths, ", "))
			}
		} else {
			fmt.Printf("- Source: Local folder - %s\n", folderPath)
			if len(excludeDirs) > 0 {
				fmt.Printf("- Exclude directories: %s\n", strings.Join(excludeDirs, ", "))
			}
			if len(excludeExts) > 0 {
				fmt.Printf("- Exclude extensions: %s\n", strings.Join(excludeExts, ", "))
			}
			if len(processExts) > 0 {
				fmt.Printf("- Process only: %s\n", strings.Join(processExts, ", "))
			}
		}

		fmt.Printf("- Chunk size: %d\n", chunkSize)
		fmt.Printf("- Chunk overlap: %d\n", overlap)
		fmt.Printf("- Chunking strategy: %s\n", chunkingStrategy)

		fmt.Print("\nCreate RAG with these settings? (y/n): ")
		confirm, _ := reader.ReadString('\n')
		confirm = strings.ToLower(strings.TrimSpace(confirm))

		if confirm != "y" && confirm != "yes" {
			fmt.Println("RAG creation cancelled.")
			return nil
		}

		// Create the RAG
		fmt.Println("\nCreating RAG...")

		// Get service provider
		provider := GetServiceProvider()
		ragService := provider.GetRagService()

		// Check that the model is available before continuing
		// This step is important to avoid errors later
		fmt.Printf("Checking if model '%s' is available...\n", modelName)
		ollamaClient := provider.GetOllamaClient()
		err = ollamaClient.CheckOllamaAndModel(modelName)
		if err != nil {
			return fmt.Errorf("model '%s' is not available: %w", modelName, err)
		}

		if useWebCrawler {
			// Use the crawler
			fmt.Printf("\nCrawling website '%s'...\n", websiteURL)

			// Create the crawler
			webCrawler, err := crawler.NewWebCrawler(websiteURL, maxDepth, concurrency, excludePaths)
			if err != nil {
				return fmt.Errorf("error initializing web crawler: %w", err)
			}

			// Set the sitemap option
			webCrawler.SetUseSitemap(useSitemap)

			// Start the crawling
			documents, err := webCrawler.CrawlWebsite()
			if err != nil {
				return fmt.Errorf("error crawling website: %w", err)
			}

			if len(documents) == 0 {
				return fmt.Errorf("no content found when crawling %s", websiteURL)
			}

			fmt.Printf("Retrieved %d pages from website. Processing content...\n", len(documents))

			// Convert documents to pointers before calling CreateTempDirForDocuments
			var docPointers []*domain.Document
			for i := range documents {
				docPointers = append(docPointers, &documents[i])
			}

			// Create a temporary directory for the documents
			tempDir := service.CreateTempDirForDocuments(docPointers)
			if tempDir != "" {
				defer service.CleanupTempDir(tempDir)
			}

			// Options for the document loader
			loaderOptions := service.DocumentLoaderOptions{
				ChunkSize:        chunkSize,
				ChunkOverlap:     overlap,
				ChunkingStrategy: chunkingStrategy,
				EnableReranker:   true,
			}

			// Create the RAG
			err = ragService.CreateRagWithOptions(modelName, ragName, tempDir, loaderOptions)
			if err != nil {
				return err
			}
		} else {
			// Use the local folder (existing code)
			loaderOptions := service.DocumentLoaderOptions{
				ExcludeDirs:      excludeDirs,
				ExcludeExts:      excludeExts,
				ProcessExts:      processExts,
				ChunkSize:        chunkSize,
				ChunkOverlap:     overlap,
				ChunkingStrategy: chunkingStrategy,
				EnableReranker:   true,
			}

			err = ragService.CreateRagWithOptions(modelName, ragName, folderPath, loaderOptions)
			if err != nil {
				return err
			}
		}

		fmt.Println("\n🎉 RAG created successfully! 🎉")
		fmt.Printf("\nYou can now use your RAG with: rlama run %s\n", ragName)

		return nil
	},
}
⋮----
// Étape 1: Nom du RAG
⋮----
// Declare modelName at the function level so it's available everywhere
var modelName string
⋮----
// Step 2: Model selection
⋮----
// Get the list of available Ollama models via the ollama list command
⋮----
// First try with ollama list without --json for better compatibility
// and capture stderr for debugging
var stdout, stderr bytes.Buffer
⋮----
// Configuration for the command execution
⋮----
// Set the OLLAMA_HOST environment variable for the command
⋮----
// Execute the command
⋮----
// Parse the output of ollama list (text format)
⋮----
var modelNames []string
⋮----
// Typical format:
// NAME             ID            SIZE    MODIFIED
// llama3           xxx...xxx     4.7 GB  X days ago
⋮----
// Skip the first line (headers)
⋮----
// Display models in our format
⋮----
// Allow the user to choose a model
⋮----
// Check if the user entered a number
var modelNumber int
modelName = "" // Initialize here too
⋮----
// The user entered a number
⋮----
// The user entered a name directly
⋮----
// If no model was selected, ask manually
⋮----
// New Step 3: Choose between local documents or website
⋮----
var folderPath string
var websiteURL string
var maxDepth, concurrency int
var excludePaths []string
var useWebCrawler bool
var useSitemap bool
⋮----
// Website crawler option
⋮----
// Ask for the website URL
⋮----
// Maximum crawl depth
⋮----
maxDepth = 2 // default value
⋮----
// Concurrency
⋮----
concurrency = 5 // default value
⋮----
// Paths to exclude
⋮----
// Ask if the user wants to use the sitemap
⋮----
// Option local folder (existing code)
⋮----
// Step 4: Chunking options
⋮----
// Add chunking strategy selection
⋮----
// Step 5: File filtering (optional)
⋮----
var excludeDirs []string
⋮----
var excludeExts []string
⋮----
var processExts []string
⋮----
// Step 6: Confirmation and creation
⋮----
// Create the RAG
⋮----
// Get service provider
⋮----
// Check that the model is available before continuing
// This step is important to avoid errors later
⋮----
// Use the crawler
⋮----
// Create the crawler
⋮----
// Set the sitemap option
⋮----
// Start the crawling
⋮----
// Convert documents to pointers before calling CreateTempDirForDocuments
var docPointers []*domain.Document
⋮----
// Create a temporary directory for the documents
⋮----
// Options for the document loader
⋮----
// Create the RAG
⋮----
// Use the local folder (existing code)
⋮----
func init()
⋮----
func ExecuteWizard(out, errOut io.Writer) error
⋮----
func NewWizardCommand() *cobra.Command
````

## File: cmd/crawl_rag.go
````go
package cmd
⋮----
import (
	"fmt"
	"strings"

	"github.com/dontizi/rlama/internal/crawler"
	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/service"
	"github.com/spf13/cobra"
)
⋮----
"fmt"
"strings"
⋮----
"github.com/dontizi/rlama/internal/crawler"
"github.com/dontizi/rlama/internal/domain"
"github.com/dontizi/rlama/internal/service"
"github.com/spf13/cobra"
⋮----
var (
	crawlMaxDepth          int
	crawlConcurrency       int
	crawlExcludePaths      []string
	crawlUseSitemap        bool
	crawlSingleURL         bool
	crawlURLsList          []string
	crawlChunkSize         int
	crawlChunkOverlap      int
	crawlChunkingStrategy  string
	crawlDisableReranker   bool
	crawlRerankerThreshold float64
	crawlRerankerWeight    float64
	crawlRerankerModel     string
)
⋮----
var crawlRagCmd = &cobra.Command{
	Use:   "crawl-rag [model] [rag-name] [website-url]",
	Short: "Create a new RAG system from a website",
	Long: `Create a new RAG system by crawling a website and indexing its content.
Example: rlama crawl-rag llama3 mysite-rag https://example.com

The crawler will try to use the sitemap.xml if available for comprehensive coverage.
It will also follow links on the pages up to the specified depth.

You can exclude certain paths and control other crawling parameters:
  rlama crawl-rag llama3 my-docs https://docs.example.com --max-depth=2
  rlama crawl-rag llama3 blog-rag https://blog.example.com --exclude-path=/archive,/tags
  rlama crawl-rag llama3 site-rag https://site.com --use-sitemap=false  # Disable sitemap search`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		modelName := args[0]
		ragName := args[1]
		websiteURL := args[2]

		// Get Ollama client with configured host and port
		ollamaClient := GetOllamaClient()
		if err := ollamaClient.CheckOllamaAndModel(modelName); err != nil {
			return err
		}

		// Create new crawler
		webCrawler, err := crawler.NewWebCrawler(websiteURL, crawlMaxDepth, crawlConcurrency, crawlExcludePaths)
		if err != nil {
			return fmt.Errorf("error initializing web crawler: %w", err)
		}

		// Define crawling options
		webCrawler.SetUseSitemap(crawlUseSitemap)
		webCrawler.SetSingleURLMode(crawlSingleURL)

		// If specific URL list, define it
		if len(crawlURLsList) > 0 {
			webCrawler.SetURLsList(crawlURLsList)
		}

		// Afficher le mode de crawling
		if len(crawlURLsList) > 0 {
			fmt.Printf("URLs list mode: crawling %d specific URLs\n", len(crawlURLsList))
		} else if crawlSingleURL {
			fmt.Println("Single URL mode: only the specified URL will be crawled (no links will be followed)")
		} else if crawlUseSitemap {
			fmt.Println("Sitemap mode enabled: will try to use sitemap.xml for comprehensive coverage")
		} else {
			fmt.Println("Standard crawling mode: will follow links to the specified depth")
		}

		// Display a message to indicate that the process has started
		fmt.Printf("Creating RAG '%s' with model '%s' by crawling website '%s'...\n",
			ragName, modelName, websiteURL)
		fmt.Printf("Max crawl depth: %d, Concurrency: %d\n", crawlMaxDepth, crawlConcurrency)

		// Start crawling
		documents, err := webCrawler.CrawlWebsite()
		if err != nil {
			return fmt.Errorf("error crawling website: %w", err)
		}

		if len(documents) == 0 {
			return fmt.Errorf("no content found when crawling %s", websiteURL)
		}

		fmt.Printf("Retrieved %d pages from website. Processing content...\n", len(documents))

		// Convertir []domain.Document en []*domain.Document
		var docPointers []*domain.Document
		for i := range documents {
			docPointers = append(docPointers, &documents[i])
		}

		// Create RAG service
		ragService := service.NewRagService(ollamaClient)

		// Set chunking options
		loaderOptions := service.DocumentLoaderOptions{
			ChunkSize:        crawlChunkSize,
			ChunkOverlap:     crawlChunkOverlap,
			ChunkingStrategy: crawlChunkingStrategy,
			EnableReranker:   !crawlDisableReranker,
			RerankerWeight:   crawlRerankerWeight,
			RerankerModel:    crawlRerankerModel,
		}

		// Create temporary directory to store crawled content
		tempDir := service.CreateTempDirForDocuments(docPointers)
		if tempDir != "" {
			// Comment this line to prevent deletion
			// defer service.CleanupTempDir(tempDir)

			// Add this to clearly display the path
			fmt.Printf("\n📁 The markdown files are located in: %s\n", tempDir)
		}

		// Create RAG system
		err = ragService.CreateRagWithOptions(modelName, ragName, tempDir, loaderOptions)
		if err != nil {
			if strings.Contains(err.Error(), "connection refused") {
				return fmt.Errorf("⚠️ Unable to connect to Ollama.\n" +
					"Make sure Ollama is installed and running.\n")
			}
			return err
		}

		// Set reranker threshold if specified
		if cmd.Flags().Changed("reranker-threshold") {
			// Load the RAG that was just created
			rag, err := ragService.LoadRag(ragName)
			if err != nil {
				return fmt.Errorf("error setting reranker threshold: %w", err)
			}

			// Set the threshold
			rag.RerankerThreshold = crawlRerankerThreshold

			// Save the updated RAG
			err = ragService.UpdateRag(rag)
			if err != nil {
				return fmt.Errorf("error updating reranker threshold: %w", err)
			}
		}

		fmt.Printf("RAG '%s' created successfully with content from %s.\n", ragName, websiteURL)

		return nil
	},
}
⋮----
// Get Ollama client with configured host and port
⋮----
// Create new crawler
⋮----
// Define crawling options
⋮----
// If specific URL list, define it
⋮----
// Afficher le mode de crawling
⋮----
// Display a message to indicate that the process has started
⋮----
// Start crawling
⋮----
// Convertir []domain.Document en []*domain.Document
var docPointers []*domain.Document
⋮----
// Create RAG service
⋮----
// Set chunking options
⋮----
// Create temporary directory to store crawled content
⋮----
// Comment this line to prevent deletion
// defer service.CleanupTempDir(tempDir)
⋮----
// Add this to clearly display the path
⋮----
// Create RAG system
⋮----
// Set reranker threshold if specified
⋮----
// Load the RAG that was just created
⋮----
// Set the threshold
⋮----
// Save the updated RAG
⋮----
func init()
⋮----
// Add local flags
⋮----
// Add reranker flags
````

## File: internal/domain/rag.go
````go
package domain
⋮----
import (
	"time"

	"github.com/dontizi/rlama/pkg/vector"
)
⋮----
"time"
⋮----
"github.com/dontizi/rlama/pkg/vector"
⋮----
// RagSystem represents a complete RAG system
type RagSystem struct {
	Name        string                      `json:"name"`
	ModelName   string                      `json:"model_name"`
	CreatedAt   time.Time                   `json:"created_at"`
	UpdatedAt   time.Time                   `json:"updated_at"`
	Description string                      `json:"description"`
	HybridStore *vector.EnhancedHybridStore // Use the hybrid store
	Documents   []*Document                 `json:"documents"`
	Chunks      []*DocumentChunk            `json:"chunks"`
	// Directory watching settings
	WatchedDir    string               `json:"watched_dir,omitempty"`
	WatchInterval int                  `json:"watch_interval,omitempty"` // In minutes, 0 means only check on use
	LastWatchedAt time.Time            `json:"last_watched_at,omitempty"`
	WatchEnabled  bool                 `json:"watch_enabled"`
	WatchOptions  DocumentWatchOptions `json:"watch_options,omitempty"`
	// Web watching settings
	WatchedURL       string          `json:"watched_url,omitempty"`
	WebWatchEnabled  bool            `json:"web_watch_enabled"`
	WebWatchInterval int             `json:"web_watch_interval,omitempty"` // In minutes
	LastWebWatchAt   time.Time       `json:"last_web_watched_at,omitempty"`
	WebWatchOptions  WebWatchOptions `json:"web_watch_options,omitempty"`
	APIProfileName   string          `json:"api_profile_name,omitempty"`  // Name of the API profile to use
	ChunkingStrategy string          `json:"chunking_strategy,omitempty"` // Type of chunking strategy used
	// Reranking settings
	RerankerEnabled   bool    `json:"reranker_enabled,omitempty"`   // Whether to use reranking
	RerankerModel     string  `json:"reranker_model,omitempty"`     // Model to use for reranking (if different from ModelName)
	RerankerWeight    float64 `json:"reranker_weight,omitempty"`    // Weight for reranker scores vs vector scores (0-1)
	RerankerThreshold float64 `json:"reranker_threshold,omitempty"` // Minimum score threshold for reranked results
	RerankerTopK      int     `json:"reranker_top_k,omitempty"`     // Default: return only top 5 results after reranking
	RerankerSilent    bool    `json:"reranker_silent,omitempty"`    // Whether to suppress warnings and output from the reranker
	// Embedding settings
	EmbeddingDimension int `json:"embedding_dimension,omitempty"` // Dimension of the embedding vectors
	
	// Vector Store Configuration
	VectorStoreType      string `json:"vector_store_type,omitempty"`      // e.g., "internal_hnsw", "qdrant"
	QdrantHost           string `json:"qdrant_host,omitempty"`
	QdrantPort           int    `json:"qdrant_port,omitempty"`            // e.g., 6333 for HTTP, 6334 for gRPC
	QdrantAPIKey         string `json:"qdrant_api_key,omitempty"`         // For Qdrant Cloud or secured instances
	QdrantCollectionName string `json:"qdrant_collection_name,omitempty"` // Typically derived from ragName
	QdrantGRPC           bool   `json:"qdrant_grpc,omitempty"`            // True to use gRPC, false for HTTP REST
}
⋮----
HybridStore *vector.EnhancedHybridStore // Use the hybrid store
⋮----
// Directory watching settings
⋮----
WatchInterval int                  `json:"watch_interval,omitempty"` // In minutes, 0 means only check on use
⋮----
// Web watching settings
⋮----
WebWatchInterval int             `json:"web_watch_interval,omitempty"` // In minutes
⋮----
APIProfileName   string          `json:"api_profile_name,omitempty"`  // Name of the API profile to use
ChunkingStrategy string          `json:"chunking_strategy,omitempty"` // Type of chunking strategy used
// Reranking settings
RerankerEnabled   bool    `json:"reranker_enabled,omitempty"`   // Whether to use reranking
RerankerModel     string  `json:"reranker_model,omitempty"`     // Model to use for reranking (if different from ModelName)
RerankerWeight    float64 `json:"reranker_weight,omitempty"`    // Weight for reranker scores vs vector scores (0-1)
RerankerThreshold float64 `json:"reranker_threshold,omitempty"` // Minimum score threshold for reranked results
RerankerTopK      int     `json:"reranker_top_k,omitempty"`     // Default: return only top 5 results after reranking
RerankerSilent    bool    `json:"reranker_silent,omitempty"`    // Whether to suppress warnings and output from the reranker
// Embedding settings
EmbeddingDimension int `json:"embedding_dimension,omitempty"` // Dimension of the embedding vectors
⋮----
// Vector Store Configuration
VectorStoreType      string `json:"vector_store_type,omitempty"`      // e.g., "internal_hnsw", "qdrant"
⋮----
QdrantPort           int    `json:"qdrant_port,omitempty"`            // e.g., 6333 for HTTP, 6334 for gRPC
QdrantAPIKey         string `json:"qdrant_api_key,omitempty"`         // For Qdrant Cloud or secured instances
QdrantCollectionName string `json:"qdrant_collection_name,omitempty"` // Typically derived from ragName
QdrantGRPC           bool   `json:"qdrant_grpc,omitempty"`            // True to use gRPC, false for HTTP REST
⋮----
// DocumentWatchOptions stores settings for directory watching
type DocumentWatchOptions struct {
	ExcludeDirs      []string `json:"exclude_dirs,omitempty"`
	ExcludeExts      []string `json:"exclude_exts,omitempty"`
	ProcessExts      []string `json:"process_exts,omitempty"`
	ChunkSize        int      `json:"chunk_size,omitempty"`
	ChunkOverlap     int      `json:"chunk_overlap,omitempty"`
	ChunkingStrategy string   `json:"chunking_strategy,omitempty"`
}
⋮----
// WebWatchOptions stores settings for web watching
type WebWatchOptions struct {
	MaxDepth         int      `json:"max_depth,omitempty"`
	Concurrency      int      `json:"concurrency,omitempty"`
	ExcludePaths     []string `json:"exclude_paths,omitempty"`
	ChunkSize        int      `json:"chunk_size,omitempty"`
	ChunkOverlap     int      `json:"chunk_overlap,omitempty"`
	ChunkingStrategy string   `json:"chunking_strategy,omitempty"`
}
⋮----
// NewRagSystem creates a new instance of RagSystem
func NewRagSystem(name, modelName string) *RagSystem
⋮----
return NewRagSystemWithDimensions(name, modelName, 1536) // Default to 1536 dimensions
⋮----
// NewRagSystemWithDimensions creates a new instance of RagSystem with specified embedding dimensions
func NewRagSystemWithDimensions(name, modelName string, dimensions int) *RagSystem
⋮----
// Handle error appropriately
⋮----
RerankerEnabled:    true,                      // Enable reranking by default
RerankerModel:      "BAAI/bge-reranker-v2-m3", // Use BGE reranker by default
RerankerWeight:     0.7,                       // Default: 70% reranker score, 30% vector similarity
RerankerTopK:       5,                         // Default: return only top 5 results after reranking
EmbeddingDimension: dimensions,                // Store the embedding dimension
VectorStoreType:    "internal",                // Default to internal vector store
⋮----
// NewRagSystemWithVectorStore creates a new instance of RagSystem with vector store configuration
func NewRagSystemWithVectorStore(name, modelName string, dimensions int, vectorStoreType, qdrantHost string, qdrantPort int, qdrantAPIKey, qdrantCollection string, qdrantGRPC bool) *RagSystem
⋮----
// Create hybrid store config
⋮----
RerankerEnabled:      true,                      // Enable reranking by default
RerankerModel:        "BAAI/bge-reranker-v2-m3", // Use BGE reranker by default
RerankerWeight:       0.7,                       // Default: 70% reranker score, 30% vector similarity
RerankerTopK:         5,                         // Default: return only top 5 results after reranking
EmbeddingDimension:   dimensions,                // Store the embedding dimension
⋮----
// AddDocument adds a document to the RAG system
func (r *RagSystem) AddDocument(doc *Document)
⋮----
// Don't use doc.Metadata if it doesn't exist
⋮----
// GetDocumentByID retrieves a document by its ID
func (r *RagSystem) GetDocumentByID(id string) *Document
⋮----
// RemoveDocument removes a document from the RAG system by its ID
func (r *RagSystem) RemoveDocument(id string) bool
⋮----
// Find the document index
var index = -1
⋮----
// Remove from the Documents slice
⋮----
// Remove from the HybridStore
⋮----
// AddChunk adds a chunk to the RAG system
func (r *RagSystem) AddChunk(chunk *DocumentChunk)
⋮----
// GetChunkByID retrieves a chunk by its ID
func (r *RagSystem) GetChunkByID(id string) *DocumentChunk
⋮----
// Search performs a hybrid search using the hybrid store
func (r *RagSystem) Search(queryVector []float32, queryText string, limit int) ([]vector.HybridSearchResult, error)
````

## File: internal/service/document_loader.go
````go
package service
⋮----
import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/dontizi/rlama/internal/domain"
)
⋮----
"bytes"
"encoding/csv"
"fmt"
"io/ioutil"
"os"
"os/exec"
"path/filepath"
"runtime"
"sort"
"strings"
"sync"
⋮----
"github.com/dontizi/rlama/internal/domain"
⋮----
// DocumentLoaderOptions defines filtering options for document loading
type DocumentLoaderOptions struct {
	ExcludeDirs      []string
	ExcludeExts      []string
	ProcessExts      []string
	ChunkSize        int
	ChunkOverlap     int
	ChunkingStrategy string  // Chunking strategy: "fixed", "semantic", "hybrid", "hierarchical"
	APIProfileName   string  // Name of the API profile to use
	EmbeddingModel   string  // Model to use for embeddings
	EnableReranker   bool    // Whether to enable reranking - now true by default
	RerankerModel    string  // Model to use for reranking
	RerankerWeight   float64 // Weight for reranker scores (0-1)
	UseONNXReranker  bool    // Whether to use ONNX reranker instead of Python
	ONNXModelDir     string  // Directory containing ONNX model files
	// Qdrant configuration
	VectorStore          string // "internal" or "qdrant"
	QdrantHost           string
	QdrantPort           int
	QdrantAPIKey         string
	QdrantCollectionName string
	QdrantGRPC           bool
}
⋮----
ChunkingStrategy string  // Chunking strategy: "fixed", "semantic", "hybrid", "hierarchical"
APIProfileName   string  // Name of the API profile to use
EmbeddingModel   string  // Model to use for embeddings
EnableReranker   bool    // Whether to enable reranking - now true by default
RerankerModel    string  // Model to use for reranking
RerankerWeight   float64 // Weight for reranker scores (0-1)
UseONNXReranker  bool    // Whether to use ONNX reranker instead of Python
ONNXModelDir     string  // Directory containing ONNX model files
// Qdrant configuration
VectorStore          string // "internal" or "qdrant"
⋮----
// NewDocumentLoaderOptions creates default document loader options with reranking enabled
func NewDocumentLoaderOptions() DocumentLoaderOptions
⋮----
EnableReranker:   true, // Enable reranking by default
RerankerWeight:   0.7,  // 70% reranker, 30% vector
RerankerModel:    "",   // Will use the same model as RAG by default
UseONNXReranker:  false, // Default to Python implementation
⋮----
VectorStore:      "internal", // Default to internal vector store
⋮----
// DocumentLoader is responsible for loading documents from the file system
type DocumentLoader struct {
	supportedExtensions map[string]bool
	extractorPath       string // Path to the external extractor
}
⋮----
extractorPath       string // Path to the external extractor
⋮----
// NewDocumentLoader creates a new instance of DocumentLoader
func NewDocumentLoader() *DocumentLoader
⋮----
// Plain text
⋮----
// Source code
⋮----
// Documents
⋮----
// We'll use pdftotext if available
⋮----
// findExternalExtractor looks for external extraction tools
func findExternalExtractor() string
⋮----
// Define platform-specific extractors
var extractors []string
⋮----
// Windows-specific extractors
⋮----
"xpdf-pdftotext.exe", // Xpdf tools for Windows
"pdftotext.exe",      // Poppler Windows
"docx2txt.exe",       // For docx files
⋮----
// Unix/Mac extractors
⋮----
"pdftotext", // Poppler-utils
"textutil",  // macOS
"catdoc",    // For .doc
"unrtf",     // For .rtf
⋮----
// LoadDocumentsFromFolderWithOptions loads documents with filtering options
func (dl *DocumentLoader) LoadDocumentsFromFolderWithOptions(folderPath string, options DocumentLoaderOptions) ([]*domain.Document, error)
⋮----
var documents []*domain.Document
var supportedFiles []string
var unsupportedFiles []string
var excludedFiles []string
⋮----
// Normalize extensions for easier comparison
⋮----
// Ensure folderPath is absolute
⋮----
// Check if the folder exists
⋮----
// Try to create the folder
⋮----
// Get information about the newly created folder
⋮----
// Preliminary file check - recursively walk the directory
⋮----
return nil // Skip this file but continue walking
⋮----
// Check if this directory should be excluded
⋮----
// Ignore hidden files (starting with .)
⋮----
// Check if the extension is in the exclude list
⋮----
// If we're only processing specific extensions
⋮----
// Display info about found files
⋮----
// Try to install dependencies if possible
⋮----
// Process supported files
⋮----
// Text extraction using multiple methods
⋮----
// Try reading as a text file
⋮----
// Check that the content is not empty
⋮----
// For PDFs, try one last method
⋮----
// Create a document with relative path for better identification
⋮----
relPath = path // Fallback to full path if relative path can't be determined
⋮----
// Use relPath for document identification, but keep the full path for file access
⋮----
doc.Name = relPath // Use relative path as the document name for better browsing
// Don't change doc.ID or doc.Path which need the absolute path
⋮----
// extractText extracts text from a file using the appropriate method based on type
func (dl *DocumentLoader) extractText(path string, ext string) (string, error)
⋮----
// Treat as a text file
⋮----
// extractFromPDF extracts text from a PDF using different methods
func (dl *DocumentLoader) extractFromPDF(path string) (string, error)
⋮----
// Method 1: Use pdftotext if available
⋮----
// Method 2: Try with other tools (pdfinfo, pdftk)
⋮----
var out []byte
⋮----
// Method 3: Last attempt, open as binary file and extract strings
⋮----
// extractFromDocument extracts text from a Word document or similar
func (dl *DocumentLoader) extractFromDocument(path string, ext string) (string, error)
⋮----
// Method 1: Use textutil on macOS
⋮----
// Method 2: Use catdoc for .doc
⋮----
// Method 3: Use unrtf for .rtf
⋮----
// Method 4: Extract strings
⋮----
// extractFromPresentation extracts text from a PowerPoint presentation
func (dl *DocumentLoader) extractFromPresentation(path string, ext string) (string, error)
⋮----
// External tools for PowerPoint are limited
⋮----
// extractFromSpreadsheet extracts text from an Excel spreadsheet
func (dl *DocumentLoader) extractFromSpreadsheet(path string, ext string) (string, error)
⋮----
// Try to use xlsx2csv for .xlsx
⋮----
// Try to use xls2csv for .xls
⋮----
// Extract strings
⋮----
// extractStringsFromBinary extracts strings from a binary file
func (dl *DocumentLoader) extractStringsFromBinary(path string) (string, error)
⋮----
// Use the 'strings' tool if available (Unix/Linux/macOS)
⋮----
// Basic implementation of 'strings' in Go
⋮----
var result strings.Builder
var currentWord strings.Builder
⋮----
if currentWord.Len() >= 4 { // Word of at least 4 characters
⋮----
// extractWithOCR attempts to extract text using OCR
func (dl *DocumentLoader) extractWithOCR(path string) (string, error)
⋮----
// Check if tesseract is available
⋮----
// Create a temporary output file
⋮----
// Determine optimal number of workers
⋮----
// For PDFs, first convert to images if possible
⋮----
// Check if pdftoppm is available
⋮----
// Convert PDF to images with parallel processing
⋮----
// First, determine the number of pages in the PDF
⋮----
// Parse page count from pdfinfo output
⋮----
// Process PDF in parallel batches
⋮----
var wg sync.WaitGroup
⋮----
semaphore <- struct{}{}        // Acquire
defer func() { <-semaphore }() // Release
⋮----
// For smaller PDFs or when pdfinfo isn't available, use the original approach
⋮----
// OCR on each image in parallel
⋮----
// Direct OCR on the file (for images)
⋮----
// Read the extracted text
⋮----
// parallelOCR processes multiple image files with tesseract in parallel
func (dl *DocumentLoader) parallelOCR(imgFiles []string, tesseractPath, outBaseDir string, numWorkers int) (string, error)
⋮----
var mutex sync.Mutex
⋮----
var wg sync.WaitGroup
⋮----
var processingErr error
⋮----
// Process each image file in parallel
⋮----
semaphore <- struct{}{}        // Acquire semaphore
defer func() { <-semaphore }() // Release semaphore
⋮----
// Read the extracted text
⋮----
// Store the result
⋮----
// Combine all text in the correct order (by filename)
⋮----
var allText strings.Builder
⋮----
// tryInstallDependencies attempts to install dependencies if necessary
func (dl *DocumentLoader) tryInstallDependencies()
⋮----
// Check if pip is available (for Python tools)
⋮----
// Try to install useful packages
⋮----
installCmd.Run() // Ignore errors
⋮----
// processContent processes the content of a document and returns chunks
func (dl *DocumentLoader) processContent(path string, content string, options DocumentLoaderOptions) []*domain.DocumentChunk
⋮----
var chunks []*domain.DocumentChunk
⋮----
// extractCSVContent extracts content from a CSV file
func (dl *DocumentLoader) extractCSVContent(path string) (string, error)
⋮----
var content strings.Builder
// Add headers as first line if present
⋮----
// Add remaining rows
⋮----
// extractExcelContent extracts content from an Excel file
func (dl *DocumentLoader) extractExcelContent(path string) (string, error)
⋮----
// First try using xlsx2csv if available
⋮----
var output bytes.Buffer
⋮----
// Fallback to Python xlsx2csv package if installed
⋮----
var output bytes.Buffer
⋮----
// extractContent extracts content from a file based on its type
func (dl *DocumentLoader) extractContent(path string) (string, error)
⋮----
// CreateRagWithOptions creates a new RAG with the specified options
func (dl *DocumentLoader) CreateRagWithOptions(options DocumentLoaderOptions) (*domain.RagSystem, error)
⋮----
// CreateTempDirForDocuments creates a temporary directory and saves crawled documents as files
func CreateTempDirForDocuments(documents []*domain.Document) string
⋮----
// Create a temporary directory
⋮----
// Save each document as a file in the temporary directory
⋮----
// Default to index-based filename
⋮----
// Try to use Path if available (more likely to exist than URL)
⋮----
// Create a safe filename from the Path
⋮----
// Trim leading/trailing dashes
⋮----
// Full path to the file
⋮----
// Write the document content to the file
⋮----
// CleanupTempDir removes a temporary directory and all its contents
func CleanupTempDir(path string)
````

## File: cmd/rag.go
````go
package cmd
⋮----
import (
	"fmt"
	"strings"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/service"
	"github.com/spf13/cobra"
)
⋮----
"fmt"
"strings"
⋮----
"github.com/dontizi/rlama/internal/client"
"github.com/dontizi/rlama/internal/service"
"github.com/spf13/cobra"
⋮----
var (
	excludeDirs          []string
	excludeExts          []string
	processExts          []string
	chunkSize            int
	chunkOverlap         int
	chunkingStrategy     string
	profileName          string
	embeddingModel       string
	ragDisableReranker   bool
	ragRerankerModel     string
	ragRerankerWeight    float64
	ragRerankerThreshold float64
	ragUseONNXReranker   bool
	ragONNXModelDir      string
	// Qdrant configuration flags
	vectorStore          string
	qdrantHost           string
	qdrantPort           int
	qdrantAPIKey         string
	qdrantCollection     string
	qdrantUseGRPC        bool
	testService          interface{} // Pour les tests
)
⋮----
// Qdrant configuration flags
⋮----
testService          interface{} // Pour les tests
⋮----
var ragCmd = &cobra.Command{
	Use:   "rag [model] [rag-name] [folder-path]",
	Short: "Create a new RAG system",
	Long: `Create a new RAG system by indexing all documents in the specified folder.
Example: rlama rag llama3.2 rag1 ./documents

The folder will be created if it doesn't exist yet.
Supported formats include: .txt, .md, .html, .json, .csv, and various source code files.

You can exclude directories or file types:
  rlama rag llama3 myproject ./code --excludedir=node_modules,dist,.git
  rlama rag llama3 mydocs ./docs --excludeext=.log,.tmp
  rlama rag llama3 specific ./mixed --processext=.md,.py,.js

Hugging Face Models:
  You can use Hugging Face GGUF models with the format:
  rlama rag hf.co/username/repository my-rag ./docs
  
  Specify quantization with:
  rlama rag hf.co/username/repository:Q4_K_M my-rag ./docs
  
OpenAI Models:
  You can use OpenAI models by setting the OPENAI_API_KEY environment variable:
  export OPENAI_API_KEY="your-api-key"
  
  Then use any OpenAI model:
  rlama rag gpt-4-turbo my-openai-rag ./docs
  rlama rag gpt-3.5-turbo my-openai-rag ./docs`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		modelName := args[0]
		ragName := args[1]
		folderPath := args[2]

		// Setup client
		llmClient, ollamaClient, err := setupLLMClient(modelName)
		if err != nil {
			return err
		}

		// Setup configuration
		loaderOptions, err := setupLoaderOptions(ragName)
		if err != nil {
			return err
		}

		// Create RAG
		ragService, err := createRagSystem(modelName, ragName, folderPath, llmClient, ollamaClient, loaderOptions)
		if err != nil {
			return handleRagCreationError(err)
		}

		// Post-creation configuration
		if err := configurePostCreation(cmd, ragService, ragName); err != nil {
			return err
		}

		fmt.Printf("RAG '%s' created successfully.\n", ragName)
		return nil
	},
}
⋮----
// Setup client
⋮----
// Setup configuration
⋮----
// Create RAG
⋮----
// Post-creation configuration
⋮----
func init()
⋮----
// Add exclusion and processing flags
⋮----
// Add flags for chunking options
⋮----
// Add reranking options - now with a flag to disable it instead
⋮----
// Add profile option
⋮----
// Add embedding model option
⋮----
// Add Qdrant configuration flags
⋮----
// Add logic to use the test service if available
⋮----
// Use the test service
⋮----
// NewRagCommand returns the rag command
func NewRagCommand() *cobra.Command
⋮----
// InjectTestService injects a test service
func InjectTestService(service interface
⋮----
// setupLLMClient creates and configures the appropriate LLM client using ServiceProvider
func setupLLMClient(modelName string) (client.LLMClient, *client.OllamaClient, error)
⋮----
// Update provider config with profile if specified
⋮----
// Get clients from provider
⋮----
// Debug: Check what type of client we got
⋮----
// Check the client and model
⋮----
// Display which client/provider is being used
⋮----
// Handle Hugging Face models (Ollama-specific)
⋮----
// Extract quantization if specified
⋮----
// Pull the model from Hugging Face
⋮----
// setupLoaderOptions creates and validates document loader options using ServiceProvider
func setupLoaderOptions(ragName string) (service.DocumentLoaderOptions, error)
⋮----
// Override ONNX configuration if specified via flags
⋮----
// Update the provider configuration
⋮----
// Create base options from configuration
⋮----
// Override with command-specific flags
⋮----
// Override chunk settings if provided via flags
if chunkSize != 1000 { // Default value check
⋮----
if chunkOverlap != 200 { // Default value check
⋮----
if chunkingStrategy != "hybrid" { // Default value check
⋮----
// Override profile and embedding settings if provided
⋮----
// Override reranker settings if provided
⋮----
if ragRerankerWeight != 0.7 { // Default value check
⋮----
if ragONNXModelDir != "./models/bge-reranker-large-onnx" { // Default value check
⋮----
// Override vector store settings if provided
if vectorStore != "internal" { // Default value check
⋮----
if qdrantHost != "localhost" { // Default value check
⋮----
if qdrantPort != 6333 { // Default value check
⋮----
// Set default collection name if not provided
⋮----
// Validate Qdrant configuration if using Qdrant
⋮----
// createRagSystem creates the RAG system with the specified configuration using ServiceProvider
func createRagSystem(modelName, ragName, folderPath string, llmClient client.LLMClient, ollamaClient *client.OllamaClient, loaderOptions service.DocumentLoaderOptions) (service.RagService, error)
⋮----
// Display a message to indicate that the process has started
⋮----
// Use the service provider to create a RAG service for the specific model
⋮----
// handleRagCreationError provides improved error messages for RAG creation failures
func handleRagCreationError(err error) error
⋮----
// Improve error messages related to Ollama
⋮----
// configurePostCreation handles post-creation configuration like reranker threshold
func configurePostCreation(cmd *cobra.Command, ragService service.RagService, ragName string) error
⋮----
// Set reranker threshold if specified
⋮----
// Load the RAG that was just created
⋮----
// Set the threshold
⋮----
// Save the updated RAG
````

## File: cmd/root.go
````go
package cmd
⋮----
import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/service"
)
⋮----
"fmt"
"time"
⋮----
"github.com/spf13/cobra"
⋮----
"github.com/dontizi/rlama/internal/client"
"github.com/dontizi/rlama/internal/service"
⋮----
const (
	Version = "0.1.36" // Current version of RLAMA
)
⋮----
Version = "0.1.36" // Current version of RLAMA
⋮----
var rootCmd = &cobra.Command{
	Use:   "rlama",
	Short: "RLAMA is a CLI tool for creating and using RAG systems with Ollama",
	Long: `RLAMA (Retrieval-Augmented Language Model Adapter) is a command-line tool 
that simplifies the creation and use of RAG (Retrieval-Augmented Generation) systems 
based on Ollama models.

Main commands:
  rag [model] [rag-name] [folder-path]    Create a new RAG system
  run [rag-name]                          Run an existing RAG system
  list                                    List all available RAG systems
  delete [rag-name]                       Delete a RAG system
  update                                  Check and install RLAMA updates

Environment variables:
  OLLAMA_HOST                            Specifies default Ollama host:port (overridden by --host and --port flags)`,
}
⋮----
// Variables to store command flags
var (
	versionFlag bool
	ollamaHost  string
	ollamaPort  string
	verbose     bool
	dataDir     string
)
⋮----
// Global service provider instance
var serviceProvider *service.ServiceProvider
⋮----
// Execute executes the root command
func Execute() error
⋮----
// GetServiceProvider returns the global service provider, creating it if necessary
func GetServiceProvider() *service.ServiceProvider
⋮----
// Override with command-line flags if provided
⋮----
var err error
⋮----
// Fallback to default config in case of error
⋮----
// GetOllamaClient returns an Ollama client configured with host and port from command flags
// Deprecated: Use GetServiceProvider().GetOllamaClient() instead
func GetOllamaClient() *client.OllamaClient
⋮----
func init()
⋮----
// Add --version flag
⋮----
// Add Ollama configuration flags
⋮----
// New flag for data directory
⋮----
// Override the Run function to handle the --version flag
⋮----
// If no arguments are provided and --version is not used, display help
⋮----
// Start the watcher daemon
⋮----
// Add this function to start the watcher daemon
func startFileWatcherDaemon()
⋮----
// Wait a bit for application initialization
⋮----
// Get services from the service provider
⋮----
// Start the daemon with a 1-minute check interval for its internal operations
// Actual RAG check intervals are controlled by each RAG's configuration
````

## File: README.md
````markdown
<!-- Social Links Navigation Bar -->
<div align="center">
  <a href="https://x.com/LeDonTizi" target="_blank">
    <img src="https://img.shields.io/badge/Twitter-1DA1F2?style=for-the-badge&logo=twitter&logoColor=white" alt="Twitter">
  </a>
  <a href="https://discord.gg/tP5JB9DR" target="_blank">
    <img src="https://img.shields.io/badge/Discord-5865F2?style=for-the-badge&logo=discord&logoColor=white" alt="Discord">
  </a>
  <a href="https://www.youtube.com/@Dontizi" target="_blank">
    <img src="https://img.shields.io/badge/YouTube-FF0000?style=for-the-badge&logo=youtube&logoColor=white" alt="YouTube">
  </a>
</div>

<br>

# RLAMA - User Guide
RLAMA is a powerful AI-driven question-answering tool for your documents that works with multiple LLM providers. It seamlessly integrates with Ollama, OpenAI, and any OpenAI-compatible endpoints (like LM Studio, VLLM, Text Generation Inference, etc.). RLAMA enables you to create, manage, and interact with Retrieval-Augmented Generation (RAG) systems tailored to your documentation needs.


[![RLAMA Demonstration](https://img.youtube.com/vi/EIsQnBqeQxQ/0.jpg)](https://www.youtube.com/watch?v=EIsQnBqeQxQ)

## Table of Contents
- [Vision & Roadmap](#vision--roadmap)
- [Installation](#installation)
- [Available Commands](#available-commands)
  - [rag - Create a RAG system](#rag---create-a-rag-system)
  - [crawl-rag - Create a RAG system from a website](#crawl-rag---create-a-rag-system-from-a-website)
  - [wizard - Create a RAG system with interactive setup](#wizard---create-a-rag-system-with-interactive-setup)
  - [watch - Set up directory watching for a RAG system](#watch---set-up-directory-watching-for-a-rag-system)
  - [watch-off - Disable directory watching for a RAG system](#watch-off---disable-directory-watching-for-a-rag-system)
  - [check-watched - Check a RAG's watched directory for new files](#check-watched---check-a-rags-watched-directory-for-new-files)
  - [web-watch - Set up website monitoring for a RAG system](#web-watch---set-up-website-monitoring-for-a-rag-system)
  - [web-watch-off - Disable website monitoring for a RAG system](#web-watch-off---disable-website-monitoring-for-a-rag-system)
  - [check-web-watched - Check a RAG's monitored website for updates](#check-web-watched---check-a-rags-monitored-website-for-updates)
  - [run - Use a RAG system](#run---use-a-rag-system)
  - [api - Start API server](#api---start-api-server)
  - [list - List RAG systems](#list---list-rag-systems)
  - [delete - Delete a RAG system](#delete---delete-a-rag-system)
  - [list-docs - List documents in a RAG](#list-docs---list-documents-in-a-rag)
  - [list-chunks - Inspect document chunks](#list-chunks---inspect-document-chunks)
  - [view-chunk - View chunk details](#view-chunk---view-chunk-details)
  - [add-docs - Add documents to RAG](#add-docs---add-documents-to-rag)
  - [crawl-add-docs - Add website content to RAG](#crawl-add-docs---add-website-content-to-rag)
  - [migrate-to-qdrant - Migrate RAG to Qdrant](#migration-between-vector-stores)
  - [migrate-to-internal - Migrate RAG to internal storage](#migration-between-vector-stores)
  - [migrate-batch - Batch migrate multiple RAGs](#migration-between-vector-stores)
  - [update-model - Change LLM model](#update-model---change-llm-model)
  - [profile - Manage API profiles](#profile---manage-api-profiles)
  - [update - Update RLAMA](#update---update-rlama)
  - [version - Display version](#version---display-version)
  - [hf-browse - Browse GGUF models on Hugging Face](#hf-browse---browse-gguf-models-on-hugging-face)
  - [run-hf - Run a Hugging Face GGUF model](#run-hf---run-a-hugging-face-gguf-model)
- [Uninstallation](#uninstallation)
- [Supported Document Formats](#supported-document-formats)
- [Troubleshooting](#troubleshooting)
- [Model Support & LLM Providers](#model-support--llm-providers)
- [Managing API Profiles](#managing-api-profiles)

## Vision & Roadmap
RLAMA aims to become the definitive tool for creating local RAG systems that work seamlessly for everyone—from individual developers to large enterprises. Here's our strategic roadmap:

### Completed Features ✅
- ✅ **Basic RAG System Creation**: CLI tool for creating and managing RAG systems
- ✅ **Document Processing**: Support for multiple document formats (.txt, .md, .pdf, etc.)
- ✅ **Document Chunking**: Advanced semantic chunking with multiple strategies (fixed, semantic, hierarchical, hybrid)
- ✅ **Vector Storage**: Local storage of document embeddings + Qdrant vector database integration
- ✅ **Production Vector Database**: Full Qdrant integration with gRPC/HTTP support, Qdrant Cloud compatibility
- ✅ **Seamless Migration Tools**: Complete migration system between internal and Qdrant storage with data integrity verification
- ✅ **Batch Operations**: Bulk migration of multiple RAGs with progress tracking and error recovery
- ✅ **Context Retrieval**: Basic semantic search with configurable context size
- ✅ **Ollama Integration**: Seamless connection to Ollama models
- ✅ **OpenAI Integration**: Full OpenAI API compatibility with profile management
- ✅ **Cross-Platform Support**: Works on Linux, macOS, and Windows
- ✅ **Easy Installation**: One-line installation script
- ✅ **API Server**: HTTP endpoints for integrating RAG capabilities in other applications
- ✅ **Web Crawling**: Create RAGs directly from websites
- ✅ **Guided RAG Setup Wizard**: Interactive interface for easy RAG creation
- ✅ **Hugging Face Integration**: Access to 45,000+ GGUF models from Hugging Face Hub
- ✅ **Advanced Reranking**: BGE reranker integration for improved search accuracy

### Small LLM Optimization (Q2 2025)
- [ ] **Prompt Compression**: Smart context summarization for limited context windows
- ✅ **Adaptive Chunking**: Dynamic content segmentation based on semantic boundaries and document structure
- ✅ **Minimal Context Retrieval**: Intelligent filtering to eliminate redundant content
- [ ] **Parameter Optimization**: Fine-tuned settings for different model sizes

### Advanced Search & Filtering (Q2 2025)
- [ ] **Enhanced Metadata Filtering**: Advanced search with document type, date, author, and custom metadata filters
- [ ] **Structured Query Language**: SQL-like queries for complex document retrieval
- [ ] **Faceted Search**: Multi-dimensional filtering with result counts
- [ ] **Similarity Thresholds**: Configurable relevance scoring and filtering

### Performance & Reliability (Q2-Q3 2025)
- [ ] **Connection Pooling**: Optimized Qdrant connections for high-throughput scenarios
- [ ] **Async Operations**: Non-blocking operations for large document imports
- [ ] **Caching Layer**: Smart caching for frequently accessed data and embeddings
- [ ] **Health Monitoring**: System health checks and performance metrics
- [ ] **Auto-retry Logic**: Exponential backoff for network failures

### Enhanced CLI & Developer Experience (Q2-Q3 2025)
- [ ] **RAG Status & Diagnostics**: `rag status`, `rag health-check`, `rag benchmark` commands
- [ ] **Performance Analytics**: Query performance metrics and optimization suggestions
- [ ] **Advanced Debugging**: Detailed logging, search result explanations, and troubleshooting tools
- [ ] **Comprehensive Testing**: Unit tests, integration tests, and performance benchmarks

### Multi-Vector Store Ecosystem (Q3 2025)
- [ ] **Additional Vector Databases**: Support for Pinecone, Weaviate, Chroma
- [ ] **Pluggable Architecture**: Easy integration of new vector store backends
- [ ] **Performance Comparisons**: Built-in benchmarking between different vector stores
- [ ] **Cross-Store Migration**: Migration tools between any supported vector databases

### User Experience Enhancements (Q3-Q4 2025)
- [ ] **Lightweight Web Interface**: Simple browser-based UI for the existing CLI backend
- [ ] **Knowledge Graph Visualization**: Interactive exploration of document connections
- [ ] **Domain-Specific Templates**: Pre-configured settings for different domains

### Enterprise Features (Q4 2025)
- [ ] **Multi-User Access Control**: Role-based permissions for team environments
- [ ] **Integration with Enterprise Systems**: Connectors for SharePoint, Confluence, Google Workspace
- [ ] **Knowledge Quality Monitoring**: Detection of outdated or contradictory information
- [ ] **System Integration API**: Webhooks and APIs for embedding RLAMA in existing workflows
- [ ] **AI Agent Creation Framework**: Simplified system for building custom AI agents with RAG capabilities

### Next-Gen Retrieval Innovations (Q1 2026)
- [ ] **Multi-Step Retrieval**: Using the LLM to refine search queries for complex questions
- [ ] **Cross-Modal Retrieval**: Support for image content understanding and retrieval
- [ ] **Feedback-Based Optimization**: Learning from user interactions to improve retrieval
- [ ] **Knowledge Graphs & Symbolic Reasoning**: Combining vector search with structured knowledge

### 🚀 **Current Status: Production-Ready with Enterprise Scaling**

RLAMA has evolved from a simple local RAG tool to a comprehensive knowledge management platform that scales from individual developers to enterprise deployments. The recent addition of **Qdrant integration** and **seamless migration tools** represents a major milestone, enabling users to:

- **Start Small**: Begin with internal storage for development and small projects
- **Scale Seamlessly**: Migrate to Qdrant for production workloads with zero data loss
- **Enterprise Ready**: Deploy to Qdrant Cloud or self-hosted instances with full feature parity
- **Future Proof**: Built-in migration paths ensure no vendor lock-in

RLAMA's core philosophy remains unchanged: to provide a simple, powerful, local RAG solution that respects privacy, minimizes resource requirements, and works seamlessly across platforms. Now with the added capability to scale to enterprise-grade vector databases when needed.

## Installation

### Prerequisites
- **For Ollama models**: [Ollama](https://ollama.ai/) installed and running
- **For OpenAI models**: OpenAI API key or API profile configured
- **For OpenAI-compatible servers**: Local server running (e.g., LM Studio, VLLM, etc.)

### Installation from terminal

```bash
curl -fsSL https://raw.githubusercontent.com/dontizi/rlama/main/install.sh | sh
```

## Tech Stack

RLAMA is built with:

- **Core Language**: Go (chosen for performance, cross-platform compatibility, and single binary distribution)
- **CLI Framework**: Cobra (for command-line interface structure)
- **LLM Integration**: Multi-provider support (Ollama, OpenAI, OpenAI-compatible endpoints)
- **Storage**: Local filesystem-based storage (JSON files for simplicity and portability)
- **Vector Search**: Custom implementation of cosine similarity for embedding retrieval

## Architecture

RLAMA follows a clean architecture pattern with clear separation of concerns:

```
rlama/
├── cmd/                  # CLI commands (using Cobra)
│   ├── root.go           # Base command
│   ├── rag.go            # Create RAG systems
│   ├── run.go            # Query RAG systems
│   └── ...
├── internal/
│   ├── client/           # External API clients
│   │   └── ollama_client.go # Ollama API integration
│   ├── domain/           # Core domain models
│   │   ├── rag.go        # RAG system entity
│   │   └── document.go   # Document entity
│   ├── repository/       # Data persistence
│   │   └── rag_repository.go # Handles saving/loading RAGs
│   └── service/          # Business logic
│       ├── rag_service.go      # RAG operations
│       ├── document_loader.go  # Document processing
│       └── embedding_service.go # Vector embeddings
└── pkg/                  # Shared utilities
    └── vector/           # Vector operations
```

## Data Flow

1. **Document Processing**: Documents are loaded from the file system, parsed based on their type, and converted to plain text.
2. **Embedding Generation**: Document text is sent to Ollama to generate vector embeddings.
3. **Storage**: The RAG system (documents + embeddings) is stored in the user's home directory (~/.rlama).
4. **Query Process**: When a user asks a question, it's converted to an embedding, compared against stored document embeddings, and relevant content is retrieved.
5. **Response Generation**: Retrieved content and the question are sent to Ollama to generate a contextually-informed response.

## Visual Representation

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Documents  │────>│  Document   │────>│  Embedding  │
│  (Input)    │     │  Processing │     │  Generation │
└─────────────┘     └─────────────┘     └─────────────┘
                                              │
                                              ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Query     │────>│  Vector     │<────│ Vector Store│
│  Response   │     │  Search     │     │ (RAG System)│
└─────────────┘     └─────────────┘     └─────────────┘
       ▲                   │
       │                   ▼
┌─────────────┐     ┌─────────────┐
│   Ollama    │<────│   Context   │
│    LLM      │     │  Building   │
└─────────────┘     └─────────────┘
```

RLAMA is designed to be lightweight and portable, focusing on providing RAG capabilities with minimal dependencies. The entire system runs locally, with the only external dependency being Ollama for LLM capabilities.

## Available Commands

You can get help on all commands by using:

```bash
rlama --help
```

### Global Flags

These flags can be used with any command:

```bash
--host string   Ollama host (default: localhost)
--port string   Ollama port (default: 11434)
```

### Custom Data Directory

RLAMA stores data in `~/.rlama` by default. To use a different location:

1. **Command-line flag** (highest priority):
   ```bash
   # Use with any command
   rlama --data-dir /path/to/custom/directory run my-rag
   ```

2. **Environment variable**:
   ```bash
   # Set the environment variable
   export RLAMA_DATA_DIR=/path/to/custom/directory
   rlama run my-rag
   ```

The precedence order is: command-line flag > environment variable > default location.

### rag - Create a RAG system

Creates a new RAG system by indexing all documents in the specified folder.

```bash
rlama rag [model] [rag-name] [folder-path]
```

**Parameters:**
- `model`: Name of the Ollama model to use (e.g., llama3, mistral, gemma) or a Hugging Face model using the format `hf.co/username/repository[:quantization]`.
- `rag-name`: Unique name to identify your RAG system.
- `folder-path`: Path to the folder containing your documents.

**Example:**

```bash
# Using a standard Ollama model
rlama rag llama3 documentation ./docs

# Using a Hugging Face model
rlama rag hf.co/bartowski/Llama-3.2-1B-Instruct-GGUF my-rag ./docs

# Using a Hugging Face model with specific quantization
rlama rag hf.co/mlabonne/Meta-Llama-3.1-8B-Instruct-abliterated-GGUF:Q5_K_M my-rag ./docs
```

#### Vector Store Options

RLAMA supports multiple vector storage backends to meet different performance and scaling needs:

**Internal Vector Store (Default)**
- File-based vector storage suitable for local development and small to medium datasets
- No external dependencies required

**Qdrant Vector Database**
- High-performance vector search engine optimized for large-scale semantic search
- Excellent for production environments and large document collections
- Supports advanced filtering and metadata search capabilities

**Using Qdrant Vector Store:**

```bash
# Create a RAG with Qdrant vector store
rlama rag llama3 docs-rag ./docs --vector-store=qdrant

# Customize Qdrant connection
rlama rag llama3 prod-rag ./docs \
  --vector-store=qdrant \
  --qdrant-host=localhost \
  --qdrant-port=6333 \
  --qdrant-collection=my-documents

# Using Qdrant Cloud with API key
rlama rag llama3 cloud-rag ./docs \
  --vector-store=qdrant \
  --qdrant-host=xyz.qdrant.cloud \
  --qdrant-port=6333 \
  --qdrant-apikey=your-api-key \
  --qdrant-grpc
```

**Qdrant Configuration Options:**
- `--vector-store`: Specify "qdrant" to use Qdrant vector database
- `--qdrant-host`: Qdrant server hostname (default: localhost)
- `--qdrant-port`: Qdrant server port (default: 6333)
- `--qdrant-apikey`: API key for Qdrant Cloud or secured instances
- `--qdrant-collection`: Collection name (defaults to RAG name)
- `--qdrant-grpc`: Use gRPC for communication (recommended for performance)

#### Migration Between Vector Stores

RLAMA provides seamless migration tools to move RAG systems between different vector storage backends without losing data.

**Individual RAG Migration:**

```bash
# Migrate from internal to Qdrant
rlama migrate-to-qdrant my-existing-rag \
  --qdrant-host=localhost \
  --qdrant-port=6333 \
  --backup

# Migrate back to internal storage
rlama migrate-to-internal my-qdrant-rag --backup

# Migrate to Qdrant Cloud
rlama migrate-to-qdrant prod-docs \
  --qdrant-host=xyz.qdrant.cloud \
  --qdrant-apikey=your-api-key \
  --qdrant-grpc \
  --backup \
  --verify
```

**Batch Migration:**

```bash
# Migrate all internal RAGs to Qdrant
rlama migrate-batch --from=internal --to=qdrant \
  --qdrant-host=production-server.com \
  --backup \
  --continue-on-error

# Migrate specific RAGs
rlama migrate-batch --from=internal --to=qdrant \
  --rags=docs,wiki,knowledge-base \
  --qdrant-host=localhost
```

**Migration Features:**
- ✅ **Data Integrity**: Automatic verification of migrated data
- ✅ **Backup Support**: Optional backup creation before migration
- ✅ **Progress Tracking**: Real-time progress for large migrations
- ✅ **Error Recovery**: Continue batch operations even if individual RAGs fail
- ✅ **Cleanup Options**: Automatic removal of old data after successful migration

### crawl-rag - Create a RAG system from a website

Creates a new RAG system by crawling a website and indexing its content.

```bash
rlama crawl-rag [model] [rag-name] [website-url]
```

**Parameters:**
- `model`: Name of the Ollama model to use (e.g., llama3, mistral, gemma).
- `rag-name`: Unique name to identify your RAG system.
- `website-url`: URL of the website to crawl and index.

**Options:**
- `--max-depth`: Maximum crawl depth (default: 2)
- `--concurrency`: Number of concurrent crawlers (default: 5)
- `--exclude-path`: Paths to exclude from crawling (comma-separated)
- `--chunk-size`: Character count per chunk (default: 1000)
- `--chunk-overlap`: Overlap between chunks in characters (default: 200)
- `--chunking-strategy`: Chunking strategy to use (options: "fixed", "semantic", "hybrid", "hierarchical", default: "hybrid")

#### Chunking Strategies

RLAMA offers multiple advanced chunking strategies to optimize document retrieval:

- **Fixed**: Traditional chunking with fixed size and overlap, respecting sentence boundaries when possible.
- **Semantic**: Intelligently splits documents based on semantic boundaries like headings, paragraphs, and natural topic shifts.
- **Hybrid**: Automatically selects the best strategy based on document type and content (markdown, HTML, code, or plain text).
- **Hierarchical**: For very long documents, creates a two-level chunking structure with major sections and sub-chunks.

The system automatically adapts to different document types:
- Markdown documents: Split by headers and sections
- HTML documents: Split by semantic HTML elements
- Code documents: Split by functions, classes, and logical blocks
- Plain text: Split by paragraphs with contextual overlap

#### Reranking Options

RLAMA includes BGE-based reranking by default to improve result quality. Two implementations are available:

**Python BGE Reranker (Default)**
- Uses the original Python FlagEmbedding library via subprocess calls
- Works out of the box with existing Python environment

**ONNX BGE Reranker (Faster)**
- Uses optimized ONNX models for **3.8x faster performance**
- Requires one-time setup but provides significant speed improvements

```bash
# Download ONNX model (one-time setup)
mkdir -p ./models
cd ./models
git clone https://huggingface.co/corto-ai/bge-reranker-large-onnx

# Use ONNX reranker for faster performance
rlama rag llama3.2 myrag ./docs --use-onnx-reranker

# Specify custom ONNX model directory
rlama rag llama3.2 myrag ./docs --use-onnx-reranker --onnx-model-dir ./models/bge-reranker-large-onnx
```

**ONNX Requirements:**
```bash
pip install onnxruntime transformers numpy
```

**Reranking Configuration Options:**
- `--disable-reranker`: Disable reranking (enabled by default)
- `--reranker-model`: Model to use for reranking (defaults to main model)
- `--reranker-weight`: Weight for reranker scores vs vector scores (0-1, default: 0.7)
- `--reranker-threshold`: Minimum score threshold for reranked results (default: 0.0)
- `--use-onnx-reranker`: Use ONNX reranker for faster performance
- `--onnx-model-dir`: Directory containing ONNX reranker model (default: ./models/bge-reranker-large-onnx)

**Performance Comparison:**
- **Python BGE**: ~7.4 seconds per query
- **ONNX BGE**: ~2.0 seconds per query (**3.8x faster**)

**Example:**

```bash
# Create a new RAG from a documentation website
rlama crawl-rag llama3 docs-rag https://docs.example.com

# Customize crawling behavior
rlama crawl-rag llama3 blog-rag https://blog.example.com --max-depth=3 --exclude-path=/archive,/tags

# Create a RAG with semantic chunking
rlama rag llama3 documentation ./docs --chunking-strategy=semantic

# Use hierarchical chunking for large documents
rlama rag llama3 book-rag ./books --chunking-strategy=hierarchical
```

### wizard - Create a RAG system with interactive setup

Provides an interactive step-by-step wizard for creating a new RAG system.

```bash
rlama wizard
```

The wizard guides you through:
- Naming your RAG
- Choosing an Ollama model
- Selecting document sources (local folder or website)
- Configuring chunking parameters
- Setting up file filtering

**Example:**

```bash
rlama wizard
# Follow the prompts to create your customized RAG
```

### watch - Set up directory watching for a RAG system

Configure a RAG system to automatically watch a directory for new files and add them to the RAG.

```bash
rlama watch [rag-name] [directory-path] [interval]
```

**Parameters:**
- `rag-name`: Name of the RAG system to watch.
- `directory-path`: Path to the directory to watch for new files.
- `interval`: Time in minutes to check for new files (use 0 to check only when the RAG is used).

**Example:**

```bash
# Set up directory watching to check every 60 minutes
rlama watch my-docs ./watched-folder 60

# Set up directory watching to only check when the RAG is used
rlama watch my-docs ./watched-folder 0

# Customize what files to watch
rlama watch my-docs ./watched-folder 30 --exclude-dir=node_modules,tmp --process-ext=.md,.txt
```

### watch-off - Disable directory watching for a RAG system

Disable automatic directory watching for a RAG system.

```bash
rlama watch-off [rag-name]
```

**Parameters:**
- `rag-name`: Name of the RAG system to disable watching.

**Example:**

```bash
rlama watch-off my-docs
```

### check-watched - Check a RAG's watched directory for new files

Manually check a RAG's watched directory for new files and add them to the RAG.

```bash
rlama check-watched [rag-name]
```

**Parameters:**
- `rag-name`: Name of the RAG system to check.

**Example:**

```bash
rlama check-watched my-docs
```

### web-watch - Set up website monitoring for a RAG system

Configure a RAG system to automatically monitor a website for updates and add new content to the RAG.

```bash
rlama web-watch [rag-name] [website-url] [interval]
```

**Parameters:**
- `rag-name`: Name of the RAG system to monitor.
- `website-url`: URL of the website to monitor.
- `interval`: Time in minutes between checks (use 0 to check only when the RAG is used).

**Example:**

```bash
# Set up website monitoring to check every 60 minutes
rlama web-watch my-docs https://example.com 60

# Set up website monitoring to only check when the RAG is used
rlama web-watch my-docs https://example.com 0

# Customize what content to monitor
rlama web-watch my-docs https://example.com 30 --exclude-path=/archive,/tags
```

### web-watch-off - Disable website monitoring for a RAG system

Disable automatic website monitoring for a RAG system.

```bash
rlama web-watch-off [rag-name]
```

**Parameters:**
- `rag-name`: Name of the RAG system to disable monitoring.

**Example:**

```bash
rlama web-watch-off my-docs
```

### check-web-watched - Check a RAG's monitored website for updates

Manually check a RAG's monitored website for new updates and add them to the RAG.

```bash
rlama check-web-watched [rag-name]
```

**Parameters:**
- `rag-name`: Name of the RAG system to check.

**Example:**

```bash
rlama check-web-watched my-docs
```

### run - Use a RAG system

Starts an interactive session to interact with an existing RAG system.

```bash
rlama run [rag-name]
```

**Parameters:**
- `rag-name`: Name of the RAG system to use.
- `--context-size`: (Optional) Number of context chunks to retrieve (default: 20)

**Example:**

```bash
rlama run documentation
> How do I install the project?
> What are the main features?
> exit
```

**Context Size Tips:**
- Smaller values (5-15) for faster responses with key information
- Medium values (20-40) for balanced performance
- Larger values (50+) for complex questions needing broad context
- Consider your model's context window limits

```bash
rlama run documentation --context-size=50  # Use 50 context chunks
```

### api - Start API server

Starts an HTTP API server that exposes RLAMA's functionality through RESTful endpoints.

```bash
rlama api [--port PORT]
```

**Parameters:**
- `--port`: (Optional) Port number to run the API server on (default: 11249)

**Example:**

```bash
rlama api --port 8080
```

**Available Endpoints:**

1. **Query a RAG system** - `POST /rag`
   ```bash
   curl -X POST http://localhost:11249/rag \
     -H "Content-Type: application/json" \
     -d '{
       "rag_name": "documentation",
       "prompt": "How do I install the project?",
       "context_size": 20
     }'
   ```

   Request fields:
   - `rag_name` (required): Name of the RAG system to query
   - `prompt` (required): Question or prompt to send to the RAG
   - `context_size` (optional): Number of chunks to include in context
   - `model` (optional): Override the model used by the RAG

2. **Check server health** - `GET /health`
   ```bash
   curl http://localhost:11249/health
   ```

**Integration Example:**
```javascript
// Node.js example
const response = await fetch('http://localhost:11249/rag', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    rag_name: 'my-docs',
    prompt: 'Summarize the key features'
  })
});
const data = await response.json();
console.log(data.response);
```

### list - List RAG systems

Displays a list of all available RAG systems.

```bash
rlama list
```

### delete - Delete a RAG system

Permanently deletes a RAG system and all its indexed documents.

```bash
rlama delete [rag-name] [--force/-f]
```

**Parameters:**
- `rag-name`: Name of the RAG system to delete.
- `--force` or `-f`: (Optional) Delete without asking for confirmation.

**Example:**

```bash
rlama delete old-project
```

Or to delete without confirmation:

```bash
rlama delete old-project --force
```

### list-docs - List documents in a RAG

Displays all documents in a RAG system with metadata.

```bash
rlama list-docs [rag-name]
```

**Parameters:**
- `rag-name`: Name of the RAG system

**Example:**

```bash
rlama list-docs documentation
```

### list-chunks - Inspect document chunks

List and filter document chunks in a RAG system with various options:

```bash
# Basic chunk listing
rlama list-chunks [rag-name]

# With content preview (shows first 100 characters)
rlama list-chunks [rag-name] --show-content

# Filter by document name/ID substring
rlama list-chunks [rag-name] --document=readme

# Combine options
rlama list-chunks [rag-name] --document=api --show-content
```

**Options:**
- `--show-content`: Display chunk content preview
- `--document`: Filter by document name/ID substring

**Output columns:**
- Chunk ID (use with view-chunk command)
- Document Source
- Chunk Position (e.g., "2/5" for second of five chunks)
- Content Preview (if enabled)
- Created Date

### view-chunk - View chunk details

Display detailed information about a specific chunk.

```bash
rlama view-chunk [rag-name] [chunk-id]
```

**Parameters:**
- `rag-name`: Name of the RAG system
- `chunk-id`: Chunk identifier from list-chunks

**Example:**

```bash
rlama view-chunk documentation doc123_chunk_0
```

### add-docs - Add documents to RAG

Add new documents to an existing RAG system.

```bash
rlama add-docs [rag-name] [folder-path] [flags]
```

**Parameters:**
- `rag-name`: Name of the RAG system
- `folder-path`: Path to documents folder

**Example:**

```bash
rlama add-docs documentation ./new-docs --exclude-ext=.tmp
```

### crawl-add-docs - Add website content to RAG

Add content from a website to an existing RAG system.

```bash
rlama crawl-add-docs [rag-name] [website-url]
```

**Parameters:**
- `rag-name`: Name of the RAG system
- `website-url`: URL of the website to crawl and add to the RAG

**Options:**
- `--max-depth`: Maximum crawl depth (default: 2)
- `--concurrency`: Number of concurrent crawlers (default: 5)
- `--exclude-path`: Paths to exclude from crawling (comma-separated)
- `--chunk-size`: Character count per chunk (default: 1000)
- `--chunk-overlap`: Overlap between chunks in characters (default: 200)

**Example:**

```bash
# Add blog content to an existing RAG
rlama crawl-add-docs my-docs https://blog.example.com

# Customize crawling behavior
rlama crawl-add-docs knowledge-base https://docs.example.com --max-depth=1 --exclude-path=/api
```

### update-model - Change LLM model

Update the LLM model used by a RAG system.

```bash
rlama update-model [rag-name] [new-model]
```

**Parameters:**
- `rag-name`: Name of the RAG system
- `new-model`: New Ollama model name

**Example:**

```bash
rlama update-model documentation deepseek-r1:7b-instruct
```

### profile - Manage API profiles

Manage API profiles for different LLM providers and endpoints.

#### profile add - Create a new profile

```bash
rlama profile add [name] [provider] [api-key] [flags]
```

**Parameters:**
- `name`: Unique name for the profile
- `provider`: Provider type (`openai` or `openai-api`)
- `api-key`: API key (use "none" for local servers without authentication)

**Flags:**
- `--base-url`: Custom base URL for OpenAI-compatible endpoints (required for `openai-api` provider)

**Examples:**

```bash
# Traditional OpenAI profile
rlama profile add openai-work openai sk-your-api-key

# LM Studio local server
rlama profile add lmstudio openai-api none --base-url http://localhost:1234/v1

# VLLM server with authentication
rlama profile add vllm openai-api your-token --base-url http://server:8000/v1
```

#### profile list - List all profiles

```bash
rlama profile list
```

#### profile delete - Delete a profile

```bash
rlama profile delete [name]
```

### update - Update RLAMA

Checks if a new version of RLAMA is available and installs it.

```bash
rlama update [--force/-f]
```

**Options:**
- `--force` or `-f`: (Optional) Update without asking for confirmation.

### version - Display version

Displays the current version of RLAMA.

```bash
rlama --version
```

or

```bash
rlama -v
```

### hf-browse - Browse GGUF models on Hugging Face

Search and browse GGUF models available on Hugging Face.

```bash
rlama hf-browse [search-term] [flags]
```

**Parameters:**
- `search-term`: (Optional) Term to search for (e.g., "llama3", "mistral")

**Flags:**
- `--open`: Open the search results in your default web browser
- `--quant`: Specify quantization type to suggest (e.g., Q4_K_M, Q5_K_M)
- `--limit`: Limit number of results (default: 10)

**Examples:**

```bash
# Search for GGUF models and show command-line help
rlama hf-browse "llama 3"

# Open browser with search results
rlama hf-browse mistral --open

# Search with specific quantization suggestion
rlama hf-browse phi --quant Q4_K_M
```

### run-hf - Run a Hugging Face GGUF model

Run a Hugging Face GGUF model directly using Ollama. This is useful for testing models before creating a RAG system with them.

```bash
rlama run-hf [huggingface-model] [flags]
```

**Parameters:**
- `huggingface-model`: Hugging Face model path in the format `username/repository`

**Flags:**
- `--quant`: Quantization to use (e.g., Q4_K_M, Q5_K_M)

**Examples:**

```bash
# Try a model in chat mode
rlama run-hf bartowski/Llama-3.2-1B-Instruct-GGUF

# Specify quantization
rlama run-hf mlabonne/Meta-Llama-3.1-8B-Instruct-abliterated-GGUF --quant Q5_K_M
```

## Uninstallation

To uninstall RLAMA:

### Removing the binary

If you installed via `go install`:

```bash
rlama uninstall
```

### Removing data

RLAMA stores its data in `~/.rlama`. To remove it:

```bash
rm -rf ~/.rlama
```

## Supported Document Formats

RLAMA supports many file formats:

- **Text**: `.txt`, `.md`, `.html`, `.json`, `.csv`, `.yaml`, `.yml`, `.xml`, `.org`
- **Code**: `.go`, `.py`, `.js`, `.java`, `.c`, `.cpp`, `.cxx`, `.h`, `.rb`, `.php`, `.rs`, `.swift`, `.kt`, `.ts`, `.tsx`, `.f`, `.F`, `.F90`, `.el`, `.svelte`
- **Documents**: `.pdf`, `.docx`, `.doc`, `.rtf`, `.odt`, `.pptx`, `.ppt`, `.xlsx`, `.xls`, `.epub`

Installing dependencies via `install_deps.sh` is recommended to improve support for certain formats.

## Troubleshooting

### Ollama is not accessible

If you encounter connection errors to Ollama:
1. Check that Ollama is running.
2. By default, Ollama must be accessible at `http://localhost:11434` or the host and port specified by the OLLAMA_HOST environment variable.
3. If your Ollama instance is running on a different host or port, use the `--host` and `--port` flags:
   ```bash
   rlama --host 192.168.1.100 --port 8000 list
   rlama --host my-ollama-server --port 11434 run my-rag
   ```
4. Check Ollama logs for potential errors.

### Text extraction issues

If you encounter problems with certain formats:
1. Install dependencies via `./scripts/install_deps.sh`.
2. Verify that your system has the required tools (`pdftotext`, `tesseract`, etc.).

### The RAG doesn't find relevant information

If the answers are not relevant:
1. Check that the documents are properly indexed with `rlama list`.
2. Make sure the content of the documents is properly extracted.
3. Try rephrasing your question more precisely.
4. Consider adjusting chunking parameters during RAG creation

### Other issues

For any other issues, please open an issue on the [GitHub repository](https://github.com/dontizi/rlama/issues) providing:
1. The exact command used.
2. The complete output of the command.
3. Your operating system and architecture.
4. The RLAMA version (`rlama --version`).

### Configuring Ollama Connection

RLAMA provides multiple ways to connect to your Ollama instance:

1. **Command-line flags** (highest priority):
   ```bash
   rlama --host 192.168.1.100 --port 8080 run my-rag
   ```

2. **Environment variable**:
   ```bash
   # Format: "host:port" or just "host"
   export OLLAMA_HOST=remote-server:8080
   rlama run my-rag
   ```

3. **Default values** (used if no other method is specified):
   - Host: `localhost`
   - Port: `11434`

The precedence order is: command-line flags > environment variable > default values.

## Advanced Usage

### Context Size Management

```bash
# Quick answers with minimal context
rlama run my-docs --context-size=10

# Deep analysis with maximum context
rlama run my-docs --context-size=50

# Balance between speed and depth
rlama run my-docs --context-size=30
```

### RAG Creation with Filtering
```bash
rlama rag llama3 my-project ./code \
  --exclude-dir=node_modules,dist \
  --process-ext=.go,.ts \
  --exclude-ext=.spec.ts
```

### Chunk Inspection
```bash
# List chunks with content preview
rlama list-chunks my-project --show-content

# Filter chunks from specific document
rlama list-chunks my-project --document=architecture
```

## Help System

Get full command help:
```bash
rlama --help
```

Command-specific help:
```bash
rlama rag --help
rlama list-chunks --help
rlama update-model --help
```

All commands support the global `--host` and `--port` flags for custom Ollama connections.

The precedence order is: command-line flags > environment variable > default values.

## Hugging Face Integration

RLAMA now supports using GGUF models directly from Hugging Face through Ollama's native integration:

### Browsing Hugging Face Models

```bash
# Search for GGUF models on Hugging Face
rlama hf-browse "llama 3"

# Open browser with search results
rlama hf-browse mistral --open
```

### Testing a Model

Before creating a RAG, you can test a Hugging Face model directly:

```bash
# Try a model in chat mode
rlama run-hf bartowski/Llama-3.2-1B-Instruct-GGUF

# Specify quantization
rlama run-hf mlabonne/Meta-Llama-3.1-8B-Instruct-abliterated-GGUF --quant Q5_K_M
```

### Creating a RAG with Hugging Face Models

Use Hugging Face models when creating RAG systems:

```bash
# Create a RAG with a Hugging Face model
rlama rag hf.co/bartowski/Llama-3.2-1B-Instruct-GGUF my-rag ./docs

# Use specific quantization
rlama rag hf.co/mlabonne/Meta-Llama-3.1-8B-Instruct-abliterated-GGUF:Q5_K_M my-rag ./docs
```

## Model Support & LLM Providers

RLAMA supports multiple LLM providers for both **text generation** and **embeddings**:

### Supported Providers

1. **Ollama** (default): Local models via Ollama server
2. **OpenAI**: Official OpenAI API  
3. **OpenAI-Compatible**: Any server implementing OpenAI API (LM Studio, VLLM, TGI, etc.)

### How Models Are Used

RLAMA uses models for two distinct purposes:

- **Text Generation (Completions)**: Answering your questions using retrieved context
- **Embeddings**: Converting documents and queries into vectors for similarity search

### Model Selection Logic

When you specify a model name, RLAMA automatically determines which provider to use:

- **OpenAI models** (e.g., `gpt-4`, `gpt-3.5-turbo`): Uses OpenAI API for completions + embeddings
- **Hugging Face models** (e.g., `hf.co/username/model`): Downloads via Ollama
- **Other models** (e.g., `llama3`, `mistral`): Uses Ollama for completions + embeddings

### Using OpenAI Models

Set your OpenAI API key:
```bash
export OPENAI_API_KEY="your-api-key"
```

Create a RAG with OpenAI model:
```bash
rlama rag gpt-4 my-rag ./documents
```

Supported OpenAI models:
- `gpt-4`, `gpt-4-turbo`, `gpt-4o`
- `gpt-3.5-turbo`
- `o3-mini` and newer models

### Using OpenAI-Compatible Endpoints

RLAMA can connect to any server that implements the OpenAI API specification, including:

- **LM Studio**: Local model serving with OpenAI API
- **VLLM**: High-performance inference server  
- **Text Generation Inference (TGI)**: Hugging Face's inference server
- **Ollama's OpenAI compatibility mode**: `ollama serve` with OpenAI endpoints
- **Any custom OpenAI-compatible server**

#### Setting Up Profiles for Custom Endpoints

Create a profile for your OpenAI-compatible server:

```bash
# For LM Studio running locally
rlama profile add lmstudio openai-api none --base-url http://localhost:1234/v1

# For VLLM server (with authentication)
rlama profile add vllm openai-api your-api-key --base-url http://your-server:8000/v1

# For remote TGI server
rlama profile add tgi openai-api dummy --base-url https://tgi.example.com/v1
```

#### Using Custom Endpoints

Create a RAG with your custom endpoint:

```bash
# Use the profile when creating a RAG
rlama rag llama-3-8b my-rag ./documents --profile lmstudio

# The model name should match what your server expects
rlama rag custom-model-name knowledge-base ./docs --profile vllm
```

#### Common OpenAI-Compatible Servers

1. **LM Studio**:
   ```bash
   # Start LM Studio with OpenAI API on default port 1234
   rlama profile add lmstudio openai-api none --base-url http://localhost:1234/v1
   rlama rag llama-3-8b-instruct my-docs ./documents --profile lmstudio
   ```

2. **VLLM**:
   ```bash
   # VLLM typically runs on port 8000
   rlama profile add vllm openai-api none --base-url http://localhost:8000/v1
   rlama rag meta-llama/Llama-3-8B-Instruct my-rag ./docs --profile vllm
   ```

3. **Ollama OpenAI Mode**:
   ```bash
   # If using Ollama's experimental OpenAI endpoints
   rlama profile add ollama-openai openai-api none --base-url http://localhost:11434/v1
   rlama rag llama3 my-rag ./docs --profile ollama-openai
   ```

#### Benefits of OpenAI-Compatible Mode

- **Unified Interface**: Same API for different inference engines
- **Easy Migration**: Switch between providers without changing RAG structure  
- **Better Performance**: Use optimized inference servers (VLLM, TGI)
- **Model Flexibility**: Access models not available through Ollama
- **Embedding Support**: Full support for both completions and embeddings

## Managing API Profiles

RLAMA allows you to create API profiles to manage multiple API keys and endpoints for different providers:

### Profile Types

- **`openai`**: Official OpenAI API profiles
- **`openai-api`**: Generic OpenAI-compatible endpoints (LM Studio, VLLM, etc.)

### Creating Profiles

#### Traditional OpenAI Profiles
```bash
# Create a profile for your OpenAI account
rlama profile add openai-work openai "sk-your-api-key"

# Create another profile for a different account
rlama profile add openai-personal openai "sk-your-personal-api-key" 
```

#### OpenAI-Compatible Endpoint Profiles
```bash
# LM Studio local server (no API key needed)
rlama profile add lmstudio openai-api none --base-url http://localhost:1234/v1

# VLLM server with authentication
rlama profile add vllm-server openai-api your-token --base-url http://192.168.1.100:8000/v1

# Remote TGI deployment
rlama profile add tgi-prod openai-api api-key --base-url https://api.mycompany.com/v1
```

### Listing Profiles

```bash
# View all available profiles with their base URLs
rlama profile list
```

Output example:
```
NAME         PROVIDER    BASE URL                  CREATED ON           LAST USED
openai-work  openai      default                   2024-01-15 10:30:25  2024-01-16 14:22:10
lmstudio     openai-api  http://localhost:1234/v1  2024-01-16 09:15:33  never
vllm-server  openai-api  http://server:8000/v1     2024-01-16 11:45:12  2024-01-16 15:30:25
```

### Deleting Profiles

```bash
# Delete a profile
rlama profile delete openai-old
```

### Using Profiles with RAGs

When creating a new RAG:

```bash
# Create a RAG with an OpenAI profile
rlama rag gpt-4 my-rag ./documents --profile openai-work

# Create a RAG with a custom endpoint
rlama rag llama-3-8b local-rag ./docs --profile lmstudio
```

When running existing RAGs:

```bash
# RAGs remember their original configuration automatically
rlama run my-rag
```

### Profile Benefits

- **Multiple Endpoints**: Manage connections to different LLM servers
- **Easy Switching**: Change between local and remote inference
- **Secure Storage**: API keys stored safely in `~/.rlama/profiles`
- **Usage Tracking**: See when profiles were last used
- **Project Organization**: Use different profiles for different projects
- **Development Workflow**: Test locally (LM Studio) → deploy remotely (VLLM)
````

## File: internal/service/rag_service.go
````go
package service
⋮----
import (
	"fmt"
	"strings"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/repository"
)
⋮----
"fmt"
"strings"
⋮----
"github.com/dontizi/rlama/internal/client"
"github.com/dontizi/rlama/internal/domain"
"github.com/dontizi/rlama/internal/repository"
⋮----
// RagService interface defines the contract for RAG operations
type RagService interface {
	CreateRagWithOptions(modelName, ragName, folderPath string, options DocumentLoaderOptions) error
	GetRagChunks(ragName string, filter ChunkFilter) ([]*domain.DocumentChunk, error)
	LoadRag(ragName string) (*domain.RagSystem, error)
	Query(rag *domain.RagSystem, query string, contextSize int) (string, error)
	AddDocsWithOptions(ragName string, folderPath string, options DocumentLoaderOptions) error
	UpdateModel(ragName string, newModel string) error
	UpdateRag(rag *domain.RagSystem) error
	UpdateRerankerModel(ragName string, model string) error
	ListAllRags() ([]string, error)
	GetOllamaClient() *client.OllamaClient
	SetPreferredEmbeddingModel(model string)
	// Directory watching methods
	SetupDirectoryWatching(ragName string, dirPath string, watchInterval int, options DocumentLoaderOptions) error
	DisableDirectoryWatching(ragName string) error
	CheckWatchedDirectory(ragName string) (int, error)
	// Web watching methods
	SetupWebWatching(ragName string, websiteURL string, watchInterval int, options domain.WebWatchOptions) error
	DisableWebWatching(ragName string) error
	CheckWatchedWebsite(ragName string) (int, error)
}
⋮----
// Directory watching methods
⋮----
// Web watching methods
⋮----
// ChunkFilter defines filtering criteria for retrieving chunks
type ChunkFilter struct {
	DocumentSubstring string
	ShowContent       bool
}
⋮----
// RagServiceImpl implements the RagService interface
type RagServiceImpl struct {
	documentLoader   *DocumentLoader
	embeddingService *EmbeddingService
	ragRepository    *repository.RagRepository
	ollamaClient     *client.OllamaClient
	rerankerService  *RerankerService
}
⋮----
// NewRagService creates a new instance of RagService
func NewRagService(ollamaClient *client.OllamaClient) RagService
⋮----
// NewRagServiceWithConfig creates a new instance of RagService with service configuration
func NewRagServiceWithConfig(ollamaClient *client.OllamaClient, config *ServiceConfig) RagService
⋮----
// Create reranker service with ONNX configuration if specified
var rerankerService *RerankerService
⋮----
// NewRagServiceWithClient creates a new instance of RagService with the specified LLM client
func NewRagServiceWithClient(llmClient client.LLMClient, ollamaClient *client.OllamaClient) RagService
⋮----
// Use the new composite service architecture
⋮----
// NewRagServiceWithEmbedding creates a new RagService with a specific embedding service
func NewRagServiceWithEmbedding(ollamaClient *client.OllamaClient, embeddingService *EmbeddingService) RagService
⋮----
// GetOllamaClient returns the Ollama client
func (rs *RagServiceImpl) GetOllamaClient() *client.OllamaClient
⋮----
// SetPreferredEmbeddingModel sets the preferred embedding model to use
func (rs *RagServiceImpl) SetPreferredEmbeddingModel(model string)
⋮----
// CreateRagWithOptions creates a new RAG system with options
func (rs *RagServiceImpl) CreateRagWithOptions(modelName, ragName, folderPath string, options DocumentLoaderOptions) error
⋮----
// Check if model is available using the correct client
// The embedding service has the right LLM client (OpenAI or Ollama)
⋮----
// Fallback to Ollama client if embedding service doesn't have a client
⋮----
// Fallback to Ollama client if embedding service is not properly configured
⋮----
// Check if the RAG already exists
⋮----
// Load documents with options
⋮----
// Detect embedding dimension
⋮----
// Create the RAG system with detected dimensions and vector store configuration
var rag *domain.RagSystem
⋮----
// Configure reranking options - enable by default
rag.RerankerEnabled = true // Always enable reranking by default
⋮----
// Only disable if explicitly set to false in options
⋮----
// Check if EnableReranker field was explicitly set
// This prevents the zero-value (false) from disabling reranking when the field isn't set
⋮----
// Set reranker model if specified, otherwise use the same model
⋮----
// Set reranker weight
⋮----
rag.RerankerWeight = 0.7 // Default to 70% reranker, 30% vector
⋮----
// Set default TopK if not already set
⋮----
rag.RerankerTopK = 5 // Default to 5 results
⋮----
// Set chunking options in WatchOptions too
⋮----
// Create chunker service
⋮----
// Process each document - chunk and generate embeddings
var allChunks []*domain.DocumentChunk
⋮----
// Add the document to the RAG
⋮----
// Chunk the document
⋮----
// Update total chunks in metadata
⋮----
// Set preferred embedding model if specified
⋮----
// Generate embeddings for all chunks
⋮----
// Add all chunks to the RAG
⋮----
// Save the RAG
⋮----
// GetRagChunks gets chunks from a RAG with filtering
func (rs *RagServiceImpl) GetRagChunks(ragName string, filter ChunkFilter) ([]*domain.DocumentChunk, error)
⋮----
// Load the RAG
⋮----
var filteredChunks []*domain.DocumentChunk
⋮----
// Apply filters
⋮----
// Apply document name filter if provided
⋮----
// LoadRag loads a RAG system
func (rs *RagServiceImpl) LoadRag(ragName string) (*domain.RagSystem, error)
⋮----
// Query performs a query on a RAG system
func (rs *RagServiceImpl) Query(rag *domain.RagSystem, query string, contextSize int) (string, error)
⋮----
// Use the embedding service's LLM client for consistency
⋮----
// The embedding service already has the right client configured
⋮----
// Generate embedding for the query
⋮----
// Use the provided context size or default value based on settings
⋮----
// Si contextSize est 0 (auto), utiliser:
// - RerankerTopK du RAG si défini
// - Sinon le TopK par défaut (5)
// - 20 si le reranking est désactivé
⋮----
contextSize = rerankerOpts.TopK // 5 par défaut
⋮----
contextSize = 20 // 20 par défaut si le reranking est désactivé
fmt.Printf("Using context size of %d (reranking disabled)\n", contextSize) // Always show this message since reranking is disabled
⋮----
// First-stage retrieval: Get initial results using vector search
// Get more results than needed for reranking
⋮----
// If reranking is enabled, retrieve more documents initially (20 or 2*contextSize, whichever is larger)
⋮----
initialRetrievalCount = contextSize * 2 // Ensure we get enough documents for reranking
⋮----
// Search for the most relevant chunks
⋮----
// Second-stage retrieval: Re-rank if enabled
var rankedResults []RankedResult
var includedDocs = make(map[string]bool)
⋮----
// Set reranker options for adaptive content-based filtering
⋮----
// Don't limit by fixed TopK but use minimum threshold
TopK:              100, // Set to a high value to avoid arbitrary limit
⋮----
RerankerModel:     "BAAI/bge-reranker-v2-m3", // Always prefer BGE reranker
ScoreThreshold:    0.3,                       // Minimum relevance threshold
⋮----
AdaptiveFiltering: true, // Enable adaptive filtering
Silent:            rag.RerankerSilent, // Use the silent setting from the RAG
⋮----
// If a specific BGE reranker model is defined in the RAG, use that one
// This allows users to choose between different BGE reranker models
⋮----
// Display the effective model being used (if not in silent mode)
⋮----
// Perform reranking with adaptive filtering
⋮----
// Track documents included after adaptive filtering
⋮----
// Show information about filtered results
⋮----
// Build the context
var context strings.Builder
⋮----
// Use the reranked results if available, otherwise use the initial results
⋮----
// Add chunk content with its metadata
⋮----
// Use original vector search results if reranking is disabled or failed
⋮----
// Add chunk content with its metadata
⋮----
// Build the prompt with better formatting and instructions for citing sources
⋮----
// Show search results to the user
⋮----
// Generate the response with the appropriate client
⋮----
// AddDocsWithOptions adds documents to a RAG with options
func (rs *RagServiceImpl) AddDocsWithOptions(ragName string, folderPath string, options DocumentLoaderOptions) error
⋮----
// Load the existing RAG system
⋮----
// Check if Ollama is available
⋮----
// Load new documents with options
⋮----
// Create chunker service with the same options as the RAG or from provided options
⋮----
// Override with provided options if specified
⋮----
// Create chunker with configured options
⋮----
// Check for duplicates
⋮----
var uniqueDocs []*domain.Document
var skippedDocs int
⋮----
// Filter out duplicate documents
⋮----
existingDocPaths[doc.Path] = true // Mark as processed to avoid future duplicates
⋮----
// Process each unique document - chunk and generate embeddings
⋮----
// Update the RAG's chunk options based on the most recent settings
⋮----
// Update reranker settings if specified in options
⋮----
// Save the updated RAG
⋮----
// UpdateModel updates the model of a RAG
func (rs *RagServiceImpl) UpdateModel(ragName string, newModel string) error
⋮----
// UpdateRag updates a RAG system
func (rs *RagServiceImpl) UpdateRag(rag *domain.RagSystem) error
⋮----
// ListAllRags lists all available RAGs
func (rs *RagServiceImpl) ListAllRags() ([]string, error)
⋮----
// SetupDirectoryWatching sets up directory watching for a RAG
func (rs *RagServiceImpl) SetupDirectoryWatching(ragName string, dirPath string, watchInterval int, options DocumentLoaderOptions) error
⋮----
// DisableDirectoryWatching disables directory watching for a RAG
func (rs *RagServiceImpl) DisableDirectoryWatching(ragName string) error
⋮----
// CheckWatchedDirectory checks a watched directory for changes
func (rs *RagServiceImpl) CheckWatchedDirectory(ragName string) (int, error)
⋮----
// SetupWebWatching sets up web watching for a RAG
func (rs *RagServiceImpl) SetupWebWatching(ragName string, websiteURL string, watchInterval int, options domain.WebWatchOptions) error
⋮----
// DisableWebWatching disables web watching for a RAG
func (rs *RagServiceImpl) DisableWebWatching(ragName string) error
⋮----
// CheckWatchedWebsite checks a watched website for changes
func (rs *RagServiceImpl) CheckWatchedWebsite(ragName string) (int, error)
⋮----
// UpdateRerankerModel updates the reranker model of a RAG
func (rs *RagServiceImpl) UpdateRerankerModel(ragName string, model string) error
⋮----
// Update the reranker model
````
