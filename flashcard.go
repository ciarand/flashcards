package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var bucketName = []byte("flashcards")

// flashcard represents a Flashcard, i.e. a complete question and answer set.
type flashcard struct {
	Question string
	Answer   string
	Class    string

	linenum int
}

type server struct {
	repo repo
}

type flashcardSlice []flashcard

type repo interface {
	all() (flashcardSlice, error)
}

func internalServerError(w http.ResponseWriter, err error) {
	printerr(err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("Something went wrong: %s", err.Error())))
}

func (h *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseGlob("layouts/*.tpl")
	if err != nil {
		internalServerError(w, err)
		return
	}

	cards, err := h.repo.all()
	if err != nil {
		internalServerError(w, err)
		return
	}

	// write to a buffer so template errors can be correctly rendered without
	// already having mucked with the response
	//
	// totally guessing on the length thing here
	//
	// TODO: maybe add some logging or w/e to auto-adjust the buffer size?
	buf := bytes.NewBuffer(make([]byte, 0, 8048))

	err = tpl.ExecuteTemplate(buf, "flashcards", struct{ Cards flashcardSlice }{cards})
	if err != nil {
		internalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	printerr(ignoreBytesWritten(io.Copy(w, buf)))
}

type fileRepo struct {
	file string
}

func (f flashcardSlice) add(c flashcard) error {
	if c.Question == "" {
		return fmt.Errorf("your card's question can't be blank (card started on line %d)", c.linenum)
	}

	if c.Answer == "" {
		return fmt.Errorf("cards need answers! (card started on line %d)", c.linenum)
	}

	f = append(f, c)

	return nil
}

func (f fileRepo) all() (flashcardSlice, error) {
	data, err := ioutil.ReadFile(f.file)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(data)

	var line string
	// assume a default size of 10, we can increase later
	collection := make(flashcardSlice, 0, 10)

	linenum := 0
	card := flashcard{linenum: linenum}

	for {
		line, err = buf.ReadString('\n')
		if err == io.EOF {
			// save the most recent card
			return append(collection, card), nil
		} else if err != nil {
			return nil, err
		}

		linenum++

		switch {
		// if the line is blank, ignore it
		case len(line) == 0:
			continue
		// if the line is a comment, ignore it
		case line[0] == '#':
			continue
		case strings.HasPrefix(line, "QUESTION: "):
			if card.Question != "" {
				collection = append(collection, card)
				card = flashcard{linenum: linenum}
			}

			card.Question = strings.TrimSpace(strings.TrimLeft(line, "QUESTION:"))

		case strings.HasPrefix(line, "ANSWER:"):
			if card.Answer != "" {
				collection = append(collection, card)
				card = flashcard{linenum: linenum}
			}

			card.Answer = strings.TrimSpace(strings.TrimLeft(line, "ANSWER:"))

		case strings.HasPrefix(line, "CLASS:"):
			if card.Class != "" {
				collection = append(collection, card)
				card = flashcard{linenum: linenum}
			}

			card.Class = strings.TrimSpace(strings.TrimLeft(line, "CLASS:"))
		}
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	mux := http.NewServeMux()

	mux.Handle("/", &server{fileRepo{"flashcards.txt"}})

	static := http.FileServer(http.Dir("./static"))
	mux.Handle("/css/", static)
	mux.Handle("/js/", static)
	mux.Handle("/fonts/", static)

	return http.ListenAndServe(":8080", mux)
}

// printerr is intended to be used for errors we can't do anything about
// (defer'd closes, write errors, etc.).
func printerr(err error) {
	if err != nil {
		fmt.Printf("Error encountered: %s\n", err)
	}
}

func printerrfn(fn func() error) {
	printerr(fn())
}

func ignoreBytesWritten(_ int64, err error) error {
	return err
}
