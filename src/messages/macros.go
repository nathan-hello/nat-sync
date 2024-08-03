package messages

import (
	"github.com/nathan-hello/nat-sync/src/utils"
)

func IsMacro(s string) []Message {
	switch s {
	case "cat":
		cat, err := New("change uri=/mnt/hdd/media/cats/exercise.mp4")
		if err != nil {
			utils.ErrorLogger.Printf("Error creating command 'change': %v\n", err)
			return nil
		}
		return cat

	case "test":
		test, err := New(
			`change uri=/mnt/hdd/media/cats/exercise.mp4;
		        wait secs=1;
		        pause;
		        wait secs=1;
		        play;
		        wait secs=1;
		        stop;
		        wait secs=1;`)
		if err != nil {
			utils.ErrorLogger.Println(err)
			return nil
		}
		return test

	case "testyt":
		testyt, err := New(
			`change uri=https://www.youtube.com/watch?v=snYu2JUqSWs;
		        wait secs=5;
		        stop;
                `)
		if err != nil {
			utils.ErrorLogger.Println(err)
			return nil
		}
		return testyt
	}
	return nil
}
