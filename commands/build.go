package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/stevvooe/continuity"
)

var (
	buildCmdConfig struct {
		format string
	}

	mediaTypeAliases = map[string]string{
		"pb":   continuity.MediaTypeManifestV0Protobuf,
		"json": continuity.MediaTypeManifestV0JSON,
	}

	BuildCmd = &cobra.Command{
		Use:   "build <root>",
		Short: "Build a manifest for the provided root",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatalln("please specify a root")
			}

			if v, ok := mediaTypeAliases[buildCmdConfig.format]; ok {
				buildCmdConfig.format = v
			}

			ctx, err := continuity.NewContext(args[0])
			if err != nil {
				log.Fatalf("error creating path context: %v", err)
			}

			m, err := continuity.BuildManifest(ctx)
			if err != nil {
				log.Fatalf("error generating manifest: %v", err)
			}

			p, err := continuity.Marshal(m, buildCmdConfig.format)
			if err != nil {
				log.Fatalf("error marshalling manifest as %s: %v",
					buildCmdConfig.format, err)
			}

			if _, err := os.Stdout.Write(p); err != nil {
				log.Fatalf("error writing to stdout: %v", err)
			}
		},
	}
)

func init() {
	BuildCmd.Flags().StringVar(&buildCmdConfig.format, "format", "pb",
		fmt.Sprintf("specify the output format of the manifest (\"pb\"|%q|\"json\"|%q)",
			continuity.MediaTypeManifestV0Protobuf, continuity.MediaTypeManifestV0JSON))
}
