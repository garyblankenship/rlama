package cmd

import (
	"fmt"

	"github.com/dontizi/rlama/internal/service"
	"github.com/spf13/cobra"
)

var (
	// Flag holders for migrate-to-internal command
	internalMigrationFlags MigrationFlags
)

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

func init() {
	rootCmd.AddCommand(migrateToInternalCmd)

	// Add common migration control flags
	AddMigrationControlFlags(migrateToInternalCmd, &internalMigrationFlags)
}