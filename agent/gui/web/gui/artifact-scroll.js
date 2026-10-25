import { markdown, ansi, partText } from "./renderers.js";

// A scroll owns only its cards and supplied element, never application controls.
export class ArtifactScroll {
  cards = new Map();
  calls = new Map();
  constructor(element, speech = null) {
    // Speech is optional: a component consumer need not construct the chat shell.
    this.speech = speech;
    this.element = element;
  }
  reset() {
    this.cards.clear();
    this.calls.clear();
    this.element.replaceChildren();
  }
  // One card shape serves every kind. The speaker closure reads current text,
  // so a card created during streaming speaks its full finalized value later.
  card(id, kind) {
    id = String(id);
    if (this.cards.has(id)) return this.cards.get(id);
    const element = document.createElement("article");
    element.className = `artifact artifact--${kind}`;
    element.dataset.id = id;
    const header = document.createElement("header");
    const label = document.createElement("span");
    label.textContent = kind;
    const speaker = document.createElement("button");
    speaker.textContent = "🔊";
    speaker.setAttribute("aria-label", "Read artifact aloud");
    const timing = document.createElement("span");
    timing.className = "artifact__timing";
    timing.hidden = true;
    const content = document.createElement("div");
    content.className = "content";
    header.append(label, timing, speaker);
    element.append(header, content);
    this.element.append(element);
    const card = { element, label, timing, content, text: "", kind, streamed: false, at: null, ms: 0 };
    speaker.hidden = !this.speech;
    speaker.addEventListener("click", () => this.speech?.enqueue(card.text));
    this.cards.set(id, card);
    return card;
  }
  // Follow the stream only when the reader was already at the bottom. Reading
  // an earlier artifact should not fight an unconditional scroll-to-bottom.
  render(card) {
    const nearEnd =
      this.element.clientHeight + this.element.scrollTop >= this.element.scrollHeight - 180;
    card.content.replaceChildren();
    // A cap bounds the normal card while the expander retains every byte.
    const shown = card.text.slice(0, 4000);
    const body = document.createElement("div");
    card.content.append(body);
    if (card.kind === "text" || ["thinking", "reasoning_summary"].includes(card.kind)) markdown(body, shown);
    else {
      const pre = document.createElement("pre");
      if (card.label.textContent === "edit_file") {
        for (const line of shown.split("\n")) {
          const row = document.createElement("div");
          row.className = line.startsWith("+")
            ? "diff-add"
            : line.startsWith("-")
              ? "diff-remove"
              : "";
          row.textContent = line;
          pre.append(row);
        }
      } else ansi(pre, shown);
      body.append(pre);
    }
    if (["thinking", "reasoning_summary"].includes(card.kind)) {
      const details = document.createElement("details");
      const summary = document.createElement("summary");
      summary.textContent = "Thinking";
      body.replaceWith(details);
      details.append(summary, body);
    }
    if (card.text.length > shown.length) {
      const details = document.createElement("details");
      const summary = document.createElement("summary");
      summary.textContent = "Full artifact";
      const full = document.createElement("pre");
      full.textContent = card.text;
      details.append(summary, full);
      card.content.append(details);
    }
    if (nearEnd) this.element.scrollTo(0, this.element.scrollHeight);
  }
  // Show only server measurements: missing legacy timestamps are not arrival time.
  paint(card) {
    const bits = [];
    if (card.at) bits.push(clock(card.at));
    if (card.ms) bits.push(elapsed(card.ms));
    card.timing.textContent = bits.join(" \u00b7 ");
    card.timing.title = card.at && card.ms
      ? `${card.at.toISOString()}, ${elapsed(card.ms)} to produce`
      : card.at
        ? card.at.toISOString()
        : "";
    card.timing.hidden = bits.length === 0;
  }
  // History replay is silent. Finals reconcile display text, while already
  // streamed speech is never re-enqueued just because a final arrives.
  receive(message) {
    const id = message.part_id;
    if (message.type === "part_delta") {
      const card = this.card(id, message.kind);
      card.text += message.chunk;
      card.streamed = true;
      this.render(card);
      if (!message.replay && ["text", "thinking", "reasoning_summary"].includes(message.kind))
        this.speech?.chunk(id, message.chunk, false, message.kind);
    } else if (
      message.type === "part_final" ||
      message.type === "user_message"
    ) {
      const part = message.part;
      const kind =
        message.type === "user_message"
          ? "user"
          : part.type === "opaque"
            ? (part.from?.surface === "responses" ? "reasoning_summary" : "thinking")
            : part.type;
      const card = this.card(id, kind);
      card.kind = kind;
      card.text = partText(part);
      card.at = stamp(message.at);
      // This is the part's own observed span, not latency since the prompt.
      card.ms = Number.isFinite(part.elapsed) && part.elapsed >= 1 ? part.elapsed : 0;
      this.render(card);
      this.paint(card);
      if (part.type === "tool_call") this.calls.set(part.call_id, card);
      if (!message.replay && ["text", "thinking", "reasoning_summary"].includes(kind))
        this.speech?.chunk(id, card.streamed ? "" : card.text, true, kind);
    } else if (message.type === "tool_dispatched") {
      const card =
        this.calls.get(message.call_id) ||
        this.card(`call-${message.call_id}`, "tool");
      this.calls.set(message.call_id, card);
      card.kind = "tool";
      card.label.textContent = message.name;
      card.text = `${message.name}\n${JSON.stringify(message.input, null, 2)}`;
      this.render(card);
      if (!message.replay) {
        this.speech?.flush();
        this.speech?.enqueue(message.name, "tool");
      }
    } else if (message.type === "tool_finished") {
      // Tool output stays silent automatically, but the same speaker button
      // still provides full on-demand access to arguments and result text.
      const card =
        this.calls.get(message.call_id) ||
        this.card(`call-${message.call_id}`, "tool");
      card.text += `\n\n${message.is_error ? "Error: " : ""}${message.result}`;
      if (message.is_error) card.element.classList.add("artifact--error");
      this.render(card);
    }
  }
}

function stamp(value) {
  if (typeof value !== "string" || !value || value === "0001-01-01T00:00:00Z") return null;
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? null : date;
}

function clock(date) {
  return date.toLocaleTimeString([], { hour12: false });
}

function elapsed(ms) {
  if (!Number.isFinite(ms) || ms < 1) return "";
  if (ms < 1000) return `${ms}ms`;
  const seconds = ms / 1000;
  if (seconds < 60) return `${seconds.toFixed(1)}s`;
  const rounded = Math.round(seconds);
  return `${Math.floor(rounded / 60)}m ${String(rounded % 60).padStart(2, "0")}s`;
}
