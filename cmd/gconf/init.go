package main

import (
    "io"
    "os"

    "github.com/getperf/getperf2/cfg"
    _ "github.com/getperf/getperf2/exp/all"
    "github.com/getperf/getperf2/gcmain"
    "github.com/pkg/errors"
    "github.com/urfave/cli"
)

var config = &cfg.CommonConfig{}

func init() {
    cmdList = append(cmdList, cli.Command{
        Name:  "init",
        Usage: "initialize config file",
        Flags: []cli.Flag{
            cli.StringFlag{
                Name:        "url",
                Usage:       "collect server url or ip",
                Destination: &config.Url,
            },
            cli.StringFlag{
                Name:        "user, u",
                Usage:       "login user name",
                Destination: &config.User,
            },
            cli.StringFlag{
                Name:        "password, p",
                Usage:       "login password",
                Destination: &config.Password,
            },
            cli.IntFlag{
                Name:        "level, l",
                Usage:       "run level",
                Destination: &config.Level,
            },
            cli.IntFlag{
                Name:        "timeout",
                Usage:       "command timeout sec",
                Destination: &config.Timeout,
            },
            cli.StringFlag{
                Name:        "server, s",
                Usage:       "collect server name",
                Destination: &config.Server,
            },
        },
        Action: func(c *cli.Context) error {
            return action(c, &initCommand{Out: os.Stdout})
        },
    })
}

type initCommand struct {
    Out io.Writer
}

func (f *initCommand) Run(c *cli.Context, base *cfg.ExportBase) error {
    controller := gcmain.NewConfigInitializer(base, config)
    if controller == nil {
        return errors.New("config initialization failed")
    }
    if err := controller.Validate(); err != nil {
        return errors.Wrap(err, "init command argument check")
    }
    if err := controller.CreateConfig(); err != nil {
        return errors.Wrap(err, "create config in init command")
    }
    // log.Printf("init base : %v\n", base)
    // log.Printf("init config : %v\n", config)
    // if templateName == "" {
    //     return errors.New("--template required. Usage: gconf --template {xxxconf} init ...")
    // }
    // exp, ok := exporter.Exporters[templateName]
    // if !ok {
    //     return fmt.Errorf("unkown template : %s.", templateName)
    // }
    // fmt.Fprintf(f.Out, exp().Config())
    // server, _ := GetHostname()
    // if c.String("server") != "" {
    //     server = c.String("server")
    // }
    // baseDir := GetBaseDir()
    // configPath := filepath.Join(baseDir, fmt.Sprintf("%s.toml", templateName))
    // if c.String("config") != "" {
    //     configPath = c.String("config")
    // }
    // fmt.Fprintf(f.Out, "CONFIGPATH :%s\n", configPath)
    // exporterServer := exporter.ServerTemplate{
    //     Url:      c.String("url"),
    //     User:     c.String("user"),
    //     Password: c.String("password"),
    //     Level:    c.Int("level"),
    //     Server:   server,
    // }
    // text := exp().Config()
    // tpl, err := template.New("config").Parse(text)
    // if err != nil {
    //     return errors.Wrap(err, "creating config")
    // }
    // nodeFile, err := CreateAndOpenFile(configPath)
    // if err != nil {
    //     return errors.Wrap(err, "creating config")
    // }
    // defer nodeFile.Close()

    // // テンプレートからサーバコンフィグファイル生成
    // err = tpl.Execute(nodeFile, exporterServer)
    // if err != nil {
    //     return errors.Wrap(err, "creating config")
    // }
    // log.Info("config created : ", configPath)

    // fmt.Fprintf(f.Out, "this is init (template=%s, url=%s)\n", templateName, c.String("url"))
    // fmt.Fprintf(f.Out, "server=%v", exporterServer)
    return nil
}

// func InitConfig() error {
//     baseDir := GetBaseDir()
//     log.Info("set base directory ", baseDir)
//     configPath := filepath.Join(baseDir, "windowsconf.toml")
//     return InitConfigFile(configPath)
// }

// func InitConfigFile(configPath string) error {
//     //     if ok, _ := CheckFile(configPath); ok {
//     //         log.Warnf("'%s' exist, Backup to '%s'", configPath, configPath+"_bak")
//     //         if err := CopyFile(configPath, configPath+"_bak"); err != nil {
//     //             return errors.Wrap(err, "backup config file")
//     //         }
//     //     }
//     //     contents := []byte(Config())
//     //     if err := ioutil.WriteFile(configPath, contents, 0664); err != nil {
//     //         return errors.Wrap(err, "write config")
//     //     }
//     return nil
// }
