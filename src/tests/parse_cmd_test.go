package tests

import (
	"reflect"
	"slices"
	"testing"

	"github.com/nathan-hello/nat-sync/src/commands"
	"github.com/nathan-hello/nat-sync/src/commands/impl"
)

func TestCmdStrings(t *testing.T) {
	happy := map[string]commands.SubCommand{
		"change --uri=asdf.com/cats --action=append --hours=23 --mins=51 --secs=12": &impl.Change{Uri: "asdf.com/cats", UriLength: 13, Action: impl.ChgAppend, Timestamp: impl.Seek{Hours: 23, Mins: 51, Secs: 12}},
		"join   --roomid=34129":                              &impl.Join{RoomId: uint16(34129)},
		"kick   --userid=2182 --isself=false --hidemsg=true": &impl.Kick{UserId: 2182, IsSelf: false, HideMsg: true},
		"pause  ":                                &impl.Pause{},
		"play   ":                                &impl.Play{},
		"seek   --hours=10 --mins=40 --secs=100": &impl.Seek{Hours: 10, Mins: 40, Secs: 100},
	}

	for k, v := range happy {
		// t.Log("testing string: ", k)
		cmd, err := commands.CmdFromString(k)
		if err != nil {
			t.Fatalf("\nCmdFromString error: %s\nstring: %s", err, k)
		}
		if !reflect.DeepEqual(cmd.Sub, v) {
			t.Fatalf("\nsubcommand from CmdFromString() does not match test case. \nstring: %s\nresult: %#v\nexpect: %#v", k, cmd.Sub, v)
		}
		t.Log("string good   : ", k)
	}

}

func TestBits(t *testing.T) {
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
		&impl.Seek{Hours: 10},
	}

	for i, v := range subs {
		t.Logf("testing %#v\n", empties[i])
		b, err := v.ToBits()
		if err != nil {
			t.Fatalf("err in subs index %d\nerr: %s", i, err)
		}
		// t.Logf("sub: %#v\ntobits:%#v\n", v, b)
		err = empties[i].FromBits(b)
		if err != nil {
			t.Fatalf("err in frombits(): %#v\n", err)
		}
		if !reflect.DeepEqual(empties[i], subs[i]) {
			t.Fatalf("s frombits does not equal expected val\nfrombits result: %#v\nexpected: %#v\n", empties[i], subs[i])
		}
		t.Logf("frombits success: %#v\n", empties[i])
	}
}

func assert(test bool, t *testing.T, msg string) {
	if !test {
		t.Fatal(msg)
	}
}

func TestEncodeCmd(t *testing.T) {
	subs := map[commands.CmdHead]commands.SubCommand{
		commands.ChangeHead: &impl.Change{Uri: "asdf.com/cats", UriLength: 13, Action: impl.ChgAppend, Timestamp: impl.Seek{Hours: 23, Mins: 51, Secs: 12}},
		commands.JoinHead:   &impl.Join{RoomId: uint16(34129)},
		commands.KickHead:   &impl.Kick{UserId: 2182, IsSelf: false, HideMsg: true},
		commands.PauseHead:  &impl.Pause{},
		commands.PlayHead:   &impl.Play{},
		commands.SeekHead:   &impl.Seek{Hours: 10},
	}

	for k, v := range subs {
		cmd := commands.Command{
			Head:    k,
			Version: commands.CurrentVersion,
			UserId:  8231,
		}
		bits, err := v.ToBits()
		if err != nil {
			t.Fatalf("err tobits(): %s", err)
		}

		cmd.Content = bits
		cmdBits, err := commands.EncodeCommand(&cmd)
		if err != nil {
			t.Fatalf("err encodecommand(): %s", err)
		}

		cmd.Sub = v // to verify that it's equal. can't put it before the encode process

		newCmd, err := commands.DecodeCommand(cmdBits)
		if err != nil {
			t.Fatalf("err decodecommand(): %s", err)
		}

		t.Logf("\n\ncmd: %#v\nnewCmd: %#v\n", cmd, newCmd)
		t.Logf("\n\ncmdSub: %#v\nnewCmd.Sub: %#v\n", cmd.Sub, newCmd.Sub)
		assert(newCmd.Head == cmd.Head, t, "head no match")
		assert(newCmd.Version == cmd.Version, t, "version no match")
		assert(newCmd.UserId == cmd.UserId, t, "userid no match")
		assert(slices.Equal(newCmd.Content, cmd.Content), t, "content no match")
		assert(reflect.DeepEqual(newCmd.Sub, cmd.Sub), t, "sub struct no match")

	}

}