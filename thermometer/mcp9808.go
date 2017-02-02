package thermometer

import (
	"github.com/alittlebrighter/embd"
	_ "github.com/alittlebrighter/embd/host/rpi"
	"github.com/alittlebrighter/embd/sensor/mcp9808"
	"github.com/alittlebrighter/thermostat/util"
)

type MCP9808 struct {
	sensor *mcp9808.MCP9808
}

func NewMCP9808() (*MCP9808, error) {
	meter := new(MCP9808)

	var err error
	bus := embd.NewI2CBus(1)
	meter.sensor, err = mcp9808.New(bus)
	if err != nil {
		return nil, err
	}

	meter.sensor.SetShutdownMode(false)
	meter.sensor.SetTempResolution(mcp9808.SixteenthC)
	meter.sensor.SetTempHysteresis(mcp9808.Zero)

	return meter, nil
}

func (meter *MCP9808) ReadTemperature() (float64, util.TemperatureUnits, error) {
	tempReading, err := meter.sensor.AmbientTemp()
	return tempReading.CelsiusDeg, util.Celsius, err
}

func (meter *MCP9808) Shutdown() {
	meter.sensor.SetShutdownMode(true)
	embd.CloseI2C()
}
