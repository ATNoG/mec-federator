package utils

import "github.com/mankings/mec-federator/internal/models"

func AddMobileCodes(target *models.MobileNetworkIds, toAdd *models.MobileNetworkIds) {
	if target == nil || toAdd == nil {
		return
	}
	if target.Mcc != toAdd.Mcc {
		return // skip mismatched MCC
	}
	mncMap := make(map[string]bool)
	for _, mnc := range target.Mncs {
		mncMap[mnc] = true
	}
	for _, mnc := range toAdd.Mncs {
		if !mncMap[mnc] {
			target.Mncs = append(target.Mncs, mnc)
		}
	}
}

func RemoveMobileCodes(target *models.MobileNetworkIds, toRemove *models.MobileNetworkIds) {
	if target == nil || toRemove == nil {
		return
	}
	if target.Mcc != toRemove.Mcc {
		return // skip mismatched MCC
	}
	removeMap := make(map[string]bool)
	for _, mnc := range toRemove.Mncs {
		removeMap[mnc] = true
	}
	var filtered []string
	for _, mnc := range target.Mncs {
		if !removeMap[mnc] {
			filtered = append(filtered, mnc)
		}
	}
	target.Mncs = filtered
}

func AddFixedCodes(existing *[]string, toAdd *[]string) *[]string {
	if existing == nil {
		existing = &[]string{}
	}
	existingMap := make(map[string]bool)
	for _, code := range *existing {
		existingMap[code] = true
	}
	for _, code := range *toAdd {
		if !existingMap[code] {
			*existing = append(*existing, code)
		}
	}
	return existing
}

func RemoveFixedCodes(existing *[]string, toRemove *[]string) *[]string {
	if existing == nil {
		return &[]string{}
	}
	removeMap := make(map[string]bool)
	for _, code := range *toRemove {
		removeMap[code] = true
	}
	var filtered []string
	for _, code := range *existing {
		if !removeMap[code] {
			filtered = append(filtered, code)
		}
	}
	return &filtered
}
