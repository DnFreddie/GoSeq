package quizz

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// FlashcardsCmd is the root command for flashcard functionality
var FlashcardsCmd = &cobra.Command{
	Use:   "flashcards",
	Short: "Manage and study with flashcards",
	Long:  `Create, edit, and study with flashcards to help memorize information.`,
}

func init() {
	app, err := NewApp(os.Stdout, os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing app: %v\n", err)
		return
	}

	// Add subcommand
	addCmd := &cobra.Command{
		Use:   "add [flashcards...]",
		Short: "Add new flashcards",
		Long: `Add new flashcards to your collection.
Each flashcard should be in the format "question:answer".`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.AddFlashcards(cmd.Context(), args)
		},
		Example: `  flashcards add "What is the capital of France?:Paris"
  flashcards add "Word:Definition" "Question:Answer"`,
	}

	// Edit subcommand
	editCmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit flashcards in your default editor",
		Long: `Open the flashcards file in your default editor.
Uses $EDITOR environment variable, defaults to 'vi'.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.EditFlashcards(cmd.Context())
		},
	}

	// Quiz subcommand
	quizCmd := &cobra.Command{
		Use:   "quiz",
		Short: "Start a flashcard quiz",
		Long: `Start an interactive quiz session.
Use '?' to see the complete answer
Use '>' to skip the current card`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.StartQuiz(cmd.Context())
		},
	}

	// Add all subcommands to FlashcardsCmd
	FlashcardsCmd.AddCommand(addCmd, editCmd, quizCmd)
}

