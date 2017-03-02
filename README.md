# Go Thermostat
This is a simple thermostat I've thrown together out of necessity (our previous one blew a capacitor and I had a raspberry pi lying around).  

## Features
Currently programmable on startup via a yaml configuration file.  After that the configuration can change via a JSON API.  Supports mode (low/high temp settings) changes that can be made based on the day of the week and a start and end time (must be within the midnight to midnight window of the day.

Note, there are two thermometer implementations.  One supports an MCP9808 temperature sensor that works over I2C and the other relies on another machine on the network providing a JSON API with the current temperature values.  Also, this works quite well with my relatively simple 5-wire central forced air HVAC system.  I haven't tested this with any other setup although it shouldn't be difficult to write another controller implementation that satisfies the controller interface.

## Road map
- let user choose thermometer implementation through configuration
- more controller implementations
- multiple thermometer support with the option of area priority in schedule (e.g. keep the temperature within the set limits in the living room during the day and focus on the temperature in the bedrooms at night)
- converge on embd library for GPIO access (it was getting colder and embd wasn't working)
- cleaner shutdown
- remote control outside of the network (Firebase?)
- convert to plugins for thermometer and controller implementations?
- security, not a priority at the moment since I only serve this over my home network which is secured

## Contributing
I'd be more than happy to accept pull requests that implement any of the stated objectives in the road map or even features I may not have thought of.  Just try to match the style or make a case for why my style is wrong.

## License
Mozilla Public License v2.0
[TLDR](https://tldrlegal.com/license/mozilla-public-license-2.0-(mpl-2))
[Full License](https://www.mozilla.org/en-US/MPL/2.0/)
