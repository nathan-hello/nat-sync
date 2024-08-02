package messages

import (
	"github.com/nathan-hello/nat-sync/src/messages/commands"
	"github.com/nathan-hello/nat-sync/src/utils"
)

func getMacro(s string) []Message {
	var msgs []Message
	switch s {
	case "test":
		change, err := commands.New("change uri=/mnt/hdd/media/cats/exercise.mp4")
		if err != nil {
			utils.ErrorLogger.Printf("Error creating command 'change': %v\n", err)
			return nil
		}
		wait1, err := commands.New("wait secs=1")
		if err != nil {
			utils.ErrorLogger.Printf("Error creating command 'wait1': %v\n", err)
			return nil
		}
		pause, err := commands.New("pause")
		if err != nil {
			utils.ErrorLogger.Printf("Error creating command 'pause': %v\n", err)
			return nil
		}
		wait2, err := commands.New("wait secs=1")
		if err != nil {
			utils.ErrorLogger.Printf("Error creating command 'wait2': %v\n", err)
			return nil
		}
		play, err := commands.New("play")
		if err != nil {
			utils.ErrorLogger.Printf("Error creating command 'play': %v\n", err)
			return nil
		}
		wait3, err := commands.New("wait secs=1")
		if err != nil {
			utils.ErrorLogger.Printf("Error creating command 'wait3': %v\n", err)
			return nil
		}
		stop, err := commands.New("stop")
		if err != nil {
			utils.ErrorLogger.Printf("Error creating command 'stop': %v\n", err)
			return nil
		}
		wait4, err := commands.New("wait secs=1")
		if err != nil {
			utils.ErrorLogger.Printf("Error creating command 'wait4': %v\n", err)
			return nil
		}
		return append(msgs, change, wait1, pause, wait2, play, wait3, stop, wait4)

	case "testyt":
		change, err := commands.New("change uri=https://www.youtube.com/watch?v=snYu2JUqSWs")
		if err != nil {
			utils.ErrorLogger.Printf("Error creating command 'change': %v\n", err)
			return nil
		}
		wait1, err := commands.New("wait secs=5")
		if err != nil {
			utils.ErrorLogger.Printf("Error creating command 'wait1': %v\n", err)
			return nil
		}
		stop, err := commands.New("stop")
		if err != nil {
			utils.ErrorLogger.Printf("Error creating command 'stop': %v\n", err)
			return nil
		}
		return append(msgs, change, wait1, stop)
	}
	return nil
}
