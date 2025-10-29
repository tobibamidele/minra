package editor

import (
	"fmt"

	"github.com/tobibamidele/minra/internal/viewport"
)

func (e *Editor) getCursorLinePercent() string {
    buf := e.bufferMgr.ActiveBuffer()
    cur := buf.Cursor()
    line := cur.Line() + 1
    total := buf.LineCount()

    if line == 1 {
        return "Top"
    }
    if line == total {
        return "Bottom"
    }

    percent := float64(line) / float64(total) * 100
    return fmt.Sprintf("%d%%", int(percent))
}

func (e *Editor) numberOfLinesAdded() (string, error) {
    if e.mode != viewport.ModeNormal && e.bufferMgr.ActiveBuffer().Modified() {
        return "", fmt.Errorf("viewport mode not normal")
    }
    prevLines := e.bufferMgr.ActiveBuffer().PreviousLineCount()
    currLines := e.bufferMgr.ActiveBuffer().LineCount()


    if prevLines > currLines {
        return fmt.Sprintf("-%d", prevLines - currLines), nil
    } else {
        return fmt.Sprintf("+%d", currLines - prevLines), nil
    }
}
