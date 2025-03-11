package domain

import (
	"fmt"
	"time"
)

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

// NewDocumentChunk creates a new chunk from a document
func NewDocumentChunk(doc *Document, content string, startPos, endPos, chunkIndex int) *DocumentChunk {
	// Generate a unique ID for the chunk
	chunkID := fmt.Sprintf("%s_chunk_%d", doc.ID, chunkIndex)
	
	// Create metadata for the chunk
	metadata := map[string]string{
		"document_name": doc.Name,
		"document_path": doc.Path,
		"content_type": doc.ContentType,
		"chunk_position": fmt.Sprintf("%d of %d", chunkIndex+1, 0), // Total will be updated later
	}
	
	return &DocumentChunk{
		ID:          chunkID,
		DocumentID:  doc.ID,
		Content:     content,
		StartPos:    startPos,
		EndPos:      endPos,
		ChunkIndex:  chunkIndex,
		CreatedAt:   time.Now(),
		Metadata:    metadata,
	}
}

// GetMetadataString returns a formatted string of the chunk's metadata
func (c *DocumentChunk) GetMetadataString() string {
	return fmt.Sprintf("Source: %s (Section %s)", 
		c.Metadata["document_name"], 
		c.Metadata["chunk_position"])
}

// UpdateTotalChunks updates the chunk position metadata with the total chunk count
func (c *DocumentChunk) UpdateTotalChunks(total int) {
	c.Metadata["chunk_position"] = fmt.Sprintf("%d of %d", c.ChunkIndex+1, total)
} 