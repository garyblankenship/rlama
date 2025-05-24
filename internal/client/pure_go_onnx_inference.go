package client

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"

	ort "github.com/yalue/onnxruntime_go"
)

// PureGoONNXInference provides pure Go ONNX inference for BGE reranking
type PureGoONNXInference struct {
	sessionMutex sync.RWMutex
	session      *ort.DynamicAdvancedSession
	modelPath    string
	initialized  bool
}

// ONNXInferenceRequest represents the input for ONNX inference
type ONNXInferenceRequest struct {
	InputIDs      [][]int64 `json:"input_ids"`
	AttentionMask [][]int64 `json:"attention_mask"`
}

// ONNXInferenceResponse represents the ONNX inference output
type ONNXInferenceResponse struct {
	Scores []float64 `json:"scores"`
}

// NewPureGoONNXInference creates a new pure Go ONNX inference client
func NewPureGoONNXInference(modelDir string) (*PureGoONNXInference, error) {
	modelPath := filepath.Join(modelDir, "model.onnx")
	
	// Set the ONNX runtime library path (absolute path from project root)
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}
	
	// Navigate to project root and construct absolute path
	projectRoot := filepath.Join(wd, "..", "..")
	libPath := filepath.Join(projectRoot, "lib", "onnxruntime-osx-arm64-1.19.0", "lib", "libonnxruntime.dylib")
	ort.SetSharedLibraryPath(libPath)
	
	client := &PureGoONNXInference{
		modelPath: modelPath,
	}
	
	// Initialize the ONNX environment and session
	if err := client.initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize ONNX inference: %w", err)
	}
	
	return client, nil
}

// initialize sets up the ONNX runtime environment and model session
func (o *PureGoONNXInference) initialize() error {
	o.sessionMutex.Lock()
	defer o.sessionMutex.Unlock()
	
	if o.initialized {
		return nil
	}
	
	// Initialize ONNX environment
	if err := ort.InitializeEnvironment(); err != nil {
		return fmt.Errorf("failed to initialize ONNX environment: %w", err)
	}
	
	// Create ONNX session with BGE model
	// Set the model directory as working directory for ONNX session to find model.onnx_data
	modelDir := filepath.Dir(o.modelPath)
	oldWd, _ := os.Getwd()
	os.Chdir(modelDir)
	defer os.Chdir(oldWd)
	
	session, err := ort.NewDynamicAdvancedSession(
		"model.onnx",                             // Use relative path from model directory
		[]string{"input_ids", "attention_mask"}, // Input names
		[]string{"logits"},                       // Output names
		nil,                                      // Use default session options
	)
	if err != nil {
		ort.DestroyEnvironment()
		return fmt.Errorf("failed to create ONNX session: %w", err)
	}
	
	o.session = session
	o.initialized = true
	
	return nil
}

// RunInference performs ONNX inference with tokenized inputs
func (o *PureGoONNXInference) RunInference(request ONNXInferenceRequest) (*ONNXInferenceResponse, error) {
	if !o.initialized {
		return nil, fmt.Errorf("ONNX inference not initialized")
	}
	
	o.sessionMutex.RLock()
	defer o.sessionMutex.RUnlock()
	
	batchSize := len(request.InputIDs)
	if batchSize == 0 {
		return &ONNXInferenceResponse{Scores: []float64{}}, nil
	}
	
	seqLength := len(request.InputIDs[0])
	scores := make([]float64, batchSize)
	
	// Process each item in the batch
	for i := 0; i < batchSize; i++ {
		score, err := o.inferenceForSingle(request.InputIDs[i], request.AttentionMask[i], seqLength)
		if err != nil {
			return nil, fmt.Errorf("inference failed for item %d: %w", i, err)
		}
		scores[i] = score
	}
	
	return &ONNXInferenceResponse{Scores: scores}, nil
}

// inferenceForSingle runs inference for a single query-document pair
func (o *PureGoONNXInference) inferenceForSingle(inputIDs, attentionMask []int64, seqLength int) (float64, error) {
	// Create input tensors
	inputShape := ort.NewShape(1, int64(seqLength))
	
	inputIDsTensor, err := ort.NewEmptyTensor[int64](inputShape)
	if err != nil {
		return 0, fmt.Errorf("failed to create input_ids tensor: %w", err)
	}
	defer inputIDsTensor.Destroy()
	
	attentionMaskTensor, err := ort.NewEmptyTensor[int64](inputShape)
	if err != nil {
		return 0, fmt.Errorf("failed to create attention_mask tensor: %w", err)
	}
	defer attentionMaskTensor.Destroy()
	
	// Create output tensor
	outputShape := ort.NewShape(1, 1)
	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		return 0, fmt.Errorf("failed to create output tensor: %w", err)
	}
	defer outputTensor.Destroy()
	
	// Fill input tensors with data
	inputIDsData := inputIDsTensor.GetData()
	attentionMaskData := attentionMaskTensor.GetData()
	
	copy(inputIDsData, inputIDs)
	copy(attentionMaskData, attentionMask)
	
	// Prepare inputs and outputs for inference
	inputs := []ort.ArbitraryTensor{inputIDsTensor, attentionMaskTensor}
	outputs := []ort.ArbitraryTensor{outputTensor}
	
	// Run inference
	if err := o.session.Run(inputs, outputs); err != nil {
		return 0, fmt.Errorf("ONNX inference failed: %w", err)
	}
	
	// Extract logits and convert to probability using sigmoid
	outputData := outputTensor.GetData()
	if len(outputData) == 0 {
		return 0, fmt.Errorf("no output data received")
	}
	
	logits := float64(outputData[0])
	
	// Apply sigmoid to convert logits to probability
	// sigmoid(x) = 1 / (1 + exp(-x))
	score := 1.0 / (1.0 + math.Exp(-logits))
	
	return score, nil
}

// Close cleans up the ONNX resources
func (o *PureGoONNXInference) Close() error {
	o.sessionMutex.Lock()
	defer o.sessionMutex.Unlock()
	
	if !o.initialized {
		return nil
	}
	
	if o.session != nil {
		o.session.Destroy()
		o.session = nil
	}
	
	if err := ort.DestroyEnvironment(); err != nil {
		return fmt.Errorf("failed to destroy ONNX environment: %w", err)
	}
	
	o.initialized = false
	return nil
}

// IsInitialized returns whether the ONNX inference is ready
func (o *PureGoONNXInference) IsInitialized() bool {
	o.sessionMutex.RLock()
	defer o.sessionMutex.RUnlock()
	return o.initialized
}