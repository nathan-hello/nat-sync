package messages

import (
	"github.com/nathan-hello/nat-sync/src/utils"
)

func IsMacro(s string) []Message {
	var roomId int64 = 1
	switch s {
	case "cat":
		cat, err := New("change roomid=1 uri=/mnt/hdd/media/cats/exercise.mp4", nil)
		if err != nil {
			utils.ErrorLogger.Printf("Error creating command 'change': %v\n", err)
			return nil
		}
		return cat

	case "test":
		test, err := New(
			`change roomid=1 uri=/mnt/hdd/media/cats/exercise.mp4;
		        wait  roomid=1 secs=1;
		        pause roomid=1;
		        wait  roomid=1 secs=1;
		        play  roomid=1;
		        wait  roomid=1 secs=1;
		        stop  roomid=1;
		        wait  roomid=1 secs=1;`, nil)
		if err != nil {
			utils.ErrorLogger.Println(err)
			return nil
		}
		return test

	case "testyt":
		testyt, err := New(
			`change roomid=1 uri=https://www.youtube.com/watch?v=snYu2JUqSWs;
		        wait roomid=1  secs=5;
		        stop roomid=1;
                `, &roomId)
		if err != nil {
			utils.ErrorLogger.Println(err)
			return nil
		}
		return testyt
	case "j":
		testyt, err := New(`join --roomid=1 username=nate;`, nil)
		if err != nil {
			utils.ErrorLogger.Println(err)
			return nil
		}
		return testyt
	}
	return nil
}
