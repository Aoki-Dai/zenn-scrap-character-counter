/**
 * content.js
 * Zennスクラップページから投稿テキストを抽出し、popup.jsからのメッセージに応答する。
 */

/**
 * Zennスクラップの投稿テキストを収集する。
 * DOM構造の変更に備えて複数のセレクタをフォールバック付きで試行する。
 * @returns {{ text: string, postCount: number }}
 */
function collectScrapText() {
  // Zennのスクラップ投稿コンテンツの候補セレクタ（優先順）
  const selectors = [
    ".ScrapComment_body__text__vTfV1", // ハッシュ付きクラス（バージョン依存）
    "[class*='ScrapComment_body']",
    "[class*='scrap-comment'] p",
    "article p",
  ];

  let nodes = [];
  for (const selector of selectors) {
    nodes = Array.from(document.querySelectorAll(selector));
    if (nodes.length > 0) break;
  }

  // テキストが取れない場合は本文全体をフォールバックとして使用
  if (nodes.length === 0) {
    const main = document.querySelector("main") || document.body;
    nodes = [main];
  }

  const texts = nodes.map((el) => el.innerText || el.textContent || "");
  const combined = texts.join("\n");

  // 投稿数：各投稿に必ず1つある <time> タグをカウント（クラス名に依存しない）
  // スクラップ外のヘッダー等の time タグを除くため main 内に絞る
  const mainEl = document.querySelector("main") || document.body;
  const postCount = mainEl.querySelectorAll("time").length;

  return { text: combined, postCount };
}

chrome.runtime.onMessage.addListener((message, _sender, sendResponse) => {
  if (message.action === "getText") {
    const result = collectScrapText();
    result.selectedText = window.getSelection().toString();
    sendResponse(result);
  }
  // sendResponse を非同期で呼ぶ場合は true を返す必要があるが、
  // ここでは同期なので省略可。念のため true を返す。
  return true;
});
