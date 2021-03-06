/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/Hex-Techs/hexctl/pkg/common/file"
	"github.com/Hex-Techs/hexctl/pkg/display"
	"github.com/Hex-Techs/hexctl/pkg/run"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run a go project",
	Long: `run a go project. For example:

Gin or other web project, it will watch *.go file and when these file changed hexctl will reload it,
you must have a main.go file in workdir and code in the directory named pkg.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			cobra.CheckErr(fmt.Errorf("unknown args %v", args))
		}
		display.Infof("%s %s %s/%s\n", Logo, version, runtime.GOOS, runtime.GOARCH)
		pwd, _ := os.Getwd()
		if !file.IsExists(run.MainFile) {
			display.Errorln("can not find main.go in", pwd)
			os.Exit(126)
		}
		initWorkDir()
		stop := make(chan bool)
		dirs, err := run.GetDirList(filepath.Join(pwd, "."))
		dirs = handlerDir(dirs)
		cobra.CheckErr(err)
		go run.NewWatcher(dirs, stop)
		time.Sleep(500 * time.Millisecond)
		go run.Reload(command, stop)
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		for range sigs {
			run.Kill()
			os.Exit(0)
		}
		// for {
		// 	select {
		// 	case <-sigs:
		// 		run.Kill()
		// 		os.Exit(0)
		// 	}
		// }
	},
}

func initWorkDir() {
	display.Successf("Create work directory %s\n", run.WorkDir)
	err := os.Mkdir(run.WorkDir, 0755)
	if err != nil {
		if !os.IsExist(err) {
			cobra.CheckErr(err)
		}
	}
}

// handlerDir remove .git and bin directory
func handlerDir(dir []string) []string {
	var res []string
	reg := regexp.MustCompile(`.+/\.`)
	for _, d := range dir {
		if !strings.HasSuffix(d, "bin") {
			if !reg.Match([]byte(d)) {
				res = append(res, d)
			}
		}
	}
	return res
}

var command []string

func init() {
	runCmd.Flags().StringSliceVarP(&command, "cmd", "", []string{}, "app command")
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
