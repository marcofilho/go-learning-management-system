package config

func GetEnvForTest(key, defaultValue string) string {
	return getEnv(key, defaultValue)
}
