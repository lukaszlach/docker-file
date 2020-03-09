package main

import (
	"io/ioutil"
	composeloader "github.com/docker/cli/cli/compose/loader"
	composetypes "github.com/docker/cli/cli/compose/types"
	dockerfileparser "github.com/moby/buildkit/frontend/dockerfile/parser"
	dockerfileinstructions "github.com/moby/buildkit/frontend/dockerfile/instructions"
	"os"
	"github.com/urfave/cli/v2"
	"log"
	"encoding/json"
	dockerfilejsonparser "github.com/keilerkonzept/dockerfile-json/pkg/dockerfile"
	"fmt"
	dockercliopts "github.com/docker/cli/opts"
	"github.com/docker/docker/builder/dockerignore"
	"github.com/docker/docker/pkg/fileutils"
	"github.com/ghodss/yaml"
)

func main() {
	app := &cli.App{
		Name:      "docker-file",
		Usage:     "Tools to handling Docker-related files",
		Copyright: "2020 Łukasz Lach https://lach.dev",
		Commands: []*cli.Command{
			{
				Name:    "compose",
				Aliases: []string{"c"},
				Usage:   "Docker Compose file handlers",
				Subcommands: []*cli.Command{
					{
						Name:    "parse",
						Aliases: []string{"p"},
						Usage:   "Parse Docker Compose YAML and convert to JSON",
						ArgsUsage: "FILE",
						Action: func(c *cli.Context) error {
							composeJson, err := loadComposeFile(c.Args().Get(0))
							if err != nil {
								fmt.Println(err)
								os.Exit(1)
							}
							fmt.Println(string(composeJson))
							return nil
						},
					},
				},
			},
			{
				Name:    "dockerfile",
				Aliases: []string{"d"},
				Usage:   "Dockerfile handlers",
				Subcommands: []*cli.Command{
					{
						Name:    "parse",
						Aliases: []string{"p"},
						Usage:   "Parse Dockerfile and convert to JSON",
						ArgsUsage: "FILE",
						Action: func(c *cli.Context) error {
							err := loadDockerfile(c.Args().Get(0))
							if err != nil {
								fmt.Println(err)
								os.Exit(1)
							}
							return nil
						},
					},
				},
			},
			{
				Name:    "env",
				Aliases: []string{"e"},
				Usage:   ".env file handlers",
				Subcommands: []*cli.Command{
					{
						Name:    "parse",
						Aliases: []string{"p"},
						Usage:   "Parse .env file and convert to JSON",
						ArgsUsage: "FILE",
						Action: func(c *cli.Context) error {
							err := loadEnvFile(c.Args().Get(0))
							if err != nil {
								fmt.Println(err)
								os.Exit(1)
							}
							return nil
						},
					},
				},
			},
			{
				Name:    "dockerignore",
				Aliases: []string{"i"},
				Usage:   ".dockerignore file handlers",
				Subcommands: []*cli.Command{
					{
						Name:    "parse",
						Aliases: []string{"p"},
						Usage:   "Parse .dockerignore file and convert to JSON",
						ArgsUsage: "FILE",
						Action: func(c *cli.Context) error {
							err := dumpIgnoreFileJson(c.Args().Get(0))
							if err != nil {
								fmt.Println(err)
								os.Exit(1)
							}
							return nil
						},
					},
					{
						Name:    "check",
						Aliases: []string{"c"},
						Usage:   "Check if file matches .dockerignore rules",
						ArgsUsage: "IGNORE_FILE FILE",
						Action: func(c *cli.Context) error {
							isIgnored, err := checkIgnoreFile(c.Args().Get(0), c.Args().Get(1))
							if err != nil {
								fmt.Println(err)
								os.Exit(1)
							}
							if isIgnored {
								fmt.Println("true")
							} else {
								fmt.Println("false")
							}
							return nil
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func loadDockerfile(filePath string) (error) {
	fileHandle, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fileHandle.Close()
	result, err := dockerfileparser.Parse(fileHandle)
	if err != nil {
		return err
	}
	_, _, err = dockerfileinstructions.Parse(result.AST)
	if err == nil {
		dockerfileJson, err2 := dockerfilejsonparser.Parse(filePath)
		if err2 == nil {
			jsonOut := json.NewEncoder(os.Stdout)
			jsonOut.Encode(dockerfileJson)
			return nil
		}
	}
	return err
}

func loadComposeFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	yamlData, err := composeloader.ParseYAML(data)
	if err != nil {
		return nil, err
	}
	_, err = composeloader.Load(composetypes.ConfigDetails{
		ConfigFiles: []composetypes.ConfigFile{
			{Config: yamlData, Filename: filePath},
		},
	})
	if err != nil {
		return nil, err
	}
	return yaml.YAMLToJSON(data)
}

func loadEnvFile(filePath string) (error) {
	envVars, err := dockercliopts.ParseEnvFile(filePath)
	if err != nil {
		return err
	}
	envVarsMap := dockercliopts.ConvertKVStringsToMapWithNil(envVars)
	jsonOut := json.NewEncoder(os.Stdout)
	jsonOut.Encode(envVarsMap)
	return nil
}

func loadIgnoreFile(filePath string) ([]string, error) {
	fileHandle, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer fileHandle.Close()
	ignorePatterns, err := dockerignore.ReadAll(fileHandle)
	if err != nil {
		return nil, err
	}
	return ignorePatterns, nil
}

func dumpIgnoreFileJson(filePath string) (error) {
	ignorePatterns, err := loadIgnoreFile(filePath)
	if err != nil {
		return err
	}
	jsonOut := json.NewEncoder(os.Stdout)
	jsonOut.Encode(ignorePatterns)
	return nil
}

func checkIgnoreFile(ignoreFilePath string, filePath string) (bool, error) {
	fileHandle, err := os.Open(ignoreFilePath)
	if err != nil {
		return false, err
	}
	defer fileHandle.Close()
	ignorePatterns, err := dockerignore.ReadAll(fileHandle)
	if err != nil {
		return false, err
	}
	ignoreMatches, err := fileutils.Matches(filePath, ignorePatterns)
	if err != nil {
		return false, err
	}
	return ignoreMatches, nil
}
