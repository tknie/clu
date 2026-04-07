/*
* Copyright 2022-2025 Thorsten A. Knieling
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
 */

package clu

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-faster/jx"
	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	start := time.Now()
	ctx := NewContext("abc", "dsddfsd")
	assert.NotNil(t, ctx)
	assert.Equal(t, "abc", ctx.user.User)
	assert.Equal(t, "abc", ctx.UserName())
	assert.Equal(t, "dsddfsd", ctx.Pass)
	assert.WithinRange(t, ctx.user.Created, start, time.Now())
}

func disableTestJXRaw(t *testing.T) {
	var e jx.Encoder
	v := 11981337726687985304.0
	x := fmt.Sprintf("%f", v)
	assert.Equal(t, "11981337726687985304.000000", x)
	e.Raw([]byte(x))
	b := e.Bytes()
	fmt.Println("Bytes:", len(b))
	assert.Equal(t, []byte{0x31, 0x31, 0x39, 0x38, 0x31, 0x33, 0x33, 0x37, 0x37, 0x32, 0x36, 0x36, 0x38, 0x37, 0x39, 0x38, 0x35, 0x34, 0x30, 0x34}, b)
	raw := jx.Raw([]byte(e.String()))
	assert.Equal(t, "11981337726687985304", string(raw))

	// abs := math.Abs(v)
	// assert.Equal(t, 11981337726687985304.0, abs)
	// fmt := byte('f')
	// b = make([]byte, 0, 32)
	// b = strconv.AppendFloat(b, v, fmt, -1, 64)
	// assert.Equal(t, []byte{0x31, 0x31, 0x39, 0x38, 0x31, 0x33, 0x33, 0x37, 0x37, 0x32, 0x36, 0x36, 0x38, 0x37, 0x39, 0x38, 0x35, 0x34, 0x30, 0x34}, b)
	// bits := math.Float64bits(v)
	// assert.Equal(t, []byte{0x31, 0x31, 0x39, 0x38, 0x31, 0x33, 0x33, 0x37, 0x37, 0x32, 0x36, 0x36, 0x38, 0x37, 0x39, 0x38, 0x35, 0x34, 0x30, 0x34}, bits)
}
