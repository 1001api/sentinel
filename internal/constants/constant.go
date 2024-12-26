// Package constants defines global variables used throughout the application.
//
// These variables SHOULD NOT be manipulated directly as they are global and shared across all packages.
// Future updates can introduce functionality to return copies of these variables for added safety,
// but this approach has overhead. For now, direct access is used to keep things simple and efficient.
package constants

var DownloadFormats = map[string]bool{
	"csv":  true,
	"xlsx": true,
	"json": true,
	"pdf":  true,
	"html": true,
}

var Intervals = map[string]int{
	"last_7_days":   7,
	"last_30_days":  30,
	"last_60_days":  60,
	"last_180_days": 180,
	"last_year":     360,
	"all_time":      -1,
}
