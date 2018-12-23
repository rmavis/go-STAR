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
	Action string `yaml:"pipe_to",omitempty`
	Editor string `yaml:"editor",omitempty`
	FilterMode string `yaml:"filter_mode",omitempty`
	PrintLines string `yaml:"print_lines",omitempty`
	SortOrder string `yaml:"sort_order",omitempty`
	Store string `yaml:"store_file",omitempty`
}

const ConfigFileName = "config.yaml"
const DefaultEditorPath = "/usr/bin/vi"
const DefaultFilterMode = "loose"
const DefaultPrintLines = "2"
const DefaultSortOrder = "desc"
const DefaultStoreFileName = "store"

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
		mergeConfigWithDefaults(&conf)
	} else {
		conf = defaultConfig()
	}

	return conf
}

// defaultConfig returns a Config filled with defaults.
func defaultConfig() Config {
	return Config{"", getEnv("EDITOR", DefaultEditorPath), DefaultFilterMode, DefaultPrintLines, DefaultSortOrder, defaultStoreFilePath()}
}

// mergeConfigWithDefaults checks each part of the given Config and
// fills in blanks with defaults.
func mergeConfigWithDefaults(conf *Config) {
	d := defaultConfig()
	conf.Action = checkAction(conf.Action, d.Action)
	conf.Editor = checkEditor(conf.Editor, d.Editor)
	conf.FilterMode = checkFilterMode(conf.FilterMode, d.FilterMode)
	conf.PrintLines = checkPrintLines(conf.PrintLines, d.PrintLines)
	conf.SortOrder = checkSortOrder(conf.SortOrder, d.SortOrder)
	conf.Store = checkStoreFile(conf.Store, d.Store)
}

// checkAction checks if the given action is valid. If so, the string
// is just returned. If not, the default action is returned.
func checkAction(_act string, def string) string {
	if strings.Contains(_act, "~") {
		_act = path.Clean(strings.Replace(_act, "~", userHome(), -1))
	}
	if (len(_act) > 0) {
		return _act
	} else {
		return def
	}
}

// checkEditor is a convenience function for getting the user's text
// editor. If the environment variable is not set, then the default
// specified above will be used.
func checkEditor(ed string, def string) string {
	if ed == "" {
		return def
	} else {
		return ed
	}
}

// checkFilterMode checks if the given filter mode is valid. If so,
// it's returned. If not, the default is returned.
func checkFilterMode(_mode string, def string) string {
	if ((_mode == "loose") || (_mode == "strict")) {
		return _mode
	} else {
		return def
	}
}

// checkPrintLines ensures that the number of lines to print is
// 1 or 2.
func checkPrintLines(num string, def string) string {
	if (num == "1" || num == "2") {
		return num
	} else {
		return def
	}
}

// checkSortOrder ensures that the sort order is valid.
func checkSortOrder(order string, def string) string {
	if (order == "asc" || order == "desc") {
		return order
	} else {
		return def
	}
}

// checkStoreFile ensures that the user's store file exists. It
// returns the given file name's absolute path.
func checkStoreFile(_path string, def string) string {
	var abs_path string

	switch {
	case _path == "":
		abs_path = def
	case strings.Contains(_path, "~"):
		abs_path = path.Clean(strings.Replace(_path, "~", userHome(), -1))
	default:
		abs_path = path.Clean(_path)
	}

	if !doesFileExist(abs_path) {
		file := createFile(abs_path)
		file.Close()
	}

	return abs_path
}

// userHome is a convenience function for getting the user's home.
func userHome() string {
	usr, err := user.Current()
	checkForError(err)

	return usr.HomeDir
}

// mergeConfigActions receives a Config and an action code and
// returns an action code. The returned action code will be a copy
// of the given code but with the intent of the Config merged in.
// Merging can only occur where the given code has 0s.
func mergeConfigActions(conf *Config, action_code []int) []int {
	act := action_code

	if act[1] == 0 {
		if len(conf.Action) > 0 {
			act[1] = 2
		} else {
			act[1] = 1
		}
	}

	if act[2] == 0 {
		if conf.FilterMode == "strict" {
			act[2] = 2
		} else {
			act[2] = 1
		}
	}

	if act[3] == 0 {
		if conf.SortOrder == "asc" {
			act[3] = 2
		} else {
			act[3] = 1
		}
	}

	if act[4] == 0 {
		if conf.PrintLines == "2" {
			act[4] = 2
		} else {
			act[4] = 1
		}
	}

	return act
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

// checkConfigFile ensures that the user's config file exists.
func checkConfigFile() {
	if config_path := configFilePath(); !doesFileExist(config_path) {
		createFile(config_path)
	}
}

// saveConfigToFile writes the given Config to the user's config
// file in the expected YAML format.
func saveConfigToFile(conf *Config) {
	file_name := configFilePath()

	file_handle, err := os.Create(file_name)
	checkForError(err)
	defer file_handle.Close()

	conf_pairs := [][]string{
		{"store_file", conf.Store},
		{"filter_mode", conf.FilterMode},
		{"pipe_to", conf.Action},
		{"editor", conf.Editor},
		{"print_lines", conf.PrintLines}}

	for _, pair := range conf_pairs {
		conf_line := []string{pair[0], ": ", pair[1], "\n"}
		_, err := file_handle.WriteString(strings.Join(conf_line, ""))
		checkForError(err)
	}
}
