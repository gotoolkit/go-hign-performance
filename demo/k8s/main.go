package main

import (
	"github.com/spf13/cobra"
	"k8s.io/client-go/1.5/tools/clientcmd"
)

func main() {
	cmd := &cobra.Command{
		Use:   "hugo",
		Short: "Hugo is a very fast static site generator",
		Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
	fs := cmd.Flags()
	overrides := &clientcmd.ConfigOverrides{}
	clientcmd.BindOverrideFlags(overrides, *fs, clientcmd.RecommendedConfigOverrideFlags(""))
}
