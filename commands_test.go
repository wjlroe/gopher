package main

import (
	. "launchpad.net/gocheck"
	"testing"
	"time"
)

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type S struct{}
var _ = Suite(&S{})

const dunno = "Dunno what you are asking me. WTF dude?"

var whatTimeTests = map[string]string{
	"what time is it?": time.LocalTime().Format(time.Kitchen),
	"WHAT TIME HUH?": time.LocalTime().Format(time.Kitchen),
	"time what": dunno,
	"lajsdlajs": time.LocalTime().Format(time.Kitchen),
}

func (s *S) TestWhatTime(c *C) {
	for input, expected := range whatTimeTests {
		actual := Command(input)
		c.Check(actual, Equals, expected)
	}
}

const greeting = "I am Gopher, I'm here to help."

var whoTests = map[string]string{
	"Who are you?": greeting,
	"WHO ARE YOU": greeting,
	"ashksdfhds f": greeting,
	"are you who": dunno,
	"skdhfskdhf": greeting,
}

func (s *S) TestWhoAreYou(c *C) {
	for input, expected := range whoTests {
		actual := Command(input)
		c.Check(actual, Equals, expected)
	}
}
