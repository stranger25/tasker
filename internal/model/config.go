package model

// Config is a struct that contains the configuration for a server.
type Config struct {
	Postgres PostgresConfig `yaml:"postgres"` // Database contains the configuration for a database.
	Redis    RedisConfig    `yaml:"redis"`    // Redis contains the configuration for a Redis server.
	Server   ServerConfig   `yaml:"server"`   // Server contains the configuration for a server.
}

// PostgresConfig is a struct that contains the configuration for a database.
type PostgresConfig struct {
	DSN string `yaml:"dsn"` // DSN is the data source name for the database.
}

// RedisConfig is a struct that contains the configuration for a Redis server.
type RedisConfig struct {
	Addr     string `yaml:"addr"`     // Addr is the address of the Redis server.
	Password string `yaml:"password"` // Password is the password for the Redis server.
	DB       int    `yaml:"db"`       // DB is the database number for the Redis server.
}

// ServerConfig is a struct that contains the configuration for a server.
type ServerConfig struct {
	Port string `yaml:"port"` // Port is the port number for the server.
	Storage string `yaml:"storage"` // Repo is the datasource for task storage.
}
