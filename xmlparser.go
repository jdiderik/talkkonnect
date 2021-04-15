/*
 * talkkonnect headless mumble client/gateway with lcd screen and channel control
 * Copyright (C) 2018-2019, Suvir Kumar <suvir@talkkonnect.com>
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * talkkonnect is the based on talkiepi and barnard by Daniel Chote and Tim Cooper
 *
 * The Initial Developer of the Original Code is
 * Suvir Kumar <suvir@talkkonnect.com>
 * Portions created by the Initial Developer are Copyright (C) Suvir Kumar. All Rights Reserved.
 *
 * Contributor(s):
 *
 * Suvir Kumar <suvir@talkkonnect.com>
 * Zoran Dimitrijevic
 * My Blog is at www.talkkonnect.com
 * The source code is hosted at github.com/talkkonnect
 *
 * xmlparser.go -> talkkonnect functionality to read from XML file and populate global variables
 */

package talkkonnect

import (
	"encoding/xml"
	"fmt"
	"github.com/jdiderik/go-openal/openal"
	"github.com/jdiderik/gumble/gumbleffmpeg"
	"golang.org/x/sys/unix"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//version and release date
const (
	talkkonnectVersion  string = "1.63.01"
	talkkonnectReleased string = "April 15 2021"
)

// Generic Global Variables
var (
	pstream               *gumbleffmpeg.Stream
	AccountCount          int  = 0
	KillHeartBeat         bool = false
	IsPlayStream          bool = false
	BackLightTime              = time.NewTicker(5 * time.Second)
	BackLightTimePtr           = &BackLightTime
	ConnectAttempts            = 0
	IsConnected           bool = false
	source                     = openal.NewSource()
	StartTime                  = time.Now()
	BufferToOpenALCounter      = 0
	AccountIndex          int  = 0
)

//account settings
var (
	Default     []bool
	Name        []string
	Server      []string
	Username    []string
	Password    []string
	Insecure    []bool
	Certificate []string
	Channel     []string
	Ident       []string
)

//software settings
var (
	OutputDevice       string = "Speaker"
	OutputDeviceShort  string
	LogFilenameAndPath string = "/var/log/talkkonnect.log"
	Logging            string = "screen"
	Loglevel           string = "info"
	Daemonize          bool
	SimplexWithMute    bool = true
	TxCounter          bool
	NextServerIndex    int = 0
)

//autoprovision settings
var (
	APEnabled    bool
	TkID         string
	URL          string
	SaveFilePath string
	SaveFilename string
)

//sound settings
var (
	EventSoundEnabled                 bool
	EventJoinedSoundFilenameAndPath   string
	EventLeftSoundFilenameAndPath     string
	EventMessageSoundFilenameAndPath  string
	AlertSoundEnabled                 bool
	AlertSoundFilenameAndPath         string
	AlertSoundVolume                  float32 = 1
	IncommingBeepSoundEnabled         bool
	IncommingBeepSoundFilenameAndPath string
	IncommingBeepSoundVolume          float32
	RogerBeepSoundEnabled             bool
	RogerBeepSoundFilenameAndPath     string
	RogerBeepSoundVolume              float32
	RepeaterToneEnabled               bool
	RepeaterToneFrequencyHz           int
	RepeaterToneDurationSec           int
	StreamSoundEnabled                bool
	StreamSoundFilenameAndPath        string
	StreamSoundVolume                 float32
)

//api settings
var (
	APIEnabled            bool
	APIListenPort         string
	APIDisplayMenu        bool
	APIChannelUp          bool
	APIChannelDown        bool
	APIMute               bool
	APICurrentVolumeLevel bool
	APIDigitalVolumeUp    bool
	APIDigitalVolumeDown  bool
	APIListServerChannels bool
	APIStartTransmitting  bool
	APIStopTransmitting   bool
	APIListOnlineUsers    bool
	APIPlayStream         bool
	APIRequestGpsPosition bool
	APIEmailEnabled       bool
	APINextServer         bool
	APIPreviousServer     bool
	APIPanicSimulation    bool
	APIScanChannels       bool
	APIDisplayVersion     bool
	APIClearScreen        bool
	APIPingServersEnabled bool
	APIRepeatTxLoopTest   bool
	APIPrintXmlConfig     bool
)


// mqtt settings
var (
	MQTTEnabled     bool = false
	Iotuuid         string
	relay1State     bool = false
	relayAllState   bool = false
	RelayPulseMills time.Duration
	TotalRelays     uint
	RelayPins       = [9]uint{}
	MQTTTopic       string
	MQTTBroker      string
	MQTTPassword    string
	MQTTUser        string
	MQTTId          string
	MQTTCleansess   bool
	MQTTQos         int
	MQTTNum         int
	MQTTPayload     string
	MQTTAction      string
	MQTTStore       string
)

// target board settings
var (
	TargetBoard string = "pc"
)

//txtimeout settings
var (
	TxTimeOutEnabled bool
	TxTimeOutSecs    int
)


//other global variables used for state tracking
var (
	txcounter         int
	togglecounter     int
	isTx              bool
	isPlayStream      bool
	CancellableStream bool = true
)

type Document struct {
	XMLName  xml.Name `xml:"document"`
	Accounts struct {
		Account []struct {
			Name          string `xml:"name,attr"`
			Default       bool   `xml:"default,attr"`
			ServerAndPort string `xml:"serverandport"`
			UserName      string `xml:"username"`
			Password      string `xml:"password"`
			Insecure      bool   `xml:"insecure"`
			Certificate   string `xml:"certificate"`
			Channel       string `xml:"channel"`
			Ident         string `xml:"ident"`
		} `xml:"account"`
	} `xml:"accounts"`
	Global struct {
		Software struct {
			Settings struct {
				OutputDevice       string `xml:"outputdevice"`
				OutputDeviceShort  string `xml:"outputdeviceshort"`
				LogFilenameAndPath string `xml:"logfilenameandpath"`
				Logging            string `xml:"logging"`
				Loglevel           string `xml:"loglevel"`
				Daemonize          bool   `xml:"daemonize"`
				CancellableStream  bool   `xml:"cancellablestream"`
				SimplexWithMute    bool   `xml:"simplexwithmute"`
				TxCounter          bool   `xml:"txcounter"`
				NextServerIndex    int    `xml:"nextserverindex"`
			} `xml:"settings"`
			Sounds struct {
				Event struct {
					Enabled                bool   `xml:"enabled,attr"`
					JoinedFilenameAndPath  string `xml:"joinedfilenameandpath"`
					LeftFilenameAndPath    string `xml:"leftfilenameandpath"`
					MessageFilenameAndPath string `xml:"messagefilenameandpath"`
				} `xml:"event"`
				Alert struct {
					Enabled         bool    `xml:"enabled,attr"`
					FilenameAndPath string  `xml:"filenameandpath"`
					Volume          float32 `xml:"volume"`
				} `xml:"alert"`
				IncommingBeep struct {
					Enabled         bool    `xml:"enabled,attr"`
					FilenameAndPath string  `xml:"filenameandpath"`
					Volume          float32 `xml:"volume"`
				} `xml:"incommingbeep"`
				RogerBeep struct {
					Enabled         bool    `xml:"enabled,attr"`
					FilenameAndPath string  `xml:"filenameandpath"`
					Volume          float32 `xml:"volume"`
				} `xml:"rogerbeep"`
				RepeaterTone struct {
					Enabled         bool `xml:"enabled,attr"`
					ToneFrequencyHz int  `xml:"tonefrequencyhz"`
					ToneDurationSec int  `xml:"tonedurationsec"`
				} `xml:"repeatertone"`
				Stream struct {
					Enabled         bool    `xml:"enabled,attr"`
					FilenameAndPath string  `xml:"filenameandpath"`
					Volume          float32 `xml:"volume"`
				} `xml:"stream"`
			} `xml:"sounds"`
			TxTimeOut struct {
				Enabled       bool `xml:"enabled,attr"`
				TxTimeOutSecs int  `xml:"txtimeoutsecs"`
			} `xml:"txtimeout"`
			API struct {
				Enabled            bool   `xml:"enabled,attr"`
				ListenPort         string `xml:"apilistenport"`
				DisplayMenu        bool   `xml:"displaymenu"`
				ChannelUp          bool   `xml:"channelup"`
				ChannelDown        bool   `xml:"channeldown"`
				Mute               bool   `xml:"mute"`
				CurrentVolumeLevel bool   `xml:"currentvolumelevel"`
				DigitalVolumeUp    bool   `xml:"digitalvolumeup"`
				DigitalVolumeDown  bool   `xml:"digitalvolumedown"`
				ListServerChannels bool   `xml:"listserverchannels"`
				StartTransmitting  bool   `xml:"starttransmitting"`
				StopTransmitting   bool   `xml:"stoptransmitting"`
				ListOnlineUsers    bool   `xml:"listonlineusers"`
				PlayStream         bool   `xml:"playstream"`
				RequestGpsPosition bool   `xml:"requestgpsposition"`
				PreviousServer     bool   `xml:"previousserver"`
				NextServer         bool   `xml:"nextserver"`
				PanicSimulation    bool   `xml:"panicsimulation"`
				ScanChannels       bool   `xml:"scanchannels"`
				DisplayVersion     bool   `xml:"displayversion"`
				ClearScreen        bool   `xml:"clearscreen"`
				RepeatTxLoopTest   bool   `xml:"repeattxlooptest"`
				PrintXmlConfig     bool   `xml:"printxmlconfig"`
				SendEmail          bool   `xml:"sendemail"`
				PingServers        bool   `xml:"pingservers"`
			} `xml:"api"`
			MQTT struct {
				MQTTEnabled   bool   `xml:"enabled,attr"`
				MQTTTopic     string `xml:"mqtttopic"`
				MQTTBroker    string `xml:"mqttbroker"`
				MQTTPassword  string `xml:"mqttpassword"`
				MQTTUser      string `xml:"mqttuser"`
				MQTTId        string `xml:"mqttid"`
				MQTTCleansess bool   `xml:"cleansess"`
				MQTTQos       int    `xml:"qos"`
				MQTTNum       int    `xml:"num"`
				MQTTPayload   string `xml:"payload"`
				MQTTAction    string `xml:"action"`
				MQTTStore     string `xml:"store"`
			} `xml:"mqtt"`
		} `xml:"software"`
		Hardware struct {
			TargetBoard string `xml:"targetboard,attr"`
		} `xml:"hardware"`
	} `xml:"global"`
}

func readxmlconfig(file string) error {
	xmlFile, err := os.Open(file)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	log.Println("info: Successfully Read file " + filepath.Base(file))
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var document Document

	err = xml.Unmarshal(byteValue, &document)
	if err != nil {
		return fmt.Errorf(filepath.Base(file) + " " + err.Error())
	}

	for i := 0; i < len(document.Accounts.Account); i++ {
		if document.Accounts.Account[i].Default == true {
			Name = append(Name, document.Accounts.Account[i].Name)
			Server = append(Server, document.Accounts.Account[i].ServerAndPort)
			Username = append(Username, document.Accounts.Account[i].UserName)
			Password = append(Password, document.Accounts.Account[i].Password)
			Insecure = append(Insecure, document.Accounts.Account[i].Insecure)
			Certificate = append(Certificate, document.Accounts.Account[i].Certificate)
			Channel = append(Channel, document.Accounts.Account[i].Channel)
			Ident = append(Ident, document.Accounts.Account[i].Ident)
			AccountCount++
		}
	}

	if AccountCount == 0 {
		FatalCleanUp("No Default Accounts Found in talkkonnect.xml File! Please Add At Least 1 Default Account in XML")
	}

	exec, err := os.Executable()
	if err != nil {
		exec = "./talkkonnect" //Hardcode our default name
	}

	// Set our default config file path (for autoprovision)
	defaultConfPath, err := filepath.Abs(filepath.Dir(file))
	if err != nil {
		FatalCleanUp("Unable to get path for config file " + err.Error())
	}

	// Set our default logging path
	//This section is pretty unix specific.. sorry if you like windows support.
	defaultLogPath := "/tmp/" + filepath.Base(exec) + ".log" // Safe assumption as it should be writable for everyone
	// First see if we can write in our CWD and use it over /tmp
	cwd, err := os.Getwd()
	if err == nil {
		cwd, err := filepath.Abs(cwd)
		if err == nil {
			if unix.Access(cwd, unix.W_OK) == nil {
				defaultLogPath = cwd + "/" + filepath.Base(exec) + ".log"
			}
		}
	}

	// Next try a file in our config path and favor it over CWD
	if unix.Access(defaultConfPath, unix.W_OK) == nil {
		defaultLogPath = defaultConfPath + "/" + filepath.Base(exec) + ".log"
	}

	// Last, see if the system talkkonnect log exists and is writeable and do that over CWD, HOME and /tmp
	if _, err := os.Stat("/var/log/" + filepath.Base(exec) + ".log"); err == nil {
		f, err := os.OpenFile("/var/log/"+filepath.Base(exec)+".log", os.O_WRONLY, 0664)
		if err == nil {
			defaultLogPath = "/var/log/" + filepath.Base(exec) + ".log"
		}
		f.Close()
	}

	// Set our default sharefile path
	defaultSharePath := "/tmp"
	dir := filepath.Dir(exec)
	//Check for soundfiles directory in various locations
	// First, check env for $GOPATH and check in the hardcoded talkkonnect/talkkonnect dir
	if os.Getenv("GOPATH") != "" {
		defaultRepo := os.Getenv("GOPATH") + "/src/github.com/jdiderik/talkkonnect"
		if stat, err := os.Stat(defaultRepo); err == nil && stat.IsDir() {
			defaultSharePath = defaultRepo
		}
	}
	// Next, check the same dir as executable for 'soundfiles'
	if stat, err := os.Stat(dir + "/soundfiles"); err == nil && stat.IsDir() {
		defaultSharePath = dir
	}
	// Last, if its in a bin directory, we check for ../share/talkkonnect/ and prioritize it if it exists
	if strings.HasSuffix(dir, "bin") {
		shareDir := filepath.Dir(dir) + "/share/" + filepath.Base(exec)
		if stat, err := os.Stat(shareDir); err == nil && stat.IsDir() {
			defaultSharePath = shareDir
		}
	}

	OutputDevice = document.Global.Software.Settings.OutputDevice
	OutputDeviceShort = document.Global.Software.Settings.OutputDeviceShort

	if len(OutputDeviceShort) == 0 {
		OutputDeviceShort = document.Global.Software.Settings.OutputDevice
	}

	LogFilenameAndPath = document.Global.Software.Settings.LogFilenameAndPath
	Logging = document.Global.Software.Settings.Logging

	if document.Global.Software.Settings.Loglevel == "trace" || document.Global.Software.Settings.Loglevel == "debug" || document.Global.Software.Settings.Loglevel == "info" || document.Global.Software.Settings.Loglevel == "warning" || document.Global.Software.Settings.Loglevel == "error" || document.Global.Software.Settings.Loglevel == "alert" {
		Loglevel = document.Global.Software.Settings.Loglevel
	}

	if strings.ToLower(Logging) != "screen" && LogFilenameAndPath == "" {
		LogFilenameAndPath = defaultLogPath
	}

	Daemonize = document.Global.Software.Settings.Daemonize
	CancellableStream = document.Global.Software.Settings.CancellableStream
	SimplexWithMute = document.Global.Software.Settings.SimplexWithMute
	TxCounter = document.Global.Software.Settings.TxCounter
	NextServerIndex = document.Global.Software.Settings.NextServerIndex

	APEnabled = document.Global.Software.AutoProvisioning.Enabled
	TkID = document.Global.Software.AutoProvisioning.TkID
	URL = document.Global.Software.AutoProvisioning.URL
	SaveFilePath = document.Global.Software.AutoProvisioning.SaveFilePath
	SaveFilename = document.Global.Software.AutoProvisioning.SaveFilename

	if APEnabled && SaveFilePath == "" {
		SaveFilePath = defaultConfPath
	}

	if APEnabled && SaveFilename == "" {
		SaveFilename = filepath.Base(exec) + ".xml" //Should default to talkkonnect.xml
	}

	EventSoundEnabled = document.Global.Software.Sounds.Event.Enabled
	EventJoinedSoundFilenameAndPath = document.Global.Software.Sounds.Event.JoinedFilenameAndPath
	EventLeftSoundFilenameAndPath = document.Global.Software.Sounds.Event.LeftFilenameAndPath
	EventMessageSoundFilenameAndPath = document.Global.Software.Sounds.Event.MessageFilenameAndPath

	if EventSoundEnabled && EventJoinedSoundFilenameAndPath == "" {
		path := defaultSharePath + "/soundfiles/events/event.wav"
		if _, err := os.Stat(path); err == nil {
			EventJoinedSoundFilenameAndPath = path
		}
	}
	if EventSoundEnabled && EventLeftSoundFilenameAndPath == "" {
		path := defaultSharePath + "/soundfiles/events/event.wav"
		if _, err := os.Stat(path); err == nil {
			EventLeftSoundFilenameAndPath = path
		}
	}
	if EventSoundEnabled && EventMessageSoundFilenameAndPath == "" {
		path := defaultSharePath + "/soundfiles/events/event.wav"
		if _, err := os.Stat(path); err == nil {
			EventMessageSoundFilenameAndPath = path
		}
	}
	AlertSoundEnabled = document.Global.Software.Sounds.Alert.Enabled
	AlertSoundFilenameAndPath = document.Global.Software.Sounds.Alert.FilenameAndPath

	if AlertSoundEnabled && AlertSoundFilenameAndPath == "" {
		path := defaultSharePath + "/soundfiles/alerts/alert.wav"
		if _, err := os.Stat(path); err == nil {
			AlertSoundFilenameAndPath = path
		}
	}

	AlertSoundVolume = document.Global.Software.Sounds.Alert.Volume

	IncommingBeepSoundEnabled = document.Global.Software.Sounds.IncommingBeep.Enabled
	IncommingBeepSoundFilenameAndPath = document.Global.Software.Sounds.IncommingBeep.FilenameAndPath

	if IncommingBeepSoundEnabled && IncommingBeepSoundFilenameAndPath == "" {
		path := defaultSharePath + "/soundfiles/rogerbeeps/Chirsp.wav"
		if _, err := os.Stat(path); err == nil {
			IncommingBeepSoundFilenameAndPath = path
		}
	}

	IncommingBeepSoundVolume = document.Global.Software.Sounds.IncommingBeep.Volume

	RogerBeepSoundEnabled = document.Global.Software.Sounds.RogerBeep.Enabled
	RogerBeepSoundFilenameAndPath = document.Global.Software.Sounds.RogerBeep.FilenameAndPath

	if RogerBeepSoundEnabled && RogerBeepSoundFilenameAndPath == "" {
		path := defaultSharePath + "/soundfiles/rogerbeeps/Chirsp.wav"
		if _, err := os.Stat(path); err == nil {
			RogerBeepSoundFilenameAndPath = path
		}
	}

	RogerBeepSoundVolume = document.Global.Software.Sounds.RogerBeep.Volume

	RepeaterToneEnabled = document.Global.Software.Sounds.RepeaterTone.Enabled
	RepeaterToneFrequencyHz = document.Global.Software.Sounds.RepeaterTone.ToneFrequencyHz
	RepeaterToneDurationSec = document.Global.Software.Sounds.RepeaterTone.ToneDurationSec

	StreamSoundEnabled = document.Global.Software.Sounds.Stream.Enabled
	StreamSoundFilenameAndPath = document.Global.Software.Sounds.Stream.FilenameAndPath

	if StreamSoundEnabled && StreamSoundFilenameAndPath == "" {
		path := defaultSharePath + "/soundfiles/alerts/stream.wav"
		if _, err := os.Stat(path); err == nil {
			StreamSoundFilenameAndPath = path
		}
	}

	StreamSoundVolume = document.Global.Software.Sounds.Stream.Volume

	TxTimeOutEnabled = document.Global.Software.TxTimeOut.Enabled
	TxTimeOutSecs = document.Global.Software.TxTimeOut.TxTimeOutSecs

	APIEnabled = document.Global.Software.API.Enabled
	APIListenPort = document.Global.Software.API.ListenPort
	APIDisplayMenu = document.Global.Software.API.DisplayMenu
	APIChannelUp = document.Global.Software.API.ChannelUp
	APIChannelDown = document.Global.Software.API.ChannelDown
	APIMute = document.Global.Software.API.Mute
	APICurrentVolumeLevel = document.Global.Software.API.CurrentVolumeLevel
	APIDigitalVolumeUp = document.Global.Software.API.DigitalVolumeUp
	APIDigitalVolumeDown = document.Global.Software.API.DigitalVolumeDown
	APIListServerChannels = document.Global.Software.API.ListServerChannels
	APIStartTransmitting = document.Global.Software.API.StartTransmitting
	APIStopTransmitting = document.Global.Software.API.StopTransmitting
	APIListOnlineUsers = document.Global.Software.API.ListOnlineUsers
	APIPlayStream = document.Global.Software.API.PlayStream
	APIRequestGpsPosition = document.Global.Software.API.RequestGpsPosition
	APIEmailEnabled = document.Global.Software.API.Enabled
	APINextServer = document.Global.Software.API.NextServer
	APIPreviousServer = document.Global.Software.API.PreviousServer
	APIPanicSimulation = document.Global.Software.API.PanicSimulation
	APIDisplayVersion = document.Global.Software.API.DisplayVersion
	APIClearScreen = document.Global.Software.API.ClearScreen
	APIPingServersEnabled = document.Global.Software.API.Enabled
	APIRepeatTxLoopTest = document.Global.Software.API.RepeatTxLoopTest
	APIPrintXmlConfig = document.Global.Software.API.PrintXmlConfig

	MQTTEnabled = document.Global.Software.MQTT.MQTTEnabled
	MQTTTopic = document.Global.Software.MQTT.MQTTTopic
	MQTTBroker = document.Global.Software.MQTT.MQTTBroker
	MQTTPassword = document.Global.Software.MQTT.MQTTPassword
	MQTTUser = document.Global.Software.MQTT.MQTTUser
	MQTTId = document.Global.Software.MQTT.MQTTId
	MQTTCleansess = document.Global.Software.MQTT.MQTTCleansess
	MQTTQos = document.Global.Software.MQTT.MQTTQos
	MQTTNum = document.Global.Software.MQTT.MQTTNum
	MQTTPayload = document.Global.Software.MQTT.MQTTPayload
	MQTTAction = document.Global.Software.MQTT.MQTTAction
	MQTTStore = document.Global.Software.MQTT.MQTTStore

	TargetBoard = document.Global.Hardware.TargetBoard

	log.Println("Successfully loaded XML configuration file into memory")

	for i := 0; i < len(document.Accounts.Account); i++ {
		if document.Accounts.Account[i].Default == true {
			log.Printf("info: Successfully Added Account %s to Index [%d]\n", document.Accounts.Account[i].Name, i)
		}
	}

	return nil
}

func modifyXMLTagServerHopping(inputXMLFile string, outputXMLFile string, nextserverindex int) {
	xmlfilein, err := os.Open(inputXMLFile)
	xmlfileout, err := os.Create(outputXMLFile)

	if err != nil {
		FatalCleanUp(err.Error())
	}

	defer xmlfilein.Close()
	defer xmlfileout.Close()
	decoder := xml.NewDecoder(xmlfilein)
	encoder := xml.NewEncoder(xmlfileout)
	encoder.Indent("", "	")

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("error: Getting token: %v\n", err)
			break
		}

		switch v := token.(type) {
		case xml.StartElement:
			if v.Name.Local == "document" {
				var document Document
				if v.Name.Local != "talkkonnect/xml" {
					err = decoder.DecodeElement(&document, &v)
					if err != nil {
						FatalCleanUp("Cannot Find XML Tag Document" + err.Error())
					}
				}
				// XML Tag to Replace
				document.Global.Software.Settings.NextServerIndex = nextserverindex

				err = encoder.EncodeElement(document, v)
				if err != nil {
					FatalCleanUp(err.Error())
				}
				continue
			}

		}

		if err := encoder.EncodeToken(xml.CopyToken(token)); err != nil {
			FatalCleanUp(err.Error())
		}
	}

	if err := encoder.Flush(); err != nil {
		FatalCleanUp(err.Error())
	} else {
		time.Sleep(2 * time.Second)
		copyFile(inputXMLFile, inputXMLFile+".bak")
		deleteFile(inputXMLFile)
		copyFile(outputXMLFile, inputXMLFile)
		c := exec.Command("reset")
		c.Stdout = os.Stdout
		c.Run()
		os.Exit(0)
	}

}
