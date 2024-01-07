package registry_client

import (
	"archive/tar"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

const (
	RegistryBaseURL = "https://registry-1.docker.io/v2/"
	AuthServiceURL  = "https://auth.docker.io/token"
)

type Manifest struct {
	SchemaVersion int     `json:"schemaVersion"`
	MediaType     string  `json:"mediaType,omitempty"`
	Config        Layer   `json:"config"`
	Layers        []Layer `json:"layers"`
}

type Layer struct {
	MediaType string `json:"mediaType,omitempty"`
	Size      int64  `json:"size"`
	Digest    string `json:"digest"`
}

type ContainerConfig struct {
	Env         []string `json:"Env"`
	Cmd         []string `json:"Cmd"`
	ArgsEscaped bool     `json:"argsEscaped"`
	// onBuild
}

type ConfigLayer struct {
	Architecture string          `json:"architecture"`
	Config       ContainerConfig `json:"config"`
	Created      string          `json:"created"`
	Os           string          `json:"os"`
	// history      []string `json:"history"`
	// rootfs rootfs `json:"rootfs"`
}

// GetTags lists tags for the specified repository.
func (c *DockerRegistryClient) GetTags(repository string) ([]string, error) {
	// ... (same as before)
	// return error
	return nil, fmt.Errorf("error querying docker registry")
}

type DockerRegistryClient struct {
	BaseURL    string
	Username   string
	Password   string
	HttpClient *http.Client
}

type AuthTokenResponse struct {
	Token       string `json:"token"`
	AccessToken string `json:"access_token"`
}

// ... (Manifest and Layer structs stay the same)

func NewDockerRegistryClient(baseURL, username, password string) *DockerRegistryClient {
	return &DockerRegistryClient{
		BaseURL:    baseURL,
		Username:   username,
		Password:   password,
		HttpClient: &http.Client{},
	}
}

func (c *DockerRegistryClient) GetAuthToken(repository string) (string, error) {
	clientId := "docker-client" // Normally this would be your client identifier

	// Construct the URL for token request; scope can be more precise, this is a wide example
	scope := fmt.Sprintf("repository:%s:pull", repository)
	url := fmt.Sprintf("%s?service=registry.docker.io&scope=%s&client_id=%s", AuthServiceURL, scope, clientId)

	// If you need to authenticate, set the Authorization header
	var req *http.Request
	var err error
	if c.Username != "" && c.Password != "" {
		auth := c.Username + ":" + c.Password
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return "", err
		}
		req.Header.Add("Authorization", "Basic "+encodedAuth)
	} else {
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return "", err
		}
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error getting auth token: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResponse AuthTokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", err
	}

	// Docker may return the access token under different names
	if tokenResponse.AccessToken != "" {
		return tokenResponse.AccessToken, nil
	}
	return tokenResponse.Token, nil
}

