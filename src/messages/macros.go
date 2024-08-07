package messages

import (
	"github.com/nathan-hello/nat-sync/src/utils"
)

func IsMacro(s string) []Message {
	var roomId int64 = 1
	switch s {
	case "cat":
		cat, err := New("roomid=1; change uri=/mnt/hdd/media/cats/exercise.mp4", nil)
		if err != nil {
			utils.ErrorLogger.Printf("Error creating command 'change': %v\n", err)
			return nil
		}
		return cat

	case "test":
		test, err := New(
			`roomid=1; change uri=/mnt/hdd/media/cats/exercise.mp4;
		        wait     secs=1;
		        pause   ;
		        wait     secs=1;
		        play    ;
		        wait     secs=1;
		        stop    ;
		        wait     secs=1;`, nil)
		if err != nil {
			utils.ErrorLogger.Println(err)
			return nil
		}
		return test

	case "testyt":
		testyt, err := New(
			`roomid=1; change uri=https://www.youtube.com/watch?v=snYu2JUqSWs;
		        wait     secs=5;
		        stop    ;
                `, &roomId)
		if err != nil {
			utils.ErrorLogger.Println(err)
			return nil
		}
		return testyt
	case "j":
		testyt, err := New(`roomid=1; join --roomname=cats username=nate;`, nil)
		if err != nil {
			utils.ErrorLogger.Println(err)
			return nil
		}
		return testyt
	}
	return nil
}
