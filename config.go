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


const EditFileInstructions = `# STAR will read this file and update its store with the new values.
# 
# An entry in the store file includes two lines in this file, so STAR
# will expect to read lines from this file in pairs that look like:
#
#   1) http://settlement.arc.nasa.gov/70sArtHiRes/70sArt/art.html
#      Tags: art, NASA, space
#
# Those parts are:
# - At the start of a line (spaces excluded) a number followed by a closing parenthesis
# - The entry, being the string that gets copied, opened, etc
# - At the start of a line (spaces excluded) the word "Tags" followed by a colon
# - The tags, being a comma-separated list
#
# You can remove entries from the store file by deleting the line
# pairs, and you can add entries by creating more.
#
# Lines that start with a pound sign will be ignored.

`



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
		file := createFile(abs_path)
		file.Close()
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



func checkEditor(ed string) string {
	if ed == "" {
		return getEnv("EDITOR", DefaultEditorPath)
	} else {
		return ed
	}
}
