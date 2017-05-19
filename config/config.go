package config

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	valid "github.com/asaskevich/govalidator"
	"github.com/go-yaml/yaml"
)

// ReadConfig reads the config used to configure the nbdserver
func ReadConfig(path string) (*Config, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("couldn't read nbdserver config: %s", err.Error())
	}

	return FromBytes(bytes)
}

// FromBytes creates a config based on given YAML 1.2 content
func FromBytes(bytes []byte) (*Config, error) {
	cfg := new(Config)

	// unmarshal the yaml content into our config,
	// which will give us basic validation guarantees
	err := yaml.Unmarshal(bytes, cfg)
	if err != nil {
		return nil, fmt.Errorf("couldn't create nbdserver config: %s", err.Error())
	}

	// now apply an extra validation layer,
	// to ensure that the values with generic types,
	// are also valid
	err = cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("couldn't create nbdserver config: %s", err.Error())
	}

	return cfg, nil
}

// Config for the nbdserver backends
type Config struct {
	StorageClusters map[string]StorageClusterConfig `yaml:"storageClusters" valid:"required"`
	Vdisks          map[string]VdiskConfig          `yaml:"vdisks" valid:"required"`
}

// Validate the Config to ensure that all required properties are present,
// and that all given properties are valid
func (cfg *Config) Validate() error {
	if _, err := valid.ValidateStruct(cfg); err != nil {
		return err
	}

	// ensure all referenced storage clusters exist, O(n^2)
	for vdiskID, vdisk := range cfg.Vdisks {
		// collect all storageClusters
		storageClusters := []string{vdisk.Storagecluster}
		if vdisk.RootStorageCluster != "" {
			storageClusters = append(
				storageClusters, vdisk.RootStorageCluster)
		}
		if vdisk.TlogStoragecluster != "" {
			storageClusters = append(
				storageClusters, vdisk.TlogStoragecluster)
		}

		// go through all storageClusters
		for clusterID := range cfg.StorageClusters {
			for index, name := range storageClusters {
				if name == clusterID {
					storageClusters = append(
						storageClusters[:index],
						storageClusters[index+1:]...)
					break
				}
			}
			if len(storageClusters) == 0 {
				break
			}
		}

		// ensure referenced (root) storage cluster exists
		if len(storageClusters) > 0 {
			return fmt.Errorf(
				"vdisk %q references unexisting storage cluster(s): %s",
				vdiskID, strings.Join(storageClusters, ", "))
		}
	}

	// valid
	return nil
}

// StorageClusterConfig defines the config for a storageCluster
type StorageClusterConfig struct {
	DataStorage     []StorageServerConfig `yaml:"dataStorage" valid:"required"`
	MetaDataStorage StorageServerConfig   `yaml:"metadataStorage" valid:"required"`
}

// ParseCSStorageServerConfigStrings allows you to parse a slice of raw dial config strings.
// Dial Config Strings are a simple format used to specify ardb connection configs
// easily as a command line argument.
// The format is as follows: `<ip>:<port>[@<db_index>][,<ip>:<port>[@<db_index>]]`,
// where the db_index is optional, and you can give multiple configs by
// seperating them with a comma.
func ParseCSStorageServerConfigStrings(dialCSConfigString string) (configs []StorageServerConfig, err error) {
	if dialCSConfigString == "" {
		return nil, nil
	}
	dialConfigStrings := strings.Split(dialCSConfigString, ",")

	// convert all connection strings into ConnectionConfigs
	for _, dialConfigString := range dialConfigStrings {
		// remove whitespace around
		dialConfigString = strings.TrimSpace(dialConfigString)

		// trailing commas are allowed
		if dialConfigString == "" {
			continue
		}

		var cfg StorageServerConfig
		parts := strings.Split(dialConfigString, "@")
		if n := len(parts); n < 2 {
			cfg.Address = dialConfigString
		} else {
			cfg.Database, err = strconv.Atoi(parts[n-1])
			if err != nil {
				err = nil // ignore actual error
				n++       // not a valid database, thus probably part of address
			}
			cfg.Address = strings.Join(parts[:n-1], "@")
		}
		if !valid.IsDialString(cfg.Address) {
			err = fmt.Errorf("%s is not a valid storage address", cfg.Address)
			return
		}
		configs = append(configs, cfg)
	}

	return
}

// StorageServerConfig defines the config for a storage server
type StorageServerConfig struct {
	Address  string `yaml:"address" valid:"dialstring,required"`
	Database int    `yaml:"db" valid:"optional"`
}

// VdiskConfig defines the config for a vdisk
type VdiskConfig struct {
	Blocksize          uint64    `yaml:"blockSize" valid:"required"`
	ReadOnly           bool      `yaml:"readOnly" valid:"optional"`
	Size               uint64    `yaml:"size" valid:"required"`
	Storagecluster     string    `yaml:"storageCluster" valid:"required"`
	RootStorageCluster string    `yaml:"rootStorageCluster" valid:"optional"`
	TlogStoragecluster string    `yaml:"tlogStorageCluster" valid:"optional"`
	Type               VdiskType `yaml:"type" valid:"required"`
}

// VdiskType represents the type of a vdisk
type VdiskType string

// MarshalYAML implements yaml.Marshaler.MarshalYAML
func (t *VdiskType) MarshalYAML() (interface{}, error) {
	switch vt := *t; vt {
	case VdiskTypeBoot, VdiskTypeDB, VdiskTypeCache:
		return string(vt), nil
	default:
		return nil, fmt.Errorf("%q is not a valid VdiskType", t)
	}
}

// UnmarshalYAML implements yaml.Unmarshaler.UnmarshalYAML
func (t *VdiskType) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var rawType string
	err = unmarshal(&rawType)
	if err != nil {
		return fmt.Errorf("%q is not a valid VdiskType", t)
	}

	vtype := VdiskType(rawType)
	switch vtype {
	case VdiskTypeBoot, VdiskTypeDB, VdiskTypeCache:
		*t = vtype
	default:
		err = fmt.Errorf("%q is not a valid VdiskType", t)
	}

	return
}

// valid vdisk types
const (
	VdiskTypeNil   = VdiskType("")
	VdiskTypeBoot  = VdiskType("boot")
	VdiskTypeDB    = VdiskType("db")
	VdiskTypeCache = VdiskType("cache")
)

func init() {
	valid.SetFieldsRequiredByDefault(true)
}
