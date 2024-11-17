package utils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/axzilla/stackpack/config"
)

func CreateFile(output *string) (*os.File, error) {
	if output == nil {
		return nil, fmt.Errorf("No output provided")
	}
	file, err := os.Create(*output)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return file, nil
}

func IsExtensionAllowed(ext string) bool {
	disallowedExts := ".png.jpg.jpeg.webp.pdf.zip"
	return !strings.Contains(disallowedExts, strings.ToLower(ext))
}

func IsValidGithubURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	if u.Host != "github.com" {
		return false
	}

	// Check for user/repo pattern
	pattern := regexp.MustCompile(`^/[^/]+/[^/]+$`)
	return pattern.MatchString(u.Path)
}

func makeRequest(remote, branch string) (*http.Response, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", remote+"/archive/refs/heads/"+branch+".zip", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", config.AppConfig.GitHubToken)
	resp, err := client.Do(req)
	return resp, nil
}

func FileSystem(remote string) (fs.FS, error) {
	if remote == "" {
		return os.DirFS("."), nil
	}

	resp, err := makeRequest(remote, "main")
	if err != nil || resp.StatusCode != 200 {
		resp, err = makeRequest(remote, "master")
		if err != nil {
			return nil, fmt.Errorf("download failed: %w", err)
		}
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read failed: %w", err)
	}

	return zip.NewReader(bytes.NewReader(data), int64(len(data)))
}

func WriteToBuffer(fileSys fs.FS, url *string, incPtn, excPtn []string, treeBuf, contentsBuf *bytes.Buffer) error {
	// Filetree header
	_, err := treeBuf.WriteString("=========================\n")
	if err != nil {
		return err
	}
	_, err = treeBuf.WriteString("Filetree\n")
	if err != nil {
		return err
	}
	_, err = treeBuf.WriteString("=========================\n")
	if err != nil {
		return err
	}

	// Iterate recursive
	err = fs.WalkDir(fileSys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Ignore some defaults
		if strings.Contains(path, ".git") || strings.Contains(path, ".DS_Store") {
			return nil
		}

		if !d.IsDir() {
			// Include files
			included := false
			for _, pattern := range incPtn {
				m, err := filepath.Match(pattern, d.Name())
				if err != nil {
					return err
				}
				if m {
					included = true
					break
				}
			}
			if !included {
				return nil
			}

			// Exclude files
			excluded := false
			for _, pattern := range excPtn {
				m, err := filepath.Match(pattern, d.Name())
				if err != nil {
					return err
				}
				if m {
					excluded = true
					break
				}
			}
			if excluded {
				return nil
			}
		}

		// Write depth indent for filetree into buffer
		depth := strings.Count(path, "/")
		for range depth {
			_, err := treeBuf.WriteString(" ")
			if err != nil {
				return err
			}
		}

		// Write file/folder names for filetree into buffer
		_, err = treeBuf.WriteString(d.Name())
		if err != nil {
			return err
		}
		if d.IsDir() {
			_, err = treeBuf.WriteString("/")
			if err != nil {
				return err
			}
		}
		_, err = treeBuf.WriteString("\n")
		if err != nil {
			return err
		}

		// Write filecontents into buffer
		ext := filepath.Ext(d.Name())
		if !d.IsDir() && IsExtensionAllowed(ext) {
			// use fs. instead os. because fs works with every filesystem not only with local one on HDD
			openedFile, err := fs.ReadFile(fileSys, path)
			if err != nil {
				return err
			}

			_, err = contentsBuf.WriteString("=========================\n")
			_, err = contentsBuf.WriteString("File: " + path + "\n")
			_, err = contentsBuf.WriteString("=========================\n")
			_, err = contentsBuf.Write(openedFile)
			_, err = contentsBuf.WriteString("\n")
			if err != nil {
				return nil
			}
		}
		return nil
	})
	return nil
}
