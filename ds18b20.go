package ds18b20

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	w1root     = "/sys/bus/w1/devices"
	w1master   = "/w1_bus_master1/w1_master_slaves"
	deviceMask = "28-"
	dataFile   = "w1_slave"
)

type Sensor struct {
	Path string
	ID   string
}

type Reading struct {
	Date  time.Time
	Value float64
}

func Sensors() ([]Sensor, error) {
	w1m, err := os.Open(filepath.Join(w1root, w1master))
	defer w1m.Close()
	if err != nil {
		return nil, err
	}

	sc := bufio.NewScanner(w1m)
	sc.Split(bufio.ScanLines)

	var sensors []Sensor
	devName := ""
	for sc.Scan() {
		devName = sc.Text()
		if strings.HasPrefix(devName, deviceMask) {
			s := Sensor{
				Path: filepath.Join(w1root, devName, dataFile),
				ID:   devName[3:]}
			sensors = append(sensors, s)
		}
	}
	if err = sc.Err(); err != nil {
		return sensors, err
	}

	return sensors, nil
}

func (s *Sensor) Reading() (*Reading, error) {
	data, err := os.Open(s.Path)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanLines)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) < 2 {
		return nil, fmt.Errorf("sensor id: %s. not enough data in file", s.ID)
	}

	if lines[0][len(lines[0])-3:] != "YES" {
		return nil, fmt.Errorf("sensor id: %s. wrong checksum", s.ID)
	}

	tempIndex := strings.LastIndexAny(lines[1], "t=")
	if tempIndex == -1 {
		return nil, fmt.Errorf("sensor id: %s. no temperature value found", s.ID)
	}
	celsius, err := strconv.ParseFloat(lines[1][tempIndex+1:], 64)
	if err != nil {
		return nil, fmt.Errorf("sensor id: %s. could not parse temperature:%s", s.ID, lines[1][tempIndex+1:])
	}

	return &Reading{Date: time.Now(), Value: celsius / 1000.0}, nil
}
