package filesearch

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Result struct {
	Name      string
	Path      string
	ParentDir string
	Score     int
}

type Engine struct {
	root  string
	files []string
}

// NewEngine creates and caches file list under root.
func NewEngine(root string) (*Engine, error) {
	e := &Engine{root: root}
	if err := e.build(); err != nil {
		return nil, err
	}
	return e, nil
}

func (e *Engine) build() error {
	var files []string
	err := filepath.WalkDir(e.root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		base := filepath.Base(p)
		if strings.HasPrefix(base, ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.IsDir() {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		return err
	}
	e.files = files
	return nil
}

// Refresh rebuilds the file cache.
func (e *Engine) Refresh() error {
	return e.build()
}

func (e *Engine) Files() []string {
	return e.files
}

// Search returns fuzzy-matched results based on filename only.
func (e *Engine) Search(q string, limit int) []Result {
	q = strings.TrimSpace(q)
	if q == "" {
		return nil
	}
	qLower := strings.ToLower(q)
	out := make([]Result, 0, 64)

	for _, p := range e.files {
		name := filepath.Base(p)
		ok, score := fzfScore(strings.ToLower(name), qLower)
		if ok {
			out = append(out, Result{
				Name:      name,
				Path:      p,
				ParentDir: filepath.Base(filepath.Dir(p)),
				Score:     score,
			})
		}
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Score == out[j].Score {
			if len(out[i].Path) == len(out[j].Path) {
				return out[i].Path < out[j].Path
			}
			return len(out[i].Path) < len(out[j].Path)
		}
		return out[i].Score > out[j].Score
	})

	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out
}

// fzfScore implements a simple FZF-like scoring.
// Returns (matched, score). Higher score is better.
func fzfScore(name, query string) (bool, int) {
	// Quick rejects
	if len(query) == 0 || len(query) > len(name) {
		return false, 0
	}

	// Exact prefix
	if strings.HasPrefix(name, query) {
		return true, 1000 + len(query)
	}

	// Full substring
	if idx := strings.Index(name, query); idx >= 0 {
		// prefer earlier occurrence and shorter name
		score := 500 + (len(query) * 2) - idx
		return true, score
	}

	// Subsequence match with bonuses for consecutive matches and start-of-word
	score := 0
	i := 0
	consec := 0
	foundAny := false
	for qi := 0; qi < len(query); qi++ {
		qc := query[qi]
		found := false
		for i < len(name) {
			nc := name[i]
			if nc == qc {
				found = true
				foundAny = true
				// base match
				score += 10
				// bonus if consecutive
				if consec > 0 {
					score += 5 * consec
				}
				// bonus if match at start or after separator
				if i == 0 || name[i-1] == '_' || name[i-1] == '-' || name[i-1] == '.' {
					score += 15
				}
				consec++
				i++
				break
			}
			// gap penalty
			consec = 0
			i++
			score -= 1
		}
		if !found {
			return false, 0
		}
	}
	// shorter name slight bonus
	score += 5 * (len(query))
	return foundAny, score
}

