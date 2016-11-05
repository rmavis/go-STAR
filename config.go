package main

import (
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"
	"gopkg.in/yaml.v2"
)



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



func userHome() string {
	usr, err := user.Current()
	checkForError(err)

	return usr.HomeDir
}



func DefaultConfigPath() string {
	return userHome() + "/.config/star";
}



func DefaultStoreFilePath() string {
	return DefaultConfigPath() + "/" + DefaultStoreFileName;
}



func ConfigFilePath() string {
	return DefaultConfigPath() + "/" + ConfigFileName;
}



func defaultConfig() Config {
	return Config{PbcopyPath, getEnv("EDITOR", DefaultEditorPath), DefaultFilterMode, DefaultStoreFilePath()}
}



// func getDefaultEditor() string {
// 	ed := os.Getenv("EDITOR")

// 	if ed == "" {
// 		return DefaultEditorPath
// 	} else {
// 		return ed
// 	}
// }

func getEnv(env_var string, _default string) string {
	ed := os.Getenv(env_var)

	if ed == "" {
		return _default
	} else {
		return ed
	}
}



func readConfig() Config {
	var conf Config
	conf_path := ConfigFilePath()

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



func checkConfig(conf *Config) {
	// conf.Action = checkAction(conf.Action)
	conf.Store = checkStoreFile(conf.Store)
	conf.FilterMode = checkFilterMode(conf.FilterMode)
}



func checkAction(_act string) string {
	// Note that `pbcopy` is the default action.
	if _act == "" {
		return PbcopyPath
	} else {
		return _act
	}
}



func checkStoreFile(_path string) string {
	var abs_path string

	if strings.Contains(_path, "~") {
		abs_path = path.Clean(strings.Replace(_path, "~", userHome(), -1))
	} else {
		abs_path = path.Clean(_path)
	}

	if !doesFileExist(abs_path) {
		createFile(abs_path)
	}

	return abs_path
}



func checkFilterMode(_mode string) string {
	if ((_mode == "loose") || (_mode == "strict")) {
		return _mode
	} else {
		return DefaultFilterMode
	}
}
