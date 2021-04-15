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
 *
 * My Blog is at www.talkkonnect.com
 * The source code is hosted at github.com/talkkonnect
 *
 *
 */

package talkkonnect

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"github.com/comail/colog"
	hd44780 "github.com/jdiderik/go-hd44780"
	"github.com/jdiderik/gpio"
	"github.com/jdiderik/gumble/gumble"
	"github.com/jdiderik/gumble/gumbleffmpeg"
	"github.com/jdiderik/gumble/gumbleutil"
	_ "github.com/jdiderik/gumble/opus"
	term "github.com/jdiderik/termbox-go"
	"github.com/jdiderik/volume-go"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	LcdText              = [4]string{"nil", "nil", "nil", "nil"}
	currentChannelID     uint32
	prevChannelID        uint32
	prevParticipantCount int    = 0
	prevButtonPress      string = "none"
	maxchannelid         uint32
	origVolume           int
	tempVolume           int
	ConfigXMLFile        string
	GPSTime              string
	GPSDate              string
	GPSLatitude          float64
	GPSLongitude         float64
	Streaming            bool
	ServerHop            bool
	HTTPServRunning      bool
	message              string
	isrepeattx           bool = true
	NowStreaming         bool
	MyLedStrip           *LedStrip
)

type Talkkonnect struct {
	Config             *gumble.Config
	Client             *gumble.Client
	Name               string
	Address            string
	Username           string
	Ident              string
	TLSConfig          tls.Config
	ConnectAttempts    uint
	Stream             *Stream
	ChannelName        string
	Daemonize          bool
	IsTransmitting     bool
	IsPlayStream       bool
	GPIOEnabled        bool
	OnlineLED          gpio.Pin
	ParticipantsLED    gpio.Pin
	TransmitLED        gpio.Pin
	HeartBeatLED       gpio.Pin
	BackLightLED       gpio.Pin
	VoiceActivityLED   gpio.Pin
	AttentionLED       gpio.Pin
	TxButton           gpio.Pin
	TxButtonState      uint
	TxToggle           gpio.Pin
	TxToggleState      uint
	UpButton           gpio.Pin
	UpButtonState      uint
	DownButton         gpio.Pin
	DownButtonState    uint
	PanicButton        gpio.Pin
	PanicButtonState   uint
	CommentButton      gpio.Pin
	CommentButtonState uint
	StreamButton       gpio.Pin
	StreamButtonState  uint
}

type ChannelsListStruct struct {
	chanID     uint32
	chanName   string
	chanParent *gumble.Channel
	chanUsers  int
}

