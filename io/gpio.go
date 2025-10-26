//go:build linux
package io

import (
	"fmt"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/raspi"
	"gobot.io/x/gobot/v2/system"
	"sync"
)

var board = raspi.NewAdaptor()
var gpio_pins = make(map[string](*gpio.DirectPinDriver))
var i2c_lock sync.Mutex

func Interrupt_init(pin string)(error) {
	fmt.Printf("Setup Interrupt Pin %s\n", pin)
	inPin, err := board.DigitalPin(pin)
	if err != nil {
		return err
	}
	err = inPin.ApplyOptions(system.WithPinDirectionInput(), system.WithPinPullUp())
	if err != nil {
		return err
	}
	gpio_pins[pin] = gpio.NewDirectPinDriver(board, pin)
	return nil
}

func Interrupt_Fired(pin string)( bool, error) {
	read, err := gpio_pins[pin].DigitalRead()
	return (read==1), err
}


func Ouput_init(pin string)(error) {
	fmt.Printf("Setup Output Pin %s\n", pin)
	inPin, err := board.DigitalPin(pin)
	if err != nil {
		return err
	}
	err = inPin.ApplyOptions(system.WithPinDirectionOutput(0))
	if err != nil {
		return err
	}
	gpio_pins[pin] = gpio.NewDirectPinDriver(board, pin)
	return nil
}

func Output_set(pin string, value byte)(error) {
	err := gpio_pins[pin].DigitalWrite(value)
	return err
}
