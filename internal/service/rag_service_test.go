package service

import (
	"testing"

	"github.com/dontizi/rlama/internal/domain"
	"github.com/stretchr/testify/assert"
)

// TestRagRerankerTopK vérifie que le reranking est configuré correctement et limite les résultats à 5 par défaut
func TestRagRerankerTopK(t *testing.T) {
	// Créer un nouveau RAG
	rag := domain.NewRagSystem("test-rerank-topk", "test-model")

	// Vérifier que les valeurs par défaut du reranking sont correctes
	assert.True(t, rag.RerankerEnabled, "Le reranking devrait être activé par défaut")
	assert.Equal(t, float64(0.7), rag.RerankerWeight, "Le poids du reranker devrait être 0.7 par défaut")
	assert.Equal(t, "test-model", rag.RerankerModel, "Le modèle du reranker devrait être le même que le RAG par défaut")
	assert.Equal(t, 5, rag.RerankerTopK, "TopK devrait être 5 par défaut")

	// Vérifier que les valeurs par défaut des options de reranking sont cohérentes
	options := DefaultRerankerOptions()
	assert.Equal(t, options.TopK, rag.RerankerTopK, "TopK dans le RAG et dans les options devraient être identiques")

	// Tester avec différentes valeurs de TopK
	testCases := []struct {
		name     string
		topK     int
		expected int
	}{
		{
			name:     "DefaultTopK",
			topK:     0, // 0 signifie utiliser la valeur par défaut
			expected: 5,
		},
		{
			name:     "CustomTopK",
			topK:     10,
			expected: 10,
		},
		{
			name:     "ZeroTopK",
			topK:     -1, // Valeur invalide, devrait utiliser le défaut du RAG
			expected: 5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Simuler la logique de Query() pour déterminer la taille du contexte
			contextSize := tc.topK

			// Si contextSize est 0 (auto), utiliser:
			// - RerankerTopK du RAG si défini
			// - Sinon le TopK par défaut (5)
			// - 20 si le reranking est désactivé
			if contextSize <= 0 {
				if rag.RerankerEnabled {
					if rag.RerankerTopK > 0 {
						contextSize = rag.RerankerTopK
					} else {
						contextSize = options.TopK // 5 par défaut
					}
				} else {
					contextSize = 20 // 20 par défaut si le reranking est désactivé
				}
			}

			// Vérifier que contextSize correspond à la valeur attendue
			assert.Equal(t, tc.expected, contextSize,
				"La taille du contexte devrait correspondre à la valeur attendue")
		})
	}

	// Tester le cas où le reranking est désactivé
	t.Run("DisabledReranking", func(t *testing.T) {
		rag.RerankerEnabled = false

		// Contextsize à 0 devrait donner 20 car le reranking est désactivé
		contextSize := 0
		if contextSize <= 0 {
			if rag.RerankerEnabled {
				if rag.RerankerTopK > 0 {
					contextSize = rag.RerankerTopK
				} else {
					contextSize = options.TopK // 5 par défaut
				}
			} else {
				contextSize = 20 // 20 par défaut si le reranking est désactivé
			}
		}

		assert.Equal(t, 20, contextSize,
			"La taille du contexte devrait être 20 par défaut si le reranking est désactivé")

		// Restaurer l'état
		rag.RerankerEnabled = true
	})
}
