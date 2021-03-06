/*
Sniperkit-Bot
- Status: analyzed
*/

// Copyright (c) 2018 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/sniperkit/snk.fork.palantir-godel-amalgomate-plugin/amalgomateplugin"
	"github.com/sniperkit/snk.fork.palantir-godel-amalgomate-plugin/amalgomateplugin/config"
)

func TestReadConfig(t *testing.T) {
	content := `
amalgomators:
  test-product:
    config: test.yml
    output-dir: test-output
    pkg: test-pkg
  next-product:
    order: 1
    config: next.yml
    output-dir: next-output
    pkg: next-pkg
  other-product:
    order: 2
    config: other.yml
    output-dir: other-output
    pkg: other-pkg
`
	var got config.Config
	err := yaml.Unmarshal([]byte(content), &got)
	require.NoError(t, err)

	wantCfg := config.Config{
		Amalgomators: config.ToAmalgomators(map[string]config.ProductConfig{
			"test-product": {
				Config:    "test.yml",
				OutputDir: "test-output",
				Pkg:       "test-pkg",
			},
			"next-product": {
				Order:     1,
				Config:    "next.yml",
				OutputDir: "next-output",
				Pkg:       "next-pkg",
			},
			"other-product": {
				Order:     2,
				Config:    "other.yml",
				OutputDir: "other-output",
				Pkg:       "other-pkg",
			},
		}),
	}
	assert.Equal(t, wantCfg, got)
}

func TestToParam(t *testing.T) {
	cfg := config.Config{
		Amalgomators: config.ToAmalgomators(map[string]config.ProductConfig{
			"test-product": {
				Order:     -1,
				Config:    "test.yml",
				OutputDir: "test-output",
				Pkg:       "test-pkg",
			},
			"next-product": {
				Config:    "next.yml",
				OutputDir: "next-output",
				Pkg:       "next-pkg",
			},
			"other-product": {
				Config:    "other.yml",
				OutputDir: "other-output",
				Pkg:       "other-pkg",
			},
		}),
	}

	wantParam := amalgomateplugin.Param{
		OrderedKeys: []string{
			"test-product",
			"next-product",
			"other-product",
		},
		Amalgomators: map[string]amalgomateplugin.ProductParam{
			"test-product": {
				Config:    "test.yml",
				OutputDir: "test-output",
				Pkg:       "test-pkg",
			},
			"next-product": {
				Config:    "next.yml",
				OutputDir: "next-output",
				Pkg:       "next-pkg",
			},
			"other-product": {
				Config:    "other.yml",
				OutputDir: "other-output",
				Pkg:       "other-pkg",
			},
		},
	}
	gotParam := cfg.ToParam()
	assert.Equal(t, wantParam, gotParam)
}
