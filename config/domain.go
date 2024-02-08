// Domain for configuration

package config

// MongoDB configuration

type MongoDbConfiguration struct {
	Database DB `json:"database"`
}
type DB struct {
	Name       string     `json:"name"`
	Connection Connection `json:"connection"`
}
type Connection struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}
