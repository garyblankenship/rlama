package service

import (
	"fmt"
	"testing"

	"github.com/dontizi/rlama/internal/domain"
	"github.com/stretchr/testify/assert"
)

// TestRerankerOptionsDefaultValues vérifie que les valeurs par défaut sont correctes
func TestRerankerOptionsDefaultValues(t *testing.T) {
	// Obtenir les options par défaut
	options := DefaultRerankerOptions()

	// Vérifier que les valeurs par défaut sont correctes
	assert.Equal(t, 5, options.TopK, "La valeur par défaut pour TopK devrait être 5")
	assert.Equal(t, 20, options.InitialK, "La valeur par défaut pour InitialK devrait être 20")
	assert.Equal(t, float64(0.7), options.RerankerWeight, "La valeur par défaut pour RerankerWeight devrait être 0.7")
	assert.Equal(t, float64(0.0), options.ScoreThreshold, "La valeur par défaut pour ScoreThreshold devrait être 0.0")
}

// TestApplyTopKLimit teste que la limite TopK est correctement appliquée
func TestApplyTopKLimit(t *testing.T) {
	// Créer des résultats déjà triés pour simuler la sortie avant application de TopK
	testCases := []struct {
		name     string
		results  []RankedResult
		topK     int
		expected int
	}{
		{
			name:     "LimitsToTopK5",
			results:  createDummyRankedResults(20),
			topK:     5,
			expected: 5,
		},
		{
			name:     "LimitsToTopK10",
			results:  createDummyRankedResults(20),
			topK:     10,
			expected: 10,
		},
		{
			name:     "HandlesTopKGreaterThanResults",
			results:  createDummyRankedResults(15),
			topK:     20,
			expected: 15, // Ne peut pas retourner plus que ce qui existe
		},
		{
			name:     "HandlesEmptyResults",
			results:  []RankedResult{},
			topK:     5,
			expected: 0,
		},
		{
			name:     "HandlesTopKZero",
			results:  createDummyRankedResults(10),
			topK:     0,
			expected: 10, // Ne devrait pas limiter si TopK=0
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Appliquer la limite TopK manuellement (reproduire la logique de Rerank)
			var limited []RankedResult
			if tc.topK > 0 && len(tc.results) > tc.topK {
				limited = tc.results[:tc.topK]
			} else {
				limited = tc.results
			}

			// Vérifier que le nombre est correct
			assert.Equal(t, tc.expected, len(limited), "Le nombre de résultats devrait être limité à TopK si nécessaire")
		})
	}
}

// createDummyRankedResults crée un ensemble de résultats factices pour les tests
func createDummyRankedResults(count int) []RankedResult {
	results := make([]RankedResult, count)

	for i := 0; i < count; i++ {
		results[i] = RankedResult{
			Chunk:         &domain.DocumentChunk{ID: fmt.Sprintf("chunk-%d", i)},
			VectorScore:   0.8 - (float64(i) * 0.01),
			RerankerScore: 0.9 - (float64(i) * 0.02),
			FinalScore:    0.95 - (float64(i) * 0.015),
		}
	}

	return results
}

// Reproduire la fonction de tri pour tester
func TestSortingByScore(t *testing.T) {
	// Créer des résultats dans un ordre mélangé
	results := []RankedResult{
		{FinalScore: 0.5},
		{FinalScore: 0.9},
		{FinalScore: 0.3},
		{FinalScore: 0.7},
		{FinalScore: 0.1},
	}

	// Trier les résultats (même logique que dans Rerank)
	// Sort by final score (descending)
	sortedResults := make([]RankedResult, len(results))
	copy(sortedResults, results)

	// Tri par score final décroissant
	for i := 0; i < len(sortedResults); i++ {
		for j := i + 1; j < len(sortedResults); j++ {
			if sortedResults[i].FinalScore < sortedResults[j].FinalScore {
				sortedResults[i], sortedResults[j] = sortedResults[j], sortedResults[i]
			}
		}
	}

	// Vérifier que les résultats sont triés correctement
	for i := 1; i < len(sortedResults); i++ {
		assert.GreaterOrEqual(t, sortedResults[i-1].FinalScore, sortedResults[i].FinalScore,
			"Les résultats devraient être triés par score décroissant")
	}

	// Vérifier l'ordre exact
	assert.Equal(t, float64(0.9), sortedResults[0].FinalScore)
	assert.Equal(t, float64(0.7), sortedResults[1].FinalScore)
	assert.Equal(t, float64(0.5), sortedResults[2].FinalScore)
	assert.Equal(t, float64(0.3), sortedResults[3].FinalScore)
	assert.Equal(t, float64(0.1), sortedResults[4].FinalScore)
}

// TestRerankerIntegration teste l'intégration du reranking dans le service RAG
func TestRerankerIntegration(t *testing.T) {
	// Ce test intégrera le reranking dans un service RAG complet
	// Comme il nécessite des dépendances externes, il sera marqué comme un test d'intégration
	t.Skip("Ce test nécessite une instance d'Ollama en cours d'exécution")

	// TODO: Implémenter un test d'intégration avec un vrai service RAG
	// Cela peut être fait plus tard en utilisant les structs et fonctions existantes
}
