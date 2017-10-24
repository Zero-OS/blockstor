package copyvdisk

import (
	"context"
	"errors"
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/zero-os/0-Disk/config"
	"github.com/zero-os/0-Disk/log"
	"github.com/zero-os/0-Disk/nbd/ardb/storage"
	"github.com/zero-os/0-Disk/nbd/nbdserver/tlog"
	"github.com/zero-os/0-Disk/tlog/copy"
	tlogserver "github.com/zero-os/0-Disk/tlog/tlogserver/server"
	cmdconfig "github.com/zero-os/0-Disk/zeroctl/cmd/config"
)

var vdiskCmdCfg struct {
	SourceConfig            config.SourceConfig
	ForceSameStorageCluster bool
	PrivKey                 string
	FlushSize               int
	JobCount                int
}

// VdiskCmd represents the vdisk copy subcommand
var VdiskCmd = &cobra.Command{
	Use:   "vdisk source_vdiskid target_vdiskid [target_clusterid]",
	Short: "Copy a vdisk",
	RunE:  copyVdisk,
}

func copyVdisk(cmd *cobra.Command, args []string) error {
	logLevel := log.InfoLevel
	if cmdconfig.Verbose {
		logLevel = log.InfoLevel
	}
	log.SetLevel(logLevel)

	// create config source
	configSource, err := config.NewSource(vdiskCmdCfg.SourceConfig)
	if err != nil {
		return err
	}
	defer configSource.Close()

	log.Info("parsing positional arguments...")

	// validate pos arg length
	argn := len(args)
	if argn < 2 {
		return errors.New("not enough arguments")
	} else if argn > 3 {
		return errors.New("too many arguments")
	}

	// store pos arguments in named variables
	sourceVdiskID, targetVdiskID := args[0], args[1]
	var targetClusterID string
	if argn == 3 {
		targetClusterID = args[2]
	}

	var sourceClusterConfig config.StorageClusterConfig
	var targetClusterConfig *config.StorageClusterConfig

	// try to read the Vdisk+NBD config of source vdisk
	sourceStaticConfig, err := config.ReadVdiskStaticConfig(configSource, sourceVdiskID)
	if err != nil {
		return fmt.Errorf(
			"couldn't read source vdisk %s's static config: %v", sourceVdiskID, err)
	}
	sourceStorageConfig, err := config.ReadNBDStorageConfig(
		configSource, sourceVdiskID)
	if err != nil {
		return fmt.Errorf(
			"couldn't read source vdisk %s's storage config: %v", sourceVdiskID, err)
	}

	sourceClusterConfig = sourceStorageConfig.StorageCluster

	// if a targetClusterID is given, check if it's not the same as the source cluster ID,
	// and if it is not, get its config
	if targetClusterID != "" {
		sourceVdiskNBDConfig, err := config.ReadVdiskNBDConfig(configSource, sourceVdiskID)
		if err != nil {
			return fmt.Errorf(
				"couldn't read source vdisk %s's storage config: %v", sourceVdiskID, err)
		}

		if sourceVdiskNBDConfig.StorageClusterID != targetClusterID {
			targetClusterConfig, err = config.ReadStorageClusterConfig(configSource, targetClusterID)
			if err != nil {
				return fmt.Errorf(
					"couldn't read target vdisk %s's storage config: %v", targetVdiskID, err)
			}
		}
	}

	// 1. copy the vdisk (meta)data

	switch stype := sourceStaticConfig.Type.StorageType(); stype {
	case config.StorageDeduped:
		err = storage.CopyDeduped(
			sourceVdiskID, targetVdiskID,
			sourceClusterConfig, targetClusterConfig)
	case config.StorageNonDeduped:
		err = storage.CopyNonDeduped(
			sourceVdiskID, targetVdiskID,
			sourceClusterConfig, targetClusterConfig)
	case config.StorageSemiDeduped:
		err = storage.CopySemiDeduped(
			sourceVdiskID, targetVdiskID,
			sourceClusterConfig, targetClusterConfig)
	default:
		err = fmt.Errorf("vdisk %s has an unknown storage type %d",
			sourceVdiskID, stype)
	}
	if err != nil {
		return err
	}

	// 2. copy the tlog data if it is needed

	err = copy.Copy(context.Background(), configSource, copy.Config{
		SourceVdiskID: sourceVdiskID,
		TargetVdiskID: targetVdiskID,
		PrivKey:       vdiskCmdCfg.PrivKey,
		FlushSize:     vdiskCmdCfg.FlushSize,
		JobCount:      vdiskCmdCfg.JobCount,
	})
	if err != nil {
		return fmt.Errorf("failed to copy/generate tlog data for vdisk `%v`: %v", targetVdiskID, err)
	}

	if !sourceStaticConfig.Type.TlogSupport() {
		return nil
	}

	// copy the tlog-specific metadata if both the source and target
	// support tlog and have enabled it
	// TODO: only try to do this if both source and target vdiskID have tlog configured
	// TODO: also fork/copy the actual tlog (meta)data, see https://github.com/zero-os/0-Disk/issues/147
	return tlog.CopyMetadata(sourceVdiskID, targetVdiskID, sourceClusterConfig, targetClusterConfig)
}

func init() {
	VdiskCmd.Long = VdiskCmd.Short + `

If no target storage cluster is given,
the storage cluster configured for the source vdisk
will also be used for the target vdisk.

If an error occured, the target vdisk should be considered as non-existent,
even though data which is already copied is not rolled back.

NOTE: by design,
  only the metadata of a deduped vdisk is copied,
  the data will be copied the first time the vdisk spins up,
  on the condition that the templateStorageCluster has been configured.

WARNING: when copying nondeduped vdisks,
  it is currently not supported that the target vdisk's data cluster
  has more or less storage servers, then the source vdisk's data cluster.
  See issue #206 for more information.
`

	VdiskCmd.Flags().Var(
		&vdiskCmdCfg.SourceConfig, "config",
		"config resource: dialstrings (etcd cluster) or path (yaml file)")
	VdiskCmd.Flags().BoolVar(
		&vdiskCmdCfg.ForceSameStorageCluster, "same", false,
		"enable flag to force copy within the same nbd servers")

	VdiskCmd.Flags().StringVar(
		&vdiskCmdCfg.PrivKey,
		"priv-key", "12345678901234567890123456789012",
		"private key")

	VdiskCmd.Flags().IntVarP(
		&vdiskCmdCfg.JobCount,
		"jobs", "j", runtime.NumCPU(),
		"the amount of parallel jobs to run the tlog generator")

	VdiskCmd.Flags().IntVar(
		&vdiskCmdCfg.FlushSize,
		"flush-size", tlogserver.DefaultConfig().FlushSize,
		"number of tlog blocks in one flush")
}
