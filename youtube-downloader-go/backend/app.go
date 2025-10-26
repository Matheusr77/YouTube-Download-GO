package backend

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx    context.Context
	Links  []string
	BinDir string
}

func NewApp() *App {
	exe, _ := os.Executable()
	dir := filepath.Dir(exe)

	// Ajusta automaticamente se estiver dentro de build/bin/
	binDir := filepath.Join(dir)
	if strings.Contains(dir, "build") {
		binDir = dir // mantém apenas a pasta do executável
	}

	return &App{
		Links:  []string{},
		BinDir: binDir,
	}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	fmt.Println("✅ Aplicativo iniciado com contexto Wails.")
}

// ======================
// 🔹 Escolher diretório
// ======================
func (a *App) EscolherDiretorio() (string, error) {
	dir, err := wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Escolha onde salvar os arquivos",
	})
	if err != nil {
		return "", err
	}
	if dir == "" {
		return "", errors.New("nenhuma pasta selecionada")
	}
	return dir, nil
}

// ======================
// 🔹 Adicionar Link
// ======================
func (a *App) AddLink(url string) string {
	url = strings.TrimSpace(url)
	if url == "" {
		return "⚠️ Link vazio!"
	}
	a.Links = append(a.Links, url)
	return fmt.Sprintf("Adicionado: %s", url)
}

// ======================
// 🔹 Limpar Links
// ======================
func (a *App) ClearLinks() {
	a.Links = []string{}
}

// ======================
// 🔹 Baixar MP3
// ======================
func (a *App) DownloadMP3(destino string) (string, error) {
	if len(a.Links) == 0 {
		return "", errors.New("nenhum link adicionado")
	}
	if destino == "" {
		return "", errors.New("nenhuma pasta selecionada")
	}

	yt := filepath.Join(a.BinDir, binName("yt-dlp"))
	ff := filepath.Join(a.BinDir, binName("ffmpeg"))

	for _, link := range a.Links {
		cmd := exec.Command(yt, "-f", "bestaudio", "-o", filepath.Join(destino, "%(title)s.%(ext)s"), link)
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("erro no yt-dlp: %v", err)
		}

		audio, _ := getLatestFile(destino)
		if audio == "" {
			continue
		}
		mp3 := strings.TrimSuffix(audio, filepath.Ext(audio)) + ".mp3"
		cmd2 := exec.Command(ff, "-i", audio, "-vn", "-ar", "44100", "-ac", "2", "-b:a", "320k", mp3)
		cmd2.Run()
		_ = os.Remove(audio)
	}

	a.ClearLinks()
	return "🎶 MP3s baixados com sucesso!", nil
}

// ======================
// 🔹 Baixar Vídeo
// ======================
func (a *App) DownloadVideo(destino string) (string, error) {
	if len(a.Links) == 0 {
		return "", errors.New("nenhum link adicionado")
	}
	if destino == "" {
		return "", errors.New("nenhuma pasta selecionada")
	}

	yt := filepath.Join(a.BinDir, binName("yt-dlp"))
	for _, link := range a.Links {
		cmd := exec.Command(yt, "-f", "bestvideo+bestaudio", "-o", filepath.Join(destino, "%(title)s.%(ext)s"), link)
		if err := cmd.Run(); err != nil {
			return "", err
		}
	}
	a.ClearLinks()
	return "📹 Vídeos baixados com sucesso!", nil
}

// ======================
// 🔹 Utilidades
// ======================
func getLatestFile(folder string) (string, error) {
	entries, err := os.ReadDir(folder)
	if err != nil {
		return "", err
	}
	var newest string
	var newestTime int64
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Unix() > newestTime {
			newestTime = info.ModTime().Unix()
			newest = filepath.Join(folder, e.Name())
		}
	}
	return newest, nil
}

func binName(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}
