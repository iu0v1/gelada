package authguard

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"sync"
	// "time"
)

// Options - structure, which is used to configure authguard.
type Options struct {
	Attempts        int
	LockoutDuration int // lockout time in seconds
	MaxLockouts     int // the maximum amount of lockouts, before ban
	BanDuration     int // subj; seconds
	ResetDuration   int // time after which to reset the number of lockouts
	BindToHost      bool
	BindToUserAgent bool

	UpdateDataDuration int // actualize records every X seconds
	UpdateDataAfter    int // actualize immediately after X updates

	SyncDuration int // sync every X seconds
	SyncAfter    int // sync immediately after X updates

	Store string // "::memory::", "/url/to/file"

	// LogDestination // TODO : make log writer

	// Exceptions

	// Backend // ql || gob || ...
}

// place for client data
type Client struct {
	ag *AuthGuard
	mu sync.Mutex
}

// reset cliet
func (c *Client) Reset() {
	// reset
	// sync
}

// AuthGuard - main struct.
type AuthGuard struct {
	options *Options
	file    *os.File
	data    []*Client
	// log // TODO : make logger

	mu sync.Mutex
}

// New - init and return new AuthGuard struct.
func New(o Options) (*AuthGuard, error) {
	ag := &AuthGuard{options: &o}

	if ag.options.Store == "::memory::" {
		ag.data = []*Client{}
	} else if ag.options.Store == "" {
		return ag, fmt.Errorf("LoginRoute not declared\n")
	} else {
		data := []*Client{}

		file, err := os.OpenFile(
			ag.options.Store,
			os.O_CREATE|os.O_RDWR|os.O_SYNC,
			700,
		)
		if err != nil {
			return ag, fmt.Errorf("error to open Store file: %v\n", err)
		}

		dec := gob.NewDecoder(file)
		if err := dec.Decode(&data); err != nil {
			return ag, fmt.Errorf("error to read Store file: %v\n", err)
		}

		ag.data = data
	}

	// TODO : check options

	return ag, nil
}

func (ag *AuthGuard) sync() error {
	ag.mu.Lock()
	defer ag.mu.Unlock()

	enc := gob.NewEncoder(ag.file)
	if err := enc.Encode(&ag.data); err != nil {
		return fmt.Errorf("error to encode Store file: %v\n", err)
	}

	return nil
}

func (ag *AuthGuard) Check(req *http.Request) bool {
	// get client
	// check client
	// sync on fail
	return true
}
