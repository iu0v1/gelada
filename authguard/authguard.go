package authguard

import (
	"encoding/gob"
	"os"
)

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
type Client struct{}

type AuthGuard struct {
	options *Options
	file    *os.File
	data    []*Client
	// log // TODO : make logger

	// TODO : create mutex?
}

func New(o Options) (*AuthGuard, error) {
	ag := &AuthGuard{options: &o}

	if ag.options.Store == "::memory::" {
		// create new data struct and put in ag.data
	} else if ag.options.Store == "" {
		// error
	} else {
		data := []*Client{}
		// open gob file
		file, err := os.OpenFile(ag.options.Store, os.O_CREATE|os.O_RDWR|os.O_SYNC, 700)
		if err != nil {
			// return error
		}

		dec := gob.NewDecoder(file)
		if err := dec.Decode(&data); err != nil {
			// return error
		}

		ag.data = data
	}

	return ag, nil
}
