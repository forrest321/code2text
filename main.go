package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type appConfig struct {
	IncludeExtensions []string
	ExcludeExtensions []string
	IgnoreDirs        []string
	OutputFile        string
}

func main() {
	config := appConfig{
		IncludeExtensions: []string{".go", ".js"},
		ExcludeExtensions: []string{},
		IgnoreDirs:        []string{".git", ".idea"},
		OutputFile:        "output.txt",
	}

	var rootCmd = &cobra.Command{
		Use:   "code2text",
		Short: "Collects code files and outputs them into a single text document.",
		Long: `A fast and flexible utility to collect code files from a directory tree
and consolidate them into a single text document for analysis.`,
		Run: func(cmd *cobra.Command, args []string) {
			start := time.Now()
			if err := processFiles(".", &config); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			duration := time.Since(start)
			fmt.Printf("Processing completed in %v\n", duration)
		},
	}

	rootCmd.PersistentFlags().StringVarP(&config.OutputFile, "output", "o", config.OutputFile, "Output file name")
	rootCmd.PersistentFlags().StringSliceVarP(&config.IncludeExtensions, "include", "i", config.IncludeExtensions, "File extensions to include")
	rootCmd.PersistentFlags().StringSliceVarP(&config.ExcludeExtensions, "exclude", "e", config.ExcludeExtensions, "File extensions to exclude")
	rootCmd.PersistentFlags().StringSliceVarP(&config.IgnoreDirs, "ignore", "g", config.IgnoreDirs, "Directories to ignore")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func processFiles(startDir string, config *appConfig) error {
	var buffer bytes.Buffer

	// Directory structure output
	dirStructure, err := generateDirStructure(startDir, config.IgnoreDirs)
	if err != nil {
		return err
	}
	buffer.WriteString(dirStructure)

	var totalLines int
	var totalSize int64
	fileCount := 0

	err = filepath.WalkDir(startDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if contains(config.IgnoreDirs, d.Name()) {
				return fs.SkipDir
			}
			return nil
		}
		ext := filepath.Ext(path)
		if contains(config.IncludeExtensions, ext) && !contains(config.ExcludeExtensions, ext) {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			lines := bytes.Count(data, []byte("\n"))
			totalLines += lines
			totalSize += int64(len(data))
			fileCount++
			buffer.WriteString(fmt.Sprintf("File: %s\n", path))
			buffer.WriteString(fmt.Sprintf("Lines: %d\n", lines))
			buffer.WriteString("-----BEGIN FILE-----\n")
			buffer.Write(data)
			buffer.WriteString("\n-----END FILE-----\n\n")
		}
		return nil
	})

	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(config.OutputFile, buffer.Bytes(), 0644); err != nil {
		return err
	}

	fmt.Printf("Total files processed: %d\n", fileCount)
	fmt.Printf("Total lines: %d\n", totalLines)
	fmt.Printf("Total size of data: %d bytes\n", totalSize)

	return nil
}

func generateDirStructure(startDir string, ignoreDirs []string) (string, error) {
	cmd := exec.Command("tree", startDir)
	if output, err := cmd.CombinedOutput(); err == nil {
		return string(output), nil
	}
	return customTree(startDir, "", ignoreDirs)
}

func customTree(dir, prefix string, ignoreDirs []string) (string, error) {
	var result strings.Builder
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if contains(ignoreDirs, file.Name()) {
			continue
		}
		if file.IsDir() {
			result.WriteString(fmt.Sprintf("%s├── %s\n", prefix, file.Name()))
			newPrefix := prefix + "│   "
			subTree, err := customTree(filepath.Join(dir, file.Name()), newPrefix, ignoreDirs)
			if err != nil {
				return "", err
			}
			result.WriteString(subTree)
		} else {
			result.WriteString(fmt.Sprintf("%s├── %s\n", prefix, file.Name()))
		}
	}
	return result.String(), nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}