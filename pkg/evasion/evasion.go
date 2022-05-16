package evasion

import (
	"fmt"
	"github.com/guervild/uru/pkg/models"
)

func GetEvasion(evasionType string) (models.ObjectModel, error) {
	if evasionType == "sleep" {
		return NewSleepEvasion(), nil
	}

	if evasionType == "hideconsole" {
		return NewHideConsoleEvasion(), nil
	}

	if evasionType == "isdomainjoined" {
		return NewIsDomainJoinedEvasion(), nil
	}

	if evasionType == "selfdelete" {
		return NewSelfDeleteEvasion(), nil
	}

	if evasionType == "ntsleep" {
		return NewNtSleepEvasion(), nil
	}

	if evasionType == "english-words" {
		return NewEnglishWordsEvasion(), nil
	}

	if evasionType == "patch" {
		return NewPatchEvasion(), nil
	}

	if evasionType == "patchetw" {
		return NewPatchEtwEvasion(), nil
	}

	if evasionType == "patchamsi" {
		return NewPatchAmsiEvasion(), nil
	}

	if evasionType == "createmutex" {
		return NewCreateMutexEvasion(), nil
	}

	return nil, fmt.Errorf("Wrong evasion type passed: evasion %s is unknown", evasionType)
}
