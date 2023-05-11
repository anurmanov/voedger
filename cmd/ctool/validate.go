/*
* Copyright (c) 2023-present Sigma-Soft, Ltd.
* @author Dmitry Molchanovsky
 */

package main

import (
	"github.com/spf13/cobra"
	"github.com/untillpro/goutils/logger"
)

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validates the configuration and status of the cluster for errors",
		RunE:  validate,
	}
}

func validate(cmd *cobra.Command, arg []string) error {
	mkCommandDirAndLogFile(cmd)

	cluster := newCluster()

	if !cluster.exists {
		logger.Error(ErrorClusterConfNotFound.Error)
		return ErrorClusterConfNotFound
	}

	err := cluster.validate()
	if err == nil {
		logger.Info("cluster configuration is ok")
	}
	return err
}