package triggers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Crosse/geneva/internal/scanner"
	"github.com/google/gopacket"
	_ "github.com/google/gopacket/layers"
)

type Trigger interface {
	Protocol() string
	Field() string
	Gas() int
	Matches(gopacket.Packet) (bool, error)
	fmt.Stringer
}

func ParseTrigger(s *scanner.Scanner) (Trigger, error) {
	if _, err := s.Expect("["); err != nil {
		return nil, err
	}

	str, err := s.Until(']')
	if err != nil {
		return nil, err
	}
	_, _ = s.Pop()

	fields := strings.Split(str, ":")
	if len(fields) < 3 {
		return nil, fmt.Errorf("invalid trigger format")
	}

	if fields[0] == "" {
		return nil, fmt.Errorf("invalid protocol")
	}

	gas := 0
	if len(fields) == 4 {
		gas, err = strconv.Atoi(fields[3])
		if err != nil {
			return nil, err
		}
	}

	var trigger Trigger
	switch strings.ToLower(fields[0]) {
	case "ip":
		trigger, err = NewIPTrigger(fields[1], fields[2], gas)
	case "tcp":
		trigger, err = NewTCPTrigger(fields[1], fields[2], gas)
	}

	if err != nil {
		return nil, err
	}

	return trigger, nil
}
