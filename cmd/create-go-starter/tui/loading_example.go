package tui

// This file contains usage examples for the LoadingModel component.
// Story 10.7 Task 4: Elegant spinners for async operations.
//
// The examples below demonstrate how to use LoadingModel and MultiLoadingModel
// for various project generation scenarios.

// Example 1: Simple single loading operation
//
// Usage:
//   loading := NewLoadingModel(LoadingInitializingGit)
//   program := tea.NewProgram(loading)
//   program.Run()
//
// When complete:
//   loading.SetSuccess() // Shows green checkmark
//   // or
//   loading.SetError()   // Shows red cross

// Example 2: Custom spinner type
//
// Usage:
//   loading := NewLoadingModelWithSpinner(
//       "Processing large dataset...",
//       spinner.Line, // Use line spinner instead of dot
//   )

// Example 3: Dynamic text updates
//
// Usage:
//   loading := NewLoadingModel("Step 1...")
//   // Later, update the text:
//   loading.SetText("Step 2...")
//   loading.SetText("Step 3...")
//   // Finally:
//   loading.SetSuccess()

// Example 4: Multiple sequential operations
//
// Usage:
//   operations := []string{
//       LoadingCreatingDirectories,
//       LoadingGeneratingFiles,
//       LoadingInitializingGit,
//       LoadingRunningGoModTidy,
//   }
//   multi := NewMultiLoadingModel(operations, true) // showCompleted=true
//
//   // As each operation completes:
//   multi.SetCurrentSuccess() // Advances to next operation
//   multi.SetCurrentSuccess() // ...
//
//   // If an operation fails:
//   multi.SetCurrentError() // Shows red cross, stops advancement

// Example 5: Integration with Model.Update() for project generation
//
// This example shows how to integrate LoadingModel into the main TUI Model
// for real-time feedback during project generation.
//
// In model.go:
//   type Model struct {
//       // ... existing fields
//       genLoading LoadingModel // Loading spinner for generation
//   }
//
// In update.go:
//   case GenerationStartedMsg:
//       m.genLoading = NewLoadingModel(LoadingGeneratingFiles)
//       return m, m.genLoading.Init()
//
//   case FileGeneratedMsg:
//       // Update loading text to show current file
//       m.genLoading.SetText(fmt.Sprintf("Creating %s...", msg.FilePath))
//       return m, nil
//
//   case GenerationCompleteMsg:
//       if msg.Success {
//           m.genLoading.SetSuccess()
//       } else {
//           m.genLoading.SetError()
//       }
//       return m, nil
//
//   default:
//       // Update the loading spinner animation
//       var cmd tea.Cmd
//       m.genLoading, cmd = m.genLoading.Update(msg)
//       return m, cmd
//
// In view.go:
//   func (m Model) viewGenerating() string {
//       var b strings.Builder
//       b.WriteString(m.genLoading.View())
//       b.WriteString("\n")
//       // ... rest of view
//       return b.String()
//   }

// Example 6: Multi-stage generation with MultiLoadingModel
//
// This example shows comprehensive project generation feedback.
//
// In model.go:
//   type Model struct {
//       // ... existing fields
//       genMultiLoading MultiLoadingModel
//   }
//
// In update.go:
//   case ConfirmGenerationMsg:
//       operations := []string{
//           LoadingCreatingDirectories,
//           LoadingGeneratingFiles,
//           LoadingInitializingGit,
//           LoadingInstallingDeps,
//           LoadingRunningGoModTidy,
//           LoadingVerifyingInstall,
//       }
//       m.genMultiLoading = NewMultiLoadingModel(operations, true)
//       return m, tea.Batch(
//           m.genMultiLoading.Init(),
//           m.generateProjectCmd(),
//       )
//
//   case StepCompletedMsg:
//       // Called by generator for each step
//       m.genMultiLoading.SetCurrentSuccess()
//       return m, nil
//
//   default:
//       // Update animation
//       var cmd tea.Cmd
//       m.genMultiLoading, cmd = m.genMultiLoading.Update(msg)
//       return m, cmd

// Example 7: Contextual messages for different operations
//
// Use the predefined constants for consistency:
//
// var contextualMessages = map[string]string{
//     "git":        LoadingInitializingGit,
//     "deps":       LoadingInstallingDeps,
//     "tidy":       LoadingRunningGoModTidy,
//     "files":      LoadingGeneratingFiles,
//     "dirs":       LoadingCreatingDirectories,
//     "config":     LoadingWritingConfig,
//     "db":         LoadingSettingUpDatabase,
//     "docker":     LoadingConfiguringDocker,
//     "test":       LoadingRunningTests,
//     "build":      LoadingBuildingProject,
//     "verify":     LoadingVerifyingInstall,
// }
//
// // Later, dynamically select message:
// operation := "git"
// loading := NewLoadingModel(contextualMessages[operation])

// Example 8: Custom colors and styles (advanced)
//
// The LoadingModel uses the global color palette from styles.go:
//   - ColorInfo (blue) for in-progress operations
//   - ColorSuccess (green) for successful completion
//   - ColorError (red) for failures
//
// To customize colors, modify styles.go:
//   InfoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#yourcolor"))
//   SuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#yourcolor"))
//   ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#yourcolor"))
