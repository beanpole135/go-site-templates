package main

import (
	"time"
)

type TimeFileDB struct {
	//Input Fields - make sure to set these when you create the object
	Basedir string

	//Internal Fields - leave these alone (state-tracking)

}

// === NewEntry() ===
// Simple creation function for making new log/file entries
func (DB *TimeFileDB) NewEntry(data interface{}, unix_time int64) error {
	if unix_time == 0 { unix_time = time.Now().Unix() } //second-accuracy
	
}

