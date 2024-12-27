package agent

import "fmt"

const Version = "2.18.0"

var Revision = "HEAD"

const HelpMessage = `
Setup:
getperf2 setup [-c <config>] [-k <sitekey>] [-p <pass>] [-u <url>]

Service control:
getperf2 [start|stop] (Windows only: [install|remove])

Manual start:
getperf2 
  -b, -background
        Agent run as background service. (default true)
  -c, -config <config>
        It performs by the specified directory.
  -f, -foreground
        Agent run as foreground service.
  -s ,-statid <statid>
        Agent run the specified category once.
`

func VersionMessage() string {
	return fmt.Sprintf("%s v%s (rev:%s)\n", cmdName, Version, Revision)
}

func PrintUsage() error {
	fmt.Printf("Getperf Cacti agent %s\n%s", VersionMessage(), HelpMessage)
	// fmt.Println(HelpMessage)
	return nil
}
