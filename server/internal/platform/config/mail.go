package config

func LoadMailConfig() (MailConfig, error) {
	port, err := GetEnvInt("SMTP_PORT", 587)
	if err != nil {
		return MailConfig{}, err
	}

	return MailConfig{
		Host: GetEnv("SMTP_HOST", ""),
		Port: port,
		User: GetEnv("SMTP_USER", ""),
		Pass: GetEnv("SMTP_PASS", ""),
		From: GetEnv("SMTP_FROM", ""),
	}, nil
}
