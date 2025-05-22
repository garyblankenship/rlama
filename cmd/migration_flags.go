package cmd

import "github.com/spf13/cobra"

// QdrantFlags holds common Qdrant connection configuration flags
type QdrantFlags struct {
	Host       *string
	Port       *int
	APIKey     *string
	Collection *string
	UseGRPC    *bool
}

// MigrationFlags holds common migration control flags
type MigrationFlags struct {
	CreateBackup         *bool
	BackupPath           *string
	VerifyAfterMigration *bool
	DeleteOldData        *bool
}

// AddQdrantFlags adds standard Qdrant connection flags to a cobra command
func AddQdrantFlags(cmd *cobra.Command, flags *QdrantFlags, collectionUsage string) {
	flags.Host = cmd.Flags().String("qdrant-host", "localhost", "Qdrant server host")
	flags.Port = cmd.Flags().Int("qdrant-port", 6333, "Qdrant server port")
	flags.APIKey = cmd.Flags().String("qdrant-apikey", "", "Qdrant API key for secured instances")
	flags.Collection = cmd.Flags().String("qdrant-collection", "", collectionUsage)
	flags.UseGRPC = cmd.Flags().Bool("qdrant-grpc", false, "Use gRPC for Qdrant communication")
}

// AddMigrationControlFlags adds standard migration control flags to a cobra command
func AddMigrationControlFlags(cmd *cobra.Command, flags *MigrationFlags) {
	flags.CreateBackup = cmd.Flags().Bool("backup", false, "Create backup before migration")
	flags.BackupPath = cmd.Flags().String("backup-path", "", "Custom backup path (default: ~/.rlama/backups)")
	flags.VerifyAfterMigration = cmd.Flags().Bool("verify", true, "Verify migration integrity after completion")
	flags.DeleteOldData = cmd.Flags().Bool("delete-old", false, "Delete old data after successful migration")
}

// AddBatchMigrationFlags adds flags specific to batch migration operations
func AddBatchMigrationFlags(cmd *cobra.Command, ragNames *[]string, continueOnError *bool) {
	cmd.Flags().StringSliceVar(ragNames, "rags", []string{}, "Specific RAG names to migrate (comma-separated, default: all RAGs)")
	cmd.Flags().BoolVar(continueOnError, "continue-on-error", false, "Continue batch migration even if individual RAGs fail")
}

// GetQdrantFlagValues returns the actual values from QdrantFlags pointers
func GetQdrantFlagValues(flags *QdrantFlags) (host string, port int, apiKey string, collection string, useGRPC bool) {
	if flags.Host != nil {
		host = *flags.Host
	}
	if flags.Port != nil {
		port = *flags.Port
	}
	if flags.APIKey != nil {
		apiKey = *flags.APIKey
	}
	if flags.Collection != nil {
		collection = *flags.Collection
	}
	if flags.UseGRPC != nil {
		useGRPC = *flags.UseGRPC
	}
	return
}

// GetMigrationFlagValues returns the actual values from MigrationFlags pointers
func GetMigrationFlagValues(flags *MigrationFlags) (createBackup bool, backupPath string, verify bool, deleteOld bool) {
	if flags.CreateBackup != nil {
		createBackup = *flags.CreateBackup
	}
	if flags.BackupPath != nil {
		backupPath = *flags.BackupPath
	}
	if flags.VerifyAfterMigration != nil {
		verify = *flags.VerifyAfterMigration
	}
	if flags.DeleteOldData != nil {
		deleteOld = *flags.DeleteOldData
	}
	return
}