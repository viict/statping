package configs

import (
	"fmt"
	"github.com/statping/statping/types/errors"
	"github.com/statping/statping/utils"
)

func (d *DbConfig) ConnectionString() string {
	var conn string
	postgresSSL := utils.Params.GetString("POSTGRES_SSLMODE")

	switch d.DbConn {
	case "memory", ":memory:":
		conn = "sqlite3"
		d.DbConn = ":memory:"
		return d.DbConn
	case "sqlite", "sqlite3":
		conn, err := findDbFile(d)
		if err != nil {
			log.Errorln(err)
		}
		d.SqlFile = conn
		log.Infof("SQL database file at: %s", d.SqlFile)
		d.DbConn = "sqlite3"
		return d.SqlFile
	case "mysql":
		host := fmt.Sprintf("%v:%v", d.DbHost, d.DbPort)
		conn = fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=UTC&time_zone=%%27UTC%%27", d.DbUser, d.DbPass, host, d.DbData)
		return conn
	case "postgres":
		conn = fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v timezone=UTC sslmode=%v", d.DbHost, d.DbPort, d.DbUser, d.DbData, d.DbPass, postgresSSL)
		return conn
	}
	return conn
}

func LoadConfigFile(directory string) (*DbConfig, error) {
	p := utils.Params
	log.Infof("Attempting to read config file at: %s/config.yml ", directory)
	utils.Params.SetConfigFile(directory + "/config.yml")

	configs := &DbConfig{
		DbConn:      p.GetString("DB_CONN"),
		DbHost:      p.GetString("DB_HOST"),
		DbUser:      p.GetString("DB_USER"),
		DbPass:      p.GetString("DB_PASS"),
		DbData:      p.GetString("DB_DATABASE"),
		DbPort:      p.GetInt("DB_PORT"),
		Project:     p.GetString("NAME"),
		Description: p.GetString("DESCRIPTION"),
		Domain:      p.GetString("DOMAIN"),
		Email:       p.GetString("EMAIL"),
		Username:    p.GetString("ADMIN_USER"),
		Password:    p.GetString("ADMIN_PASS"),
		Location:    utils.Directory,
		SqlFile:     p.GetString("SQL_FILE"),
	}
	log.WithFields(utils.ToFields(configs)).Debugln("read config file: " + directory + "/config.yml")

	if configs.DbConn == "" {
		return configs, errors.New("Starting in setup mode")
	}
	return configs, nil
}
