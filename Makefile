include $(GOROOT)/src/Make.inc

TARG=campfire
GOFILES=rooms.go\
	users.go\
	sites.go\
	messages.go\
	commands.go\
	campfire.go\


include $(GOROOT)/src/Make.cmd
