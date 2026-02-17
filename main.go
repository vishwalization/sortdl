package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	green  = "\033[32m"
	cyan   = "\033[36m"
	yellow = "\033[33m"
	gray   = "\033[90m"
	reset  = "\033[0m"
)

var categories = map[string][]string{
	"Images":    {".jpg", ".jpeg", ".png", ".gif", ".webp", ".heic", ".tiff", ".bmp"},
	"Documents": {".pdf", ".doc", ".docx", ".txt", ".md", ".rtf", ".xlsx", ".xls", ".pptx", ".ppt", ".csv"},
	"Videos":    {".mp4", ".mov", ".avi", ".mkv", ".webm", ".flv", ".wmv"},
	"Music":     {".mp3", ".wav", ".m4a", ".flac", ".aac"},
	"Archives":  {".zip", ".rar", ".tar", ".gz", ".7z", ".dmg"},
}

func getCategory(ext string) string {
	ext = strings.ToLower(ext)
	for cat, exts := range categories {
		for _, e := range exts {
			if ext == e {
				return cat
			}
		}
	}
	return "Others"
}

func main() {
	dryRun := flag.Bool("dry-run", false, "Show what would happen without moving files")
	flag.Parse()

	home, _ := os.UserHomeDir()
	downloads := filepath.Join(home, "Downloads")

	fmt.Println(cyan + "🔧 sortdl is organizing your Downloads folder..." + reset)

	filesMoved := 0
	stats := make(map[string]int)

	entries, err := os.ReadDir(downloads)
	if err != nil {
		fmt.Println("Error opening Downloads folder")
		return
	}

	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == ".DS_Store" {
			continue
		}

		filePath := filepath.Join(downloads, entry.Name())
		ext := filepath.Ext(entry.Name())
		category := getCategory(ext)

		targetFolder := filepath.Join(downloads, category)
		targetPath := filepath.Join(targetFolder, entry.Name())

		// Handle duplicate names
		if _, err := os.Stat(targetPath); err == nil {
			base := strings.TrimSuffix(entry.Name(), ext)
			targetPath = filepath.Join(targetFolder, base+"_1"+ext)
		}

		stats[category]++

		if *dryRun {
			fmt.Printf("%sWould move:%s %s → %s/%s\n", yellow, reset, entry.Name(), category, entry.Name())
		} else {
			os.MkdirAll(targetFolder, 0755)
			os.Rename(filePath, targetPath)
			filesMoved++
		}
	}

	fmt.Println(green + "✅ Done!" + reset)
	if *dryRun {
		fmt.Println(gray + "This was a dry-run. Remove --dry-run to actually move files." + reset)
	} else {
		fmt.Printf("Moved %d files\n", filesMoved)
	}

	fmt.Println("\nSummary:")
	for cat := range categories {
		if stats[cat] > 0 {
			fmt.Printf("  %s%s:%s %d files\n", yellow, cat, reset, stats[cat])
		}
	}
	if stats["Others"] > 0 {
		fmt.Printf("  %sOthers:%s %d files\n", yellow, reset, stats["Others"])
	}
}
