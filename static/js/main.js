window.onload = () => {
    copyButton();
}

document.addEventListener("DOMContentLoaded", () => {
    activeMenuItem();
})

function copyButton() {
    const pre = document.querySelectorAll("pre");

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
