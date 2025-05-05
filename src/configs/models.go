package configs

type Config struct {
	Server ConfigServer
	Logger ConfigLogger
}

type ConfigServer struct {
	Host string
	Port int
}

type ConfigLogger struct {
	Type  string
	Level string
}