func (c *DockerRegistryClient) GetManifest(repository, tag string) (*Manifest, error) {
	token, err := c.GetAuthToken(repository)
	if err != nil {
		return nil, err
	}

	url := c.BaseURL + repository + "/manifests/" + tag
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error querying docker registry: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// print raw manifest
	// fmt.Println("Raw manifest:", string(body))

	var manifest Manifest
	if err := json.Unmarshal(body, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}
func (c *DockerRegistryClient) DownloadLayer(repository, digest string) (string, error) {
	// check if layer exists
	layerDir := "layers"
	layerFile := filepath.Join(layerDir, digest)
	if _, err := os.Stat(layerFile); err == nil {
		// fmt.Println("Layer already downloaded:", layerFile)
		return layerFile, nil
	}

	// Get the auth token to access the repository
	token, err := c.GetAuthToken(repository)
	if err != nil {
		return "", err
	}

	// Construct the URL for the layer
	layerURL := fmt.Sprintf("%s%s/blobs/%s", c.BaseURL, repository, digest)

	// Create an HTTP request for the layer
	req, err := http.NewRequest("GET", layerURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Execute the request
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download layer %s: %s", digest, resp.Status)
	}

	// Create the directory where the layer will be stored

	if err := os.MkdirAll(layerDir, 0755); err != nil {
		return "", err
	}

	// Create a file to write the layer to
	file, err := os.Create(layerFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Write the layer content to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println("Layer downloaded:", layerFile)
	return layerFile, nil
}

// extractLayer extracts a GZIP'd tar archive (a Docker layer) into the specified destination directory.
func extractLayer(filePath, destination string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	// Iterate through the files in the tar archive
	for {
		header, err := tarReader.Next()
		switch {
		case err == io.EOF:
			return nil // End of archive
		case err != nil:
			return err // Other error
		case header == nil:
			continue // Next file
		}

		// Construct the path where the file should be created
		target := filepath.Join(destination, header.Name)

		// if target exists delete it first
		if tar.TypeReg == header.Typeflag || tar.TypeRegA == header.Typeflag {
			if _, err := os.Stat(target); err == nil {
				if err := os.Remove(target); err != nil {
					return err
				}
			}
		}

		// Handle the different types of files in the tar archive
		switch header.Typeflag {
		case tar.TypeDir:
			// Create a directory
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}

		case tar.TypeReg, tar.TypeRegA:
			// Handle regular files
			if err := extractFile(tarReader, header, target); err != nil {
				return err
			}

		case tar.TypeSymlink:
			// Handle symlinks
			// if exists remove old symlink
			if _, err := os.Stat(target); err == nil {
				if err := os.Remove(target); err != nil {
					return err
				}
			}

			link := header.Linkname
			if err := os.Symlink(link, target); err != nil {
				return err
			}
		}
	}
}

// Extract a regular file from a tar archive.
func extractFile(tarReader *tar.Reader, header *tar.Header, filepath string) error {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(header.Mode))
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy file data from tar archive to the file
	if _, err := io.CopyN(file, tarReader, header.Size); err != nil {
		return err
	}

	return nil
}

func ExtractAndAssembleImage(client *DockerRegistryClient, repository, tag, rootFsPath string) (*ConfigLayer, error) {
	// Ensure the root filesystem directory exists
	if _, err := os.Stat(rootFsPath); err == nil {
		if err := os.RemoveAll(rootFsPath); err != nil {
			return nil, err
		}
	}

	if err := os.MkdirAll(rootFsPath, 0755); err != nil {
		return nil, err
	}

	// Get the manifest for the specific image tag
	manifest, err := client.GetManifest(repository, tag)
	if err != nil {
		return nil, err
	}

	// print manifest
	// fmt.Println("Manifest:", manifest)

	fmt.Printf("Found %d layers\n", len(manifest.Layers))
	accumulatedSize := int64(0)
	for i, layer := range manifest.Layers {
		fmt.Printf("Layer (%d) %s (%d MB)\n", i, layer.Digest, layer.Size/1024/1024)
		accumulatedSize += layer.Size
	}
	fmt.Printf("Total size: %d MB\n", accumulatedSize/1024/1024)

	// Download all layers defined in the manifest
	for i, layer := range manifest.Layers {
		// Download layer

		fmt.Printf("Downloading layer (%d) %s...\n", i, layer.Digest)
		layerFile, err := client.DownloadLayer(repository, layer.Digest)
		if err != nil {
			return nil, fmt.Errorf("error downloading layer %s: %w", layer.Digest, err)
		}

		fmt.Printf("Extracting layer (%d) %s into %s\n", i, layer.Digest, rootFsPath)
		// Extract layer
		if err := extractLayer(layerFile, rootFsPath); err != nil {
			return nil, fmt.Errorf("error extracting layer %s: %w", layer.Digest, err)
		}

		// Clean up the downloaded layer file
		// if err := os.Remove(layerFile); err != nil {
		// 	return fmt.Errorf("error cleaning up layer file %s: %w", layerFile, err)
		// }
	}

	// Extract the image configuration
	println("Extracting image configuration...", manifest.Config.Digest)
	configFile, err := client.DownloadLayer(repository, manifest.Config.Digest)
	if err != nil {
		return nil, fmt.Errorf("error downloading image configuration: %w", err)
	}

	config := ConfigLayer{}
	fileContent, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("error reading config file %s: %w", configFile, err)
	}
	if err := json.Unmarshal(fileContent, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file %s: %w", configFile, err)
	}

	println("Successfully extracted image and configuration")
	return &config, nil
}
