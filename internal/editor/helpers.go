package editor

import (
	"fmt"
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
