document.addEventListener("DOMContentLoaded", function() {
    const fetchAnswersButton = document.getElementById("fetchAnswersButton");
    const surveyItem = document.getElementById("pnl_SurveyItem");

    function displayLoading() {
        const loadingDiv = document.createElement("div");
        loadingDiv.className = "loading";
        loadingDiv.textContent = "Loading...";
        surveyItem.appendChild(loadingDiv);
    }

    function displayError(message) {
        const errorDiv = document.createElement("div");
        errorDiv.className = "error";
        errorDiv.textContent = "Error: " + message;
        surveyItem.appendChild(errorDiv);
    }

    function displayAnswers(answers) {
        console.log("Received answers:", answers);  // 受け取った回答をコンソールに出力
        answers.forEach((answer, index) => {
            const textarea = document.getElementById(`tbx_${index + 1}`);
            if (textarea) {
                textarea.value = answer.answer;
            }
        });
    }

    function fetchAnswers() {
        const htmlContent = document.documentElement.outerHTML;

        surveyItem.innerHTML = ''; // Clear previous content
        displayLoading();

        fetch('http://localhost:8080/getAnswers', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ html: htmlContent })
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok ' + response.statusText);
            }
            return response.json();
        })
        .then(data => {
            surveyItem.innerHTML = ''; // Clear loading message
            displayAnswers(data);
        })
        .catch(error => {
            surveyItem.innerHTML = ''; // Clear loading message
            displayError(error.message);
        });
    }

    fetchAnswersButton.addEventListener("click", fetchAnswers);
});