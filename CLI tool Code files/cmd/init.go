package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/telugusmasher2010-collab/project-cli-tool/internal/generator"
	"github.com/telugusmasher2010-collab/project-cli-tool/internal/output"
	"github.com/telugusmasher2010-collab/project-cli-tool/internal/templates"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Scaffold a new project interactively",
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		projectName, err := promptInput(reader, "Project name: ", validateProjectName)
		if err != nil {
			return err
		}

		t, err := selectTemplate(reader)
		if err != nil {
			return err
		}

		outputPath, err := promptInput(reader, fmt.Sprintf("Output path (default: ./%s): ", projectName), validatePath)
		if err != nil {
			return err
		}
		if outputPath == "" {
			outputPath = "./" + projectName
		}

		output.Infof("Scaffolding %s with %s template...", projectName, t.Name)
		s := output.NewSpinner("Generating project...")
		s.Start()

		vars := generator.NewVariables()
		vars.Set("ProjectName", projectName)
		vars.Set("GoModule", "github.com/user/"+projectName)

		gen := generator.New(outputPath, vars, generator.Options{Overwrite: false})
		if err := gen.Generate(t.Name); err != nil {
			s.Stop()
			output.Cross(fmt.Sprintf("Failed: %v", err))
			return err
		}

		s.Stop()
		output.Check(fmt.Sprintf("Project %s created at %s", projectName, outputPath))
		output.Infof("Run: cd %s && npm install && git init", outputPath)
		return nil
	},
}

func promptInput(reader *bufio.Reader, prompt string, validate func(string) error) (string, error) {
	for {
		fmt.Print(prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		input = strings.TrimSpace(input)
		if validate != nil {
			if err := validate(input); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				continue
			}
		}
		return input, nil
	}
}

func selectTemplate(reader *bufio.Reader) (*templates.TemplateInfo, error) {
	available := templates.List()
	fmt.Println("\nAvailable templates:")
	for i, t := range available {
		fmt.Printf("  %d. %s — %s\n", i+1, t.Name, t.Description)
	}

	for {
		fmt.Print("\nSelect template (1-" + fmt.Sprintf("%d", len(available)) + "): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		input = strings.TrimSpace(input)

		var idx int
		if _, err := fmt.Sscanf(input, "%d", &idx); err != nil || idx < 1 || idx > len(available) {
			fmt.Fprintf(os.Stderr, "Error: enter a number between 1 and %d\n", len(available))
			continue
		}
		return &available[idx-1], nil
	}
}

func validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}
	for _, r := range name {
		if !isAlphaNum(r) && r != '-' && r != '_' {
			return fmt.Errorf("project name must only contain letters, numbers, hyphens, and underscores")
		}
	}
	return nil
}

func validatePath(path string) error {
	if path == "" {
		return nil
	}
	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		return fmt.Errorf("path exists and is not a directory")
	}
	return nil
}

func isAlphaNum(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

func init() {
	rootCmd.AddCommand(initCmd)
}
