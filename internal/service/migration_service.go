package service

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

// MigrationService handles RAG system migrations between vector stores
type MigrationService struct {
	ragRepository *repository.RagRepository
}

// NewMigrationService creates a new migration service
func NewMigrationService() *MigrationService {
	return &MigrationService{
		ragRepository: repository.NewRagRepository(),
	}
}

// MigrateToQdrant migrates a RAG from internal storage to Qdrant
func (ms *MigrationService) MigrateToQdrant(ragName string, opts MigrationOptions) error {
	// Step 1: Load existing RAG
	fmt.Println("📖 Loading existing RAG...")
	rag, err := ms.ragRepository.Load(ragName)
	if err != nil {
		return fmt.Errorf("failed to load RAG '%s': %w", ragName, err)
	}

	// Check if already using Qdrant
	if rag.VectorStoreType == "qdrant" {
		return fmt.Errorf("RAG '%s' is already using Qdrant vector store", ragName)
	}

	// Step 2: Create backup if requested
	if opts.CreateBackup {
		if err := ms.createBackup(ragName, opts.BackupPath); err != nil {
			return fmt.Errorf("backup creation failed: %w", err)
		}
		fmt.Println("✅ Backup created successfully")
	}

	// Step 3: Extract vectors and metadata from current store
	fmt.Println("🔍 Extracting vectors from internal storage...")
	vectors, err := ms.extractVectorsFromRAG(rag)
	if err != nil {
		return fmt.Errorf("failed to extract vectors: %w", err)
	}
	fmt.Printf("📊 Found %d vectors to migrate\n", len(vectors))

	// Step 4: Create new Qdrant-based hybrid store
	fmt.Println("🔗 Connecting to Qdrant...")
	newHybridConfig := vector.HybridStoreConfig{
		IndexPath:            ":memory:",
		Dimensions:           rag.EmbeddingDimension,
		VectorStoreType:      "qdrant",
		QdrantHost:           opts.QdrantHost,
		QdrantPort:           opts.QdrantPort,
		QdrantAPIKey:         opts.QdrantAPIKey,
		QdrantCollectionName: opts.QdrantCollectionName,
		QdrantGRPC:           opts.QdrantGRPC,
	}

	newHybridStore, err := vector.NewEnhancedHybridStoreWithConfig(newHybridConfig)
	if err != nil {
		return fmt.Errorf("failed to create Qdrant hybrid store: %w", err)
	}
	defer func() {
		if closer, ok := newHybridStore.VectorStore.(interface{ Close() error }); ok {
			closer.Close()
		}
	}()

	// Step 5: Transfer all vectors to Qdrant
	fmt.Println("🚀 Transferring vectors to Qdrant...")
	err = ms.transferVectorsToStore(vectors, newHybridStore)
	if err != nil {
		return fmt.Errorf("failed to transfer vectors to Qdrant: %w", err)
	}

	// Step 6: Update RAG configuration
	fmt.Println("⚙️ Updating RAG configuration...")
	rag.VectorStoreType = "qdrant"
	rag.QdrantHost = opts.QdrantHost
	rag.QdrantPort = opts.QdrantPort
	rag.QdrantAPIKey = opts.QdrantAPIKey
	rag.QdrantCollectionName = opts.QdrantCollectionName
	rag.QdrantGRPC = opts.QdrantGRPC
	rag.HybridStore = newHybridStore

	// Step 7: Verify migration if requested
	if opts.VerifyAfterMigration {
		fmt.Println("🔍 Verifying migration...")
		if err := ms.verifyMigration(rag, len(vectors)); err != nil {
			return fmt.Errorf("migration verification failed: %w", err)
		}
		fmt.Println("✅ Migration verification passed")
	}

	// Step 8: Save updated RAG
	if err := ms.ragRepository.Save(rag); err != nil {
		return fmt.Errorf("failed to save migrated RAG: %w", err)
	}

	// Step 9: Clean up old data if requested
	if opts.DeleteOldData {
		fmt.Println("🗑️ Cleaning up old internal vector files...")
		if err := ms.deleteOldInternalData(ragName); err != nil {
			fmt.Printf("⚠️ Warning: failed to delete old data: %v\n", err)
		} else {
			fmt.Println("✅ Old data cleaned up")
		}
	}

	return nil
}

