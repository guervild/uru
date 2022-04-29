package tampering

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/guervild/uru/pkg/common"

	"github.com/josephspurrier/goversioninfo"
)

type FileProperties struct {
	CompanyName         string
	InternalName        string
	FileDescription     string
	FileVersion         string
	LegalCopyright      string
	OriginalFilename    string
	ProductVersionPatch int
	ProductVersionMajor int
	ProductVersionMinor int
	ProductName         string
	ProductVersion      string
	FileVersionMajor    int
	FileVersionMinor    int
	FileVersionPatch    int
	FileVersionBuild    int
}

func BuildFromJson(filename, arch, dirpath string) (string, error) {

	if err := common.CheckIfFileExists(filename); err != nil {
		return "", err
	}

	configFilepath, _ := filepath.Abs(filename)
	configData, err := ioutil.ReadFile(configFilepath)

	if err != nil {
		return "", err
	}

	vi := &goversioninfo.VersionInfo{}

	if err := vi.ParseJSON(configData); err != nil {
		return "", fmt.Errorf("File Properties error could not parse the .json file: %v", err)
	}

	name := vi.StringFileInfo.InternalName

	vi.VarFileInfo.Translation.LangID = goversioninfo.LangID(1033)
	vi.VarFileInfo.Translation.CharsetID = goversioninfo.CharsetID(1200)

	vi.Build()
	vi.Walk()

	fileout := "resource_windows.syso"

	if dirpath != "" || len(dirpath) > 0 {
		fileout = filepath.Join(dirpath, filepath.Base(fileout))
	}

	if err := vi.WriteSyso(fileout, arch); err != nil {
		return "", fmt.Errorf("Error writing syso: %v", err)
	}

	return name, nil
}