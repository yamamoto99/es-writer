/// <reference types="chrome"/>

function openProfileForm() {
  chrome.tabs.create({ url: chrome.runtime.getURL('src/pages/profile.html') });
}

export default openProfileForm;
