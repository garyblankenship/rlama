package cmd

import (
	"fmt"

	"github.com/dontizi/rlama/internal/service"
	"github.com/spf13/cobra"
)

var (
	// Flag holders for migrate-to-qdrant command
	qdrantFlags    QdrantFlags
	migrationFlags MigrationFlags
)

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

func init() {
	rootCmd.AddCommand(migrateToQdrantCmd)

	// Add common flag sets
	AddQdrantFlags(migrateToQdrantCmd, &qdrantFlags, "Qdrant collection name (defaults to RAG name)")
	AddMigrationControlFlags(migrateToQdrantCmd, &migrationFlags)
}