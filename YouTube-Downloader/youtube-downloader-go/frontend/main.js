import * as App from "./wailsjs/go/backend/App";
import { OpenDirectoryDialog } from "./wailsjs/runtime";

async function addLink() {
  const url = document.getElementById("url").value.trim();
  if (!url) {
    log("Digite um link válido!");
    return;
  }

  try {
    const msg = await App.AddLink(url);
    log(msg);
    document.getElementById("links").value += url + "\n";
    document.getElementById("url").value = "";
  } catch (err) {
    log("Erro ao adicionar link: " + err);
  }
}

async function baixarMP3() {
  try {
    const destino = await OpenDirectoryDialog();
    if (!destino) return;
    log("Baixando MP3...");
    const result = await App.DownloadMP3(destino);
    log(result);
  } catch (err) {
    log("Erro ao baixar MP3: " + err);
  }
}

async function baixarVideo() {
  try {
    const destino = await OpenDirectoryDialog();
    if (!destino) return;
    log("Baixando vídeos...");
    const result = await App.DownloadVideo(destino);
    log(result);
  } catch (err) {
    log("Erro ao baixar vídeo: " + err);
  }
}

function log(text) {
  const div = document.getElementById("log");
  div.innerText += text + "\n";
  div.scrollTop = div.scrollHeight;
}

// expõe as funções ao HTML
window.addLink = addLink;
window.baixarMP3 = baixarMP3;
window.baixarVideo = baixarVideo;
