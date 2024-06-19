/// <reference types="chrome"/>

function openProfileForm() {
  console.log("openProfile called")
  chrome.tabs.create({ url: chrome.runtime.getURL("../pages/profile.html") });
}

export default openProfileForm;
