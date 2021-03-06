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
	"os"

	"github.com/Hex-Techs/hexctl/pkg/display"
	"github.com/Hex-Techs/hexctl/pkg/kc"
	"github.com/spf13/cobra"
)

var kubeconfig string
var src string

// kcCmd represents the kc command
var kcCmd = &cobra.Command{
	Use:              "kc",
	TraverseChildren: true,
	Short:            "manage your kubeconfig and context",
	Long: `kc helps you manage kubeconfig files and contexts,
it will switch context or show current context.

- show current context
  hexctl kc show

you must have kubectl command already.`,
}

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "switch your kube context",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			cmd.Help()
			os.Exit(1)
		}
		kc.Switch(kubeconfig)
	},
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show your current kube context",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			cmd.Help()
			os.Exit(1)
		}
		kc.Show(kubeconfig)
	},
}

var nsCmd = &cobra.Command{
	Use:   "ns",
	Short: "switch your current kube context default namespace",
	Run: func(cmd *cobra.Command, args []string) {
		var ns string
		if len(args) == 1 {
			ns = args[0]
		} else if len(args) == 0 {
			ns = "default"
		} else {
			display.Errorln("error: you must give a namespace by the current context cluster")
			os.Exit(1)
		}
		kc.Namespace(kubeconfig, ns)
	},
}

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "merge tow kubeconfig file in ~/.kube/config",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			cmd.Help()
			os.Exit(1)
		}
		kc.Merge(src, kubeconfig)
	},
}

func init() {
	mergeCmd.Flags().StringVarP(&src, "src", "s", "", "Specify the kubeconfig file to merge (required)")
	kcCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "", "", "Specify the kubeconfig file to modify, default ~/.kube/config")

	mergeCmd.MarkFlagRequired("src")

	kcCmd.AddCommand(switchCmd)
	kcCmd.AddCommand(showCmd)
	kcCmd.AddCommand(nsCmd)
	kcCmd.AddCommand(mergeCmd)
	rootCmd.AddCommand(kcCmd)
}
