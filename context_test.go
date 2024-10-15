/*
* Copyright 2022-2024 Thorsten A. Knieling
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
	"testing"
	"time"

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
