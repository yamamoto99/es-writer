/// <reference types="chrome"/>

function openProfileForm() {
  console.log("openProfile called")
  chrome.tabs.create({ url: chrome.runtime.getURL("../tabs/profile.html") })
}

export default openProfileForm
