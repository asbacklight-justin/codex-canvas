package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App exposes the small, typed API used by the desktop UI.
type App struct {
	ctx     context.Context
	service *CodexService
}

func NewApp() *App {
	return &App{service: NewCodexService()}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetStatus() AppStatus {
	return a.service.Status()
}

// SelectImage uses the native Finder picker and returns a local preview. The
// image never leaves this Mac.

func (a *App) SelectImage(language string) (ImageSelection, error) {
	title := "Choose a Codex background"
	filterName := "Images (PNG, JPEG, WebP, GIF)"
	if language == "zh" {
		title = "选择 Codex 背景图片"
		filterName = "图片（PNG、JPEG、WebP、GIF）"
	}
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
		Filters: []runtime.FileFilter{{
			DisplayName: filterName,
			Pattern:     "*.png;*.jpg;*.jpeg;*.webp;*.gif",
		}},
	})
	if err != nil {
		return ImageSelection{}, err
	}
	if path == "" {
		return ImageSelection{Cancelled: true}, nil
	}
	return loadImageSelection(path)
}

func (a *App) ListPresets() ([]PresetBackground, error) {
	return listBuiltinPresets()
}

func (a *App) SelectPreset(id string) (ImageSelection, error) {
	return a.service.selectBuiltinPreset(id)
}

func (a *App) ApplyTheme(theme ThemeConfig) (ActionResult, error) {
	return a.service.Apply(theme)
}

func (a *App) RestoreOfficial() (ActionResult, error) {
	return a.service.Restore()
}

func (a *App) RevealDataFolder() error {
	dir, err := a.service.dataDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	runtime.BrowserOpenURL(a.ctx, "file://"+filepath.ToSlash(dir))
	return nil
}

func loadImageSelection(path string) (ImageSelection, error) {
	info, err := os.Stat(path)
	if err != nil {
		return ImageSelection{}, fmt.Errorf("无法读取图片：%w", err)
	}
	if info.Size() > maxImageBytes {
		return ImageSelection{}, fmt.Errorf("图片不能超过 16 MB")
	}
	mime, ok := imageMIME(filepath.Ext(path))
	if !ok {
		return ImageSelection{}, fmt.Errorf("暂不支持 %s 格式", strings.TrimPrefix(filepath.Ext(path), "."))
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return ImageSelection{}, fmt.Errorf("无法读取图片：%w", err)
	}
	return ImageSelection{
		Path:       path,
		Name:       filepath.Base(path),
		Size:       info.Size(),
		PreviewURL: makeDataURL(mime, data),
	}, nil
}
