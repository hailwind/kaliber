module github.com/mwat56/kaliber

go 1.14

require (
	github.com/NYTimes/gziphandler v1.1.1
	github.com/mwat56/apachelogger v1.4.9
	github.com/mwat56/cssfs v0.2.1
	github.com/mwat56/errorhandler v1.1.6
	github.com/mwat56/ini v1.3.8
	github.com/mwat56/jffs v0.1.0
	github.com/mwat56/kaliber/db v0.0.0-20200413173747-0c15ddb66a49
	github.com/mwat56/passlist v1.3.1
	github.com/mwat56/sessions v0.3.9
	github.com/mwat56/whitespace v0.2.1
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	golang.org/x/crypto v0.0.0-20200429183012-4b2356b1ed79 // indirect
	golang.org/x/sys v0.0.0-20200509044756-6aff5f38e54f // indirect
)

replace github.com/mwat56/kaliber/db => ./db