func Init(file string, ServerIndex string) {
	err := term.Init()
	if err != nil {
		FatalCleanUp("Cannot Initialize Terminal Error: " + err.Error())
	}
	defer term.Close()

	colog.Register()
	colog.SetOutput(os.Stdout)

	ConfigXMLFile = file
	err = readxmlconfig(ConfigXMLFile)
	if err != nil {
		message := err.Error()
		FatalCleanUp(message)
	}

	if Logging == "screen" {
		colog.SetFlags(log.Ldate | log.Ltime)
	}

	if Logging == "screenwithlineno" || Logging == "screenandfilewithlineno" {
		colog.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	}

	switch Loglevel {
	case "trace":
		colog.SetMinLevel(colog.LTrace)
		log.Println("info: Loglevel Set to Trace")
	case "debug":
		colog.SetMinLevel(colog.LDebug)
		log.Println("info: Loglevel Set to Debug")
	case "info":
		colog.SetMinLevel(colog.LInfo)
		log.Println("info: Loglevel Set to Info")
	case "warning":
		colog.SetMinLevel(colog.LWarning)
		log.Println("info: Loglevel Set to Warning")
	case "error":
		colog.SetMinLevel(colog.LError)
		log.Println("info: Loglevel Set to Error")
	case "alert":
		colog.SetMinLevel(colog.LAlert)
		log.Println("info: Loglevel Set to Alert")
	default:
		colog.SetMinLevel(colog.LInfo)
		log.Println("info: Default Loglevel unset in XML config automatically loglevel to Info")
	}


	if APEnabled {
		log.Println("info: Contacting http Provisioning Server Pls Wait")
		err := autoProvision()
		time.Sleep(5 * time.Second)
		if err != nil {
			FatalCleanUp("Error from AutoProvisioning Module " + err.Error())
		} else {
			log.Println("info: Loading XML Config")
			ConfigXMLFile = file
			readxmlconfig(ConfigXMLFile)
		}
	}

	if NextServerIndex > 0 {
		AccountIndex = NextServerIndex
	} else {
		AccountIndex, err = strconv.Atoi(ServerIndex)
	}

	b := Talkkonnect{
		Config:      gumble.NewConfig(),
		Name:        Name[AccountIndex],
		Address:     Server[AccountIndex],
		Username:    Username[AccountIndex],
		Ident:       Ident[AccountIndex],
		ChannelName: Channel[AccountIndex],
		Daemonize:   Daemonize,
	}

	if MQTTEnabled == true {
		log.Printf("info: Attempting to Contact MQTT Server")
		log.Printf("info: MQTT Broker      : %s\n", MQTTBroker)
		log.Printf("info: Subscribed topic : %s\n", MQTTTopic)
		go b.mqttsubscribe()
	} else {
		log.Printf("info: MQTT Server Subscription Disabled in Config")
	}

	if len(b.Username) == 0 {
		buf := make([]byte, 6)
		_, err := rand.Read(buf)
		if err != nil {
			FatalCleanUp("Cannot Generate Random Number Error " + err.Error())
		}

		buf[0] |= 2
		b.Config.Username = fmt.Sprintf("talkkonnect-%02x%02x%02x%02x%02x%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
	} else {
		b.Config.Username = Username[AccountIndex]
	}

	b.Config.Password = Password[AccountIndex]

	if Insecure[AccountIndex] {
		b.TLSConfig.InsecureSkipVerify = true
	}
	if Certificate[AccountIndex] != "" {
		cert, err := tls.LoadX509KeyPair(Certificate[AccountIndex], Certificate[AccountIndex])
		if err != nil {
			FatalCleanUp("Certificate Error " + err.Error())
		}
		b.TLSConfig.Certificates = append(b.TLSConfig.Certificates, cert)
	}

	if APIEnabled && !HTTPServRunning {
		go func() {
			http.HandleFunc("/", b.httpAPI)

			if err := http.ListenAndServe(":"+APIListenPort, nil); err != nil {
				FatalCleanUp("Problem Starting HTTP API Server " + err.Error())
			}
		}()
	}

	b.ClientStart()
	IsConnected = false

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	exitStatus := 0

	<-sigs
	b.CleanUp()
	os.Exit(exitStatus)
}

func (b *Talkkonnect) ClientStart() {
	f, err := os.OpenFile(LogFilenameAndPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.Println("info: Trying to Open File ", LogFilenameAndPath)
	if err != nil {
		FatalCleanUp("Problem Opening talkkonnect.log file " + err.Error())
	}

	if TargetBoard == "rpi" {
		if !LedStripEnabled {
			b.LEDOffAll()
		}
	}

	if Logging == "screenandfile" {
		log.Println("info: Logging is set to: ", Logging)
		wrt := io.MultiWriter(os.Stdout, f)
		colog.SetFlags(log.Ldate | log.Ltime)
		colog.SetOutput(wrt)
	}

	if Logging == "screenandfilewithlineno" {
		log.Println("info: Logging is set to: ", Logging)
		wrt := io.MultiWriter(os.Stdout, f)
		colog.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		colog.SetOutput(wrt)
	}

	b.Config.Attach(gumbleutil.AutoBitrate)
	b.Config.Attach(b)

	log.Printf("info: [%d] Default Mumble Accounts Found in XML config\n", AccountCount)

	if TargetBoard == "rpi" {
		log.Println("info: Target Board Set as RPI (gpio enabled) ")
		b.initGPIO()
		if LedStripEnabled {
			MyLedStrip, _ = NewLedStrip()
			log.Printf("info: Led Strip %v %s\n", MyLedStrip.buf, MyLedStrip.display)
		}
	} else {
		log.Println("info: Target Board Set as PC (gpio disabled) ")
	}

	if (TargetBoard == "rpi" && LCDBackLightTimerEnabled == true) && (OLEDEnabled == true || LCDEnabled == true) {

		log.Println("info: Backlight Timer Enabled by Config")
		BackLightTime = *BackLightTimePtr
		BackLightTime = time.NewTicker(LCDBackLightTimeoutSecs * time.Second)

		go func() {
			for {
				<-BackLightTime.C
				log.Printf("debug: LCD Backlight Ticker Timed Out After %d Seconds", LCDBackLightTimeoutSecs)
				LCDIsDark = true
				if LCDInterfaceType == "parallel" {
					b.LEDOff(b.BackLightLED)
				}
				if LCDInterfaceType == "i2c" {
					lcd := hd44780.NewI2C4bit(LCDI2CAddress)
					if err := lcd.Open(); err != nil {
						log.Println("error: Can't open lcd: " + err.Error())
						return
					}
					lcd.ToggleBacklight()
				}
				if OLEDEnabled == true && OLEDInterfacetype == "i2c" {
					Oled.DisplayOff()
					LCDIsDark = true
				}
			}
		}()
	} else {
		log.Println("debug: Backlight Timer Disabled by Config")
	}

	talkkonnectBanner("\u001b[44;1m") // add blue background to banner reference https://www.lihaoyi.com/post/BuildyourownCommandLinewithANSIescapecodes.html#background-colors

	err = volume.Unmute(OutputDevice)

	if err != nil {
		log.Println("error: Unable to Unmute ", err)
	} else {
		log.Println("debug: Speaker UnMuted Before Connect to Server")
	}

	if TTSEnabled && TTSTalkkonnectLoaded {
		err := playWavLocal(TTSTalkkonnectLoadedFilenameAndPath, TTSVolumeLevel)
		if err != nil {
			log.Println("error: PlayWavLocal(TTSTalkkonnectLoadedFilenameAndPath) Returned Error: ", err)
		}
	}

	b.Connect()

	pstream = gumbleffmpeg.New(b.Client, gumbleffmpeg.SourceFile(""), 0)

	if (HeartBeatEnabled) && (TargetBoard == "rpi") {
		HeartBeat := time.NewTicker(time.Duration(PeriodmSecs) * time.Millisecond)

		go func() {
			for range HeartBeat.C {
				timer1 := time.NewTimer(time.Duration(LEDOnmSecs) * time.Millisecond)
				timer2 := time.NewTimer(time.Duration(LEDOffmSecs) * time.Millisecond)
				<-timer1.C
				if HeartBeatEnabled {
					b.LEDOn(b.HeartBeatLED)
				}
				<-timer2.C
				if HeartBeatEnabled {
					b.LEDOff(b.HeartBeatLED)
				}
				if KillHeartBeat == true {
					HeartBeat.Stop()
				}

			}
		}()
	}

	if BeaconEnabled {
		BeaconTicker := time.NewTicker(time.Duration(BeaconTimerSecs) * time.Second)

		go func() {
			for range BeaconTicker.C {
				IsPlayStream = true
				b.playIntoStream(BeaconFilenameAndPath, BVolume)
				IsPlayStream = false
				log.Println("info: Beacon Enabled and Timed Out Auto Played File ", BeaconFilenameAndPath, " Into Stream")
			}
		}()
	}

	b.BackLightTimer()

	if LCDEnabled == true {
		b.LEDOn(b.BackLightLED)
		LCDIsDark = false
	}

	if OLEDEnabled == true {
		Oled.DisplayOn()
		LCDIsDark = false
	}

	if AudioRecordEnabled == true {

		if AudioRecordOnStart == true {

			if AudioRecordMode != "" {

				if AudioRecordMode == "traffic" {
					log.Println("info: Incoming Traffic will be Recorded with sox")
					AudioRecordTraffic()
					if TargetBoard == "rpi" {
						if LCDEnabled == true {
							LcdText = [4]string{"nil", "nil", "nil", "Traffic Recording ->"} // 4
							LcdDisplay(LcdText, LCDRSPin, LCDEPin, LCDD4Pin, LCDD5Pin, LCDD6Pin, LCDD7Pin, LCDInterfaceType, LCDI2CAddress)
						}
						if OLEDEnabled == true {
							oledDisplay(false, 6, 1, "Traffic Recording") // 6
						}
					}
				}
				if AudioRecordMode == "ambient" {
					log.Println("info: Ambient Audio from Mic will be Recorded with sox")
					AudioRecordAmbient()
					if TargetBoard == "rpi" {
						if LCDEnabled == true {
							LcdText = [4]string{"nil", "nil", "nil", "Mic Recording ->"} // 4
							LcdDisplay(LcdText, LCDRSPin, LCDEPin, LCDD4Pin, LCDD5Pin, LCDD6Pin, LCDD7Pin, LCDInterfaceType, LCDI2CAddress)
						}
						if OLEDEnabled == true {
							oledDisplay(false, 6, 1, "Mic Recording") // 6
						}
					}
				}
				if AudioRecordMode == "combo" {
					log.Println("info: Both Incoming Traffic and Ambient Audio from Mic will be Recorded with sox")
					AudioRecordCombo()
					if TargetBoard == "rpi" {
						if LCDEnabled == true {
							LcdText = [4]string{"nil", "nil", "nil", "Combo Recording ->"} // 4
							LcdDisplay(LcdText, LCDRSPin, LCDEPin, LCDD4Pin, LCDD5Pin, LCDD6Pin, LCDD7Pin, LCDInterfaceType, LCDI2CAddress)
						}
						if OLEDEnabled == true {
							oledDisplay(false, 6, 1, "Combo Recording") //6
						}
					}
				}

			}

		}
	}

keyPressListenerLoop:
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {
			case term.KeyEsc:
				log.Println("error: ESC Key is Invalid")
				reset()
				break keyPressListenerLoop
			case term.KeyDelete:
				b.cmdDisplayMenu()
			case term.KeyF1:
				b.cmdChannelUp()
			case term.KeyF2:
				b.cmdChannelDown()
			case term.KeyF3:
				b.cmdMuteUnmute("toggle")
			case term.KeyF4:
				b.cmdCurrentVolume()
			case term.KeyF5:
				b.cmdVolumeUp()
			case term.KeyF6:
				b.cmdVolumeDown()
			case term.KeyF7:
				b.cmdListServerChannels()
			case term.KeyF8:
				b.cmdStartTransmitting()
			case term.KeyF9:
				b.cmdStopTransmitting()
			case term.KeyF10:
				b.cmdListOnlineUsers()
			case term.KeyF11:
				b.cmdPlayback()
			case term.KeyF12:
				b.cmdGPSPosition()
			case term.KeyCtrlC:
				talkkonnectAcknowledgements("\u001b[44;1m") // add blue background to banner reference https://www.lihaoyi.com/post/BuildyourownCommandLinewithANSIescapecodes.html#background-colors
				b.cmdQuitTalkkonnect()
			case term.KeyCtrlD:
				b.cmdDebugStacktrace()
			case term.KeyCtrlE:
				b.cmdSendEmail()
			case term.KeyCtrlF:
				b.cmdConnPreviousServer()
			case term.KeyCtrlI: // New. Audio Recording. Traffic
				b.cmdAudioTrafficRecord()
			case term.KeyCtrlJ: // New. Audio Recording. Mic
				b.cmdAudioMicRecord()
			case term.KeyCtrlK: // New/ Audio Recording. Combo
				b.cmdAudioMicTrafficRecord()
			case term.KeyCtrlL:
				b.cmdClearScreen()
			case term.KeyCtrlO:
				b.cmdPingServers()
			case term.KeyCtrlN:
				b.cmdConnNextServer()
			case term.KeyCtrlP:
				b.cmdPanicSimulation()
			case term.KeyCtrlG:
				b.cmdPlayRepeaterTone()
			case term.KeyCtrlR:
				b.cmdRepeatTxLoop()
			case term.KeyCtrlS:
				b.cmdScanChannels()
			case term.KeyCtrlT:
				b.cmdThanks()
			case term.KeyCtrlV:
				b.cmdShowUptime()
			case term.KeyCtrlU:
				b.cmdDisplayVersion()
			case term.KeyCtrlX:
				b.cmdDumpXMLConfig()
			default:
				if ev.Ch != 0 {
					log.Println("error: Invalid Keypress ASCII ", ev.Ch, "Press <DEL> for Menu")
				} else {
					log.Println("error: Key Not Mapped, Press <DEL> for menu", ev.Ch)
				}
			}
		case term.EventError:
			FatalCleanUp("Terminal Error " + err.Error())
		}
	}
}
