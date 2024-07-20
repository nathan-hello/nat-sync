package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"net"
	"os"
	"strconv"
	"strings"
)

// Config is generated when the program launches.
// It should never be called until InitConfig() is done.
var config FullConfig
var initialized bool

type Room struct {
	ServerIp          string `json:"server_ip"`
	RoomName          string `json:"room_name"`
	Username          string `json:"username"`
	EncryptedPassword string `json:"encrypted_password"`
	Autoconnect       bool   `json:"connect_automatically"`
}

type Client struct {
	Rooms []Room `json:"rooms"`
}

type Server struct {
	Port string `json:"port"`
}

type argsCommandLine struct {
	StartClient bool
	ConnectIp   string
	Username    string
	StartServer bool
	ServerPort  int
}

type argsConfigFile struct {
	Client Client `json:"client"`
	Server Server `json:"server"`
}

type FullConfig struct {
	StartClient bool
	ConnectIp   string
	ConnectPort int
	Username    string
	StartServer bool
	ServerPort  int
}

func ReadConfig(configPath string) *argsConfigFile {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil
	}
	var conf argsConfigFile
	if err := json.Unmarshal(file, &conf); err != nil {
		return nil
	}

	return &conf
}

func readCmdLineArgs(flags *flag.FlagSet) *argsCommandLine {
	asdf := &argsCommandLine{
		StartClient: true,
		Username:    randomString(8),
		ConnectIp:   "127.0.0.1",
		StartServer: false,
		ServerPort:  1412,
	}
	var args []string
	for _, v := range args {
		v = strings.ToLower(v)
		v = strings.TrimPrefix(v, "-")
		v = strings.TrimPrefix(v, "-")
		switch {
		case strings.HasPrefix(v, "connect="):
			flag, _ := strings.CutPrefix(v, "connect=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			ok := net.ParseIP(flag)
			if ok == nil {
				UserLogger.Printf("ip is not valid. (did you try adding :<port> without --port?) ip: %s", flag)
			}
			asdf.ConnectIp = flag
		case strings.HasPrefix(v, "no-client"):
			asdf.StartClient = false
		case strings.HasPrefix(v, "start-server"):
			asdf.StartServer = true
		case strings.HasPrefix(v, "server-port"):
			flag, _ := strings.CutPrefix(v, "server-port=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			i, err := strconv.Atoi(flag)
			if err != nil {
				UserLogger.Printf("port is not valid. port: %s", flag)
			}
			asdf.ServerPort = i

		}

	}
	return asdf
}

func randomString(length int) string {
	randomBytes := make([]byte, length)
	rand.Read(randomBytes)
	randomString := base64.URLEncoding.EncodeToString(randomBytes)[:length]
	return randomString
}

func Config() *FullConfig {
	return &config
}

func readConfigFile() {

}
