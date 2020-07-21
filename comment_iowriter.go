package reviewdog

import (
	"context"
	"fmt"
	"io"
	"strings"
)

var _ CommentService = &RawCommentWriter{}

// RawCommentWriter is comment writer which writes results to given writer
// without any formatting.
type RawCommentWriter struct {
	w io.Writer
}

func NewRawCommentWriter(w io.Writer) *RawCommentWriter {
	return &RawCommentWriter{w: w}
}

func (s *RawCommentWriter) Post(_ context.Context, c *Comment) error {
	_, err := fmt.Fprintln(s.w, strings.Join(c.Result.Lines, "\n"))
	return err
}

var _ CommentService = &UnifiedCommentWriter{}

// UnifiedCommentWriter is comment writer which writes results to given writer
// in one of following unified formats.
//
// Format:
//   - <file>: [<tool name>] <message>
//   - <file>:<lnum>: [<tool name>] <message>
//   - <file>:<lnum>:<col>: [<tool name>] <message>
// where <message> can be multiple lines.
type UnifiedCommentWriter struct {
	w io.Writer
}

func NewUnifiedCommentWriter(w io.Writer) *UnifiedCommentWriter {
	return &UnifiedCommentWriter{w: w}
}

func (mc *UnifiedCommentWriter) Post(_ context.Context, c *Comment) error {
	loc := c.Result.Diagnostic.GetLocation()
	s := loc.GetPath()
	start := loc.GetRange().GetStart()
	if start.GetLine() > 0 {
		s += fmt.Sprintf(":%d", start.GetLine())
		if start.GetColumn() > 0 {
			s += fmt.Sprintf(":%d", start.GetColumn())
		}
	}
	s += fmt.Sprintf(": [%s] %s", c.ToolName, c.Body)
	_, err := fmt.Fprintln(mc.w, s)
	return err
}
