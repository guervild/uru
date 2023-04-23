package evasion

import (
	"fmt"
	"github.com/guervild/uru/pkg/evasion/c"
	"github.com/guervild/uru/pkg/evasion/go"
	"github.com/guervild/uru/pkg/models"
)

func GetEvasion(evasionType string, langType string) (models.ObjectModel, error) {

	switch langType {
	case "go":
		if evasionType == "sleep" {
			return _go.NewSleepEvasion(), nil
		}

		if evasionType == "hideconsole" {
			return _go.NewHideConsoleEvasion(), nil
		}

		if evasionType == "isdomainjoined" {
			return _go.NewIsDomainJoinedEvasion(), nil
		}

		if evasionType == "selfdelete" {
			return _go.NewSelfDeleteEvasion(), nil
		}

		if evasionType == "ntsleep" {
			return _go.NewNtSleepEvasion(), nil
		}

		if evasionType == "english_words" {
			return _go.NewEnglishWordsEvasion(), nil
		}

		if evasionType == "patch" {
			return _go.NewPatchEvasion(), nil
		}

		if evasionType == "patchetw" {
			return _go.NewPatchEtwEvasion(), nil
		}

		if evasionType == "patchamsi" {
			return _go.NewPatchAmsiEvasion(), nil
		}

		if evasionType == "createmutex" {
			return _go.NewCreateMutexEvasion(), nil
		}

		if evasionType == "refreshdll" {
			return _go.NewRefreshDllEvasion(), nil
		}

		return nil, fmt.Errorf("Wrong evasion type passed: evasion %s is unknown", evasionType)

	case "c":
		switch evasionType {
		case "sleep":
			return c.NewCSleepEvasion(), nil
		case "dllforward":
			return c.NewCDllForwardEvasion(), nil
		default:
			break
		}
		return nil, fmt.Errorf("Wrong evasion type passed: evasion %s is unknown", evasionType)
	case "rust":
		return nil, fmt.Errorf("Wrong evasion type passed: evasion %s is unknown", evasionType)
	}
	return nil, fmt.Errorf("Wrong langtype: %s is unknown", langType)
}
