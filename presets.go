package main

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed presets/*.jpg
var builtinPresetFS embed.FS

type PresetBackground struct {
	ID         string `json:"id"`
	PreviewURL string `json:"previewUrl"`
}

var builtinPresetIDs = []string{"technology", "anime", "city", "nature", "animal"}

func presetAssetPath(id string) (string, bool) {
	for _, allowed := range builtinPresetIDs {
		if id == allowed {
			return "presets/" + id + ".jpg", true
		}
	}
	return "", false
}

func listBuiltinPresets() ([]PresetBackground, error) {
	result := make([]PresetBackground, 0, len(builtinPresetIDs))
	for _, id := range builtinPresetIDs {
		assetPath, _ := presetAssetPath(id)
		data, err := builtinPresetFS.ReadFile(assetPath)
		if err != nil {
			return nil, fmt.Errorf("read preset %s: %w", id, err)
		}
		result = append(result, PresetBackground{ID: id, PreviewURL: makeDataURL("image/jpeg", data)})
	}
	return result, nil
}

func (s *CodexService) selectBuiltinPreset(id string) (ImageSelection, error) {
	assetPath, ok := presetAssetPath(id)
	if !ok {
		return ImageSelection{}, errors.New("unknown built-in background")
	}
	data, err := builtinPresetFS.ReadFile(assetPath)
	if err != nil {
		return ImageSelection{}, fmt.Errorf("read built-in background: %w", err)
	}
	dir, err := s.dataDir()
	if err != nil {
		return ImageSelection{}, err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return ImageSelection{}, err
	}
	destination := filepath.Join(dir, "builtin-"+id+".jpg")
	temp := destination + ".tmp"
	if err := os.WriteFile(temp, data, 0o600); err != nil {
		return ImageSelection{}, err
	}
	if err := os.Rename(temp, destination); err != nil {
		return ImageSelection{}, err
	}
	return ImageSelection{
		Path: destination, Name: id, Size: int64(len(data)),
		PreviewURL: makeDataURL("image/jpeg", data),
	}, nil
}
