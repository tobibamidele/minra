package viewport

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tobibamidele/minra/internal/cursor"
	"github.com/tobibamidele/minra/internal/syntax"
	"github.com/tobibamidele/minra/pkg/utils"
	// "github.com/tobibamidele/minra/pkg/utils"
)

// Render renders the viewport
func (v *Viewport) Render(highlighter *syntax.Highlighter, cur *cursor.Cursor, mode Mode) string {
	var b strings.Builder
	
	startLine := v.scrollY
	endLine := v.scrollY + v.height
	if endLine > v.buffer.LineCount() {
		endLine = v.buffer.LineCount()
	}
	
	currentLineStyle := lipgloss.NewStyle().Background(lipgloss.Color("236"))
	lineNumStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	activeLineNumStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("236")).
		Foreground(lipgloss.Color("220")).
		Bold(true)
	
	for lineNum := startLine; lineNum < endLine; lineNum++ {
		line := v.buffer.Line(lineNum)
		isCursorLine := lineNum == cur.Line()
		
		// Line number
		if v.lineNumbers {
			lineNumStr := fmt.Sprintf("%4d ", lineNum+1)
			if isCursorLine {
				b.WriteString(activeLineNumStyle.Render(lineNumStr))
			} else {
				b.WriteString(lineNumStyle.Render(lineNumStr))
			}
		}
		
		// Expand tabs
		displayLine := v.expandTabs(line)
		
		// Apply syntax highlighting
		if highlighter != nil {
			displayLine = highlighter.Highlight(displayLine)
			// displayLine = utils.StripANSI(displayLine)
		}
		
		// Handle horizontal scrolling
		visibleStart := v.scrollX
		visibleEnd := v.scrollX + v.width
		visibleLine := utils.SafeSliceANSI(displayLine, visibleStart, visibleEnd)
		
		// var visibleLine string
		// if visibleStart < len(displayLine) {
		// 	if visibleEnd > len(displayLine) {
		// 		visibleEnd = len(displayLine)
		// 	}
		// 	visibleLine = displayLine[visibleStart:visibleEnd]
		// }
		
		// Render cursor if on this line
		if isCursorLine && (mode == ModeInsert || mode == ModeNormal) {
			displayCursorPos := v.calculateDisplayCol(cur)
			relativeCursorPos := displayCursorPos - v.scrollX
			plainLine := utils.StripANSI(visibleLine)
			if relativeCursorPos > len(plainLine) {
				relativeCursorPos = len(plainLine)
			}
			
			if relativeCursorPos >= 0 && relativeCursorPos <= len(visibleLine) {
				before := visibleLine[:relativeCursorPos]
				
				var cursorStyle lipgloss.Style
				if mode == ModeInsert {
					cursorStyle = lipgloss.NewStyle().
						Background(lipgloss.Color("230")).
						Foreground(lipgloss.Color("0"))
				} else {
					cursorStyle = lipgloss.NewStyle().
						Background(lipgloss.Color("240")).
						Foreground(lipgloss.Color("230"))
				}
				
				if relativeCursorPos < len(visibleLine) {
					cursorChar := string(visibleLine[relativeCursorPos])
					after := ""
					if relativeCursorPos+1 < len(visibleLine) {
						after = visibleLine[relativeCursorPos+1:]
					}
					visibleLine = before + cursorStyle.Render(cursorChar) + after
				} else {
					visibleLine = visibleLine + cursorStyle.Render(" ")
				}
			}
		}
		
		// Apply current line highlight
		if isCursorLine {
			visibleLine = currentLineStyle.Render(visibleLine)
		}
		
		b.WriteString(visibleLine)
		b.WriteString("\n")
	}
	
	// Fill remaining lines
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
