package main

import (
	"fmt"

	"google.golang.org/grpc"

	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/vars"

	"github.com/pkg/errors"

	"github.com/euforia/thrap/config"
	"github.com/euforia/thrap/core"
	"github.com/euforia/thrap/thrapb"
	"gopkg.in/urfave/cli.v2"
)

var (
	errRemoteRequired = errors.New("thrap remote required")
)

func newCLI() *cli.App {
	cli.VersionPrinter = func(ctx *cli.Context) {
		fmt.Println(ctx.App.Version)
	}

	app := &cli.App{
		Name:     "thrap",
		HelpName: "thrap",
		Version:  version(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "thrap-addr",
				Usage:   "thrap registry address",
				EnvVars: []string{"THRAP_ADDR"},
			},
		},
		Commands: []*cli.Command{
			commandConfigure(),
			commandIdentity(),
			commandAgent(),
			commandStack(),
			commandPack(),
			commandVersion(),
		},
	}

	app.HideVersion = true

	return app
}

func commandVersion() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Show version",
		Action: func(ctx *cli.Context) error {
			fmt.Println(version())
			return nil
		},
	}
}

func commandConfigure() *cli.Command {
	return &cli.Command{
		Name:  "configure",
		Usage: "Configure global settings",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:   vars.VcsID,
				Usage:  "version control `provider`",
				Value:  "github",
				Hidden: true,
			},
			&cli.StringFlag{
				Name:  vars.VcsUsername,
				Usage: "version control `username`",
			},
			&cli.StringFlag{
				Name:  "data-dir",
				Usage: "data `directory`",
				Value: "~/" + consts.WorkDir,
			},
			&cli.BoolFlag{
				Name:  "no-prompt",
				Usage: "do not prompt for input",
			},
		},
		Action: func(ctx *cli.Context) error {
			opts := core.ConfigureOptions{
				VCS: &config.VCSConfig{
					ID:       ctx.String(vars.VcsID),
					Username: ctx.String(vars.VcsUsername),
				},
				DataDir:  ctx.String("data-dir"),
				NoPrompt: ctx.Bool("no-prompt"),
			}

			// Only configures things that are not configured
			return core.ConfigureGlobal(opts)
		},
	}
}

var errNotConfigured = errors.New("thrap not configured. Try running 'thrap configure'")

func newThrapClient(ctx *cli.Context) (thrapb.ThrapClient, error) {
	// Check remote addr
	remoteAddr := ctx.String("thrap-addr")
	if remoteAddr == "" {
		return nil, errRemoteRequired
	}

	cc, err := grpc.Dial(remoteAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return thrapb.NewThrapClient(cc), nil
}
