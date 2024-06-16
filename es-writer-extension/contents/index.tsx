import type { PlasmoCSConfig } from "plasmo"

export const config: PlasmoCSConfig = {
  matches: ["<all_urls>"]
}

async function indexContents() {
  console.log("indexContents called")  // 関数呼び出しの確認用ログ
  try {
    const res = await fetch("http://localhost:8080/getAnswers", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ html: document.documentElement.outerHTML })
    })
    
    if (res.ok) {
      const answers = await res.json()
      console.log("Received answers:", answers) // 受け取ったデータをコンソールに出力
      replaceTextareaText(answers)
    } else {
      console.error("Network response was not ok", res.statusText)
    }
  } catch (error) {
    console.error("Error fetching answers:", error)
  }
}

console.log(document.documentElement.outerHTML)

function replaceTextareaText(answers: any) {
  const allTextareas = document.getElementsByTagName("textarea")
  Array.from(allTextareas).forEach((textarea, index) => {
    if (answers[index]) {
      textarea.value = answers[index].answer
    }
  })
}

export { indexContents }

document.addEventListener("DOMContentLoaded", () => {
  const fetchAnswersButton = document.getElementById("fetchAnswersButton")
  if (fetchAnswersButton) {
    fetchAnswersButton.addEventListener("click", indexContents)
    console.log("Event listener added to fetchAnswersButton")
  } else {
    console.error("fetchAnswersButton not found")
  }
})