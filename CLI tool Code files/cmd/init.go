package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/telugusmasher2010-collab/project-cli-tool/internal/generator"
	"github.com/telugusmasher2010-collab/project-cli-tool/internal/output"
)

type templateOption struct {
	Name        string
	Description string
	Stack       string
}

var templates = []templateOption{
	{Name: "tauri-llm", Description: "Tauri v2 + Rust + React + local LLM sidecar", Stack: "Tauri/Rust/React"},
	{Name: "whatsapp-bot", Description: "Node.js + Baileys + SQLite + Fastify", Stack: "Node.js"},
	{Name: "expense-splitter", Description: "Flutter + Dart + Supabase + UPI", Stack: "Flutter/Dart"},
	{Name: "next-webapp", Description: "Next.js 15 + Prisma + Tailwind + Auth", Stack: "Next.js"},
	{Name: "react-native-map", Description: "React Native + Expo + MapLibre", Stack: "React Native"},
	{Name: "cli-go", Description: "Minimal Go CLI with cobra", Stack: "Go"},
}

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

func selectTemplate(reader *bufio.Reader) (*templateOption, error) {
	fmt.Println("\nAvailable templates:")
	for i, t := range templates {
		fmt.Printf("  %d. %s — %s (%s)\n", i+1, t.Name, t.Description, t.Stack)
	}

	for {
		fmt.Print("\nSelect template (1-" + fmt.Sprintf("%d", len(templates)) + "): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		input = strings.TrimSpace(input)

		var idx int
		if _, err := fmt.Sscanf(input, "%d", &idx); err != nil || idx < 1 || idx > len(templates) {
			fmt.Fprintf(os.Stderr, "Error: enter a number between 1 and %d\n", len(templates))
			continue
		}
		return &templates[idx-1], nil
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
