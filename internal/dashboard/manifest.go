package dashboard

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Manifest defines the structure of a dashboard plugin.
type Manifest struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Entry       string            `json:"entry"` // e.g., "dist/index.js"
	Config      map[string]string `json:"config,omitempty"`
	Permissions []string          `json:"permissions,omitempty"`
}

// Registry manages the loading and listing of plugins.
type Registry struct {
	pluginDir string
	plugins   map[string]Manifest
}

func NewRegistry(pluginDir string) *Registry {
	return &Registry{
		pluginDir: pluginDir,
		plugins:   make(map[string]Manifest),
	}
}

// LoadPlugins reads all manifests from the plugin directory.
func (r *Registry) LoadPlugins() error {
	entries, err := os.ReadDir(r.pluginDir)
	if err != nil {
		return fmt.Errorf("failed to read plugin directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			manifestPath := filepath.Join(r.pluginDir, entry.Name(), "plugin.json")
			if _, err := os.Stat(manifestPath); err == nil {
				manifest, err := r.loadManifest(manifestPath)
				if err != nil {
					continue // Log error and skip
				}
				r.plugins[manifest.ID] = manifest
			}
		}
	}
	return nil
}

func (r *Registry) loadManifest(path string) (Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, err
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return Manifest{}, err
	}
	return m, nil
}

func (r *Registry) ListPlugins() []Manifest {
	list := make([]Manifest, 0, len(r.plugins))
	for _, m := range r.plugins {
		list = append(list, m)
	}
	return list
}
