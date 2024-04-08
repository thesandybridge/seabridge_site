window.onload = () => {
    copyButton();
    terminal();
    showTerm();
}

document.addEventListener("DOMContentLoaded", () => {
    activeMenuItem();
})

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

function processCommands(cmdInput, messagesDiv) {
    const menu_items = [...document.querySelectorAll('#main-nav a')].map(i => i.outerText);
    const commandText = cmdInput.value.trim();
    const args = parseArguments(commandText);
    appendToTarget(messagesDiv, `> ${commandText}`);

    switch (args[0]) {
        case 'clear':
            messagesDiv.innerHTML = '';
            break;
        case 'ls':
            appendToTarget(messagesDiv, menu_items.join(" "));
            break;
        case 'cd':
            if (args.length === 2) {
                const inputPath = args[1].toLowerCase();
                let newPath;

                if (inputPath === '/') {
                    newPath = window.location.origin;
                } else if (inputPath === 'home') {
                    newPath = window.location.origin;
                } else if (inputPath === '..') {
                    const pathParts = window.location.pathname.split('/');
                    const filteredParts = pathParts.filter(Boolean).slice(0, -1);
                    newPath = `${window.location.origin}/${filteredParts.join('/')}`;
                } else if (inputPath.startsWith('/')) {
                    newPath = window.location.origin + inputPath;
                } else {
                    const currPath = window.location.pathname.endsWith('/') ?
                        window.location.pathname.slice(0, -1) : window.location.pathname;
                    newPath = `${window.location.origin}${currPath}/${inputPath}`;
                }

                window.location.href = newPath;
            } else if (args.length === 1) {
                window.location.href = window.location.origin;
            } else {
                appendToTarget(target, "Invalid args");
            }
            break;
        case 'echo':
            if (args.length === 2) {
                appendToTarget(messagesDiv, args[1]);
            } else {
                appendToTarget(messagesDiv, "Invalid args");
            }
            break;
        case 'contact':
            appendToTarget(messagesDiv, "Nice try, it doesn't work yet... check me out on GitHub");
            break;
        case 'help':
            appendToTarget(messagesDiv, "Thanks for visiting my site! Here are some commands you can try...\n\tcd <arg>\n\techo <arg>\n\tclear\n\thelp\n\tcontact <your_email> <your_message>");
            break;
        default:
            appendToTarget(messagesDiv, "Invalid command");
    }

    cmdInput.remove();
    createNewInput(messagesDiv);
}

function appendToTarget(messagesDiv, message) {
    const newMessage = document.createElement('pre');
    newMessage.classList.add('ignore');
    newMessage.textContent = message;
    messagesDiv.appendChild(newMessage);
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
}

function createNewInput(messagesDiv) {
    const wrapper = document.querySelector('.cmd-wrapper');
    const newCmdInput = document.createElement('input');
    newCmdInput.type = 'text';
    newCmdInput.id = 'cmd';
    newCmdInput.autofocus = true;
    newCmdInput.spellcheck = false;

    wrapper.append(newCmdInput);

    newCmdInput.focus();

    newCmdInput.addEventListener('keydown', (event) => {
        if (event.key === 'Enter') {
            processCommands(newCmdInput, messagesDiv);
        }
    });
}


function positionCursorForNextLine(cmdInput) {
    cmdInput.focus();
}

function parseArguments(input) {
    const args = [];
    let currentArg = '';
    let inQuotes = false;

    for (const char of input) {
        if (char === '"') {
            inQuotes = !inQuotes;
        } else if (!inQuotes && char === ' ') {
            if (currentArg) {
                args.push(currentArg);
                currentArg = '';
            }
        } else {
            currentArg += char;
        }
    }

    if (currentArg) {
        args.push(currentArg);
    }

    return args;
}

function terminal() {
    const messagesDiv = document.querySelector('.msg');
    const initialCmdInput = document.querySelector('#cmd');

    initialCmdInput.addEventListener('keydown', (event) => {
        if (event.key === 'Enter') {
            processCommands(initialCmdInput, messagesDiv);
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

