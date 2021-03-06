/*
Sniperkit-Bot
- Status: analyzed
*/

// Copyright (c) 2018 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package amalgomateplugin

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/palantir/amalgomate/amalgomate"
	"github.com/palantir/godel/pkg/dirchecksum"
	"github.com/pkg/errors"
)

const indentLen = 2

func Run(param Param, verify bool, stdout io.Writer) error {
	var verifyFailedKeys []string
	verifyFailedErrors := make(map[string]string)
	verifyFailedFn := func(name, errStr string) {
		verifyFailedKeys = append(verifyFailedKeys, name)
		verifyFailedErrors[name] = errStr
	}

	sortedKeys := param.OrderedKeys
	if len(sortedKeys) == 0 {
		for k := range param.Amalgomators {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)
	}
	for _, k := range param.OrderedKeys {
		val := param.Amalgomators[k]
		cfg, err := amalgomate.LoadConfig(val.Config)
		if err != nil {
			return errors.Wrapf(err, "failed to read amalgomate configuration for %s", k)
		}

		if verify {
			if _, err := os.Stat(val.OutputDir); os.IsNotExist(err) {
				verifyFailedFn(k, fmt.Sprintf("output directory %s does not exist", val.OutputDir))
				continue
			}

			originalChecksums, err := dirchecksum.ChecksumsForMatchingPaths(val.OutputDir, nil)
			if err != nil {
				return errors.Wrapf(err, "failed to compute original checksums")
			}

			newChecksums, err := dirchecksum.ChecksumsForDirAfterAction(val.OutputDir, func(dir string) error {
				return amalgomate.Run(cfg, dir, val.Pkg)
			})
			if err != nil {
				return errors.Wrapf(err, "amalgomate verify failed for %s", k)
			}

			if diff := originalChecksums.Diff(newChecksums); len(diff.Diffs) > 0 {
				verifyFailedFn(k, diff.String())
				continue
			}
			continue
		}
		if err := amalgomate.Run(cfg, val.OutputDir, val.Pkg); err != nil {
			return errors.Wrapf(err, "amalgomate failed for %s", k)
		}
	}
	if verify && len(verifyFailedKeys) > 0 {
		fmt.Fprintf(stdout, "amalgomator output differs from what currently exists: %v\n", verifyFailedKeys)
		for _, currKey := range verifyFailedKeys {
			fmt.Fprintf(stdout, "%s%s:\n", strings.Repeat(" ", indentLen), currKey)
			for _, currErrLine := range strings.Split(verifyFailedErrors[currKey], "\n") {
				fmt.Fprintf(stdout, "%s%s\n", strings.Repeat(" ", indentLen*2), currErrLine)
			}
		}
		return fmt.Errorf("")
	}
	return nil
}
