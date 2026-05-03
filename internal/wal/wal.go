package wal

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"
)

type Entry struct {
	Command string `json:"cmd"`
	Key     string `json:"key"`
	Value   string `json:"value"`
}

type WAL struct {
	mu     sync.Mutex
	file   *os.File
	writer *bufio.Writer
}

func Open(path string) (*WAL, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &WAL{
		file:   f,
		writer: bufio.NewWriter(f),
	}, nil
}

func (w *WAL) Append(entry Entry) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	line, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	line = append(line, '\n')
	if _, err := w.writer.Write(line); err != nil {
		return err
	}
	if err := w.writer.Flush(); err != nil {
		return err
	}
	return w.file.Sync()
}

func (w *WAL) Replay(fn func(Entry)) error {
	if _, err := w.file.Seek(0, 0); err != nil {
		return err
	}

	scanner := bufio.NewScanner(w.file)
	for scanner.Scan() {
		var entry Entry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			return err
		}
		fn(entry)
	}
	return scanner.Err()
}

func (w *WAL) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.writer.Flush()
	return w.file.Close()
}
