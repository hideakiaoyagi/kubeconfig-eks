package main

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

func getCredential(c *cli.Context) error {
	// Check command-line arguments.
	eksregion := c.String("region")
	if eksregion == "" {
		return errors.New("argument error: 'region' is empty")
	}
	ekscluster := c.String("name")
	if ekscluster == "" {
		return errors.New("argument error: 'name' is empty")
	}
	configpath := c.String("config")
	if configpath == "" {
		return errors.New("argument error: 'config' is empty")
	}
	configKeyNameCluster := c.String("config-key-cluster")
	if configKeyNameCluster == "" {
		return errors.New("argument error: 'config-key-cluster' is empty")
	}
	configKeyNameUser := c.String("config-key-user")
	if configKeyNameUser == "" {
		return errors.New("argument error: 'config-key-user' is empty")
	}
	configKeyNameContext := c.String("config-key-context")
	if configKeyNameContext == "" {
		return errors.New("argument error: 'config-key-context' is empty")
	}

	// Check config file path.
	usr, _ := user.Current()
	configpath = strings.Replace(configpath, "~", usr.HomeDir, 1)
	configpath, _ = filepath.Abs(configpath)

	// Read the information from EKS.
	eks := NewEksService(eksregion)
	if err := eks.ReadParameters(ekscluster); err != nil {
		return err
	}

	// Read yaml data from config file.
	config := NewYamlFile(configpath)
	if err := config.ReadYamlFile(); err != nil {
		return err
	}

	// Set parameters from EKS information to yaml data.
	if err := config.SetParamCluster(configKeyNameCluster, eks.param.endpoint, eks.param.cadata); err != nil {
		return err
	}
	if err := config.SetParamUser(configKeyNameUser, eks.param.clustername); err != nil {
		return err
	}
	if err := config.SetParamContext(configKeyNameContext, configKeyNameCluster, configKeyNameUser); err != nil {
		return err
	}
	if err := config.SetParamCurrentContext(configKeyNameContext); err != nil {
		return err
	}

	// Write yaml data to config file.
	if err := config.WriteYamlFile(); err != nil {
		return err
	}

	if config.isnew {
		fmt.Printf("[Success] kubeconfig file '%s' has been generated.\n", config.path)
	} else {
		fmt.Printf("[Success] kubeconfig file '%s' has been updated.\n", config.path)
	}

	return nil
}

func main() {
	app := cli.NewApp()

	app.Name = "kubeconfig-eks"
	app.Version = "0.0.1"
	app.Usage = "generate/update kubeconfig file from Amazon EKS cluster informaton"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "~/.kube/config",
			Usage: "specify kubeconfig `FILE`",
		},
		cli.StringFlag{
			Name:  "name, n",
			Usage: "cluster name (required)",
		},
		cli.StringFlag{
			Name:  "region, r",
			Usage: "cluster region (required)",
		},
		cli.StringFlag{
			Name:  "config-key-cluster",
			Value: "eks-cluster",
			Usage: "specify 'cluster' key-name in config file",
		},
		cli.StringFlag{
			Name:  "config-key-user",
			Value: "eks-user",
			Usage: "specify 'user' key-name in config file",
		},
		cli.StringFlag{
			Name:  "config-key-context",
			Value: "eks-cluster",
			Usage: "specify 'context' key-name in config file",
		},
	}
	app.Action = getCredential

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("[Error] %s\n", err)
		os.Exit(1)
	}
}
