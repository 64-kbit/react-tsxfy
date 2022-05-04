package main

import (
	"fmt"
	"github.com/thatisuday/commando"
	"io/fs"
	"io/ioutil"
	"os"
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
		AddFlag("path,p", "Specify Path of File/Folder", commando.String, nil).
		SetAction(PrintJsFiles())
	// parse command-line arguments from the STDIN
	commando.Parse(nil)
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
		fmt.Println("Listing JavaScript Files contents...")
		dir, _ := flags["path"].GetString()
		files := GetJavaScriptFiles(dir)
		for _, file := range files {
			fmt.Print(file.Path)
			fmt.Print("Size: ", file.File.Size())
			fmt.Println("\n")
		}
	}

}
