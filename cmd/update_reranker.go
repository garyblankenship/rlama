package cmd

import (
	"fmt"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/service"
	"github.com/spf13/cobra"
)

var updateRerankerCmd = &cobra.Command{
	Use:   "update-reranker [rag-name]",
	Short: "Met à jour le modèle de reranking d'un RAG existant",
	Long:  `Configure un RAG existant pour utiliser le modèle BGE Reranker par défaut.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		updateReranker(args[0])
	},
}

func init() {
	rootCmd.AddCommand(updateRerankerCmd)
}

func updateReranker(ragName string) {
	// Charger le service RAG
	ollamaClient := client.NewDefaultOllamaClient()
	ragService := service.NewRagService(ollamaClient)

	// Mettre à jour le modèle de reranking
	err := ragService.UpdateRerankerModel(ragName, "BAAI/bge-reranker-v2-m3")
	if err != nil {
		fmt.Printf("Erreur lors de la mise à jour du modèle de reranking: %v\n", err)
		return
	}

	fmt.Printf("✅ RAG '%s' mis à jour pour utiliser le modèle BGE Reranker.\n", ragName)
}
