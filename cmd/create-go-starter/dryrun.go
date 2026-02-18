package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// getFilesForTemplate returns the list of files that would be generated for the
// given template/database/observability combination, without writing anything.
// Used by runDryRun to preview files and by tests to validate file lists.
func getFilesForTemplate(projectPath, projectName, template, database, observabilityLevel string) []FileGenerator {
	switch template {
	case TemplateFull:
		return buildFullFileList(projectPath, projectName, database, observabilityLevel)
	case TemplateMinimal:
		return buildMinimalFileList(projectPath, projectName, database)
	case TemplateGraphQL:
		return buildGraphQLFileList(projectPath, projectName, database)
	default:
		return buildFullFileList(projectPath, projectName, database, observabilityLevel)
	}
}

// runDryRun previews what files would be created without generating any.
// It satisfies AC#1-#6 from Story 10.2:
//   - AC#1: Lists all files with relative paths
//   - AC#2: No files are created on the filesystem
//   - AC#3: Summary includes file count, directory count, template, database, observability
//   - AC#4: Only files for the specified template/database are listed
//   - AC#5: Returns nil error (exit code 0) with "Dry-run completed" message
//   - AC#6: Warns if project directory already exists but continues preview
func runDryRun(projectName, template, database, observabilityLevel string) error {
	projectPath := projectName

	// AC#6: Warn if the project directory already exists (does not block dry-run)
	if _, err := os.Stat(projectPath); err == nil {
		fmt.Printf(Yellow("⚠️  Warning: Directory '%s' already exists.\n"), projectName)
		fmt.Println(Yellow("   In real mode, this would cause an error."))
		fmt.Println(Yellow("   Continuing dry-run preview..."))
		fmt.Println()
	}

	// AC#1, AC#4: Get the file list for this template/database combination
	files := getFilesForTemplate(projectPath, projectName, template, database, observabilityLevel)

	// AC#3: Collect unique directories for the summary
	dirSet := make(map[string]bool)
	for _, f := range files {
		dirSet[filepath.Dir(f.Path)] = true
	}

	// Sort directories for consistent output
	dirs := make([]string, 0, len(dirSet))
	for d := range dirSet {
		dirs = append(dirs, d)
	}
	sort.Strings(dirs)

	// Display dry-run header
	fmt.Println("══════════════════════════════════════════════════════")
	fmt.Println("  DRY-RUN PREVIEW — No files will be created")
	fmt.Println("══════════════════════════════════════════════════════")
	fmt.Printf("  Project:       %s\n", projectName)
	fmt.Printf("  Template:      %s\n", template)
	fmt.Printf("  Database:      %s\n", database)
	fmt.Printf("  Observability: %s\n", observabilityLevel)
	fmt.Println()

	// AC#1: Display each file with relative path prefix
	fmt.Printf("Files that would be created (%d files):\n\n", len(files))
	for _, f := range files {
		fmt.Printf("  📄 %s\n", f.Path)
	}
	fmt.Println()

	// AC#3: Display unique directories
	fmt.Printf("Directories that would be created (%d unique):\n", len(dirs))
	for _, d := range dirs {
		fmt.Printf("  📁 %s/\n", d)
	}
	fmt.Println()

	// AC#5: Footer message confirming no files were created
	fmt.Println("══════════════════════════════════════════════════════")
	fmt.Println("  ✅ Dry-run completed. No files created.")
	fmt.Println("  Run without --dry-run to generate the project.")
	fmt.Println("══════════════════════════════════════════════════════")

	return nil
}
