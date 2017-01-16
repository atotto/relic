/*
 * Copyright (c) SAS Institute Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package shared

import (
	"errors"
	"fmt"
	"os"

	"gerrit-pdt.unx.sas.com/tools/relic.git/config"
	"github.com/spf13/cobra"
)

var ArgConfig string
var CurrentConfig *config.Config
var argVersion bool

var RootCmd = &cobra.Command{
	Use:              "relic",
	PersistentPreRun: showVersion,
	RunE:             bailUnlessVersion,
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&ArgConfig, "config", "c", "", "Configuration file")
	RootCmd.PersistentFlags().BoolVar(&argVersion, "version", false, "Show version and exit")
}

func showVersion(cmd *cobra.Command, args []string) {
	if argVersion {
		fmt.Printf("relic version %s\n", config.Version)
		os.Exit(0)
	}
}

func bailUnlessVersion(cmd *cobra.Command, args []string) error {
	if !argVersion {
		return errors.New("Expected a command")
	}
	return nil
}

func InitConfig() error {
	usedDefault := false
	if ArgConfig == "" {
		ArgConfig = config.DefaultConfig()
		if ArgConfig == "" {
			return errors.New("--config not specified")
		}
		usedDefault = true
	}
	config, err := config.ReadFile(ArgConfig)
	if err != nil {
		if os.IsNotExist(err) && usedDefault {
			return fmt.Errorf("--config not specified and default config at %s does not exist", ArgConfig)
		}
		return err
	}
	CurrentConfig = config
	return nil
}

func Main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}