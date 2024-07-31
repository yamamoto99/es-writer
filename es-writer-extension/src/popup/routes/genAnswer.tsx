/// <reference types="chrome"/>
import { api_endpoint } from "../../contents/index"

async function genAnswer() {
  return new Promise<void>((resolve, reject) => {
    chrome.tabs.query({ active: true, currentWindow: true }, async (tabs) => {
      if (tabs[0] && tabs[0].id !== undefined) {
        chrome.tabs.sendMessage(
          tabs[0].id,
          { action: "getHTML" },
          async (response) => {
            if (response && response.html) {
              const html_source = response.html
              console.log("html loaded")
              try {
                const apiResponse = await fetch(
                  api_endpoint + "/app/generate/generateAnswers",
                  {
                    method: "POST",
                    headers: {
                      "Content-Type": "application/json"
                    },
                    body: JSON.stringify({ html: html_source })
                  }
                )
                if (!apiResponse.ok) {
                  console.error(
                    "Network response was not ok",
                    apiResponse.statusText
                  )
                  reject(new Error("Network response was not ok"))
                  return
                }
                const answers = await apiResponse.json()
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

function replaceTextareaText(answers: any) {
  chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
    if (tabs[0] && tabs[0].id !== undefined) {
      chrome.tabs.sendMessage(tabs[0].id, {
        action: "replaceTextareas",
        answers: answers
      })
    }
  })
}

export default genAnswer
