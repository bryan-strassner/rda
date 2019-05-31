package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var InputFile string
var OutputDir string

func init() {
	rootCmd.AddCommand(flattenCmd)
	wd, _ := os.Getwd()
	flattenCmd.PersistentFlags().StringVarP(&OutputDir, "output", "o", wd, "Full path of the root directory as output")
	flattenCmd.MarkPersistentFlagFilename("output")
}

type YF struct {
	Schema   string
	Metadata struct {
		Name               string
		LayeringDefinition struct {
			Layer string
		} `yaml:"layeringDefinition"`
	}
	Data struct{}
}

var flattenCmd = &cobra.Command{
	Use:   "flatten",
	Short: "Flatten a set of rendered documents to folders",
	Long: `Converts a set of rendered documents to a set of folders arranged
as folders per schema`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Resolve input and output, such that the input file name is either
		// absolute or relative and ends up being absolute

		//TODO: encapsulate things into funcs
		inputFileName := args[0]
		if !strings.HasPrefix(args[0], "/") {
			wd, _ := os.Getwd()
			inputFileName = path.Join(wd, args[0])
		}
		_, inFileName := path.Split(inputFileName)
		inFileName = strings.Split(inFileName, ".")[0]
		outputFileName := path.Join(cmd.Flag("output").Value.String(), inFileName)
		fmt.Printf("Flatten Input is %s\n", inputFileName)
		fmt.Printf("Flatten Output is %s\n", outputFileName)

		file, _ := os.Open(inputFileName)
		file.Close()
		yamlIn, err := ioutil.ReadFile(inputFileName)
		if err != nil {
			fmt.Printf("Yaml file could not be read #%v ", err)
		}

		decoder := yaml.NewDecoder(bytes.NewReader(yamlIn))
		// use a Node for demarshalling so that we can keep the whole doc
		// Decode to a YF struct to allow for directing the output
		var yn yaml.Node
		var yf YF

		for decoder.Decode(&yn) == nil {
			err := yn.Decode(&yf)
			if err != nil {
				fmt.Printf("#%v", err)
			}
			fmt.Printf("%s with name: %s and layer: %s\n", yf.Schema, yf.Metadata.Name, yf.Metadata.LayeringDefinition.Layer)
			//TODO: Create directories based on the schema
			//TODO: Write files based on the name-layer
			yba, err := yaml.Marshal(&yn)
			if err != nil {
				fmt.Printf("#%v", err)
			}
			fmt.Println(string(yba))
		}
		fmt.Println("Done.")

	},
}
