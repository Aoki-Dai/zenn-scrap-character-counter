/**
 * popup.js
 * Wasmの読み込み、content.jsへのメッセージ送信、結果表示を担う。
 */

async function loadWasm() {
  const go = new Go(); // wasm_exec.js で定義されるグローバルクラス

  // MV3では fetch() の Response を直接 instantiateStreaming に渡せないため、
  // ArrayBuffer に変換してから WebAssembly.instantiate() を使用する。
  const response = await fetch(chrome.runtime.getURL("main.wasm"));
  const buffer = await response.arrayBuffer();
  const { instance } = await WebAssembly.instantiate(buffer, go.importObject);

  go.run(instance); // Goのmain()を非同期実行（<-done でブロックし続ける）
}

function setStatus(msg) {
  document.getElementById("status").textContent = msg;
}

function setError(msg) {
  document.getElementById("error").textContent = msg;
  setStatus("");
}

function displaySelectionResults(selectedText) {
  const section = document.getElementById("selection-section");
  if (selectedText === "") {
    section.style.display = "none";
  } else {
    const stats = analyzeText(selectedText); // eslint-disable-line no-undef
    document.getElementById("sel-total").textContent = stats.total.toLocaleString();
    document.getElementById("sel-japanese").textContent = stats.japanese.toLocaleString();
    document.getElementById("sel-noSpace").textContent = stats.noSpace.toLocaleString();
    section.style.display = "block";
  }
}

function displayResults(stats, postCount) {
  document.getElementById("total").textContent = stats.total.toLocaleString();
  document.getElementById("japanese").textContent =
    stats.japanese.toLocaleString();
  document.getElementById("noSpace").textContent =
    stats.noSpace.toLocaleString();
  document.getElementById("postCount").textContent = postCount.toLocaleString();
  setStatus("");
}

async function main() {
  try {
    setStatus("Wasm を読み込み中...");
    await loadWasm();

    setStatus("ページのテキストを取得中...");

    // アクティブなタブにメッセージを送信してテキストを取得
    const [tab] = await chrome.tabs.query({
      active: true,
      currentWindow: true,
    });

    if (
      !tab ||
      !tab.url ||
      !tab.url.match(/^https:\/\/zenn\.dev\/[^/]+\/scraps\//)
    ) {
      setError("Zenn のスクラップページを開いてください。");
      return;
    }

    let response;
    try {
      response = await chrome.tabs.sendMessage(tab.id, { action: "getText" });
    } catch (e) {
      setError(
        "ページと通信できませんでした。リロードして再試行してください。",
      );
      return;
    }

    if (!response || !response.text) {
      setError("テキストを取得できませんでした。");
      return;
    }

    setStatus("文字数を集計中...");

    // Go Wasm の analyzeText をグローバルから呼び出す
    const stats = analyzeText(response.text); // eslint-disable-line no-undef
    displayResults(stats, response.postCount || 0);
    displaySelectionResults(response.selectedText || "");
  } catch (err) {
    setError(`エラー: ${err.message}`);
  }
}

document.addEventListener("DOMContentLoaded", main);
