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

function processCommands(buffer, target, e) {
    e.preventDefault();

    const lines = buffer.split('\n');
    let ignore = false;
    const menu_items = [...document.querySelectorAll('#main-nav a')].map(i => i.outerText);

    lines.forEach((line, index) => {
        if (index === lines.length - 1) {
            const args = parseArguments(line);

            switch (args[0]) {
                case 'clear':
                    target.value = '';
                    ignore = true;
                    break;
                case 'ls':
                    appendToTarget(target, menu_items.join(" "))
                    break;
                case 'cd':
                    if (args.length === 2) {
                        const path = args[1];
                        switch (path.toLowerCase()) {
                            case '/':
                                window.location.href = path.toLowerCase();
                                break;
                            case 'home':
                                window.location.href = '/';
                                break;
                            case '..':
                                window.location.href = '/';
                                break
                            default:
                                window.location.href = '/'+path.toLowerCase();
                        }
                    } else if (args.length === 1) {
                        window.location.href = '/';
                    } else {
                        appendToTarget(target, "Invalid args");
                    }
                    break;
                case 'echo':
                    if (args.length === 2) {
                        appendToTarget(target, args[1]);
                    } else {
                        appendToTarget(target, "Invalid args");
                    }
                    break;
                case 'contact':
                    appendToTarget(target, "Nice try, it doesn't work yet... check me out on github");
                    break;
                case 'help':
                    appendToTarget(target, "Thanks for visiting my site! Here are some commands you can try...\n\tcd <arg>\n\techo <arg>\n\tclear\n\thelp\n\tcontact <your_email> <your_message>");
                    break;
                default:
                    appendToTarget(target, "");
                    ignore = true
            }
        }
    });

    positionCursorForNextLine(target, ignore);
}

function appendToTarget(target, message) {
    target.value += (target.value ? "\n" : "") + message;
}

function positionCursorForNextLine(target, ignore) {
    if (ignore) return
    target.value += '\n';
    target.scrollTop = target.scrollHeight;
    const end = target.value.length;
    target.setSelectionRange(end, end);
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
    const term = document.querySelector('#terminal');
    if (!term) return

    term.addEventListener('keydown', (event) => {
        if (event.key === 'Enter') {
            processCommands(term.value, term, event);
        }
    })

}


function showTerm() {
    const wrapper = document.querySelector('.terminal-wrapper');
    const toggle = document.querySelector('#term-toggle');
    const textarea = document.querySelector('#terminal');

    if (!toggle) return

    toggle.addEventListener('click', () => {
        wrapper.classList.toggle('visible');
        textarea.value = '';
    });

    textarea.addEventListener('click', (event) => {
        event.stopPropagation();
    });

    wrapper.addEventListener('click', () => {
        wrapper.classList.remove('visible');
        textarea.value = '';
    });

    document.addEventListener('keydown', (event) => {
        if (event.key === 'Escape') {
            wrapper.classList.remove('visible');
            textarea.value = '';
        }
    });
}

