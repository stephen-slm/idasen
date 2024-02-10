package config

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var defaultConfig = Configuration{
	ConnectionAddress: "",
	LocalName:         "",
	StandHeight:       1.12,
	SitHeight:         0.74,
}

type Configuration struct {
	//  ConnectionAddress changes per a device (Mac, Windows, Linux) but this is
	// the foundational device used to connect to the desk for triggering the
	// standing and sitting.
	ConnectionAddress string `json:"connection_address" yaml:"connection_address"`

	// LocalName defines the localised name for the connected device. Used for
	// displaying if and when the user uses the configuration window.
	LocalName string `json:"local_name" yaml:"local_name"`

	// StandHeight is the configured stand height for the desk.
	StandHeight float64 `json:"stand_height" yaml:"stand_height"`

	// SitHeight is the configured sit height for the desk.
	SitHeight float64 `json:"sit_height" yaml:"sit_height"`
}

// Load attempts to pull the configuration from the given absolute path.
//
// No configuration changes will happen if an error occurred during the loading
// stage.
func Load(absolutePath string) (*Configuration, error) {
	log.WithField("path", absolutePath).Debug("loading configuration")

	var configuration Configuration
	_, err := os.Stat(absolutePath)

	// No reason to continue with the load operation and should use the default
	// values if the file does not exist here. Otherwise the following is just
	// going to fail anyway.
	if os.IsNotExist(err) {
		return &defaultConfig, nil
	}

	yamlFile, err := os.ReadFile(absolutePath)

	if err != nil {
		return &configuration, fmt.Errorf("failed to read file from disk, %w", err)
	}

	if err = yaml.Unmarshal(yamlFile, &configuration); err != nil {
		return &configuration, fmt.Errorf("failed to parse file contents as yaml, %w", err)
	}

	return &configuration, nil
}

// Save attempts to save the configuration to the given absolute path.
func (c *Configuration) Save(absolutePath string) error {
	log.
		WithField("path", absolutePath).
		WithField("configuration", c).
		Debug("saving configuration")

	bytes, err := yaml.Marshal(c)

	if err != nil {
		return fmt.Errorf("failed to marhsal configuration, %w", err)
	}

	file, err := os.Create(absolutePath)
	if err != nil {
		return fmt.Errorf("failed to create file, %w", err)
	}

	if _, err = file.Write(bytes); err != nil {
		return fmt.Errorf("failed to write file contents, %w", err)
	}

	return nil
}
