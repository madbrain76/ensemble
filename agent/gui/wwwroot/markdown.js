"use strict";
// Presentation only: callers retain original Markdown for History and speech.
// Never enable raw HTML or plugins without reviewing this security boundary.
const messageMarkdown = window.markdownit({ html: false, linkify: false, typographer: false });
messageMarkdown.validateLink = url => /^https?:\/\//i.test(url);
messageMarkdown.renderer.rules.image = (tokens, index) =>
    messageMarkdown.utils.escapeHtml(`[Image: ${tokens[index].content}]`);
messageMarkdown.renderer.rules.link_open = (tokens, index, options, env, self) => {
    tokens[index].attrSet("rel", "noopener noreferrer");
    return self.renderToken(tokens, index, options);
};
function markdown(parent, text) {
    const openCode = new Set([...parent.querySelectorAll("details[data-code-index][open]")]
        .map(node => node.dataset.codeIndex));
    const tableScroll = [...parent.querySelectorAll(".markdown-table")].map(node => node.scrollLeft);
    // This HTML is generated only by the pinned parser with raw HTML disabled,
    // restrictive links, and image tokens replaced by escaped text.
    parent.innerHTML = messageMarkdown.render(text);
    for (const [index, table] of [...parent.querySelectorAll("table")].entries()) {
        const wrapper = document.createElement("div");
        wrapper.className = "markdown-table";
        wrapper.tabIndex = 0;
        wrapper.setAttribute("role", "region");
        wrapper.setAttribute("aria-label", "Scrollable table");
        table.before(wrapper);
        wrapper.append(table);
        wrapper.scrollLeft = tableScroll[index] || 0;
        for (const cell of table.querySelectorAll("thead th")) cell.scope = "col";
    }
    for (const [index, pre] of [...parent.querySelectorAll("pre")].entries()) {
        const code = pre.querySelector("code");
        if (!code) continue;
        const source = code.textContent;
        code.replaceChildren();
        highlight(code, source);
        if (source.length <= 1800 && source.split("\n").length <= 24) continue;
        const details = document.createElement("details");
        details.dataset.codeIndex = String(index);
        details.open = openCode.has(String(index));
        const summary = document.createElement("summary");
        summary.textContent = "Expand complete code";
        const preview = document.createElement("pre");
        preview.textContent = source.slice(0, 1800).split("\n").slice(0, 24).join("\n") + "\n…";
        preview.hidden = details.open;
        pre.before(preview, details);
        details.append(summary, pre);
        details.ontoggle = () => { preview.hidden = details.open; };
    }
    // A fenced block is content the operator was meant to take away, so it becomes
    // a copyable box. The payload lives in a read-only field and stays literal, so
    // a copy carries what the model wrote, including lines a long block collapsed.
    for (const pre of [...parent.querySelectorAll("pre")]) {
        if (!pre.querySelector("code")) continue;
        // Anchor on the collapsible wrapper when there is one: the copy control has
        // to stay visible even while the block itself is folded away.
        const anchor = pre.closest("details") || pre;
        if (anchor.closest(".clipboard-block")) continue;
        // A fence contributes one trailing newline; the copy must be what the model
        // wrote, not what the markup added, and a collapsed preview shows "…".
        const payload = pre.textContent.replace(/\n…$/, "").replace(/\n$/, "");
        if (!payload.trim()) continue;
        const block = document.createElement("div");
        block.className = "clipboard-block";
        const field = document.createElement("textarea");
        field.readOnly = true;
        field.className = "clipboard-payload";
        field.tabIndex = -1;
        field.setAttribute("aria-label", "Copyable text");
        field.value = payload;
        const button = document.createElement("button");
        button.type = "button";
        button.className = "btn btn-quiet btn-sm clipboard-copy";
        button.textContent = "Copy text";
        button.onclick = async () => {
            const reset = () => { button.textContent = "Copy text"; };
            let copied = false;
            try {
                if (navigator.clipboard?.writeText) {
                    await navigator.clipboard.writeText(payload);
                    copied = true;
                }
            } catch (error) { copied = false; }
            if (!copied) {
                // Clipboard write needs a transient selection; keep it invisible and
                // restore the scroll so the visible block never jumps.
                const top = field.scrollTop;
                field.select();
                field.setSelectionRange(0, payload.length);
                field.scrollTop = top;
                try { copied = document.execCommand("copy"); } catch (error) { copied = false; }
            }
            // Say what happened, then hand focus back to the control that answered.
            button.textContent = copied ? "Copied" : "Select the text and press Ctrl+C";
            button.focus();
            clearTimeout(block.timer);
            block.timer = setTimeout(reset, copied ? 1600 : 6000);
        };
        const actions = document.createElement("div");
        actions.className = "clipboard-actions";
        actions.append(button);
        anchor.before(block);
        block.append(field, actions, anchor);
    }
}
