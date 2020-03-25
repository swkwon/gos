package glog

type NetworkConfig struct {
	Host string `json:"host"`
}

type FileConfig struct {
	Path     string `json:"path"`
	FileName string `json:"file_name"`
	Rotation string `json:"rotation"`
}

type Config struct {
	Type     string         `json:"type"`
	LogLevel string         `json:"log_level"`
	File     *FileConfig    `json:"file"`
	TCP      *NetworkConfig `json:"tcp"`
	UDP      *NetworkConfig `json:"udp"`
}
