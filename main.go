/*
Sniperkit-Bot
- Status: analyzed
*/

// Copyright (c) 2018 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/palantir/godel/framework/pluginapi/v2/pluginapi"

	"github.com/sniperkit/snk.fork.palantir-godel-amalgomate-plugin/cmd"
)

func main() {
	if ok := pluginapi.InfoCmd(os.Args, os.Stdout, cmd.PluginInfo); ok {
		return
	}
	os.Exit(cmd.Execute())
}
