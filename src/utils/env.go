package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// Config is generated when the program launches.
// It should never be called until InitConfig() is done.
var initialized bool

type ServerRoom struct {
	Id       int64
	Name     string
	Password string
	Clients  map[int64]Client
}

type Client struct {
	Id   int64
	Name string
	Conn net.Conn
}

type ClientSavedRoom struct {
	ServerIp          string `json:"server_ip"`
	RoomId            int64  `json:"room_id"`
	RoomName          string `json:"room_name"`
	Username          string `json:"username"`
	EncryptedPassword string `json:"encrypted_password"`
	Autoconnect       bool   `json:"connect_automatically"`
}

type Server struct {
	Port  string       `json:"port"`
	Rooms []ServerRoom `json:"rooms"`
}

type argsCommandLine struct {
	StartClient bool
	ConnectIp   string
	Username    string
	StartServer bool
	ServerPort  string
}

type argsConfigFile struct {
	Client Client `json:"client"`
	Server Server `json:"server"`
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

func ParseArgs() (*argsCommandLine, error) {
	asdf := &argsCommandLine{}
	fs := flag.NewFlagSet("nat-sync", flag.ContinueOnError)

	connect := fs.String("connect", "", "IP address to connect to")
	connectShort := fs.String("c", "", "IP address to connect to (shorthand)")

	username := fs.String("username", "", "Username")
	usernameShort := fs.String("u", "", "Username (shorthand)")

	noClient := fs.Bool("no-client", false, "Disable client start")
	noClientShort := fs.Bool("n", false, "Disable client start (shorthand)")

	startServer := fs.Bool("start-server", false, "Start the server")
	startServerShort := fs.Bool("s", false, "Start the server (shorthand)")

	serverPort := fs.String("server-port", "", "Port for the server")
	serverPortShort := fs.String("p", "", "Port for the server (shorthand)")

	help := fs.Bool("help", false, "Print this menu")
	helpShort := fs.Bool("h", false, "Print this menu")

	if *help || *helpShort || len(fs.Args()) == 0 {
		fs.PrintDefaults()
		return nil, nil
	}

	// Use the value from either the long or short flag
	connectVal := *connect
	if *connectShort != "" {
		connectVal = *connectShort
	}

	var ip string
	if connectVal != "" {
		flag := strings.Trim(connectVal, "\"")
		ok := net.ParseIP(flag)
		if ok == nil {
			fmt.Printf("IP is not valid. (did you try adding :<port> without --port?) IP: %s\n", flag)
		}
		ip = flag
	}

	serverPortVal := *serverPort
	if *serverPortShort != "" {
		serverPortVal = *serverPortShort
	}

	user := *username
	if *usernameShort != "" {
		user = *usernameShort
	}

	var sp string
	if serverPortVal != "" {
		flag := strings.Trim(serverPortVal, "\"")
		flag = strings.TrimPrefix(flag, ":")
		_, err := strconv.Atoi(flag)
		if err != nil {
			return nil, fmt.Errorf("port is not valid. port: %s", flag)
		}

		sp = flag
	}

	asdf.ServerPort = sp
	asdf.StartClient = !*noClient && !*noClientShort
	asdf.StartServer = *startServer || *startServerShort
	asdf.ConnectIp = ip
	asdf.Username = user

	return asdf, nil
}

func randomString(length int) string {
	randomBytes := make([]byte, length)
	rand.Read(randomBytes)
	randomString := base64.URLEncoding.EncodeToString(randomBytes)[:length]
	return randomString
}

func readConfigFile() {

}
