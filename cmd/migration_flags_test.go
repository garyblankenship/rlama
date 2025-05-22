package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestQdrantFlags(t *testing.T) {
	cmd := &cobra.Command{}
	var flags QdrantFlags

	// Add Qdrant flags
	AddQdrantFlags(cmd, &flags, "Test collection usage")

	// Test flag registration
	flagSet := cmd.Flags()
	
	// Check that all expected flags are registered
	expectedFlags := []string{
		"qdrant-host",
		"qdrant-port", 
		"qdrant-apikey",
		"qdrant-collection",
		"qdrant-grpc",
	}

	for _, flagName := range expectedFlags {
		flag := flagSet.Lookup(flagName)
		if flag == nil {
			t.Errorf("Expected flag %s not found", flagName)
		}
	}

	// Test default values by parsing empty args
	err := cmd.ParseFlags([]string{})
	if err != nil {
		t.Fatalf("Failed to parse empty flags: %v", err)
	}

	// Check default values
	host, port, apiKey, collection, useGRPC := GetQdrantFlagValues(&flags)
	
	if host != "localhost" {
		t.Errorf("Expected default host 'localhost', got '%s'", host)
	}
	if port != 6333 {
		t.Errorf("Expected default port 6333, got %d", port)
	}
	if apiKey != "" {
		t.Errorf("Expected empty default apiKey, got '%s'", apiKey)
	}
	if collection != "" {
		t.Errorf("Expected empty default collection, got '%s'", collection)
	}
	if useGRPC != false {
		t.Errorf("Expected default useGRPC false, got %v", useGRPC)
	}
}

func TestMigrationFlags(t *testing.T) {
	cmd := &cobra.Command{}
	var flags MigrationFlags

	// Add migration flags
	AddMigrationControlFlags(cmd, &flags)

	// Test flag registration
	flagSet := cmd.Flags()
	
	expectedFlags := []string{
		"backup",
		"backup-path",
		"verify",
		"delete-old",
	}

	for _, flagName := range expectedFlags {
		flag := flagSet.Lookup(flagName)
		if flag == nil {
			t.Errorf("Expected flag %s not found", flagName)
		}
	}

	// Test default values
	err := cmd.ParseFlags([]string{})
	if err != nil {
		t.Fatalf("Failed to parse empty flags: %v", err)
	}

	createBackup, backupPath, verify, deleteOld := GetMigrationFlagValues(&flags)
	
	if createBackup != false {
		t.Errorf("Expected default createBackup false, got %v", createBackup)
	}
	if backupPath != "" {
		t.Errorf("Expected empty default backupPath, got '%s'", backupPath)
	}
	if verify != true {
		t.Errorf("Expected default verify true, got %v", verify)
	}
	if deleteOld != false {
		t.Errorf("Expected default deleteOld false, got %v", deleteOld)
	}
}

func TestFlagValueParsing(t *testing.T) {
	cmd := &cobra.Command{}
	var qdrantFlags QdrantFlags
	var migrationFlags MigrationFlags

	AddQdrantFlags(cmd, &qdrantFlags, "Test collection")
	AddMigrationControlFlags(cmd, &migrationFlags)

	// Test parsing custom values
	testArgs := []string{
		"--qdrant-host=example.com",
		"--qdrant-port=9999",
		"--qdrant-apikey=test-key",
		"--qdrant-collection=test-collection",
		"--qdrant-grpc",
		"--backup",
		"--backup-path=/custom/path",
		"--verify=false",
		"--delete-old",
	}

	err := cmd.ParseFlags(testArgs)
	if err != nil {
		t.Fatalf("Failed to parse test flags: %v", err)
	}

	// Verify Qdrant flag values
	host, port, apiKey, collection, useGRPC := GetQdrantFlagValues(&qdrantFlags)
	if host != "example.com" {
		t.Errorf("Expected host 'example.com', got '%s'", host)
	}
	if port != 9999 {
		t.Errorf("Expected port 9999, got %d", port)
	}
	if apiKey != "test-key" {
		t.Errorf("Expected apiKey 'test-key', got '%s'", apiKey)
	}
	if collection != "test-collection" {
		t.Errorf("Expected collection 'test-collection', got '%s'", collection)
	}
	if !useGRPC {
		t.Errorf("Expected useGRPC true, got %v", useGRPC)
	}

	// Verify migration flag values
	createBackup, backupPath, verify, deleteOld := GetMigrationFlagValues(&migrationFlags)
	if !createBackup {
		t.Errorf("Expected createBackup true, got %v", createBackup)
	}
	if backupPath != "/custom/path" {
		t.Errorf("Expected backupPath '/custom/path', got '%s'", backupPath)
	}
	if verify {
		t.Errorf("Expected verify false, got %v", verify)
	}
	if !deleteOld {
		t.Errorf("Expected deleteOld true, got %v", deleteOld)
	}
}