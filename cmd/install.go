/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"html/template"
	"log"
	"os"

	"github.com/spf13/cobra"
)

type Operator struct {
	Name            string
	Namespace       string
	Source          string
	DefaultChannel  string
	Csv             string
	TargetNamespace string
	Crd             string
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [operator]",
	Short: "Install operator",
	Long: `Install operator
	Examples needed here`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		operator := args[0]
		stdout, _ := cmd.Flags().GetBool("stdout")
		namespace, source, defaultchannel, csv, _, target_namespace, crd := get_operator(operator)
		if stdout == true {
			templateStr := `{{ if eq .Namespace "openshift-operators" }}
apiVersion: v1
kind: Namespace
metadata:
  labels:
    openshift.io/cluster-monitoring: "true"
  name: {{ .Namespace }}
---
apiVersion: operators.coreos.com/v1
kind: OperatorGroup
metadata:
  name: {{ .Name }}-operatorgroup
  namespace: {{ .Namespace }}
spec:
  targetNamespaces:
  - {{ .Namespace }}
---
{{ end }}
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: {{ .Name }}-subscription
  namespace: {{ .Namespace }}
spec:
  channel: "{{ .DefaultChannel }}"
  name: {{ .Name }}
  source: {{ .Source }}
  sourceNamespace: openshift-marketplace
`
			t := template.New("Template")
			tpl, err := t.Parse(templateStr)
			if err != nil {
				log.Fatalln(err)
			}
			operatordata := Operator{
				Name:            operator,
				Namespace:       namespace,
				Source:          source,
				DefaultChannel:  defaultchannel,
				Csv:             csv,
				TargetNamespace: target_namespace,
				Crd:             crd,
			}
			err = tpl.Execute(os.Stdout, operatordata)
			if err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolP("stdout", "s", false, "Print to stdout")
}
