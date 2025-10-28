//go:build linux
package io

import (
	"fmt"
	"gobot.io/x/gobot/v2/platforms/raspi"
	"gobot.io/x/gobot/v2/system"
	"sync"
)

var board = raspi.NewAdaptor()
var i2c_lock sync.Mutex

func Raspi_init()(error) {
	err := 	board.Connect()
	return err
}

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
	return nil
}

func Interrupt_Fired(pin string)( bool, error) {
	read, err := board.DigitalRead(pin)
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
	return nil
}

func Output_set(pin string, value byte)(error) {
	err := board.DigitalWrite(pin, value)
	return err
}
