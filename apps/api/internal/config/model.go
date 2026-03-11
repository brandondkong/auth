package config

type Config struct {
	DatabaseUrl				string		`env:"DATABASE_URL"`
	JwtRefreshSigningKey	string		`env:"JWT_REFRESH_SIGNING_KEY"`
	JwtAccessSigningKey		string		`env:"JWT_ACCESS_SIGNING_KEY"`
}

