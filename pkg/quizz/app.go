package quizz

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Config struct {
	FlashcardsDir  string
	FlashcardsFile string
	Editor         string
	Colors         map[string]string
}

type Flashcard struct {
	Question string
	Answer   string
}

type App struct {
	config Config
	writer io.Writer
	reader io.Reader
}

const (
	ColorOrange = "\033[38;5;214m"
	ColorGreen  = "\033[38;5;114m"
	ColorRed    = "\033[38;5;203m"
	ColorBlue   = "\033[38;5;110m"
	ColorYellow = "\033[38;5;226m"
	ColorReset  = "\033[0m"
)

var (
	ErrNoFlashcards      = errors.New("no flashcards found")
	ErrInvalidFormat     = errors.New("invalid flashcard format")
	ErrInvalidEditor     = errors.New("invalid editor configuration")
	ErrDirectoryCreation = errors.New("failed to create flashcards directory")
)

func NewApp(w io.Writer, r io.Reader) (*App, error) {
	config, err := newConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create config: %w", err)
	}

	return &App{
		config: config,
		writer: w,
		reader: r,
	}, nil
}

func newConfig() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("failed to get home directory: %w", err)
	}

	flashcardsDir := filepath.Join(home, "github.com/DnFreddie/Notes/content/flashcards/")
	return Config{
		FlashcardsDir:  flashcardsDir,
		FlashcardsFile: filepath.Join(flashcardsDir, "german.md"),
		Editor:         os.Getenv("EDITOR"),
		Colors: map[string]string{
			"orange": ColorOrange,
			"green":  ColorGreen,
			"red":    ColorRed,
			"blue":   ColorBlue,
			"yellow": ColorYellow,
			"reset":  ColorReset,
		},
	}, nil
}

func (a *App) ensureDirectory() error {
	if err := os.MkdirAll(a.config.FlashcardsDir, 0755); err != nil {
		return fmt.Errorf("%w: %v", ErrDirectoryCreation, err)
	}
	return nil
}

func (a *App) colorize(color, text string) string {
	return fmt.Sprintf("%s%s%s", a.config.Colors[color], text, a.config.Colors["reset"])
}

func (a *App) AddFlashcards(ctx context.Context, cards []string) error {
	if err := a.ensureDirectory(); err != nil {
		return err
	}

	f, err := os.OpenFile(a.config.FlashcardsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open flashcards file: %w", err)
	}
	defer f.Close()

	for _, card := range cards {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("operation cancelled: %w", err)
		}

		question, answer, err := parseFlashcard(card)
		if err != nil {
			fmt.Fprintln(a.writer, a.colorize("red", fmt.Sprintf("Invalid flashcard: %s - %v", card, err)))
			continue
		}

		line := fmt.Sprintf("%s : %s\n", question, answer)
		if _, err := f.WriteString(line); err != nil {
			return fmt.Errorf("failed to write flashcard: %w", err)
		}

		fmt.Fprintln(a.writer, a.colorize("green", fmt.Sprintf("Added: %s", line)))
	}
	return nil
}

func parseFlashcard(card string) (question, answer string, err error) {
	parts := strings.Split(card, ":")
	if len(parts) != 2 {
		return "", "", ErrInvalidFormat
	}
	question = strings.TrimSpace(parts[0])
	answer = strings.TrimSpace(parts[1])
	if question == "" || answer == "" {
		return "", "", ErrInvalidFormat
	}
	return question, answer, nil
}

func (a *App) EditFlashcards(ctx context.Context) error {
	if err := a.ensureDirectory(); err != nil {
		return err
	}

	editor := a.config.Editor
	if editor == "" {
		editor = "vi"
	}

	if !isValidEditor(editor) {
		return ErrInvalidEditor
	}

	cmd := exec.CommandContext(ctx, editor, a.config.FlashcardsFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func isValidEditor(editor string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(editor)
}

func (a *App) LoadFlashcards() ([]Flashcard, error) {
	file, err := os.Open(a.config.FlashcardsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open flashcards file: %w", err)
	}
	defer file.Close()

	var cards []Flashcard
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		question, answer, err := parseFlashcard(scanner.Text())
		if err != nil {
			continue
		}
		cards = append(cards, Flashcard{Question: question, Answer: answer})
	}

	if len(cards) == 0 {
		return nil, ErrNoFlashcards
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read flashcards: %w", err)
	}

	return cards, nil
}

func (a *App) StartQuiz(ctx context.Context) error {
	cards, err := a.LoadFlashcards()
	if err != nil {
		return fmt.Errorf("failed to load flashcards: %w", err)
	}

	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})

	clearScreen()
	for _, card := range cards {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("quiz cancelled: %w", err)
		}
		if err := a.processFlashcard(ctx, card); err != nil {
			return fmt.Errorf("failed to process flashcard: %w", err)
		}
	}

	fmt.Fprintln(a.writer, a.colorize("green", "Congratulations! You've completed the quiz!"))
	return nil
}

func (a *App) processFlashcard(ctx context.Context, card Flashcard) error {
	fmt.Fprintln(a.writer, a.colorize("orange", fmt.Sprintf("Question: %s", card.Question)))

	hintLength := 0
	reader := bufio.NewReader(a.reader)

	for hintLength <= len(card.Answer) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			hint := card.Answer[:hintLength]
			fmt.Fprintln(a.writer, a.colorize("blue", fmt.Sprintf("Hint: %s", hint)))
			fmt.Fprint(a.writer, a.colorize("blue", "> "))

			input, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			input = strings.TrimSpace(input)

			switch input {
			case "?":
				clearLines(2)
				fmt.Fprintln(a.writer, a.colorize("yellow", fmt.Sprintf("Complete: %s", card.Answer)))
				continue
			case ">":
				clearLines(2)
				fmt.Fprintln(a.writer, a.colorize("yellow", fmt.Sprintf("Skipped. Answer was: %s", card.Answer)))
				time.Sleep(2 * time.Second)
				clearScreen()
				return nil
			}

			if strings.EqualFold(input, card.Answer) {
				clearLines(2)
				fmt.Fprintln(a.writer, a.colorize("green", "Correct!"))
				time.Sleep(time.Second)
				clearLines(2)
				fmt.Fprintln(a.writer, a.colorize("orange", fmt.Sprintf("Question: %s\nAnswer: %s", card.Question, card.Answer)))
				time.Sleep(time.Second)
				clearScreen()
				return nil
			}

			matchLength := findMatchLength(input, card.Answer)
			hintLength = matchLength + 1

			clearLines(2)
			if hintLength <= len(card.Answer) {
				fmt.Fprintln(a.writer, a.colorize("red", "Incorrect, try again!"))
				time.Sleep(500 * time.Millisecond)
				clearLines(1)
			}
		}
	}
	return nil
}

func findMatchLength(input, answer string) int {
	input = strings.ToLower(input)
	answer = strings.ToLower(answer)

	var i int
	for i = 0; i < len(input) && i < len(answer); i++ {
		if input[i] != answer[i] {
			break
		}
	}
	return i
}

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func clearLines(count int) {
	for i := 0; i < count; i++ {
		fmt.Print("\033[F\033[2K")
	}
}
