package main

import (
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"
	"gopkg.in/yaml.v2"
)



// Config is a structure that contains keys corresponding to each
// valid key in the config file. The default values are the YAML
// commands used to parse/marshal the YAML -- those defaults are
// replaced by the values from the YAML file.
type Config struct {
	Action string `yaml:"pipe_to"`
	Editor string `yaml:"editor"`
	FilterMode string `yaml:"filter_mode"`
	Store string `yaml:"file_name"`  // ,omitempty
}



const ConfigFileName = "config.yaml"
const DefaultEditorPath = "/usr/bin/vi"
const DefaultFilterMode = "loose"
const DefaultStoreFileName = "store"
const OpenPath = "/usr/bin/open"
const PagerPath = "/usr/bin/less"
const PbcopyPath = "/usr/bin/pbcopy"





// defaultConfigPath returns the path to the directory containing
// the user's config file.
func defaultConfigPath() string {
	return userHome() + "/.config/star";
}



// configFilePath returns the path to the user's config file.
func configFilePath() string {
	return defaultConfigPath() + "/" + ConfigFileName;
}



// defaultStoreFilePath returns the path to the user's default store
// file. The default should only be used when an alternative location
// is not specified in the config file.
func defaultStoreFilePath() string {
	return defaultConfigPath() + "/" + DefaultStoreFileName;
}



// readConfig checks for the user's config file. If it exists, then
// it will be read and transformed into a Config. If it doesn't, then
// the default Config will be returned instead.
func readConfig() Config {
	var conf Config
	conf_path := configFilePath()

	if doesFileExist(conf_path) {
		cont, err := ioutil.ReadFile(conf_path)
		checkForError(err)
		yaml.Unmarshal(cont, &conf)
	} else {
		conf = defaultConfig()
	}

	checkConfig(&conf)

	return conf
}



// defaultConfig returns a Config filled with defaults.
func defaultConfig() Config {
	return Config{PbcopyPath, getEnv("EDITOR", DefaultEditorPath), DefaultFilterMode, defaultStoreFilePath()}
}



// checkConfig checks the required parts of the user's Config.
func checkConfig(conf *Config) {
	// conf.Action = checkAction(conf.Action)
	conf.Store = checkStoreFile(conf.Store)
	conf.FilterMode = checkFilterMode(conf.FilterMode)
}



// checkAction checks if the given action is valid. If so, the string
// is just returned. If not, the default action is returned.
func checkAction(_act string) string {
	// Note that `pbcopy` is the default action.
	if _act == "" {
		return PbcopyPath
	} else {
		return _act
	}
}



// checkStoreFile ensures that the user's store file exists. It
// returns the given file name's absolute path.
func checkStoreFile(_path string) string {
	var abs_path string

	if strings.Contains(_path, "~") {
		abs_path = path.Clean(strings.Replace(_path, "~", userHome(), -1))
	} else {
		abs_path = path.Clean(_path)
	}

	if !doesFileExist(abs_path) {
		file := createFile(abs_path)
		file.Close()
	}

	return abs_path
}



// checkFilterMode checks if the given filter mode is valid. If so,
// it's returned. If not, the default is returned.
func checkFilterMode(_mode string) string {
	if ((_mode == "loose") || (_mode == "strict")) {
		return _mode
	} else {
		return DefaultFilterMode
	}
}



// checkEditor is a convenience function for getting the user's text
// editor. If the environment variable is not set, then the default
// specified above will be used.
func checkEditor(ed string) string {
	if ed == "" {
		return getEnv("EDITOR", DefaultEditorPath)
	} else {
		return ed
	}
}



// userHome is a convenience function for getting the user's home.
func userHome() string {
	usr, err := user.Current()
	checkForError(err)

	return usr.HomeDir
}



// getEnv checks if the given environment variable is set. If so, its
// value is returned. If not, then the given default is returned.
func getEnv(env_var string, _default string) string {
	ed := os.Getenv(env_var)

	if ed == "" {
		return _default
	} else {
		return ed
	}
}
