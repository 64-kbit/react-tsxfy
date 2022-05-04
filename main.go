package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/thatisuday/commando"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	commando.SetExecutableName("react-tsxfy").SetVersion("v1.0.0").SetDescription("tsxAjs is a tool to generate ocnvert ja")

	commando.
		Register(nil).
		AddArgument("list", "Show list of all files and folders", ""). // required
		AddFlag("verbose,V", "display log information ", commando.Bool, nil).
		AddFlag("help,h", "display help information", commando.Bool, nil).
		AddFlag("path,p", "Specify Path of File/Folder", commando.String, nil).
		AddFlag("files,f", "display files", commando.Bool, false).
		AddFlag("size,s", "Display File Size", commando.Bool, false).
		SetAction(ListDirectoryContents())

	commando.Register("show").
		AddArgument("file", "Show file", ""). // required
		AddFlag("verbose,V", "display log information ", commando.Bool, nil).
		AddFlag("help,h", "display help information", commando.Bool, nil).
		AddFlag("path,p", "Specify Path of File/Folder", commando.String, "./").
		SetAction(PrintJsFiles())

	commando.Register("rename").SetShortDescription("Renames JavaScript file to TypeScript file").
		AddFlag("verbose,V", "display log information ", commando.Bool, nil).
		AddFlag("help,h", "display help information", commando.Bool, nil).
		AddFlag("path,p", "Specify Path of File/Folder", commando.String, "./").
		SetAction(RenameJsToTs())
	// parse command-line arguments from the STDIN
	commando.Parse(nil)
}

func RenameJsToTs() func(map[string]commando.ArgValue, map[string]commando.FlagValue) {
	return func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
		path, err := flags["path"].GetString()
		if err != nil {
			panic(err)
		}
		verbose, _ := flags["verbose"].GetBool()
		if verbose {
			c := color.New(color.FgBlue).Add(color.Bold)
			_, _ = c.Println("path:", path)

		}
		files := GetJavaScriptFiles(path)
		if len(files) == 0 {
			c := color.New(color.FgRed).Add(color.Bold)
			_, _ = c.Println("No Files are Found on path:", path)
			os.Exit(1)
		}
		for _, file := range files {
			if verbose {
				c := color.New(color.FgBlue).Add(color.Bold)
				_, _ = c.Println("file:", file)
			}
			newFile := strings.Replace(file.Path, ".jsx", ".tsx", 1)
			newFile = strings.Replace(newFile, ".js", ".tsx", 1)
			newFile = strings.Replace(newFile, ".react.", ".", 1)
			if verbose {
				c := color.New(color.FgGreen).Add(color.Bold).Add(color.Underline)
				_, _ = c.Println("newFile:", newFile)
			}
			err := os.Rename(file.Path, newFile)
			if err != nil {
				c := color.New(color.FgRed).Add(color.Bold).Add(color.Underline)
				_, _ = c.Println("Error:", err)
			}

		}
	}
}

func ScanFilesInDir(dir string) []fs.FileInfo {
	var files []fs.FileInfo
	var err error
	if files, err = ioutil.ReadDir(dir); err != nil {
		panic(err)
	}
	return files
}

func PrintDirContents(dir string, showFiles bool, showSize bool) {
	var contents = ScanFilesInDir(dir)
	for _, content := range contents {
		if showFiles {
			fmt.Print("-", dir+"/"+content.Name())
		}
		if showSize {
			fmt.Println("(", content.Size(), ")")
		}

		if content.IsDir() {
			fmt.Print("/ Directory ", content.Name())
			PrintDirContents(dir+"/"+content.Name(), showFiles, showSize)
			fmt.Println("\n----::----")
		}
	}
}

type customFileInfo struct {
	Path string
	File fs.FileInfo
}

func GetJavaScriptFiles(dir string) []customFileInfo {
	var files []fs.FileInfo
	var err error
	if files, err = ioutil.ReadDir(dir); err != nil {
		panic(err)
		return nil
	}
	var jsFiles []customFileInfo
	for _, file := range files {
		if file.IsDir() {
			jsFiles = append(jsFiles, GetJavaScriptFiles(dir+"/"+file.Name())...)
		}
		fl := customFileInfo{Path: dir + "/" + file.Name(), File: file}
		if file.Name()[len(file.Name())-3:] == ".js" {
			jsFiles = append(jsFiles, fl)
		}
	}
	return jsFiles
}

func ListDirectoryContents() func(map[string]commando.ArgValue, map[string]commando.FlagValue) {
	return func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
		fmt.Println("Listing directory contents...")
		dir, _ := flags["path"].GetString()
		fileInfo, err := os.Stat(dir)
		if err != nil {
			panic("File/Folder Name doesnot exitss ")
			return // exit
		}
		fileOly, _ := flags["files"].GetBool()
		sizeOly, _ := flags["size"].GetBool()
		if fileInfo.IsDir() {
			PrintDirContents(dir, fileOly, sizeOly)
		} else {
			fmt.Print(fileInfo.Name())
			fmt.Print("Size: ", fileInfo.Size())
			fmt.Println(" ")
		}
	}
}

func PrintJsFiles() func(map[string]commando.ArgValue, map[string]commando.FlagValue) {
	return func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {

		c := color.New(color.FgBlue).Add(color.Underline)
		_, err := c.Println("Listing JavaScript Files contents...")
		if err != nil {
			return
		}
		dir, _ := flags["path"].GetString()
		files := GetJavaScriptFiles(dir)
		for _, file := range files {
			fmt.Printf(file.Path)
			fmt.Printf("Size: %d \n", file.File.Size())
		}
	}

}