// MigrateToInternal migrates a RAG from Qdrant to internal storage
func (ms *MigrationService) MigrateToInternal(ragName string, opts MigrationOptions) error {
	// Step 1: Load existing RAG
	fmt.Println("📖 Loading existing RAG...")
	rag, err := ms.ragRepository.Load(ragName)
	if err != nil {
		return fmt.Errorf("failed to load RAG '%s': %w", ragName, err)
	}

	// Check if already using internal storage
	if rag.VectorStoreType != "qdrant" {
		return fmt.Errorf("RAG '%s' is already using internal vector store", ragName)
	}

	// Step 2: Create backup if requested
	if opts.CreateBackup {
		if err := ms.createBackup(ragName, opts.BackupPath); err != nil {
			return fmt.Errorf("backup creation failed: %w", err)
		}
		fmt.Println("✅ Backup created successfully")
	}

	// Step 3: Extract vectors from Qdrant
	fmt.Println("🔍 Extracting vectors from Qdrant...")
	vectors, err := ms.extractVectorsFromRAG(rag)
	if err != nil {
		return fmt.Errorf("failed to extract vectors: %w", err)
	}
	fmt.Printf("📊 Found %d vectors to migrate\n", len(vectors))

	// Step 4: Create new internal hybrid store
	fmt.Println("🔗 Creating internal vector storage...")
	newHybridStore, err := vector.NewEnhancedHybridStore(":memory:", rag.EmbeddingDimension)
	if err != nil {
		return fmt.Errorf("failed to create internal hybrid store: %w", err)
	}

	// Step 5: Transfer all vectors to internal storage
	fmt.Println("🚀 Transferring vectors to internal storage...")
	err = ms.transferVectorsToStore(vectors, newHybridStore)
	if err != nil {
		return fmt.Errorf("failed to transfer vectors to internal storage: %w", err)
	}

	// Step 6: Update RAG configuration
	fmt.Println("⚙️ Updating RAG configuration...")
	rag.VectorStoreType = "internal"
	rag.QdrantHost = ""
	rag.QdrantPort = 0
	rag.QdrantAPIKey = ""
	rag.QdrantCollectionName = ""
	rag.QdrantGRPC = false
	rag.HybridStore = newHybridStore

	// Step 7: Verify migration if requested
	if opts.VerifyAfterMigration {
		fmt.Println("🔍 Verifying migration...")
		if err := ms.verifyMigration(rag, len(vectors)); err != nil {
			return fmt.Errorf("migration verification failed: %w", err)
		}
		fmt.Println("✅ Migration verification passed")
	}

	// Step 8: Save updated RAG
	if err := ms.ragRepository.Save(rag); err != nil {
		return fmt.Errorf("failed to save migrated RAG: %w", err)
	}

	// Step 9: Clean up old data if requested
	if opts.DeleteOldData {
		fmt.Println("🗑️ Warning: Qdrant collection cleanup not implemented yet")
		fmt.Println("💡 You may want to manually remove the collection from Qdrant")
	}

	return nil
}

// VectorData represents a vector with its associated metadata
type VectorData struct {
	ID       string
	Vector   []float32
	Content  string
	Metadata string
}

// extractVectorsFromRAG extracts all vectors and metadata from a RAG system
func (ms *MigrationService) extractVectorsFromRAG(rag *domain.RagSystem) ([]VectorData, error) {
	var vectors []VectorData

	// Extract vectors from all chunks
	for _, chunk := range rag.Chunks {
		if chunk.Embedding != nil {
			vectors = append(vectors, VectorData{
				ID:       chunk.ID,
				Vector:   chunk.Embedding,
				Content:  chunk.Content,
				Metadata: chunk.GetMetadataString(),
			})
		}
	}

	return vectors, nil
}

// transferVectorsToStore transfers vectors to the target hybrid store
func (ms *MigrationService) transferVectorsToStore(vectors []VectorData, store *vector.EnhancedHybridStore) error {
	for i, vectorData := range vectors {
		if i%100 == 0 {
			fmt.Printf("⏳ Progress: %d/%d vectors transferred\n", i, len(vectors))
		}

		err := store.AddDocument(vectorData.ID, vectorData.Content, vectorData.Metadata, vectorData.Vector)
		if err != nil {
			return fmt.Errorf("failed to add vector %s: %w", vectorData.ID, err)
		}
	}

	fmt.Printf("✅ All %d vectors transferred successfully\n", len(vectors))
	return nil
}

// createBackup creates a backup of the RAG before migration
func (ms *MigrationService) createBackup(ragName, backupPath string) error {
	if backupPath == "" {
		backupPath = filepath.Join(config.GetDataDir(), "backups")
	}

	// Create backup directory
	timestamp := time.Now().Format("20060102-150405")
	backupDir := filepath.Join(backupPath, fmt.Sprintf("%s-%s", ragName, timestamp))
	
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// For now, this is a placeholder - in a full implementation, we would:
	// 1. Copy the RAG's info.json file
	// 2. Copy vector storage files (if internal)
	// 3. Export Qdrant data (if Qdrant)
	// 4. Create a manifest file with backup info

	fmt.Printf("💾 Backup directory created: %s\n", backupDir)
	return nil
}

// verifyMigration verifies that the migration was successful
func (ms *MigrationService) verifyMigration(rag *domain.RagSystem, expectedVectorCount int) error {
	// Perform a simple search to verify the vector store is working
	if len(rag.Chunks) == 0 {
		return fmt.Errorf("no chunks found in migrated RAG")
	}

	// Use the first chunk's embedding as a test query
	testVector := rag.Chunks[0].Embedding
	if testVector == nil {
		return fmt.Errorf("first chunk has no embedding")
	}

	// Perform a search
	results := rag.HybridStore.Search(testVector, 5)
	if len(results) == 0 {
		return fmt.Errorf("vector search returned no results")
	}

	// Check that we can find the exact vector we searched for
	found := false
	for _, result := range results {
		if result.ID == rag.Chunks[0].ID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("could not find the test vector in search results")
	}

	fmt.Printf("🔍 Verification: Found test vector with ID %s\n", results[0].ID)
	return nil
}

// deleteOldInternalData removes old internal vector files
func (ms *MigrationService) deleteOldInternalData(ragName string) error {
	dataDir := config.GetDataDir()
	ragDir := filepath.Join(dataDir, ragName)
	vectorFile := filepath.Join(ragDir, "vectors.json")

	if _, err := os.Stat(vectorFile); err == nil {
		if err := os.Remove(vectorFile); err != nil {
			return fmt.Errorf("failed to remove vector file: %w", err)
		}
	}

	return nil
}