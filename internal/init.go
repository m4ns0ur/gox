package internal

import (
	"flag"

	"github.com/maprost/gox/gxcfg"
	"github.com/maprost/gox/internal/args"
	"github.com/maprost/gox/internal/log"
)

type initCommand struct {
	hdd  *bool
	log  *string
	file *string
}

func InitCommand() args.SubCommand {
	return &initCommand{}
}

func (cmd *initCommand) Name() string {
	return "init"
}

func (cmd *initCommand) DefineFlags(fs *flag.FlagSet) {
	cmd.hdd = fs.Bool("-hdd", false, "use ")
	cmd.log = args.LogFlag(fs)
	cmd.file = args.FileFlag(fs)
}

func (cmd *initCommand) Run() {
	var err error
	cfgFile := "config.gox"

	log.Info("Init go project.")

	// load config file
	err = gxcfg.InitConfig(cfgFile, gxcfg.DatabaseAccessLink)
	if err != nil {
		log.Fatal("Can't init config: ", err.Error())
	}

}
