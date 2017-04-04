// Copyright 2017 Paul BCode All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package displayRGB

import (
	"log"
	"testing"
	"time"

	"golang.org/x/exp/io/i2c"
)

func TestInitLCD(t *testing.T) {
	d, err := Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0x3f)
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	d.Write("Hello François")
	time.Sleep(2 * time.Minute)

	d.Write("Hello François\n Hollande")
}
