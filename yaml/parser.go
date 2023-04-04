package yaml_parser

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

/*
Cmd		"ECHO 123"
When	"5min" -> in 5minutes from now | "23:12" next time local clock hits 23:12
RepeatIn	"1h"	-> every 1 hour | undefined -> no repeats
RepeatTimes
*/
type processYaml struct {
	Cmd         string
	When        string
	RepeatIn    string `yaml:"repeatin"`
	RepeatTimes int
}

type ProcessRepeat struct {
	Name		string
	Cmd         string
	When        time.Time
	RepeatIn    time.Duration
	RepeatTimes int
}

type ProcessOnce struct {
	Name		string
	Cmd  		string
	When 		time.Time
}

func ParseCmds(path string) ([]ProcessOnce, []ProcessRepeat) {
	data := parseYaml(path)
	processOnce, processRepeat := parseToProcess(data)
	return processOnce, processRepeat
}

func parseYaml(path string) map[string]processYaml {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	data := make(map[string]processYaml)
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func parseToProcess(ps map[string]processYaml) ([]ProcessOnce, []ProcessRepeat) {
	procsOnce := []ProcessOnce{}
	procsRep := []ProcessRepeat{}
	for key, p := range ps {
		when, repeatIn := parseWhenRepeat(p.When, p.RepeatIn)
		if repeatIn == 0 {
			procsOnce = append(procsOnce, ProcessOnce{
				Name:	key,
				Cmd:  p.Cmd,
				When: when,
			})
		} else {
			procsRep = append(procsRep, ProcessRepeat{
				Name:		key,
				Cmd:         p.Cmd,
				When:        when,
				RepeatIn:    repeatIn,
				RepeatTimes: p.RepeatTimes,
			})
		}
	}
	return procsOnce, procsRep
}

func parseWhenRepeat(whenStr, repeatStr string) (when time.Time, repeatIn time.Duration) {
	/*PARSE when Argument*/
	repeatIn, err := time.ParseDuration(repeatStr)
	if err != nil {
		if repeatStr == "" {
			repeatIn = time.Duration(0)
		} else {
			log.Fatal(err)
		}
	}
	
	/*PARSE RepeatIn Argument*/
	dur, err := time.ParseDuration(whenStr)
	// parse regular pattern "1m" "12h" "1h15m30s"
	if err == nil {
		when = time.Now().Add(dur)
		return when, repeatIn
	}
	// time in past/no-input-made -> instantly start it by moving it in the past
	if whenStr == "" {
		return time.Now().Add(-time.Hour), repeatIn

	}
	// try to parse alternative clock-like pattern "12:30" "22:1" "8:00"
	dur, err = parseClockDuration(whenStr)
	if err != nil {
		log.Fatal(err)
	}
	when = time.Now().Truncate(24 * time.Hour) // 00:00 of today
	when = when.Add(dur)                       // hh:mm on today
	if when.Compare(time.Now()) == -1 {        // check if that time already happened today
		when.Add(time.Hour * 24)
	}
	return when, repeatIn
}

func parseClockDuration(str string) (time.Duration, error) {
	a := strings.Split(str, ":")
	if len(a) != 2 {
		return time.Duration(0), fmt.Errorf("ERROR: failed at parseClock() cant parse Time")
	}
	hours, err := strconv.ParseUint(a[0], 10, 64)
	if err != nil {
		return time.Duration(0), fmt.Errorf("ERROR: failed at parseClock() cant parse Time")
	}
	mins, err := strconv.ParseUint(a[1], 10, 64)
	if hours > 23 || mins > 59 {
		return time.Duration(0), fmt.Errorf("ERROR: failed at parseClock() time must between 0:00 and 23:59")
	}
	return (time.Hour*time.Duration(hours) + time.Minute*time.Duration(mins)), nil
}
