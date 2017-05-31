package gxcfg

type DBDocker struct {
	Image     string
	Container string
	DiscSpace string
}

type Database struct {
	driver   string
	dbname   string
	user     string
	password string
	host     string
	port     string
	Docker   DBDocker
}

func (d Database) Driver() string {
	return d.driver
}

func (d Database) Database() string {
	return d.dbname
}

func (d Database) Host() string {
	return d.host
}

func (d Database) Port() string {
	return d.port
}

func (d Database) Username() string {
	return d.user
}

func (d Database) Password() string {
	return d.password
}

type CfgDocker struct {
	Image     string
	Container string

	// Important for building the docker image
	// Insert public content, like website templates or image path.
	Volumes []string

	ProjectPath string // -> /go/src/github.com/maprost/gox
}

type Config struct {
	Name            string
	Port            string
	FullProjectPath string // -> /home/maprost/golang/src/github.com/maprost/gox
	ProjectPath     string // -> src/github.com/maprost/gox
	Database        []Database
	Docker          CfgDocker
	Clients         map[string]string
}

var singleton *Config = nil

func InitConfig(file string, databaseAccess DatabaseAccess, configSearch bool) error {
	cfg, err := loadConfig(file, databaseAccess, configSearch)
	if err != nil {
		return err
	}
	singleton = &cfg
	return nil
}

func GetConfig() *Config {
	return singleton
}
