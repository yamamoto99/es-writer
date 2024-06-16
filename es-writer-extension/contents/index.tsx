// import { useState } from "react"
import type { PlasmoCSConfig } from 'plasmo'

export const config: PlasmoCSConfig = {
  matches: ["<all_urls>"]
}

async function indexContents() {
  // const [data, setData] = useState('');

  const res = await fetch(
    'http://localhost:8080/getAnswers',
    {
      method: 'POST',
      mode: 'no-cors',
      body: document.body.innerHTML
    }
  );

  if (res.ok) {
    replaceTextareaText(res.json());
  }
}

console.log(document.body.innerHTML);

function replaceTextareaText(answers: any) {
  const allTextareas = document.getElementsByTagName('textarea');
  Array.from(allTextareas).forEach((Textarea, index) => {
    Textarea.value = answers[index].Answer;
  });
}

export { indexContents }