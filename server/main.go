// package comment
package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/washtubs/evplanner"

	logging "github.com/op/go-logging"
)

// TODO: move these to a main?
var log *logging.Logger
var logBackend *logging.LogBackend
var LeveledLogBackend logging.LeveledBackend

var stdStore evplanner.Store

type serverConfig struct {
	port *int
}

// All dependency management will be done here. I'm hoping I just need
// a flat object graph of singletons
func configureSingletonDependencies() {
	stdStore = new(evplanner.InMemoryStore)
}

func main() {
	startServer()
}

func startServer() {
	var (
		port = flag.Int("port", 8145, "port number for the HTTP server")
	)
	flag.Parse()

	configureSingletonDependencies()

	conf := serverConfig{port: port}

	mux := http.NewServeMux()
	mux.HandleFunc("/read", handleRead)
	mux.HandleFunc("/write", handleWrite)
	mux.HandleFunc("/lock", handleLock)
	mux.HandleFunc("/unlock", handleUnlock)

	server := http.Server{
		Addr:    ":" + strconv.Itoa(*conf.port),
		Handler: mux,
	}

	log.Debugf("Starting ev-planner management server at port %d", *conf.port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

var forbiddenResp http.Response = http.Response{
	StatusCode: 403,
}

var errorResp http.Response = http.Response{
	StatusCode: 500,
}

var httpUnlock chan bool = nil

func writeResponse(w http.ResponseWriter, resp http.Response) {
	w.WriteHeader(resp.StatusCode)
	resp.Write(w)
}

func handleLock(w http.ResponseWriter, req *http.Request) {
	if httpUnlock != nil {
		prematureLockResp := http.Response{StatusCode: evplanner.PrematureLock}
		writeResponse(w, prematureLockResp)
		return
	}

	done := make(chan bool, 1)

	timeoutResp := http.Response{StatusCode: evplanner.LockTimeout}
	timeoutTicker := time.NewTicker(time.Duration(30) * time.Second)

	go func() {
		defer stdStore.UnlockForModification()
		stdStore.LockForModification()
		done <- true

		// after the lock is obtained make it so that
		httpUnlock = make(chan bool, 1)
		lockMaxDurationExceededTicker := time.NewTicker(time.Duration(10) * time.Second)
		select {
		case <-lockMaxDurationExceededTicker.C:
			log.Debug("Exceeded max duration")
		case <-httpUnlock:
			log.Debug("Received unlock")
		}
		httpUnlock = nil
		log.Debug("nil httpUnlock")

		// unlock
	}()

	select {
	case <-done:
		return
	case <-timeoutTicker.C:
		writeResponse(w, timeoutResp)
	}

}

func handleUnlock(w http.ResponseWriter, req *http.Request) {
	if httpUnlock == nil {
		prematureUnlockResp := http.Response{StatusCode: evplanner.PrematureUnlock}
		writeResponse(w, prematureUnlockResp)
	}

	go func() {
		httpUnlock <- true
	}()
}

func handleRead(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		writeResponse(w, forbiddenResp)
		return
	}

	w.Write([]byte(stdStore.Read().Serialize()))
}

func handleWrite(w http.ResponseWriter, req *http.Request) {
	if req.Method != "PUT" {
		writeResponse(w, forbiddenResp)
		return
	}

	if !stdStore.IsLockedForModification() {
		lockedResp := http.Response{StatusCode: evplanner.WriteErrorNotLocked}
		writeResponse(w, lockedResp)
		return
	}

	defer req.Body.Close()
	bytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("Failed to read the request body")
		writeResponse(w, errorResp)
	}

	stdStore.Write(evplanner.PlaceholderFromString(string(bytes)))
}

func init() {
	log = logging.MustGetLogger("notification-tracker")
	logBackend = logging.NewLogBackend(os.Stderr, "", 0)
	LeveledLogBackend = logging.AddModuleLevel(logBackend)
	LeveledLogBackend.SetLevel(logging.DEBUG, "")
	logging.SetBackend(LeveledLogBackend)

	logging.SetFormatter(logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	))
}
