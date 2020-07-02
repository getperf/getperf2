package agent

// import (
// 	"bytes"
// 	"context"
// 	"os"
// 	"testing"
// 	// "github.com/getperf/getperf2"
// )

// func TestRun(t *testing.T) {
// 	got := new(bytes.Buffer)
// 	// err := getperf2.Run(context.Background(), []string{"-version"}, os.Stdout, os.Stderr)
// 	err := Run(context.Background(), []string{"-version"}, got, os.Stderr)
// 	t.Log(got)
// 	t.Log(len(got.String()))
// 	if err != nil && len(got.String()) < 27 {
// 		t.Error(`Run("detartrated") = false`)
// 	}
// }

// const cmdName = "getperf2"

// // Run the getperf2
// func Run(ctx context.Context, argv []string, outStream, errStream io.Writer) error {
// 	log.SetOutput(errStream)
// 	fs := flag.NewFlagSet(
// 		fmt.Sprintf("%s (v%s rev:%s)", cmdName, version, revision), flag.ContinueOnError)
// 	fs.SetOutput(errStream)
// 	ver := fs.Bool("version", false, "display version")
// 	if err := fs.Parse(argv); err != nil {
// 		return err
// 	}
// 	if *ver {
// 		return printVersion(outStream)
// 	}
// 	return nil
// }

// func printVersion(out io.Writer) error {
// 	_, err := fmt.Fprintf(out, "%s v%s (rev:%s)\n", cmdName, version, revision)
// 	return err
// }
