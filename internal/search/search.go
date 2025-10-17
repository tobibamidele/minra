package search

import (
	"strings"

	"github.com/tobibamidele/minra/internal/buffer"
)

// Result represents a search result
type Result struct {
	Line   int
	Column int
	Length int
}

// Engine performs text search
type Engine struct {
	query         string
	caseSensitive bool
	regex         bool
	results       []Result
	currentIdx    int
}

// NewEngine creates a new search engine
func NewEngine() *Engine {
	return &Engine{
		caseSensitive: false,
		regex:         false,
		results:       make([]Result, 0),
		currentIdx:    -1,
	}
}

// SetQuery sets the search query
func (e *Engine) SetQuery(query string) {
	e.query = query
	e.results = make([]Result, 0)
	e.currentIdx = -1
}

// SetCaseSensitive sets case sensitivity
func (e *Engine) SetCaseSensitive(sensitive bool) {
	e.caseSensitive = sensitive
}

// Search performs search on buffer
func (e *Engine) Search(buf *buffer.Buffer) []Result {
	e.results = make([]Result, 0)

	if e.query == "" {
		return e.results
	}

	searchQuery := e.query
	if !e.caseSensitive {
		searchQuery = strings.ToLower(searchQuery)
	}

	for lineNum := 0; lineNum < buf.LineCount(); lineNum++ {
		line := buf.Line(lineNum)
		searchLine := line
		if !e.caseSensitive {
			searchLine = strings.ToLower(searchLine)
		}

		col := 0
		for {
			idx := strings.Index(searchLine[col:], searchQuery)
			if idx == -1 {
				break
			}

			e.results = append(e.results, Result{
				Line:   lineNum,
				Column: col + idx,
				Length: len(e.query),
			})
			col += idx + 1
		}
	}

	if len(e.results) > 0 {
		e.currentIdx = 0
	}

	return e.results
}

// Next returns the next result
func (e *Engine) Next() *Result {
	if len(e.results) == 0 {
		return nil
	}

	e.currentIdx = (e.currentIdx + 1) % len(e.results)
	return &e.results[e.currentIdx]
}

// Previous returns previous result
func (e *Engine) Previous() *Result {
	if len(e.results) == 0 {
		return nil
	}

	e.currentIdx = (e.currentIdx - 1 + len(e.results)) % len(e.results)
	return &e.results[e.currentIdx]
}

// Current returns the current result
func (e *Engine) Current() *Result {
	if e.currentIdx < 0 || e.currentIdx >= len(e.results) {
		return nil
	}
	return &e.results[e.currentIdx]
}

// Results returns all results
func (e *Engine) Results() []Result {
	return e.results
}

// Count returns nummber of results
func (e *Engine) Count() int {
	return len(e.results)
}
