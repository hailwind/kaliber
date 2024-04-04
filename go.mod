module github.com/mwat56/kaliber

go 1.21

require (
	github.com/NYTimes/gziphandler v1.1.1
	github.com/mwat56/apachelogger v1.6.3
	github.com/mwat56/cssfs v0.2.7
	github.com/mwat56/errorhandler v1.1.10
	github.com/mwat56/ini v1.5.4
	github.com/mwat56/jffs v0.1.4
	github.com/mwat56/kaliber/db v0.0.0-20230309215109-1c5b592b38b8
	github.com/mwat56/passlist v1.3.10
	github.com/mwat56/sessions v0.3.15
	github.com/mwat56/whitespace v0.2.5
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
)

require (
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/term v0.18.0 // indirect
)

replace (
	github.com/mwat56/apachelogger => ../apachelogger
	github.com/mwat56/cssfs => ../cssfs
	github.com/mwat56/errorhandler => ../errorhandler
	github.com/mwat56/hashtags => ../hashtags
	github.com/mwat56/ini => ../ini
	github.com/mwat56/jffs => ../jffs
	github.com/mwat56/kaliber/db => ./db
	github.com/mwat56/passlist => ../passlist
	github.com/mwat56/sessions => ../sessions
	github.com/mwat56/uploadhandler => ../uploadhandler
	github.com/mwat56/whitespace => ../whitespace
)
