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

const fileModeForDir = os.FileMode(int(0744))
const fileModeForCreated = os.FileMode(int(0664))

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

		inputFileName := args[0]
		if !strings.HasPrefix(args[0], "/") {
			wd, _ := os.Getwd()
			inputFileName = path.Join(wd, args[0])
		}
		_, inFileName := path.Split(inputFileName)
		inFileName = strings.Split(inFileName, ".")[0]
		outputFileDir := path.Join(cmd.Flag("output").Value.String(), inFileName)
		fmt.Printf("Flatten Input is %s\n", inputFileName)
		fmt.Printf("Flatten Output is %s\n", outputFileDir)

		decoder := yaml.NewDecoder(bytes.NewReader(readYamlInput(inputFileName)))
		// use a Node for demarshalling so that we can keep the whole doc
		// Decode to a YF struct to allow for directing the output
		var yn yaml.Node

		for decoder.Decode(&yn) == nil {
			var yf YF
			err := yn.Decode(&yf)
			if err != nil {
				fmt.Printf("#%v\n", err)
			}
			fmt.Printf("%s with name: %s and layer: %s\n", yf.Schema, yf.Metadata.Name, yf.Metadata.LayeringDefinition.Layer)
			yba, err := yaml.Marshal(&yn)
			if err != nil {
				fmt.Printf("#%v \n", err)
			}
			writeYamlToFolder(outputFileDir, yf.Schema, yf.Metadata.Name, yf.Metadata.LayeringDefinition.Layer, yba)
		}
		fmt.Println("Done.")
	},
}

// Reads the input file into a []byte
func readYamlInput(inputFileName string) (yamlIn []byte) {
	yamlIn, err := ioutil.ReadFile(inputFileName)
	if err != nil {
		fmt.Printf("Yaml file could not be read #%v \n", err)
	}
	return yamlIn
}

// Convert a input string schema to a folder name by converting all "/" characters to "."
func schemaToFolderName(schema string) (folder string) {
	if strings.TrimSpace(schema) == "" {
		folder = "unknown"
	} else {
		folder = strings.ReplaceAll(schema, "/", ".")
	}
	return folder
}

// Create directories using the base path and schema name
func createFolderFromSchema(baseFolderPath string, schema string) {
	os.MkdirAll(path.Join(baseFolderPath, schemaToFolderName(schema)), fileModeForDir)
}

// Creates a file name of <docname>-<doclayer>
func docnameToFileName(docname string, doclayer string) (name string) {
	name = strings.Join([]string{docname, doclayer}, "-")
	return name
}

// Writes a []byte to baseFolderPath/schema(translated)/docname(translated)
func writeYamlToFolder(baseFolderPath string, schema string, docname string, doclayer string, yamlBytes []byte) {
	createFolderFromSchema(baseFolderPath, schema)
	filename := path.Join(baseFolderPath, schemaToFolderName(schema), docnameToFileName(docname, doclayer))
	fullYaml := append(append([]byte("---\n"), yamlBytes...), []byte("...\n")...)
	err := ioutil.WriteFile(filename, fullYaml, fileModeForCreated)
	if err != nil {
		fmt.Printf("#%v \n", err)
	} else {
		fmt.Printf("Wrote %s\n", filename)
	}
}
