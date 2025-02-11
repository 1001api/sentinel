import { marked } from "marked";

interface History {
    role: string
    content: string
}

// history storage
const storage = window.localStorage;
const HISTORY_STORAGE = "HISTORY_STORAGE";
const MESSAGE_USER = "user";
const MESSAGE_ASSISTANT = "assistant";
let histories: History[] = [];

if (storage.getItem(HISTORY_STORAGE)) {
    histories = JSON.parse(storage.getItem(HISTORY_STORAGE) || "[]");
} else {
    storage.setItem(HISTORY_STORAGE, JSON.stringify([]));
}

const projectID = document.getElementById("project-id")?.textContent;
const id = projectID ? JSON.parse(projectID) : null;

// Get DOM elements
const chatContainer = document.querySelector<HTMLDivElement>('#chat-container');
const chatInput = document.querySelector<HTMLInputElement>('#chat-input');
const chatSubmit = document.querySelector<HTMLButtonElement>('#chat-submit');
const chatClearBtn = document.querySelector<HTMLButtonElement>('#chat-clear-btn');

function parseSSEData(data: string) {
    if (!data.startsWith('data: ')) return null;
    // Remove 'data: ' prefix and clean up the quoted content
    let content = data.slice(6)
        .replace(/^"/, '') // Remove leading quote
        .replace(/"$/, ''); // Remove trailing quote
    return content;
}

function cleanText(text: string) {
    return text
        .replace(/\\n/g, '\n') // Convert \n to actual newlines
        .trim();
}

function scrollToBottom() {
    if (!chatContainer) return;

    const scrollOptions: ScrollToOptions = {
        top: chatContainer.scrollHeight,
        behavior: "smooth",
    };

    chatContainer.scrollTo(scrollOptions);
}

function createMessageBubble(message: string, isUser = false) {
    const bubbleWrapper = document.createElement('div');
    bubbleWrapper.className = `flex w-full mb-4 ${isUser ? 'justify-end' : 'justify-start'}`;

    const bubble = document.createElement('div');
    bubble.className = `max-w-[70%] font-inherit p-3 rounded-lg overflow-x-auto whitespace-pre-wrap ${isUser
        ? 'bg-blue-500 text-white text-sm rounded-br-none'
        : 'bg-gray-200 dark:bg-gray-600 text-sm text-gray-900 dark:text-white rounded-bl-none'
        }`;

    bubble.textContent = message;
    bubbleWrapper.appendChild(bubble);
    return bubbleWrapper;
}

function addMessage(message: string, isUser = false, isStreaming = false) {
    if (!chatContainer) return;

    const bubble = createMessageBubble(message, isUser);
    chatContainer.appendChild(bubble);

    if (isStreaming) {
        const textElem = bubble.querySelector("div");
        if (!textElem) return;
        handleStreamingResponse(textElem);
    }
}

chatSubmit?.addEventListener('click', handleSubmit);
chatInput?.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') handleSubmit();
});
chatClearBtn?.addEventListener('click', () => {
    if (!chatContainer) return;

    storage.setItem(HISTORY_STORAGE, JSON.stringify([]));
    histories = [];
    chatContainer.innerHTML = "";
})

function handleSubmit() {
    if (!chatInput || !chatSubmit) return;

    const message = chatInput.value.trim();
    if (message === '') return;

    // Disable input and button while processing
    chatInput.disabled = true;
    chatSubmit.disabled = true;

    // add user message
    addMessage(message, true);

    // add AI message
    addMessage("Generating response...", false, true);

    // add message to history
    saveMessage(MESSAGE_USER, message);

    // clear input
    chatInput.value = '';

    // Re-enable input and button
    chatInput.disabled = false;
    chatSubmit.disabled = false;
    chatInput.focus();
}

async function handleStreamingResponse(bubble: HTMLElement) {
    try {
        const response = await fetch(
            `/api/ai/stream/summary?query=${chatInput?.value}&projectId=${id}&provider=openai`,
            {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    history: histories,
                })
            }
        );
        if (!response.ok || !response.body) {
            throw new Error(response.status.toString());
        }

        const reader = response.body.getReader();
        const decoder = new TextDecoder();
        let accumulatedText = '';
        let buffer = '';

        while (true) {
            const { value, done } = await reader.read();

            if (done) break;

            // Decode the chunk
            buffer += decoder.decode(value, { stream: true });

            // Process each line in the buffer
            const lines = buffer.split('\n');
            buffer = lines.pop() || '';

            for (const line of lines) {
                const content = parseSSEData(line);

                if (content !== null) {
                    accumulatedText += content;
                    const cleanedText = cleanText(accumulatedText);
                    const markdown = await marked.parse(cleanedText);
                    bubble.innerHTML = markdown;

                    scrollToBottom();
                }
            }
        }

        saveMessage(MESSAGE_ASSISTANT, accumulatedText);
    } catch (error) {
        console.error('Streaming error:', error);

        // if error is too much request (rate limited)
        if (error.message === "429") {
            bubble.innerHTML = "<strong>[Oopss, you are using DEMO account]</strong><br/>Requests are limited to only 5 requests per minute. Please wait for a while...";
            scrollToBottom();
            return;
        }

        bubble.textContent = "Sorry, there was an error processing your request.";
    }
}

function saveMessage(role: string, message: string) {
    const obj: History = { role: role, content: message };
    histories.push(obj);
    storage.setItem(HISTORY_STORAGE, JSON.stringify(histories));
}

addMessage("Hello! How can I help you today?");

if (histories) {
    for (const v of histories) {
        addMessage(v.content, v.role === "user", false);
    }
}
