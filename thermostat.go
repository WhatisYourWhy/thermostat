package main

import (
	"log"
	"time"

	"github.com/alittlebrighter/thermostat/controller"
	tmeter "github.com/alittlebrighter/thermostat/thermometer"
	"github.com/alittlebrighter/thermostat/util"
)

const (
	lowTemp  float64 = 69
	highTemp         = 80

	errorTolerance = 3

	tempCheckInterval = 1 * time.Minute
)

type Thermostat struct {
	Modes                 `json:"modes"`
	DefaultMode           string                `json:"defaultMode"`
	Schedule              []*ScheduleEvent      `json:"schedule"`
	Overshoot             float64               `json:"overshoot"`
	PollInterval          int                   `json:"pollInterval"`
	MaxErrors, errorCount uint8                 `json:"maxErrors"`
	UnitPreference        util.TemperatureUnits `json:"unitPreference"`
	control               controller.Controller
	thermometer           tmeter.Thermometer
	Events                *util.RingBuffer `json:"events"`
}

type Modes map[string]*Window

type Window struct {
	LowTemp  float64 `json:"low"`
	HighTemp float64 `json:"high"`
}

type ScheduleEvent struct {
	Days     []time.Weekday `json:"days"`
	ModeName string         `json:"mode"`
	Start    util.ClockTime `json:"start"`
	End      util.ClockTime `json:"end"`
}

func (stat *Thermostat) CurrentTemperatureWindow(t time.Time) *Window {
	for _, spec := range stat.Schedule {
		if _, ok := stat.Modes[spec.ModeName]; !ok {
			continue
		}

		dayMatch := false
		for _, day := range spec.Days {
			if t.Weekday() == day {
				dayMatch = true
				break
			}
		}
		if !dayMatch {
			continue
		}

		switch {
		case t.Hour() < spec.Start.Hour() || t.Hour() > spec.End.Hour():
			fallthrough
		case t.Hour() == spec.Start.Hour() && t.Minute() < spec.Start.Minute():
			fallthrough
		case t.Hour() == spec.End.Hour() && t.Minute() > spec.End.Minute():
			continue
		default:
			return stat.Modes[spec.ModeName]
		}
	}

	return stat.Modes[stat.DefaultMode]
}

func (stat *Thermostat) ProcessTemperatureReading(temp float64, units util.TemperatureUnits) {
	if units == util.Celsius && stat.UnitPreference != util.Celsius {
		temp = util.TempCToF(temp)
	} else if units == util.Fahrenheit && stat.UnitPreference != util.Fahrenheit {
		temp = util.TempFToC(temp)
	}

	window := stat.CurrentTemperatureWindow(time.Now())

	log.Printf("Current Temperature (%s): %f\n", stat.UnitPreference[0], temp)
	switch {
	case (stat.control.Direction() == controller.Heating && temp > window.LowTemp+stat.Overshoot) || (stat.control.Direction() == controller.Cooling && temp < window.HighTemp-stat.Overshoot):
		log.Println("turning OFF")
		stat.control.Off()
	case temp < window.LowTemp:
		log.Println("turning on HEAT")
		stat.control.Heat()
	case temp > window.HighTemp:
		log.Println("turning on COOL")
		stat.control.Cool()
	default:
		log.Println("doing NOTHING")
	}

	stat.Events.Add(&util.EventLog{AmbientTemperature: temp, Units: stat.UnitPreference, Direction: stat.control.Direction()})
}

func (stat *Thermostat) HandleError() {
	stat.errorCount++

	if stat.errorCount > stat.MaxErrors {
		stat.control.Off()
		stat.errorCount = 0
	}
}

func (stat *Thermostat) Run(cancel <-chan bool) {
	// we want to do something right away
	temp, units, err := stat.thermometer.ReadTemperature()
	if err != nil {
		log.Println("Error reading Temperature: " + err.Error())
		stat.HandleError()
	} else {
		stat.ProcessTemperatureReading(temp, units)
	}

	ticker := time.NewTicker(time.Duration(stat.PollInterval) * time.Minute)
	for {
		select {
		case <-ticker.C:
			temp, units, err := stat.thermometer.ReadTemperature()
			if err != nil {
				log.Println("Error reading Temperature: " + err.Error())
				stat.HandleError()
			} else {
				stat.ProcessTemperatureReading(temp, units)
			}
		case <-cancel:
			return
		}
	}
}
