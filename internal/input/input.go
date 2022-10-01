package input

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/chrisgavin/ipman/internal/types"
	"github.com/chrisgavin/ipman/internal/validation"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const extension = "yaml"
const self = "_self." + extension
const currentVersion = 1

var ErrUnsupportedVersion = errors.New("Unsupported input version.")

func readFile[T interface{}](path string, destination T) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	if err := decoder.Decode(destination); err != nil {
		return validation.NewValidationError(path, "File could not be parsed.", err)
	}

	return nil
}

func ReadInput(path string) (*types.Input, error) {
	rootPath := filepath.Join(path, self)
	input := types.Input{}
	if err := readFile(rootPath, &input); err != nil {
		return nil, err
	}
	input.Path = rootPath
	if input.Version != currentVersion {
		return nil, ErrUnsupportedVersion
	}

	networkPaths, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, networkPathInfo := range networkPaths {
		fullNetworkPath := filepath.Join(path, networkPathInfo.Name())
		if !networkPathInfo.IsDir() {
			if networkPathInfo.Name() == self {
				continue
			}
			return nil, errors.Errorf("Unexpected file at %s.", fullNetworkPath) // TODO
		}
		networkPath := filepath.Join(fullNetworkPath, self)
		network := types.Network{}
		if err := readFile(networkPath, &network); err != nil {
			return nil, err
		}
		network.Path = networkPath
		network.Name = networkPathInfo.Name()
		sitePaths, err := ioutil.ReadDir(fullNetworkPath)
		if err != nil {
			return nil, err
		}
		for _, sitePathInfo := range sitePaths {
			fullSitePath := filepath.Join(fullNetworkPath, sitePathInfo.Name())
			if !sitePathInfo.IsDir() {
				if sitePathInfo.Name() == self {
					continue
				}
				return nil, errors.Errorf("Unexpected file at %s.", fullSitePath)
			}
			sitePath := filepath.Join(fullSitePath, self)
			site := types.Site{}
			if err := readFile(sitePath, &site); err != nil {
				return nil, err
			}
			site.Path = sitePath
			site.Name = sitePathInfo.Name()
			poolPaths, err := ioutil.ReadDir(fullSitePath)
			if err != nil {
				return nil, err
			}
			for _, poolPathInfo := range poolPaths {
				fullPoolPath := filepath.Join(fullSitePath, poolPathInfo.Name())
				if !poolPathInfo.IsDir() {
					if poolPathInfo.Name() == self {
						continue
					}
					return nil, errors.Errorf("Unexpected file at %s.", fullPoolPath)
				}
				poolPath := filepath.Join(fullPoolPath, self)
				pool := types.Pool{}
				if err := readFile(poolPath, &pool); err != nil {
					return nil, err
				}
				pool.Path = poolPath
				pool.Name = poolPathInfo.Name()
				hostPaths, err := ioutil.ReadDir(fullPoolPath)
				if err != nil {
					return nil, err
				}
				for _, hostPathInfo := range hostPaths {
					fullHostPath := filepath.Join(fullPoolPath, hostPathInfo.Name())
					if hostPathInfo.IsDir() {
						return nil, errors.Errorf("Unexpected directory at %s.", fullHostPath)
					}
					if hostPathInfo.Name() == self {
						continue
					}
					host := types.Host{}
					if err := readFile(fullHostPath, &host); err != nil {
						return nil, err
					}
					host.Path = fullHostPath
					host.Name = strings.Split(hostPathInfo.Name(), ".")[0]
					pool.Hosts = append(pool.Hosts, host)
				}
				site.Pools = append(site.Pools, pool)
			}
			network.Sites = append(network.Sites, site)
		}
		input.Networks = append(input.Networks, network)
	}

	return &input, nil
}
