package cmd

import (
	"fmt"
	"strings"

	"github.com/dontizi/rlama/internal/repository"
	"github.com/dontizi/rlama/internal/service"
	"github.com/spf13/cobra"
)

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

func init() {
	rootCmd.AddCommand(migrateBatchCmd)

	// Required flags
	migrateBatchCmd.Flags().StringVar(&batchFromVectorStore, "from", "", "Source vector store type (internal, qdrant)")
	migrateBatchCmd.Flags().StringVar(&batchToVectorStore, "to", "", "Target vector store type (internal, qdrant)")
	migrateBatchCmd.MarkFlagRequired("from")
	migrateBatchCmd.MarkFlagRequired("to")

	// Add common flag sets
	AddQdrantFlags(migrateBatchCmd, &batchQdrantFlags, "Qdrant collection name (will be prefixed if --qdrant-collection-prefix is set)")
	AddMigrationControlFlags(migrateBatchCmd, &batchMigrationFlags)
	AddBatchMigrationFlags(migrateBatchCmd, &batchRagNames, &batchContinueOnError)

	// Batch-specific Qdrant flag
	migrateBatchCmd.Flags().StringVar(&batchQdrantCollectionPrefix, "qdrant-collection-prefix", "", "Prefix for Qdrant collection names (default: use RAG names)")
}

// filterRagsByVectorStore returns only RAGs that use the specified vector store type
func filterRagsByVectorStore(ragNames []string, vectorStoreType string) ([]string, error) {
	ragRepo := repository.NewRagRepository()
	var filtered []string

	for _, ragName := range ragNames {
		rag, err := ragRepo.Load(ragName)
		if err != nil {
			fmt.Printf("⚠️ Warning: Failed to load RAG '%s': %v\n", ragName, err)
			continue
		}

		// Check vector store type (default to internal if not set)
		ragVectorStore := rag.VectorStoreType
		if ragVectorStore == "" {
			ragVectorStore = "internal"
		}

		if ragVectorStore == vectorStoreType {
			filtered = append(filtered, ragName)
		}
	}

	return filtered, nil
}

// getCollectionName generates the Qdrant collection name
func getCollectionName(ragName, prefix string) string {
	if prefix == "" {
		return ragName
	}
	return fmt.Sprintf("%s-%s", prefix, ragName)
}