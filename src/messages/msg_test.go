package messages

import (
	"reflect"
	"slices"
	"testing"

	"github.com/nathan-hello/nat-sync/src/messages/impl"
	"github.com/nathan-hello/nat-sync/src/utils"
)

var roomId int64 = 1

func TestCmdStrings(t *testing.T) {
	utils.InitLogger()
	happy := map[string]impl.Command{
		"1 change --uri=asdf.com/cats --action=append --hours=23 --mins=51 --secs=12": &impl.Change{Uri: "asdf.com/cats", UriLength: 13, Action: impl.ChgAppend, Timestamp: impl.Seek{Hours: 23, Mins: 51, Secs: 12}},
		"1 kick   --userid=2182 --isself=false --hidemsg=true":                        &impl.Kick{UserId: 2182, IsSelf: false, HideMsg: true},
		"1 seek --hours=10 --mins=40 --secs=100":                                      &impl.Seek{Hours: 10, Mins: 40, Secs: 100},
		"1 pause":                                                                     &impl.Pause{},
		"1 play ":                                                                     &impl.Play{},
	}

	for k, v := range happy {
		cmd := Message{}
		t.Log("testing string: ", k)
		err := cmd.TextUnmarshaller([]byte(k))

		if err != nil {
			t.Fatalf("\nNewCmdFromString error: %s\nstring: %s", err, k)
		}

		// [0] because we are assuming that the above cmds don't have ; in them for multi commands
		if !reflect.DeepEqual(cmd.Sub, v) {
			t.Fatalf("\nsubcommand from NewCmdFromString() does not match test case. \nstring: %s\nresult: %#v\nexpect: %#v", k, cmd.Sub, v)
		}
		t.Log("string good   : ", k)
	}

}

func TestBits(t *testing.T) {
	utils.InitLogger()
	empties := []impl.Command{
		&impl.Change{},
		// &impl.Join{},
		&impl.Kick{},
		&impl.Pause{},
		&impl.Play{},
		&impl.Seek{},
	}

	subs := []impl.Command{
		&impl.Change{Uri: "asdf.com/cats", UriLength: 13, Action: impl.ChgAppend, Timestamp: impl.Seek{Hours: 23, Mins: 51, Secs: 12}},
		// &impl.Join{RoomId: int64(34129)},
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
		err = empties[i].New(b)
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
	utils.InitLogger()
	subs := []impl.Command{
		&impl.Change{Uri: "asdf.com/cats", UriLength: 13, Action: impl.ChgAppend, Timestamp: impl.Seek{Hours: 23, Mins: 51, Secs: 12}},
		// &impl.Join{RoomId: int64(34129)},
		&impl.Kick{UserId: 2182, IsSelf: false, HideMsg: true},
		&impl.Pause{},
		&impl.Play{},
		&impl.Seek{Hours: 10},
	}

	for _, v := range subs {
		h, _ := getHeadFromString(v.GetHead())
		cmd := Message{
			Head:    h,
			Version: utils.CurrentVersion,
		}
		bits, err := v.ToBits()
		if err != nil {
			t.Fatalf("err tobits(): %s", err)
		}

		cmd.Content = bits

		cmdBits, err := cmd.MarshalBinary()
		if err != nil {
			t.Fatalf("err encodecommand(): %s", err)
		}

		cmd.Sub = v // to verify that it's equal. can't put it before the encode process

		newCmd := Message{}
		err = newCmd.UnmarshalBinary(cmdBits)
		if err != nil {
			t.Fatalf("err decodecommand(): %s", err)
		}

		t.Logf("\n\ncmd: %#v\nnewCmd: %#v\n", cmd, newCmd)
		t.Logf("\n\ncmdSub: %#v\nnewCmd.Sub: %#v\n", cmd.Sub, newCmd.Sub)

		assert(newCmd.Head == cmd.Head, t, "head no match")
		assert(newCmd.Version == cmd.Version, t, "version no match")
		assert(slices.Equal(newCmd.Content, cmd.Content), t, "content no match")
		assert(reflect.DeepEqual(newCmd.Sub, cmd.Sub), t, "sub struct no match")

	}

}
