package viewport

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/tobibamidele/minra/internal/cursor"
	"github.com/tobibamidele/minra/internal/syntax"
	"github.com/tobibamidele/minra/internal/syntax/matchers"
	"github.com/tobibamidele/minra/internal/ui"
	"github.com/tobibamidele/minra/pkg/utils"
)

func (v *Viewport) renderBinaryNotice() string {
    notice := "[Binary file — viewing disabled]"
    style := lipgloss.NewStyle().
        Foreground(lipgloss.Color("240")).
        Background(ui.ColorBackground).
        Width(v.Width()).
        Height(v.Height()).
        Align(lipgloss.Center, lipgloss.Center)

    return style.Render(notice)
}


func (v *Viewport) Render(highlighter *syntax.Highlighter, cur *cursor.Cursor, mode Mode) string {
	if v.isBinary {
		v.renderBinaryNotice()
	}
	var b strings.Builder

	startLine := v.scrollY
	endLine := v.scrollY + v.height
	if endLine > v.buffer.LineCount() {
		endLine = v.buffer.LineCount()
	}

	// --- Styles ---
	currentLineStyle := lipgloss.NewStyle().Background(lipgloss.Color("236"))
	lineNumStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#232e33"))
	activeLineNumStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("236")).
		Foreground(lipgloss.Color("220")).
		Bold(true)

	bracketStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("220")).
		Background(lipgloss.Color("238")).
		Bold(true)

	for lineNum := startLine; lineNum < endLine; lineNum++ {
		rawLine := v.buffer.Line(lineNum)
		isCursorLine := lineNum == cur.Line()

		// --- Line numbers ---
		if v.lineNumbers {
			lineNumStr := fmt.Sprintf("%4d ", lineNum+1)
			if isCursorLine {
				b.WriteString(activeLineNumStyle.Render(lineNumStr))
			} else {
				b.WriteString(lineNumStyle.Render(lineNumStr))
			}
		}

		// Expand tabs before syntax highlighting
		displayLine := v.expandTabs(rawLine)

		// --- Bracket matching on raw text ---
		bracketA, bracketB := -1, -1
		if isCursorLine && (mode == ModeInsert || mode == ModeNormal) {
			runes := []rune(rawLine)
			if cur.Col() >= 0 && cur.Col() < len(runes) && matchers.IsBracket(runes[cur.Col()]) {
				match := matchers.FindMatchingBracket(rawLine, cur.Col())
				if match != -1 {
					bracketA, bracketB = cur.Col(), match
				}
			}
		}

		// --- Apply syntax highlighting ---
		if highlighter != nil {
			displayLine = highlighter.Highlight(displayLine)
		}

		// --- Highlight matching brackets (ANSI safe) ---
		if bracketA != -1 && bracketB != -1 {
			displayLine = applyBracketHighlight(displayLine, bracketA, bracketB, bracketStyle)
		}

		// --- Scroll horizontally (ANSI safe) ---
		visibleLine := utils.SafeSliceANSI(displayLine, v.scrollX, v.scrollX+v.width)

		// --- Draw cursor ---
		if isCursorLine && (mode == ModeInsert || mode == ModeNormal) {
			displayCursorPos := v.calculateDisplayCol(cur)
			relativeCursorPos := displayCursorPos - v.scrollX
			plain := utils.StripANSI(visibleLine)
			runeCount := utf8.RuneCountInString(plain)
			if relativeCursorPos > runeCount {
				relativeCursorPos = runeCount
			}

			before := utils.SafeSliceANSI(visibleLine, 0, relativeCursorPos)
			after := utils.SafeSliceANSI(visibleLine, relativeCursorPos+1, utils.VisibleWidth(visibleLine))

			var cursorStyle lipgloss.Style
			if mode == ModeInsert {
				cursorStyle = ui.ActiveCursorStyle
			} else {
				cursorStyle = ui.InactiveCursorStyle
			}

			cursorChar := " "
			if relativeCursorPos < runeCount {
				runes := []rune(plain)
				cursorChar = string(runes[relativeCursorPos])
			}

			visibleLine = before + cursorStyle.Render(cursorChar) + after
		}

		if isCursorLine {
			visibleLine = currentLineStyle.Width(v.Width()).Render(visibleLine)
		} else { // Render lines with background
			visibleLine = lipgloss.NewStyle().
				Background(ui.ColorBackground).
				Width(v.Width()).
				Render(visibleLine)
		}

		b.WriteString(visibleLine)
		b.WriteString("\n")
	}

	// --- Fill empty space ---
	for i := endLine - startLine; i < v.height; i++ {
		if v.lineNumbers {
			b.WriteString(lineNumStyle.Render("   ~ "))
		} else {
			b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("~"))
		}
		b.WriteString("\n")
	}

	return b.String()
}

// --- Highlight helper ---

func applyBracketHighlight(line string, idx1, idx2 int, style lipgloss.Style) string {
	runes := []rune(utils.StripANSI(line))
	if idx1 > idx2 {
		idx1, idx2 = idx2, idx1
	}
	if idx1 < 0 || idx2 >= len(runes) {
		return line
	}

	// Convert visible indices to ANSI-safe slices
	beforeFirst := utils.SafeSliceANSI(line, 0, idx1)
	first := utils.SafeSliceANSI(line, idx1, idx1+1)
	middle := utils.SafeSliceANSI(line, idx1+1, idx2)
	second := utils.SafeSliceANSI(line, idx2, idx2+1)
	after := utils.SafeSliceANSI(line, idx2+1, utils.VisibleWidth(line))

	// Apply bracket style
	first = style.Render(first)
	second = style.Render(second)

	return beforeFirst + first + middle + second + after
}

func (v *Viewport) expandTabs(line string) string {
	var result strings.Builder
	col := 0
	for _, ch := range line {
		if ch == '\t' {
			spacesToAdd := v.tabSize - (col % v.tabSize)
			for i := 0; i < spacesToAdd; i++ {
				result.WriteRune(' ')
			}
			col += spacesToAdd
		} else {
			result.WriteRune(ch)
			col++
		}
	}
	return result.String()
}
