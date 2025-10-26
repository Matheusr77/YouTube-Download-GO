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
)

type App struct {
	Links  []string
	BinDir string
}

func NewApp() *App {
	exe, _ := os.Executable()
	dir := filepath.Dir(exe)
	return &App{
		Links:  []string{},
		BinDir: filepath.Join(dir, "bin"),
	}
}

func (a *App) AddLink(url string) string {
	url = strings.TrimSpace(url)
	if url == "" {
		return "Link vazio!"
	}
	a.Links = append(a.Links, url)
	return fmt.Sprintf("Adicionado: %s", url)
}

func (a *App) ClearLinks() {
	a.Links = []string{}
}

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
		os.Remove(audio)
	}
	a.ClearLinks()
	return "MP3s baixados com sucesso!", nil
}

func (a *App) DownloadVideo(destino string) (string, error) {
	if len(a.Links) == 0 {
		return "", errors.New("nenhum link adicionado")
	}
	yt := filepath.Join(a.BinDir, binName("yt-dlp"))
	for _, link := range a.Links {
		cmd := exec.Command(yt, "-f", "bestvideo+bestaudio", "-o", filepath.Join(destino, "%(title)s.%(ext)s"), link)
		if err := cmd.Run(); err != nil {
			return "", err
		}
	}
	a.ClearLinks()
	return "Vídeos baixados com sucesso!", nil
}

func getLatestFile(folder string) (string, error) {
	files, err := os.ReadDir(folder)
	if err != nil {
		return "", err
	}
	var newest string
	var newestTime int64
	for _, f := range files {
		if info, err := f.Info(); err == nil {
			if info.ModTime().Unix() > newestTime {
				newestTime = info.ModTime().Unix()
				newest = filepath.Join(folder, f.Name())
			}
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

func (a *App) Startup(ctx context.Context) {
	// Chamado quando o app inicia
}
