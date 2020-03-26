package glog

// NetworkConfig ...
type NetworkConfig struct {
	Host string `json:"host"`
}

// FileConfig ...
type FileConfig struct {
	Path     string `json:"path"`
	FileName string `json:"file_name"`
	Rotation string `json:"rotation"`
}

// Config ...
type Config struct {
	Type           string         `json:"type"`
	Format         string         `json:"format"`
	DateTimeFormat string         `json:"datetime_format"`
	LogLevel       string         `json:"log_level"`
	File           *FileConfig    `json:"file"`
	TCP            *NetworkConfig `json:"tcp"`
	UDP            *NetworkConfig `json:"udp"`
	Sub            []*Config      `json:"sub_logger,omitempty"`
}
