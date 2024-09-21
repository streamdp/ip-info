package domain

type Config struct {
	HttpPort              int
	GrpcPort              int
	GrpcUseReflection     bool
	DatabaseUrl           string
	HttpWriteTimeout      int
	HttpReadTimeout       int
	HttpReadHeaderTimeout int
	Version               string
}

func (c *Config) IsValid() bool {
	if c.HttpPort < 0 || c.HttpPort > 65535 {
		return false
	}
	if c.GrpcPort < 0 || c.GrpcPort > 65535 {
		return false
	}
	if c.DatabaseUrl == "" {
		return false
	}
	return true
}
