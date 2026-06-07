package config

import (
	"fmt"
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

	llm, err := loadLLMConfig()
	if err != nil {
		return AppConfig{}, err
	}

	runtimeDirs, err := resolveRuntimePaths(appHome, fileCfg)
	if err != nil {
		return AppConfig{}, err
	}
	dbURL := getenv("DATABASE_URL", buildDefaultMySQLDSN())
	redisDB, err := getenvInt("REDIS_DB", 0)
	if err != nil {
		return AppConfig{}, fmt.Errorf("invalid REDIS_DB: %w", err)
	}

	return AppConfig{
		LLM:                        llm,
		ServerAddr:                 getenv("SERVER_ADDR", ":8080"),
		AllowedOrigin:              getenv("ALLOWED_ORIGIN", "http://localhost:5173"),
		DatabaseURL:                dbURL,
		RedisAddr:                  getenv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword:              getenv("REDIS_PASSWORD", ""),
		RedisDB:                    redisDB,
		JWTSecret:                  getenv("JWT_SECRET", "nano-cc-local-secret"),
		AppHome:                    appHome,
		BuiltinSkillsDir:           runtimeDirs.builtinSkillsDir,
		CommandBinDir:              runtimeDirs.commandBinDir,
		CommandScriptDir:           runtimeDirs.commandScriptDir,
		SystemPromptPath:           runtimeDirs.systemPromptPath,
		WorkspaceRoot:              runtimeDirs.workspaceRoot,
		WebAllowedTools:            parseCSVList(getenv("WEB_ALLOWED_TOOLS", firstNonEmpty(fileCfg.WebAllowedTools, "load_skill,bash,read_file,write_file,edit_file,todo_write"))),
		BashAllowOutsideWorkspace:  getenvBool("BASH_ALLOW_OUTSIDE_WORKSPACE", fileCfg.BashAllowOutsideWorkspace),
		BashAllowDangerousCommands: getenvBool("BASH_ALLOW_DANGEROUS_COMMANDS", fileCfg.BashAllowDangerousCommands),
		CookieName:                 getenv("SESSION_COOKIE_NAME", "nano_cc_session"),
		SessionTTLMinutes:          getenvIntOrDefault("SESSION_TTL_MINUTES", 60*24*7),
	}, nil
}
