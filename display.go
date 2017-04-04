// Copyright 2017 Paul BCode All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package displayRGB

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/exp/io/i2c"
	"golang.org/x/exp/io/i2c/driver"
)

const (
	LCD_CHR = 1 // Mode - Sending data
	LCD_CMD = 0 // Mode - Sending command

	LCD_BACKLIGHT_ON  = 0x08 // On
	LCD_BACKLIGHT_OFF = 0x00 //Off

	LCD_LINE_1 = 0x80 // LCD RAM address for the 1st line
	LCD_LINE_2 = 0xC0 // LCD RAM address for the 2nd line

	TIME_PULSE = (5 * time.Millisecond)
	TIME_DELAY = (5 * time.Millisecond)

	ENABLE = 4

	LCD_WIDTH = 16 //Size of the LCD
)

type Device struct {
	lcd       *i2c.Device
	backlight int
}

// Open opens a connection the the RGB backlight display.
// Once display is no longer in-use, it should be closed by Close.
func Open(o driver.Opener, lcdAddr int) (*Device, error) {
	lcd, err := i2c.Open(o, lcdAddr)
	if err != nil {
		return nil, fmt.Errorf("cannot open LCD device: %v", err)
	}
	d := &Device{lcd: lcd, backlight: LCD_BACKLIGHT_ON}
	err = d.InitLCD()
	if err != nil {
		return nil, err
	}
	return d, nil
}
func (d *Device) InitLCD() error {
	d.LcdByte(0x33, LCD_CMD) // 110011 Initialise
	d.LcdByte(0x32, LCD_CMD) // 110010 Initialise
	d.LcdByte(0x06, LCD_CMD) // 000110 Cursor move direction
	d.LcdByte(0x0C, LCD_CMD) // 001100 Display On,Cursor Off, Blink Off
	d.LcdByte(0x28, LCD_CMD) // 101000 Data length, number of lines, font size
	d.LcdByte(0x01, LCD_CMD) // 000001 Clear display
	return nil
}

func (d *Device) Write(msg string) {
	msg1 := ""
	msg2 := ""
	// Test si \n
	msgTmp := strings.Split(msg, "\n")
	if len(msgTmp) > 1 {
		msg1 = msgTmp[0]
		msg2 = msgTmp[1]
	} else {
		msg1 = msg
	}
	d.LcdString(msg1, LCD_LINE_1)
	d.LcdString(msg2, LCD_LINE_2)
}

func (d *Device) LcdString(msg string, line int) {
	d.LcdByte(line, LCD_CMD)
	space := " "
	maxi := 0
	for i, carac := range msg {
		maxi = i
		if i > LCD_WIDTH {
			goto endmsg
		}
		d.LcdByte(int(carac), LCD_CHR)
	}
	for ; maxi <= LCD_WIDTH; maxi++ {
		d.LcdByte(int(space[0]), LCD_CHR)
	}
endmsg:
}

func (d *Device) LcdByte(bits int, mode int) {
	bits_high := mode | (bits & 0xF0) | d.backlight
	bits_low := mode | ((bits << 4) & 0xF0) | d.backlight

	d.lcd.Write([]byte{byte(bits_high)})
	d.ToggleEnable(bits_high)

	d.lcd.Write([]byte{byte(bits_low)})
	d.ToggleEnable(bits_low)
}

func (d *Device) ToggleEnable(bits int) {
	time.Sleep(TIME_DELAY)
	d.lcd.Write([]byte{byte(bits | ENABLE)})
	time.Sleep(TIME_PULSE)
	d.lcd.Write([]byte{byte(bits & (0xFF - ENABLE))})
	time.Sleep(TIME_DELAY)
}
func (d *Device) Close() error {
	return d.lcd.Close()
}

func (d *Device) Clear() {
}

func (d *Device) BackLight_OFF() {
	d.backlight = LCD_BACKLIGHT_OFF
	d.LcdByte(0x00, LCD_CMD)
}
func (d *Device) BackLight_ON() {
	d.backlight = LCD_BACKLIGHT_ON
	d.LcdByte(0x00, LCD_CMD)
}
