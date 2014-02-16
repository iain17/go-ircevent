package irc

import (
	//	"github.com/thoj/go-ircevent"
	"testing"
	"time"
)

func TestConnection(t *testing.T) {
	irccon := IRC("go-eventirc", "go-eventirc")
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	err := irccon.Connect("irc.freenode.net:6667")
	if err != nil {
		t.Fatal("Can't connect to freenode.")
	}
	irccon.AddCallback("001", func(e *Event) { irccon.Join("#go-eventirc") })

	irccon.AddCallback("366", func(e *Event) {
		irccon.Privmsg("#go-eventirc", "Test Message\n")
		irccon.Nick("go-eventnewnick")
	})
	irccon.AddCallback("NICK", func(e *Event) {
		irccon.Quit()
		if irccon.nickcurrent == "go-eventnewnick" {
			t.Fatal("Nick change did not work!")
		}
	})
	irccon.Loop()
}

func TestConnectionSSL(t *testing.T) {
	irccon := IRC("go-eventirc", "go-eventirc")
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	irccon.UseTLS = true
	err := irccon.Connect("irc.freenode.net:7000")
	if err != nil {
		t.Fatal("Can't connect to freenode.")
	}
	irccon.AddCallback("001", func(e *Event) { irccon.Join("#go-eventirc") })

	irccon.AddCallback("366", func(e *Event) {
		irccon.Privmsg("#go-eventirc", "Test Message\n")
		time.Sleep(2 * time.Second)
		irccon.Quit()
	})

	irccon.Loop()
}

func TestRemoveCallback(t *testing.T) {
	irccon := IRC("go-eventirc", "go-eventirc")
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true

	done := make(chan int, 10)

	irccon.AddCallback("TEST", func(e *Event) { done <- 1 })
	id := irccon.AddCallback("TEST", func(e *Event) { done <- 2 })
	irccon.AddCallback("TEST", func(e *Event) { done <- 3 })

	// Should remove callback at index 1
	irccon.RemoveCallback("TEST", id)

	irccon.RunCallbacks(&Event{
		Code: "TEST",
	})

	var results []int

	results = append(results, <-done)
	results = append(results, <-done)

	if len(results) != 2 || !(results[0] == 1 && results[1] == 3) {
		t.Error("Callback 2 not removed")
	}
}

func TestWildcardCallback(t *testing.T) {
	irccon := IRC("go-eventirc", "go-eventirc")
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true

	done := make(chan int, 10)

	irccon.AddCallback("TEST", func(e *Event) { done <- 1 })
	irccon.AddCallback("*", func(e *Event) { done <- 2 })

	irccon.RunCallbacks(&Event{
		Code: "TEST",
	})

	var results []int

	results = append(results, <-done)
	results = append(results, <-done)

	if len(results) != 2 || !(results[0] == 1 && results[1] == 2) {
		t.Error("Wildcard callback not called")
	}
}

func TestClearCallback(t *testing.T) {
	irccon := IRC("go-eventirc", "go-eventirc")
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true

	done := make(chan int, 10)

	irccon.AddCallback("TEST", func(e *Event) { done <- 0 })
	irccon.AddCallback("TEST", func(e *Event) { done <- 1 })
	irccon.ClearCallback("TEST")
	irccon.AddCallback("TEST", func(e *Event) { done <- 2 })
	irccon.AddCallback("TEST", func(e *Event) { done <- 3 })

	irccon.RunCallbacks(&Event{
		Code: "TEST",
	})

	var results []int

	results = append(results, <-done)
	results = append(results, <-done)

	if len(results) != 2 || !(results[0] == 2 && results[1] == 3) {
		t.Error("Callbacks not cleared")
	}
}

func TestIRCemptyNick(t *testing.T) {
	irccon := IRC("", "go-eventirc")
	irccon = nil
	if nil != irccon {
		t.Error("empty nick didn't result in error")
		t.Fail()
	}
/*
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	irccon.UseTLS = true
	err := irccon.Connect("irc.freenode.net:7000")
	if err != nil {
		t.Fatal("Can't connect to freenode.")
	}
	irccon.AddCallback("001", func(e *Event) { irccon.Join("#go-eventirc") })

	irccon.AddCallback("366", func(e *Event) {
		irccon.Privmsg("#go-eventirc", "Test Message\n")
		time.Sleep(2 * time.Second)
		irccon.Quit()
	})

	irccon.Loop()
*/
}
 
func TestIRCemptyUser(t *testing.T) {
	irccon := IRC("go-eventirc", "")
	if nil != irccon {
		t.Error("empty user didn't result in error")
	}
}


func TestHasConnectionValues0(t *testing.T) {
	irccon := IRC("go-eventirc", "go-eventirc")
	if nil == irccon {
		t.Error("creating IRC struct failed")
	}

	irccon.server = "foo"
	if false == irccon.hasConnectionValues() {
		t.Error("valid struct not detected as such")
	}
}

func TestHasConnectionValues1(t *testing.T) {
	irccon := IRC("go-eventirc", "go-eventirc")
	if nil == irccon {
		t.Error("creating IRC struct failed")
	}
	irccon.server = "foo"

	irccon.Version = ""
	if irccon.hasConnectionValues() {
		t.Error("empty 'Version' not detected")
	}
}

func TestHasValidValues2(t *testing.T) {
	irccon := IRC("go-eventirc", "go-eventirc")
	if nil == irccon {
		t.Error("creating IRC struct failed")
	}
	irccon.server = "foo"

	irccon.nick = ""
	if irccon.hasConnectionValues() {
		t.Error("empty 'nick' not detected")
	}
}

func TestHasConnectionValues3(t *testing.T) {
	irccon := IRC("go-eventirc", "go-eventirc")
	if nil == irccon {
		t.Error("creating IRC struct failed")
	}
	irccon.server = "foo"

	irccon.user = ""
	if irccon.hasConnectionValues() {
		t.Error("empty 'user' not detected")
	}
}

func TestHasConnectionValues4(t *testing.T) {
	irccon := IRC("go-eventirc", "go-eventirc")
	if nil == irccon {
		t.Error("creating IRC struct failed")
	}
	irccon.server = "foo"

	irccon.Log = nil
	if irccon.hasConnectionValues() {
		t.Error("nil pointer 'Log' not detected")
	}
}

func TestHasConnectionValues5(t *testing.T) {
	irccon := IRC("go-eventirc", "go-eventirc")
	if nil == irccon {
		t.Error("creating IRC struct failed")
	}

	if irccon.hasConnectionValues() {
		t.Error("empty 'server' not detected")
	}
}
