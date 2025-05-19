package main

import (
	"os"
	"path/filepath"

	"github.com/KimMachineGun/automemlimit/memlimit"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/ostrace/runner/common"
	cli_helpers "github.com/ostrace/runner/helpers/cli"
	"github.com/ostrace/runner/log"
	"gitlab.com/gitlab-org/labkit/fips"

	_ "github.com/ostrace/runner/cache/azure"
	_ "github.com/ostrace/runner/cache/gcs"
	_ "github.com/ostrace/runner/cache/gcsv2"
	_ "github.com/ostrace/runner/cache/s3"
	_ "github.com/ostrace/runner/cache/s3v2"
	_ "github.com/ostrace/runner/commands"
	_ "github.com/ostrace/runner/commands/fleeting"
	_ "github.com/ostrace/runner/commands/helpers"
	_ "github.com/ostrace/runner/executors/custom"
	_ "github.com/ostrace/runner/executors/docker"
	_ "github.com/ostrace/runner/executors/docker/autoscaler"
	_ "github.com/ostrace/runner/executors/docker/machine"
	_ "github.com/ostrace/runner/executors/instance"
	_ "github.com/ostrace/runner/executors/kubernetes"
	_ "github.com/ostrace/runner/executors/parallels"
	_ "github.com/ostrace/runner/executors/shell"
	_ "github.com/ostrace/runner/executors/ssh"
	_ "github.com/ostrace/runner/executors/virtualbox"
	_ "github.com/ostrace/runner/helpers/secrets/resolvers/akeyless"
	_ "github.com/ostrace/runner/helpers/secrets/resolvers/azure_key_vault"
	_ "github.com/ostrace/runner/helpers/secrets/resolvers/gcp_secret_manager"
	_ "github.com/ostrace/runner/helpers/secrets/resolvers/vault"
	_ "github.com/ostrace/runner/shells"
)

func init() {
	_, _ = maxprocs.Set()
	memlimit.SetGoMemLimitWithEnv()
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			// log panics forces exit
			if _, ok := r.(*logrus.Entry); ok {
				os.Exit(1)
			}
			panic(r)
		}
	}()

	fips.Check()

	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Usage = "a OsTrace Runner"
	app.Version = common.AppVersion.ShortLine()
	cli.VersionPrinter = common.AppVersion.Printer
	app.Authors = []cli.Author{
		{
			Name:  "ostrace",
			Email: "dev.sulaiman@icloud.com",
		},
	}
	app.Commands = common.GetCommands()
	app.CommandNotFound = func(context *cli.Context, command string) {
		logrus.Fatalln("Command", command, "not found.")
	}

	cli_helpers.InitCli()
	cli_helpers.LogRuntimePlatform(app)
	cli_helpers.SetupCPUProfile(app)
	cli_helpers.FixHOME(app)
	cli_helpers.WarnOnBool(os.Args)

	log.ConfigureLogging(app)

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
