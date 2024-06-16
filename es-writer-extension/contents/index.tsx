import type { PlasmoCSConfig } from "plasmo"

export const config: PlasmoCSConfig = {
  matches: ["<all_urls>"]
}

function indexContents() {

  console.log("indexContents called")  // 関数呼び出しの確認用ログ

  fetch("http://localhost:8080/getAnswers", {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({ html: document.documentElement.outerHTML })
  })
    .then(res => {
      if (!res.ok) {
        console.error("Network response was not ok", res.statusText)
      }
    })
    .then(answers => {
      console.log("Received answers:", answers) // 受け取ったデータをコンソールに出力
      replaceTextareaText(answers)
    });

    return <></>
}

function replaceTextareaText(answers: any) {
  const allTextareas = document.getElementsByTagName("textarea")
  Array.from(allTextareas).forEach((textarea, index) => {
    if (answers[index]) {
      textarea.value = answers[index].answer
    }
  })
}

export default indexContents

document.addEventListener("DOMContentLoaded", () => {
  const fetchAnswersButton = document.getElementById("fetchAnswersButton")
  if (fetchAnswersButton) {
    fetchAnswersButton.addEventListener("click", indexContents)
    console.log("Event listener added to fetchAnswersButton")
  } else {
    console.error("fetchAnswersButton not found")
  }
})