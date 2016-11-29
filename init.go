package main





// makeInitializer returns a main action function that will
// initialize the user's config and store files. If the user
// specifies a file path on the command line as the first argument,
// that file will be used for their store file.
func makeInitializer(terms []string) func() {
	init := func() {
		checkConfigFile()
		conf := readConfig()
		conf.Store = checkStoreFile(terms[0])
		mergeConfigWithDefaults(&conf)
		saveConfigToFile(&conf)
	}

	return init
}
