package utils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"strings"

	"github.com/mankings/mec-federator/internal/models/dto"
	"gopkg.in/yaml.v2"
)

func GetDescriptorData(tgzData []byte) (map[string]interface{}, error) {
	buffer := bytes.NewReader(tgzData)

	gzReader, err := gzip.NewReader(buffer)
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	var yamlFiles []string
	var yamlFileHeader *tar.Header

	// First pass: find all .yaml files
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if strings.HasSuffix(header.Name, ".yaml") {
			yamlFiles = append(yamlFiles, header.Name)
			yamlFileHeader = header
		}
	}

	if len(yamlFiles) == 0 {
		return nil, errors.New("no .yaml file found in the archive")
	}

	if len(yamlFiles) > 1 {
		return nil, errors.New("more than one .yaml file found in the archive")
	}

	// Reset readers to extract the yaml file
	buffer = bytes.NewReader(tgzData)
	gzReader, err = gzip.NewReader(buffer)
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()

	tarReader = tar.NewReader(gzReader)

	// Second pass: extract the yaml file
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if header.Name == yamlFileHeader.Name {
			yamlContent, err := io.ReadAll(tarReader)
			if err != nil {
				return nil, err
			}

			var result map[string]interface{}
			err = yaml.Unmarshal(yamlContent, &result)
			if err != nil {
				return nil, err
			}

			return result, nil
		}
	}

	return nil, errors.New("yaml file not found during extraction")
}

func ValidateDescriptorData(descriptorData map[string]interface{}) (dto.NewAppPkg, error) {
	var appPkg dto.NewAppPkg

	mec_appd := descriptorData["mec-appd"].(map[interface{}]interface{})

	appPkg.AppdId = mec_appd["id"].(string)
	appPkg.Name = mec_appd["name"].(string)
	appPkg.InfoName = mec_appd["info-name"].(string)
	appPkg.Provider = mec_appd["provider"].(string)
	appPkg.Version = mec_appd["version"].(string)
	appPkg.Description = mec_appd["description"].(string)
	appPkg.MecVersion = mec_appd["mec-version"].(string)

	return appPkg, nil
}
