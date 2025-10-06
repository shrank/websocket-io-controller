//go:build linux
package io

import (
	"fmt"
	"log"
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/raspi"
	"gobot.io/x/gobot/v2/system"
	"time"
	"sync"
)

board := raspi.NewAdaptor()
gpio_pins := make(map[string](*DirectPinDriver))
i2c_lock := sync.Mutex()

func Interrupt_init(pin string)(error) {
	fmt.Printf("Setup Interrupt Pin %s\n", pin)
	p, err := board.DigitalPin(pin)
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
	p, err := board.DigitalPin(pin)
	if err != nil {
		return err
	}
	err = inPin.ApplyOptions(system.WithPinDirectionOutput())
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
