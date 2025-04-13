let state = {
    socket: null,
    userId: null,
    quiz: null,
    currentQuestionIndex: 0,
    selectedOption: null,
    timerInterval: null,
    timeRemaining: 30
};

const MESSAGE_TYPES = {
    JOIN: 'join',
    JOINED: 'joined',
    ANSWER: 'answer',
    RESULT: 'result',
    LEADERBOARD: 'leaderboard',
    ERROR: 'error'
};

function getElement(id) {
    return document.getElementById(id);
}

// Connect to WebSocket
function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;

    state.socket = new WebSocket(wsUrl);

    state.socket.onopen = () => console.log('WebSocket connected');
    state.socket.onmessage = handleSocketMessage;
    state.socket.onerror = (error) => console.error('WebSocket error:', error);
    state.socket.onclose = () => console.log('WebSocket disconnected');
}

// Send message through WebSocket
function sendMessage(type, payload) {
    if (!state.socket || state.socket.readyState !== WebSocket.OPEN) {
        console.error('WebSocket not connected');
        return;
    }

    state.socket.send(JSON.stringify({ type, payload }));
}

// Handle incoming WebSocket messages
function handleSocketMessage(event) {
    const message = JSON.parse(event.data);
    console.log('Received message:', message);

    switch (message.type) {
        case MESSAGE_TYPES.JOINED:
            handleJoined(message.payload);
            break;
        case MESSAGE_TYPES.LEADERBOARD:
            updateLeaderboard(message.payload);
            break;
        case MESSAGE_TYPES.RESULT:
            showResult(message.payload);
            break;
        case MESSAGE_TYPES.ERROR:
            showError(message.payload.message);
            break;
    }
}

// Handle successful join
function handleJoined(payload) {
    state.userId = payload.userId;

    fetch(`/api/quiz/${payload.quizId}`)
        .then(response => response.json())
        .then(data => {
            state.quiz = data;
            initializeQuiz();
            showQuestion(0);
        })
        .catch(error => {
            console.error('Error fetching quiz:', error);
            showError('Failed to load quiz. Please try again.');
        });
}

function initializeQuiz() {
    getElement('quizTitle').textContent = `Quiz ${state.quiz.id}`;
    getElement('userNameDisplay').textContent = getElement('userName').value;

    getElement('joinForm').classList.add('hidden');
    getElement('quizSection').classList.remove('hidden');
}

function showError(message) {
    alert(`Error: ${message}`);
}

function showQuestion(index) {
    state.currentQuestionIndex = index;
    state.selectedOption = null;

    if (index >= state.quiz.questions.length) {
        getElement('questionText').textContent = 'Quiz completed!';
        getElement('options').innerHTML = '';
        return;
    }

    const question = state.quiz.questions[index];
    getElement('questionText').textContent = question.text;

    const optionsContainer = getElement('options');
    optionsContainer.innerHTML = '';

    question.options.forEach((option, i) => {
        const optionElement = document.createElement('div');
        optionElement.className = 'option';
        optionElement.textContent = option;
        optionElement.dataset.index = i;
        optionElement.addEventListener('click', () => selectOption(optionElement, i));
        optionsContainer.appendChild(optionElement);
    });

    getElement('result').classList.add('hidden');
    getElement('nextButton').classList.add('hidden');

    // Start timer
    startTimer(question);
}

function startTimer(question) {
    state.timeRemaining = question.timeLimit || 30;
    const timerElement = getElement('timer');
    timerElement.textContent = state.timeRemaining;

    clearInterval(state.timerInterval);
    state.timerInterval = setInterval(() => {
        state.timeRemaining--;
        timerElement.textContent = state.timeRemaining;

        if (state.timeRemaining <= 0) {
            clearInterval(state.timerInterval);
            if (state.selectedOption === null) {
                submitAnswer(-1);
            }
        }
    }, 1000);
}

function selectOption(element, index) {
    document.querySelectorAll('.option').forEach(opt => {
        opt.classList.remove('selected');
    });

    element.classList.add('selected');
    state.selectedOption = index;
}

// Submit answer
function submitAnswer(answerIndex) {
    clearInterval(state.timerInterval);

    const question = state.quiz.questions[state.currentQuestionIndex];

    sendMessage(MESSAGE_TYPES.ANSWER, {
        quizId: state.quiz.id,
        userId: state.userId,
        questionId: question.id,
        answer: answerIndex
    });
}

// Show result
function showResult(resultData) {
    const resultElement = getElement('result');
    resultElement.classList.remove('hidden', 'correct', 'incorrect');

    if (resultData.correct) {
        resultElement.classList.add('correct');
        resultElement.textContent = 'Correct! +10 points';
    } else {
        resultElement.classList.add('incorrect');
        resultElement.textContent = 'Incorrect!';
    }

    getElement('userScore').textContent = resultData.score;

    if (state.currentQuestionIndex < state.quiz.questions.length - 1) {
        getElement('nextButton').classList.remove('hidden');
    } else {
        setTimeout(() => {
            getElement('questionText').textContent = 'Quiz completed!';
            getElement('options').innerHTML = '';
            resultElement.classList.add('hidden');
        }, 2000);
    }
}

// Update leaderboard
function updateLeaderboard(entries) {
    const leaderboardEntries = getElement('leaderboardEntries');
    leaderboardEntries.innerHTML = '';

    entries.forEach((entry, index) => {
        const entryElement = document.createElement('div');
        entryElement.className = 'leaderboard-entry';

        const rank = document.createElement('span');
        rank.textContent = `${index + 1}.`;

        const name = document.createElement('span');
        name.textContent = entry.userName;

        const joinedAt = document.createElement('span');
        joinedAt.textContent = new Date(entry.joinedAt).toLocaleString();

        const lastActive = document.createElement('span');
        lastActive.textContent = new Date(entry.lastActive).toLocaleString();

        const score = document.createElement('span');
        score.textContent = `${entry.score} points`;

        entryElement.appendChild(rank);
        entryElement.appendChild(name);
        entryElement.appendChild(joinedAt);
        entryElement.appendChild(lastActive);
        entryElement.appendChild(score);

        if (entry.userId === state.userId) {
            entryElement.style.fontWeight = 'bold';
        }

        leaderboardEntries.appendChild(entryElement);
    });
}

function nextQuestion() {
    showQuestion(state.currentQuestionIndex + 1);
}

function setupEventListeners() {
    // Join button
    getElement('joinButton').addEventListener('click', () => {
        const quizId = getElement('quizId').value.trim();
        const userName = getElement('userName').value.trim();

        if (!quizId || !userName) {
            showError('Please enter both Quiz ID and your name');
            return;
        }

        connectWebSocket();

        setTimeout(() => {
            sendMessage(MESSAGE_TYPES.JOIN, {
                quizId: quizId,
                userName: userName
            });
        }, 500);
    });

    getElement('nextButton').addEventListener('click', nextQuestion);

    getElement('options').addEventListener('click', (event) => {
        if (event.target.classList.contains('option') && !getElement('result').classList.contains('hidden')) {
            return;
        }

        if (event.target.classList.contains('option')) {
            const index = parseInt(event.target.dataset.index);
            submitAnswer(index);
        }
    });
}

function init() {
    setupEventListeners();
}

document.addEventListener('DOMContentLoaded', init); 