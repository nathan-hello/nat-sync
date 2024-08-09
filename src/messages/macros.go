package messages

import "github.com/nathan-hello/nat-sync/src/utils"

var macros map[string][]byte = map[string][]byte{
	"room": []byte("0 room action=create name=asdf"),
	"test": []byte(`1 change uri=/mnt/hdd/media/cats/exercise.mp4;
		        1 wait secs=1;
		        1 pause;
		        1 wait secs=1;
		        1 play;
		        1 wait secs=1;
		        1 stop;
		        1 wait secs=1;`),
	"youtube": []byte(`1 change uri=https://www.youtube.com/watch?v=snYu2JUqSWs;
		         1 wait     secs=5;
		         1 stop    ;`),
	"j": []byte("0 join --roomname=cats username=nate;"),
}

func IsMacro(s string) []Message {
	text, ok := macros[s]
	if !ok {
		return nil
	}

	msgs, err := NewMulti(text)
	if err != nil {
		utils.ErrorLogger.Printf("macro was found but did not work. macro: %s, err: %s", text, err)
		return nil
	}

	return msgs
}
