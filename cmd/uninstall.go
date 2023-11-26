package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"k8s.io/klog/v2"
	"os"
	"strings"
)

const (
	operatorUninstallDesc = `
 'uninstall' command deletes the Nineinfra platform along with all the dependencies.`
	operatorUninstallExample = `  kubectl nine uninstall`
)

type operatorUninstallCmd struct {
	out    io.Writer
	errOut io.Writer
	output bool
}

func newUninstallCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	o := &operatorUninstallCmd{out: out, errOut: errOut}

	cmd := &cobra.Command{
		Use:     "uninstall",
		Short:   "Uninstall the NineInfra",
		Long:    operatorUninstallDesc,
		Example: operatorUninstallExample,
		Args:    cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.run(out)
			if err != nil {
				klog.Warning(err)
				return err
			}
			return nil
		},
	}
	cmd = DisableHelp(cmd)
	return cmd
}

// run deletes the Nineinfra to Kubernetes cluster.
func (o *operatorUninstallCmd) run(writer io.Writer) error {

	exist, cl := CheckNineClusterExist("", "")
	if exist {
		fmt.Printf("Error: NineClusters Exists! Please delete these NineClusters firstly!\n")
		PrintClusterList(cl)
		os.Exit(1)
	}
	if err := InitHelm(); err != nil {
		fmt.Printf("Error: %v \n", err)
		os.Exit(1)
	}

	path, _ := rootCmd.Flags().GetString(kubeconfig)

	parameters := []string{}
	if path != "" {
		parameters = append([]string{"--kubeconfig", path}, parameters...)
	}
	flags := strings.Join(parameters, " ")

	for _, v := range DefaultChartList {
		err := HelmUnInstall(v, "", DefaultNamespace, flags)
		if err != nil {
			fmt.Printf("Error: %v \n", err)
			os.Exit(1)
		}
	}

	if err := RemoveHelmRepo(DefaultHelmRepo); err != nil {
		fmt.Printf("Error: %v \n", err)
		os.Exit(1)
	}

	if err := DeleteIfExist(DefaultNamespace, flags); err != nil {
		fmt.Printf("Error: %v \n", err)
		os.Exit(1)
	}
	return nil
}