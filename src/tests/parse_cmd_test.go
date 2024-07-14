package tests

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/nathan-hello/nat-sync/src/commands"
	"github.com/nathan-hello/nat-sync/src/commands/impl"
)

func TestCmdStringParsing(t *testing.T) {
	happy := map[string]commands.SubCommand{
		"change --uri=asdf.com/cats --action=append --hours=23 --mins=51 --secs=12": &impl.Change{Uri: "asdf.com/cats", UriLength: 13, Action: impl.ChgAppend, Timestamp: impl.Seek{Hours: 23, Mins: 51, Secs: 12}},
		"join   --roomid=34129":                              &impl.Join{RoomId: uint16(34129)},
		"kick   --userid=2182 --isself=false --hidemsg=true": &impl.Kick{UserId: 2182, IsSelf: false, HideMsg: true},
		"pause  ":                                &impl.Pause{},
		"play   ":                                &impl.Play{},
		"seek   --hours=20 --mins=40 --secs=100": &impl.Seek{Hours: 20, Mins: 40, Secs: 100},
	}

	for k, v := range happy {
		// fmt.Println("testing string: ", k)
		cmd, err := commands.CmdFromString(k)
		if err != nil {
			t.Fatalf("\nCmdFromString error: %s\nstring: %s", err, k)
		}
		if !reflect.DeepEqual(cmd.Sub, v) {
			t.Fatalf("\nsubcommand from CmdFromString() does not match test case. \nstring: %s\nresult: %#v\nexpect: %#v", k, cmd.Sub, v)
		}
		fmt.Println("string good   : ", k)
	}

}

func TestFromBits(t *testing.T) {
	empties := []commands.SubCommand{
		&impl.Change{},
		&impl.Join{},
		&impl.Kick{},
		&impl.Pause{},
		&impl.Play{},
		&impl.Seek{},
	}

	subs := []commands.SubCommand{
		&impl.Change{Uri: "asdf.com/cats", UriLength: 13, Action: impl.ChgAppend, Timestamp: impl.Seek{Hours: 23, Mins: 51, Secs: 12}},
		&impl.Join{RoomId: uint16(34129)},
		&impl.Kick{UserId: 2182, IsSelf: false, HideMsg: true},
		&impl.Pause{},
		&impl.Play{},
		&impl.Seek{Hours: 20, Mins: 40, Secs: 100},
	}

	for i, v := range subs {
		fmt.Printf("testing %#v\n", empties[i])
		b, err := v.ToBits()
		if err != nil {
			t.Fatalf("err in subs index %d\nerr: %s", i, err)
		}
		// fmt.Printf("sub: %#v\ntobits:%#v\n", v, b)
		err = empties[i].FromBits(b)
		if err != nil {
			t.Fatalf("err in frombits(): %#v\n", err)
		}
		if !reflect.DeepEqual(empties[i], subs[i]) {
			t.Fatalf("s frombits does not equal expected val\nfrombits result: %#v\nexpected: %#v\n", empties[i], subs[i])
		}
		// fmt.Printf("frombits success: %#v\n", empties[i])
	}
}
