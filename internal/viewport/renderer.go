package viewport

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/tobibamidele/minra/internal/cursor"
	"github.com/tobibamidele/minra/internal/syntax"
	"github.com/tobibamidele/minra/internal/syntax/matchers"
	"github.com/tobibamidele/minra/pkg/utils"
)

// // Render renders the viewport
// func (v *Viewport) Render(highlighter *syntax.Highlighter, cur *cursor.Cursor, mode Mode) string {
// 	var b strings.Builder
//
// 	startLine := v.scrollY
// 	endLine := v.scrollY + v.height
// 	if endLine > v.buffer.LineCount() {
// 		endLine = v.buffer.LineCount()
// 	}
//
// 	currentLineStyle := lipgloss.NewStyle().Background(lipgloss.Color("236"))
// 	lineNumStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
// 	activeLineNumStyle := lipgloss.NewStyle().
// 		Background(lipgloss.Color("236")).
// 		Foreground(lipgloss.Color("220")).
// 		Bold(true)
//
// 	// Style for highlighting matching brackets
// 	bracketStyle := lipgloss.NewStyle().
// 		Foreground(lipgloss.Color("220")).
// 		Background(lipgloss.Color("237")).
// 		Bold(true)
//
// 	for lineNum := startLine; lineNum < endLine; lineNum++ {
// 		rawLine := v.buffer.Line(lineNum)
// 		isCursorLine := lineNum == cur.Line()
//
// 		// --- Line number ---
// 		if v.lineNumbers {
// 			lineNumStr := fmt.Sprintf("%4d ", lineNum+1)
// 			if isCursorLine {
// 				b.WriteString(activeLineNumStyle.Render(lineNumStr))
// 			} else {
// 				b.WriteString(lineNumStyle.Render(lineNumStr))
// 			}
// 		}
//
// 		// --- Expand tabs before highlighting ---
// 		displayLine := v.expandTabs(rawLine)
//
// 		// --- Apply syntax highlighting ---
// 		if highlighter != nil {
// 			displayLine = highlighter.Highlight(displayLine)
// 		}
//
// 		// --- Highlight matching brackets if cursor is on one ---
// 		if isCursorLine && (mode == ModeNormal || mode == ModeInsert) {
// 			runes := []rune(rawLine)
// 			cursorCol := cur.Col()
// 			if cursorCol >= 0 && cursorCol < len(runes) && matchers.IsBracket(runes[cursorCol]) {
// 				matchIdx := matchers.FindMatchingBracket(rawLine, cursorCol)
// 				if matchIdx != -1 {
// 					displayLine = applyBracketHighlight(displayLine, cursorCol, matchIdx, bracketStyle)
// 				}
// 			}
// 		}
//
// 		// --- Handle horizontal scrolling safely (ANSI aware) ---
// 		visibleStart := v.scrollX
// 		visibleEnd := v.scrollX + v.width
// 		visibleLine := utils.SafeSliceANSI(displayLine, visibleStart, visibleEnd)
//
// 		// --- Draw cursor ---
// 		if isCursorLine && (mode == ModeInsert || mode == ModeNormal) {
// 			displayCursorPos := v.calculateDisplayCol(cur)
// 			relativeCursorPos := displayCursorPos - v.scrollX
//
// 			plainLine := utils.StripANSI(visibleLine)
// 			runeCount := utf8.RuneCountInString(plainLine)
// 			if relativeCursorPos > runeCount {
// 				relativeCursorPos = runeCount
// 			}
//
// 			if relativeCursorPos >= 0 {
// 				before := utils.SafeSliceANSI(visibleLine, 0, relativeCursorPos)
// 				after := utils.SafeSliceANSI(visibleLine, relativeCursorPos+1, utils.VisibleWidth(visibleLine))
//
// 				var cursorStyle lipgloss.Style
// 				if mode == ModeInsert {
// 					cursorStyle = lipgloss.NewStyle().
// 						Background(lipgloss.Color("230")). // pale yellow
// 						Foreground(lipgloss.Color("0"))    // black text
// 				} else {
// 					cursorStyle = lipgloss.NewStyle().
// 						Background(lipgloss.Color("240")). // gray block
// 						Foreground(lipgloss.Color("230"))
// 				}
//
// 				cursorChar := " "
// 				if relativeCursorPos < runeCount {
// 					runes := []rune(plainLine)
// 					cursorChar = string(runes[relativeCursorPos])
// 				}
//
// 				visibleLine = before + cursorStyle.Render(cursorChar) + after
// 			}
// 		}
//
// 		// --- Highlight entire line (in both normal and insert mode) ---
// 		if isCursorLine {
// 			visibleLine = currentLineStyle.Render(visibleLine)
// 		}
//
// 		b.WriteString(visibleLine)
// 		b.WriteString("\n")
// 	}
//
// 	// --- Fill remaining lines ---
// 	for i := endLine - startLine; i < v.height; i++ {
// 		if v.lineNumbers {
// 			b.WriteString(lineNumStyle.Render("   ~ "))
// 		} else {
// 			b.WriteString(lipgloss.NewStyle().
// 				Foreground(lipgloss.Color("240")).Render("~"))
// 		}
// 		b.WriteString("\n")
// 	}
//
// 	return b.String()
// }

func (v *Viewport) Render(highlighter *syntax.Highlighter, cur *cursor.Cursor, mode Mode) string {
	var b strings.Builder

	startLine := v.scrollY
	endLine := v.scrollY + v.height
	if endLine > v.buffer.LineCount() {
		endLine = v.buffer.LineCount()
	}

	// --- Styles ---
	currentLineStyle := lipgloss.NewStyle().Background(lipgloss.Color("236"))
	lineNumStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
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

			cursorStyle := lipgloss.NewStyle()
			if mode == ModeInsert {
				cursorStyle = cursorStyle.Background(lipgloss.Color("230")).Foreground(lipgloss.Color("0"))
			} else {
				cursorStyle = cursorStyle.Background(lipgloss.Color("240")).Foreground(lipgloss.Color("230"))
			}

			cursorChar := " "
			if relativeCursorPos < runeCount {
				runes := []rune(plain)
				cursorChar = string(runes[relativeCursorPos])
			}

			visibleLine = before + cursorStyle.Render(cursorChar) + after
		}

		if isCursorLine {
			// currentLineStyle = currentLineStyle.Foreground(lipgloss.Color("230"))
			visibleLine = currentLineStyle.Width(v.Width()).Render(visibleLine)
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

