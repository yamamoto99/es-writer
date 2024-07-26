/// <reference types="chrome"/>

import { api_endpoint } from "../../contents/index"

async function genAnswer() {
  return new Promise<void>((resolve, reject) => {
    chrome.tabs.query({ active: true, currentWindow: true }, async (tabs) => {
      if (tabs[0] && tabs[0].id !== undefined) {
        chrome.scripting.executeScript(
          {
            target: { tabId: tabs[0].id },
            func: getActiveTabHTML
          },
          async (result) => {
            if (result && result[0]) {
              const html_source = result[0].result
              console.log("html loaded")

              try {
                const response = await fetch(
                  api_endpoint + "/app/generate/generateAnswers",
                  {
                    method: "POST",
                    headers: {
                      "Content-Type": "application/json"
                    },
                    body: JSON.stringify({ html: html_source })
                  }
                )

                if (!response.ok) {
                  console.error(
                    "Network response was not ok",
                    response.statusText
                  )
                  reject(new Error("Network response was not ok"))
                  return
                }

                const answers = await response.json()
                console.log("Received answers:", answers)
                replaceTextareaText(answers)
                resolve()
              } catch (error) {
                console.error("Fetch error:", error)
                reject(error)
              }
            } else {
              console.error("Failed to get active tab HTML.")
              reject(new Error("Failed to get active tab HTML."))
            }
          }
        )
      } else {
        console.error("No active tab found or tab ID is undefined.")
        reject(new Error("No active tab found or tab ID is undefined."))
      }
    })
  })
}

function getActiveTabHTML() {
  return document.documentElement.outerHTML
}

function replaceTextareaText(answers: any) {
  chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
    if (tabs[0] && tabs[0].id !== undefined) {
      chrome.scripting.executeScript({
        target: { tabId: tabs[0].id },
        func: replaceTextareas,
        args: [answers]
      })
    }
  })
}

function replaceTextareas(answers: any) {
  const allTextareas = document.getElementsByTagName("textarea")
  Array.from(allTextareas).forEach((textarea, index) => {
    if (answers[index]) {
      textarea.value = answers[index].answer
    }
  })
}

export default genAnswer
