package downloader

import (
	tea "github.com/charmbracelet/bubbletea"
	"io"
	"os"
)

var P *tea.Program

type ProgressMsg float64

type progressErrMsg struct{ err error }

type ProgressWriter struct {
	Total      int
	Downloaded int
	File       *os.File
	Reader     io.Reader
	OnProgress func(float64)
}

func (pw *ProgressWriter) Start() {
	// TeeReader calls pw.Write() each time a new response is received
	_, err := io.Copy(pw.File, io.TeeReader(pw.Reader, pw))
	if err != nil {
		P.Send(progressErrMsg{err})
	}
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	pw.Downloaded += len(p)
	if pw.Total > 0 && pw.OnProgress != nil {
		pw.OnProgress(float64(pw.Downloaded) / float64(pw.Total))
	}
	return len(p), nil
}
