"use strict";
// Appearance belongs to this origin/browser profile, not to the selected Agent.
class BrowserTheme {
    constructor() {
        // v2: v1 was written by the seed path below while the server default was
        // "dark", so every profile that ever loaded the page recorded a choice it
        // never made and could no longer leave. A stored v1 value cannot be told
        // apart from a deliberate pick, so it is dropped rather than trusted.
        this.key = "ensemble.theme.v2";
        this.legacyKey = "ensemble.theme.v1";
        this.preference = null;
        this.seeded = null;
        this.system = matchMedia("(prefers-color-scheme: dark)");
        try {
            this.preference = this.valid(localStorage.getItem(this.key));
            localStorage.removeItem(this.legacyKey);
        } catch { /* Session-only when storage is blocked. */ }
        this.system.addEventListener("change", () => this.render());
        window.addEventListener("storage", event => {
            // sessionStorage also emits storage events; it is not our preference.
            try { if (event.storageArea !== localStorage) return; } catch { return; }
            if (event.key !== this.key && event.key !== null) return;
            // Removal means use the browser's system preference, not server state.
            this.preference = this.valid(event.newValue) || "system";
            this.render();
        });
        document.querySelector("#theme-toggle").onclick = () => this.select(
            document.documentElement.dataset.bsTheme === "dark" ? "light" : "dark");
        document.querySelector("#theme").onchange = event => this.select(event.target.value);
        this.render();
    }
    valid(value) { return ["dark", "light", "system"].includes(value) ? value : null; }
    // The browser owns appearance; the server field is a legacy default that no
    // current client can even patch (Settings.cs: "current controls use browser
    // storage"). Applying it for this frame is kind; recording it as the user's
    // choice is not, because it then outranks the operating system forever.
    seed(legacyTheme) {
        if (this.preference !== null) return;
        this.seeded = this.valid(legacyTheme);
        this.render();
    }
    select(value) {
        if (!this.valid(value)) return;
        this.preference = value;
        try { localStorage.setItem(this.key, value); } catch { /* Keep this tab usable without durable storage. */ }
        this.render();
    }
    render() {
        const preference = this.preference || this.seeded || "system";
        const dark = preference === "dark" || (preference === "system" && this.system.matches);
        document.documentElement.dataset.theme = preference;
        document.documentElement.dataset.bsTheme = dark ? "dark" : "light";
        const button = document.querySelector("#theme-toggle");
        button.disabled = false;
        button.querySelector("use").setAttribute("href", dark ? "#icon-moon" : "#icon-sun");
        button.setAttribute("aria-label", dark ? "Switch to light theme" : "Switch to dark theme");
        button.title = dark ? "Dark theme; switch to light theme" : "Light theme; switch to dark theme";
        document.querySelector("#theme").value = preference;
    }
}
const browserTheme = new BrowserTheme();
