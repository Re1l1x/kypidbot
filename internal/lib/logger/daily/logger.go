package daily

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Writer struct {
	dir  string
	file *os.File
	date string
	mu   sync.Mutex
	env  string
}

func NewLogsWriter(dir string, env string) *Writer {
	return &Writer{dir: dir, env: env}
}

func (w *Writer) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	today := time.Now().Format("2006-01-02")

	if w.file == nil || w.date != today {
		if w.file != nil {
			w.file.Close()
		}

		os.MkdirAll(w.dir, 0755)

		path := filepath.Join(w.dir, "kypidbot-"+w.env+"-"+today+".log")

		f, err := os.OpenFile(path,
			os.O_CREATE|os.O_APPEND|os.O_WRONLY,
			0644,
		)
		if err != nil {
			return 0, err
		}

		w.file = f
		w.date = today
	}

	return w.file.Write(p)
}
