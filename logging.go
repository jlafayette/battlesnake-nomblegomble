package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"sync"

	"github.com/jlafayette/battlesnake-go/t"
)

var (
	logger *GameLogger
)

func init() {
	logger = NewGameLogger("gamelogs")
}

type GameLogger struct {
	rootDir string
	lookup  map[string]*FileLogger
}

func NewGameLogger(rootDir string) *GameLogger {
	_, err := os.Stat(rootDir)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(rootDir, 0755)
		if errDir != nil {
			log.Fatalf("ERROR: Failed to create log directory %s, %s", rootDir, err)
		}
	}
	lookup := make(map[string]*FileLogger)
	return &GameLogger{
		rootDir: rootDir,
		lookup:  lookup,
	}
}

func (g *GameLogger) newGame(id string) error {
	fileLogger, err := NewFileLogger(path.Join(g.rootDir, id+".log"))
	if err != nil {
		return err
	}
	g.lookup[id] = fileLogger
	return nil
}

func (g *GameLogger) removeGame(id string) {
	if fileLogger, ok := g.lookup[id]; ok {
		err := fileLogger.Close()
		if err != nil {
			log.Printf("WARNING: Failed to close log file for game id %s, %s", id, err)
		}
		delete(g.lookup, id)
	}
}

func (g *GameLogger) log(id string, msg []byte) {
	if fileLogger, ok := g.lookup[id]; ok {
		_, err := fileLogger.Write(msg)
		if err != nil {
			log.Printf("WARNING: Failed to write to log for game id %s, %s", id, err)
		}
		_, err = fileLogger.Write([]byte("\n"))
		if err != nil {
			log.Printf("WARNING: Failed to write to log for game id %s, %s", id, err)
		}
	}
}

type FileLogger struct {
	lock     sync.Mutex
	filename string // should be set to the actual filename
	fp       *os.File
}

// Make a new FileLogger. Return nil if error occurs during setup.
func NewFileLogger(filename string) (*FileLogger, error) {
	w := &FileLogger{filename: filename}
	err := w.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open logfile %s, %s", filename, err)
	}
	return w, nil
}

// Create a file.
func (w *FileLogger) Open() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.fp, err = os.Create(w.filename)
	return
}

func (w *FileLogger) Close() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	// Close existing file if open
	if w.fp != nil {
		err = w.fp.Close()
		w.fp = nil
	}
	return
}

// Write satisfies the io.Writer interface.
func (w *FileLogger) Write(output []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.fp.Write(output)
}

func init_game_log(id string) {
	err := logger.newGame(id)
	if err != nil {
		log.Printf("ERROR: Failed to initialize logs for game %s, %s", id, err)
	}
}

func close_game_log(id string) {
	logger.removeGame(id)
}

func log_move(state *t.GameState) {
	data, err := json.Marshal(state)
	if err != nil {
		log.Printf("ERROR: Failed to marshal game state for logs, %s", err)
	}
	logger.log(state.Game.ID, data)
}

func log_move_response(id string, response *t.BattlesnakeMoveResponse) {
	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("ERROR: Failed to marshal respone state for logs, %s", err)
	}
	logger.log(id, data)
}
