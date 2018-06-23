package main

import (
	"fmt"
	"os"

	"github.com/euforia/thrap/manifest"
	"github.com/euforia/thrap/utils"
	"github.com/euforia/thrap/vcs"
	"gopkg.in/urfave/cli.v2"
)

func commandStack() *cli.Command {
	return &cli.Command{
		Name:  "stack",
		Usage: "Stack operations",
		Subcommands: []*cli.Command{
			commandStackBuild(),
			commandStackDeploy(),
			commandStackInit(),
			commandStackRegister(),
			commandStackValidate(),
			commandStackVersion(),
		},
	}
}

func commandStackValidate() *cli.Command {
	return &cli.Command{
		Name:      "validate",
		Usage:     "Validate a manifest",
		ArgsUsage: "[path to manifest]",
		Action: func(ctx *cli.Context) error {
			mfile := ctx.Args().Get(0)
			mf, err := manifest.LoadManifest(mfile)
			if err != nil {
				return err
			}

			rpath, err := utils.GetLocalPath("")
			if err != nil {
				return err
			}

			mf.Version = vcs.GetRepoVersion(rpath).String()
			errs := mf.Validate()
			if errs != nil {
				return utils.FlattenErrors(errs)

			}

			writeHCLManifest(mf, os.Stdout)

			return nil
		},
	}
}

func commandStackVersion() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Show stack version",
		Action: func(ctx *cli.Context) error {
			lpath, err := utils.GetLocalPath("")
			if err == nil {
				fmt.Println(vcs.GetRepoVersion(lpath))
			}
			return err
		},
	}
}