package golang

import (
	"io/ioutil"

	"github.com/maprost/gox/gxarg"
	"github.com/maprost/gox/gxcfg"
	"github.com/maprost/gox/internal/docker"
	"github.com/maprost/gox/internal/log"
	"github.com/maprost/gox/internal/shell"
)

func CompileInDocker(cfg *gxcfg.Config) error {
	err := RemoveDockerContainer(cfg)
	if err != nil {
		return err
	}

	log.Info("Build project ", cfg.Name, " in docker container.")
	dock := docker.NewRunBuilder(cfg.Docker.Container, cfg.Docker.Image)

	// add project
	dock.Value(cfg.FullProjectPath, cfg.Docker.ProjectPath)

	// add command
	dock.Execute("cd " + cfg.Docker.ProjectPath +
		" && go fmt ./..." +
		" && go build -o " + BinaryGxName() +
		" && chmod o+w " + BinaryGxName())

	_, err = dock.Run(log.LevelInfo)
	if err != nil {
		return err
	}

	return RemoveDockerContainer(cfg)
}

func CompileBinary() (err error) {
	_, err = shell.Command("go", "fmt", "./...")
	if err != nil {
		return
	}

	_, err = shell.Command("go", "get", "-d")
	if err != nil {
		return
	}

	_, err = shell.Stream(log.LevelInfo, "go", "build", "-o", BinaryName())
	return
}

func TestInDocker(cfgFile string, cfg *gxcfg.Config) error {
	err := RemoveDockerContainer(cfg)
	if err != nil {
		return err
	}

	log.Info("Test project ", cfg.Name, " in docker container.")

	dock := docker.NewRunBuilder(cfg.Docker.Container, cfg.Docker.Image)

	// add project
	dock.Value(cfg.FullProjectPath, cfg.Docker.ProjectPath)

	// add database
	for _, db := range cfg.Database {
		dock.Link(db.Docker.Container, db.Docker.Container)
	}

	// add command
	dock.Execute("cd " + cfg.Docker.ProjectPath +
		" && touch " + gxcfg.FileInsideDockerContainer +
		" && chmod o+w " + gxcfg.FileInsideDockerContainer +
		" && go test -cover ./... -args -" + gxarg.Config + "=" + cfgFile)

	_, err = dock.Run(log.LevelInfo)
	// delete FileInsideDockerContainer file
	defer shell.Command("rm", gxcfg.FileInsideDockerContainer)

	if err != nil {
		return err
	}

	return RemoveDockerContainer(cfg)
}

func BuildDockerImage(cfgFile string) error {
	cfg := gxcfg.GetConfig()
	log.Info("Build docker image: ", cfg.Docker.Container)

	err := docker.RemoveImage(cfg.Docker.Container)
	if err != nil {
		return err
	}

	fileContent := "From " + cfg.Docker.Image + "\n\n" +
		"COPY " + BinaryGxName() + " " + cfg.Docker.ProjectPath + "\n\n" +
		"COPY " + cfgFile + " " + cfg.Docker.ProjectPath + "\n\n" +
		"RUN touch " + gxcfg.FileInsideDockerContainer + " && mv " + gxcfg.FileInsideDockerContainer + " " + cfg.Docker.ProjectPath + " \n\n"

	// add volume
	for _, v := range cfg.Docker.Volumes {
		fileContent += "COPY " + v + " " + cfg.Docker.ProjectPath + "/" + v + "\n\n"
	}

	// add entry point
	fileContent += "ENTRYPOINT [\"" + cfg.Docker.ProjectPath + "/" + BinaryGxName() + " -" + gxarg.Config + "=" + cfgFile + "\"]" + "\n"
	err = ioutil.WriteFile("DockerFile", []byte(fileContent), 0644)
	if err != nil {
		return err
	}

	_, err = shell.Stream(log.LevelInfo, "docker", "build", "-t", cfg.Docker.Container, "-f", "./DockerFile", ".")
	return err
}

func RemoveDockerContainer(cfg *gxcfg.Config) error {
	return docker.StopAndRemove(cfg.Docker.Container)
}

func PullDockerImage() error {
	return docker.Pull(gxcfg.GetConfig().Docker.Image)
}

func BinaryGxName() string {
	return BinaryName() + "_gx"
}

func BinaryName() string {
	return gxcfg.GetConfig().Name
}
