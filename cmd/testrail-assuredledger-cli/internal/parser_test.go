//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/testrail-cli/LICENSE.md.

package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseMessage(t *testing.T) {
	t.Run("simple 1", func(t *testing.T) {
		inp := "one caller=assured-ledger/ledger-core/virtual/statemachine/logger_step.go:57 Component=sm"
		length, obj, err := parseMessage(inp, 0)

		assert.Equal(t, 3, length)
		assert.Equal(t, "one", string(obj))
		assert.NoError(t, err)
	})

	t.Run("simple 2", func(t *testing.T) {
		inp := "one two caller=assured-ledger/ledger-core/virtual/statemachine/logger_step.go:57 Component=sm"
		length, obj, err := parseMessage(inp, 0)

		assert.Equal(t, 7, length)
		assert.Equal(t, "one two", string(obj))
		assert.NoError(t, err)
	})

	t.Run("simple 3", func(t *testing.T) {
		inp := "caller=assured-ledger/ledger-core/virtual/statemachine/logger_step.go:57 Component=sm"
		length, obj, err := parseMessage(inp, 0)

		assert.Equal(t, 0, length)
		assert.Equal(t, "", string(obj))
		assert.NoError(t, err)
	})

	t.Run("simple 4", func(t *testing.T) {
		inp := "  caller=assured-ledger/ledger-core/virtual/statemachine/logger_step.go:57 Component=sm"
		length, obj, err := parseMessage(inp, 0)

		assert.Equal(t, 1, length)
		assert.Equal(t, " ", string(obj))
		assert.NoError(t, err)
	})

	t.Run("simple 5", func(t *testing.T) {
		inp := "   caller=assured-ledger/ledger-core/virtual/statemachine/logger_step.go:57 Component=sm"
		length, obj, err := parseMessage(inp, 0)

		assert.Equal(t, 2, length)
		assert.Equal(t, "  ", string(obj))
		assert.NoError(t, err)
	})
}

func Test_parseKeyString(t *testing.T) {
	t.Run("simple 1", func(t *testing.T) {
		t.Parallel()

		inp := "caller=assured-ledger/ledger-core/virtual/statemachine/logger_step.go:57 Component=sm"
		length, obj, err := parseKeyString(inp)

		assert.Equal(t, 6, length)
		assert.Equal(t, "caller", string(obj))
		assert.NoError(t, err)
	})

	t.Run("simple 2", func(t *testing.T) {
		t.Parallel()

		inp := "caller"
		_, _, err := parseKeyString(inp)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), " 2")
	})

	t.Run("simple 3", func(t *testing.T) {
		t.Parallel()

		inp := "cal ler=123"
		_, _, err := parseKeyString(inp)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), " 1")
	})
}

func Test_parseJsonValueString(t *testing.T) {
	fld := "assured-ledger/ledger-core/virtual/statemachine/logger_step.go:57"

	t.Run("simple 1", func(t *testing.T) {
		inp := "caller=assured-ledger/ledger-core/virtual/statemachine/logger_step.go:57 Component=sm"
		length, obj, err := parseJsonValueString(inp[7:])

		assert.Equal(t, 65, length)
		assert.Equal(t, fld, string(obj))
		assert.NoError(t, err)
	})

	t.Run("simple 2", func(t *testing.T) {
		inp := "caller= "
		length, obj, err := parseJsonValueString(inp[7:])

		assert.Equal(t, 0, length)
		assert.Equal(t, "", string(obj))
		assert.NoError(t, err)
	})

	t.Run("harder", func(t *testing.T) {
		inp := "caller=assured-ledger/ledger-core/virtual/statemachine/logger_step.go:57"
		length, obj, err := parseJsonValueString(inp[7:])

		assert.Equal(t, 65, length)
		assert.Equal(t, fld, string(obj))
		assert.NoError(t, err)
	})

	t.Run("quotes 1", func(t *testing.T) {
		inp := "caller=\"assured-ledger/ledger-core/virtual/statemachine/logger_step.go:57\""
		length, obj, err := parseJsonValueString(inp[7:])

		assert.Equal(t, 66, length)
		assert.Equal(t, fld, string(obj))
		assert.NoError(t, err)
	})

	t.Run("quotes 2", func(t *testing.T) {
		inp := "caller=\"assured-ledger/ledger-core/virtual/statemachine/logger_step.go:57\" Component=sm"
		length, obj, err := parseJsonValueString(inp[7:])

		assert.Equal(t, 66, length)
		assert.Equal(t, fld, string(obj))
		assert.NoError(t, err)
	})

	t.Run("quotes 3", func(t *testing.T) {
		inp := "caller=\"\""
		length, obj, err := parseJsonValueString(inp[7:])

		assert.Equal(t, 1, length)
		assert.Equal(t, "", string(obj))
		assert.NoError(t, err)
	})

	t.Run("quotes 4", func(t *testing.T) {
		inp := "caller=\""
		length, obj, err := parseJsonValueString(inp[7:])

		assert.Equal(t, 0, length)
		assert.Equal(t, []byte(nil), obj)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), " 5")
	})

	t.Run("quotes 5", func(t *testing.T) {
		inp := "caller=\"123"
		length, obj, err := parseJsonValueString(inp[7:])

		assert.Equal(t, 0, length)
		assert.Equal(t, []byte(nil), obj)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), " 5")
	})

	t.Run("quotes 6", func(t *testing.T) {
		inp := "caller=\"123"
		length, obj, err := parseJsonValueString(inp[7:])

		assert.Equal(t, 0, length)
		assert.Equal(t, []byte(nil), obj)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), " 5")
	})

	t.Run("quotes 7", func(t *testing.T) {
		inp := "caller=\"\t\""
		length, obj, err := parseJsonValueString(inp[7:])

		assert.Equal(t, 2, length)
		assert.Len(t, obj, 1)
		assert.NoError(t, err)
	})

	t.Run("quotes 8", func(t *testing.T) {
		inp := "caller=\"\\u00\""
		_, _, err := parseJsonValueString(inp[7:])

		assert.Error(t, err)
		assert.Contains(t, err.Error(), " 2")
	})

	t.Run("quotes 9", func(t *testing.T) {
		inp := "caller=\"\\z\""
		_, _, err := parseJsonValueString(inp[7:])

		assert.Error(t, err)
		assert.Contains(t, err.Error(), " 4")
	})

	t.Run("quotes 10", func(t *testing.T) {
		inp := "caller=\"\\"
		_, _, err := parseJsonValueString(inp[7:])

		assert.Error(t, err)
		assert.Contains(t, err.Error(), " 1")
	})

	t.Run("quotes 11", func(t *testing.T) {
		inp := "caller=\"\\u00"
		_, _, err := parseJsonValueString(inp[7:])

		assert.Error(t, err)
		assert.Contains(t, err.Error(), " 2")
	})

	t.Run("quotes 12", func(t *testing.T) {
		inp := "caller=\"\\u001\""
		_, _, err := parseJsonValueString(inp[7:])

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid syntax")
	})

	t.Run("quotes 13", func(t *testing.T) {
		inp := "caller=\"\\u002b\""
		_, obj, err := parseJsonValueString(inp[7:])

		assert.Equal(t, "+", string(obj))
		assert.NoError(t, err)
	})
}