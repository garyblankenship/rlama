package vector

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

// QdrantStore implements VectorStoreInterface using Qdrant vector database
type QdrantStore struct {
	client         qdrant.PointsClient
	collections    qdrant.CollectionsClient
	conn           *grpc.ClientConn
	collectionName string
	dims           uint64
	timeout        time.Duration
}

// Ensure QdrantStore implements VectorStoreInterface
var _ VectorStoreInterface = (*QdrantStore)(nil)

// NewQdrantStore creates and configures a new Qdrant client and store
func NewQdrantStore(host string, port int, collectionName string, dims int, apiKey string, useGRPC bool, createCollectionIfNotExists bool) (*QdrantStore, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var conn *grpc.ClientConn
	var err error

	// Setup gRPC connection options
	var dialOpts []grpc.DialOption
	if apiKey != "" {
		// For Qdrant Cloud or secured instances, typically use TLS
		config := &tls.Config{}
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(config)))
		fmt.Printf("Warning: API key provided - using TLS. Ensure proper authentication is configured.\n")
	} else {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err = grpc.DialContext(ctxTimeout, addr, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Qdrant at %s: %w", addr, err)
	}

	pointsClient := qdrant.NewPointsClient(conn)
	collectionsClient := qdrant.NewCollectionsClient(conn)

	store := &QdrantStore{
		client:         pointsClient,
		collections:    collectionsClient,
		conn:           conn,
		collectionName: collectionName,
		dims:           uint64(dims),
		timeout:        10 * time.Second,
	}

	if createCollectionIfNotExists {
		err := store.ensureCollectionExists(ctxTimeout)
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("failed to ensure Qdrant collection '%s': %w", collectionName, err)
		}
	}

	return store, nil
}

// ensureCollectionExists creates the collection if it doesn't exist
func (s *QdrantStore) ensureCollectionExists(ctx context.Context) error {
	listResp, err := s.collections.List(ctx, &qdrant.ListCollectionsRequest{})
	if err != nil {
		return fmt.Errorf("failed to list Qdrant collections: %w", err)
	}

	for _, coll := range listResp.GetCollections() {
		if coll.GetName() == s.collectionName {
			fmt.Printf("Qdrant collection '%s' already exists.\n", s.collectionName)
			return nil
		}
	}

	fmt.Printf("Creating Qdrant collection '%s' with %d dimensions...\n", s.collectionName, s.dims)
	_, err = s.collections.Create(ctx, &qdrant.CreateCollection{
		CollectionName: s.collectionName,
		VectorsConfig: &qdrant.VectorsConfig{
			Config: &qdrant.VectorsConfig_Params{
				Params: &qdrant.VectorParams{
					Size:     s.dims,
					Distance: qdrant.Distance_Cosine,
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create Qdrant collection '%s': %w", s.collectionName, err)
	}

	fmt.Printf("Qdrant collection '%s' created successfully.\n", s.collectionName)
	return nil
}

// Add implements VectorStoreInterface - adds a vector without payload
func (s *QdrantStore) Add(id string, vector []float32) {
	err := s.UpsertPointWithPayload(id, vector, nil)
	if err != nil {
		fmt.Printf("Error adding point %s to Qdrant: %v\n", id, err)
	}
}

// UpsertPointWithPayload adds or updates a point with its vector and payload
func (s *QdrantStore) UpsertPointWithPayload(id string, vector []float32, payload map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	qdrantPayload := make(map[string]*qdrant.Value)
	if payload != nil {
		for k, v := range payload {
			switch val := v.(type) {
			case string:
				qdrantPayload[k] = &qdrant.Value{Kind: &qdrant.Value_StringValue{StringValue: val}}
			case int:
				qdrantPayload[k] = &qdrant.Value{Kind: &qdrant.Value_IntegerValue{IntegerValue: int64(val)}}
			case int64:
				qdrantPayload[k] = &qdrant.Value{Kind: &qdrant.Value_IntegerValue{IntegerValue: val}}
			case float64:
				qdrantPayload[k] = &qdrant.Value{Kind: &qdrant.Value_DoubleValue{DoubleValue: val}}
			case bool:
				qdrantPayload[k] = &qdrant.Value{Kind: &qdrant.Value_BoolValue{BoolValue: val}}
			default:
				qdrantPayload[k] = &qdrant.Value{Kind: &qdrant.Value_StringValue{StringValue: fmt.Sprintf("%v", v)}}
			}
		}
	}

	points := []*qdrant.PointStruct{
		{
			Id: &qdrant.PointId{
				PointIdOptions: &qdrant.PointId_Uuid{Uuid: id},
			},
			Vectors: &qdrant.Vectors{
				VectorsOptions: &qdrant.Vectors_Vector{
					Vector: &qdrant.Vector{Data: vector},
				},
			},
			Payload: qdrantPayload,
		},
	}

	_, err := s.client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: s.collectionName,
		Points:         points,
		Wait:           proto.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("failed to upsert point %s to Qdrant: %w", id, err)
	}
	return nil
}

// Search performs a vector search in Qdrant
func (s *QdrantStore) Search(query []float32, limit int) []SearchResult {
	return s.SearchWithFilter(query, limit, nil)
}

// SearchWithFilter performs a vector search with optional payload filtering
func (s *QdrantStore) SearchWithFilter(query []float32, limit int, filter *qdrant.Filter) []SearchResult {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	searchRequest := &qdrant.SearchPoints{
		CollectionName: s.collectionName,
		Vector:         query,
		Limit:          uint64(limit),
		WithPayload: &qdrant.WithPayloadSelector{
			SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true},
		},
		WithVectors: &qdrant.WithVectorsSelector{
			SelectorOptions: &qdrant.WithVectorsSelector_Enable{Enable: false},
		},
		Filter: filter,
	}

	searchResults, err := s.client.Search(ctx, searchRequest)
	if err != nil {
		fmt.Printf("Error searching Qdrant collection '%s': %v\n", s.collectionName, err)
		return []SearchResult{}
	}

	results := make([]SearchResult, 0, len(searchResults.GetResult()))
	for _, hit := range searchResults.GetResult() {
		var originalID string
		if hit.GetId().GetUuid() != "" {
			originalID = hit.GetId().GetUuid()
		} else {
			fmt.Printf("Warning: Qdrant hit ID is not UUID: %v\n", hit.GetId())
			continue
		}

		results = append(results, SearchResult{
			ID:    originalID,
			Score: float64(hit.GetScore()),
		})
	}
	return results
}

// Remove deletes a point from Qdrant
func (s *QdrantStore) Remove(id string) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	_, err := s.client.Delete(ctx, &qdrant.DeletePoints{
		CollectionName: s.collectionName,
		Points: &qdrant.PointsSelector{
			PointsSelectorOneOf: &qdrant.PointsSelector_Points{
				Points: &qdrant.PointsIdsList{
					Ids: []*qdrant.PointId{
						{PointIdOptions: &qdrant.PointId_Uuid{Uuid: id}},
					},
				},
			},
		},
		Wait: proto.Bool(true),
	})
	if err != nil {
		fmt.Printf("Error removing point %s from Qdrant: %v\n", id, err)
	}
}

// Save is a no-op for QdrantStore as Qdrant server handles persistence
func (s *QdrantStore) Save(path string) error {
	return nil
}

// Load is a no-op for QdrantStore as connection is established at construction
func (s *QdrantStore) Load(path string) error {
	return nil
}

// Close closes the gRPC connection to Qdrant
func (s *QdrantStore) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}