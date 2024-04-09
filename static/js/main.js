window.onload = () => {
    copyButton();
    showTerm();
}

document.addEventListener("DOMContentLoaded", () => {
    activeMenuItem();
    document.body.addEventListener('htmx:afterRequest', function(event) {
        if (event.detail.requestConfig.verb === "post") {
            document.getElementById('cmd').value = '';
        }
    });

    document.body.addEventListener('htmx:beforeSwap', function(event) {
        try {
            const response = JSON.parse(event.detail.xhr.responseText);
            event.preventDefault()
            switch (response.action) {
                case "clear":
                    document.querySelector('.msg').innerHTML = '';
                    break;
                case "open-url":
                    window.open(response.url, "_blank");
                    break;
                case "navigate":
                    window.location.href = response.url;
                    break;
                case "rotate":
                    const main = document.querySelector("body")
                    main.style = "transform: rotate(45deg)";
                    break;
                case "malware":
                    does_something_dangerous();
                    break;
            }
        } catch (e) {
        }
    });

    document.body.addEventListener('htmx:afterSwap', function() {
        scrollToBottom();
        focusInput();
    });

})

function does_something_dangerous() {
    const term = document.querySelector('.terminal-wrapper');
    term.style = "display: none;";
    const originalContent = document.querySelector('body').innerHTML;

    const body = document.querySelector('body');

    const maxX = window.innerWidth;
    const maxY = window.innerHeight;

    setInterval(() => {
        const newDiv = document.createElement('div');
        newDiv.classList.add('replicated-content');

        const randomX = Math.random() * maxX;
        const randomY = Math.random() * maxY;
        newDiv.style.left = randomX + 'px';
        newDiv.style.top = randomY + 'px';

        newDiv.innerHTML = originalContent;

        body.appendChild(newDiv);
    }, 200)
}

function scrollToBottom() {
    const terminalWrapper = document.querySelector('#terminal');
    terminalWrapper.scrollTop = terminalWrapper.scrollHeight;
}

function focusInput() {
    const inputField = document.getElementById('cmd');
    if (inputField) {
        inputField.focus();
    }
}


function copyButton() {
    const posts = document.querySelector(".post");

    if (!posts) return;

    const pre = document.querySelectorAll("pre:not(.ascii)");

    if (pre.length == 0) return;

    const svg = '<svg xmlns="http://www.w3.org/2000/svg" class="button-copy" aria-hidden="true" height="16" viewBox="0 0 16 16" version="1.1" width="16" data-view-component="true"><path fill-rule="evenodd" d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"></path><path fill-rule="evenodd" d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"></path></svg>';

    pre.forEach((code, idx) => {
        code.setAttribute("id", `code-${idx}`);
        const button = document.createElement("button");
        button.innerHTML = svg;
        button.setAttribute("class", "copy-button");
        button.setAttribute("data-clipboard-target", `#code-${idx}`);
        code.appendChild(button);

        button.addEventListener("click", () => {
            navigator.clipboard.writeText(code.innerText).then(() => {
                button.innerHTML = "Copied!";
                button.classList.add("copied");
                setTimeout(() => {
                    button.innerHTML = svg;
                    button.classList.remove("copied");
                }, 1000)
            })
        })
    })

}

function activeMenuItem() {
    const currentPath = window.location.pathname;
    const links = document.querySelectorAll('#main-nav a');

    links.forEach(function(link) {
        const href = link.getAttribute('href');
        if (href === currentPath) {
            link.classList.add('active');
        }
    });
}

function showTerm() {
    const wrapper = document.querySelector('.terminal-wrapper');
    const toggle = document.querySelector('#term-toggle');

    if (!toggle) return;

    toggle.addEventListener('click', () => {
        wrapper.classList.toggle('visible');
        const currentInput = wrapper.querySelector('input');
        if (currentInput && wrapper.classList.contains('visible')) {
            currentInput.value = '';
            currentInput.focus();
        }
    });

    wrapper.addEventListener('click', (event) => {
        if (event.target.tagName.toLowerCase() === 'input') {
            event.stopPropagation();
        } else {
            wrapper.classList.remove('visible');
            const currentInput = wrapper.querySelector('input');
            if (currentInput) {
                currentInput.value = '';
            }
        }
    });

    document.addEventListener('keydown', (event) => {
        if (event.key === 'Escape') {
            wrapper.classList.remove('visible');
            const currentInput = wrapper.querySelector('input');
            if (currentInput) {
                currentInput.value = '';
            }
        } else if (event.ctrlKey && (event.key === 'k' || event.key === 'K')) {
            event.preventDefault();
            wrapper.classList.toggle('visible');
            const currentInput = wrapper.querySelector('input');
            if (currentInput && wrapper.classList.contains('visible')) {
                currentInput.value = '';
                currentInput.focus();
            }
        }
    });
}


