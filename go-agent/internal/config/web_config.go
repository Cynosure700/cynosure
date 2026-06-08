package config

import (
	"os"
)

func LoadWebConfig() (AppConfig, error) {
	fileCfg, err := loadConfigFile()
	if err != nil && !os.IsNotExist(err) {
		return AppConfig{}, err
	}

	appHome, err := resolveAppHome(fileCfg)
	if err != nil {
		return AppConfig{}, err
	}

	llm, err := loadLLMConfig(fileCfg)
	if err != nil {
		return AppConfig{}, err
	}

	runtimeDirs, err := resolveRuntimePaths(appHome, fileCfg)
	if err != nil {
		return AppConfig{}, err
	}

	dbURL := firstNonEmpty(getenv("DATABASE_URL"))
	if dbURL == "" {
		dbURL = buildMySQLDSN(fileCfg)
	}

	esAddresses := parseCSVList(firstNonEmpty(getenv("ES_ADDRESSES"), fileCfg.ESAddresses, "http://1.12.217.28:9200"))

	return AppConfig{
		LLM:                        llm,
		ServerAddr:                 firstNonEmpty(fileCfg.ServerAddr, ":8080"),
		AllowedOrigin:              firstNonEmpty(fileCfg.AllowedOrigin, "http://localhost:5173"),
		DatabaseURL:                dbURL,
		RedisAddr:                  firstNonEmpty(fileCfg.RedisAddr, "1.12.217.28:6379"),
		RedisPassword:              firstNonEmpty(getenv("REDIS_PASSWORD"), "213140"),
		RedisDB:                    fileCfg.RedisDB,
		ESAddresses:                esAddresses,
		ESUsername:                 firstNonEmpty(getenv("ES_USERNAME"), fileCfg.ESUsername),
		ESPassword:                 getenv("ES_PASSWORD"),
		JWTSecret:                  firstNonEmpty(getenv("JWT_SECRET"), "nano-cc-local-secret"),
		AppHome:                    appHome,
		BuiltinSkillsDir:           runtimeDirs.builtinSkillsDir,
		CommandBinDir:              runtimeDirs.commandBinDir,
		CommandScriptDir:           runtimeDirs.commandScriptDir,
		SystemPromptPath:           runtimeDirs.systemPromptPath,
		WorkspaceRoot:              runtimeDirs.workspaceRoot,
		LogsDir:                    runtimeDirs.logsDir,
		WebAllowedTools:            parseCSVList(firstNonEmpty(fileCfg.WebAllowedTools, "load_skill,bash,read_file,write_file,edit_file,todo_write,update_memory")),
		BashAllowOutsideWorkspace:  fileCfg.BashAllowOutsideWorkspace,
		BashAllowDangerousCommands: fileCfg.BashAllowDangerousCommands,
		CookieName:                 firstNonEmpty(fileCfg.SessionCookieName, "nano_cc_session"),
		SessionTTLMinutes:          intOrDefault(fileCfg.SessionTTLMinutes, 60*24*7),
	}, nil
}
