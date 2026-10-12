# Preface

A coding agent can write a great deal of code without learning why any of it
should have that shape. This book asks a harder question: can an experienced
engineer put enough of that judgment into a book that another agent can build
the system, discover what the teaching missed, and help improve the next edition?

The second edition is the result of trying it. Codex coordinated a fresh student
implementation of Ensemble, chapter by chapter, with code reviews and real use.
The completed implementation then supplied the evidence for revising the prose
and graders. Some instructions worked. Some omitted facts a student needed.
Some tests accepted the wrong thing, and some rejected a correct implementation.
Those discoveries belong in the book a reader gets next.

We call this an **agentic codebook**: a technical book whose exercises build a
working program, with executable checks alongside the teaching. The larger
ambition is a loop in which the program helps produce the next edition. This
edition has reached working software and a revised book. The move from Codex
into the new Ensemble has begun; it has not yet produced Edition 3.

## An agent for one engineer

Bill Cox's first-edition account starts with a bet. After forty years writing
compilers and chip-design tools, he told his team he could build a better proof
of concept than the coding agent they were using in two weeks. His manager asked
him to prove it. He built the demonstration, spent a week deciding whether to
keep it, deleted every line, and built it again properly.

The deletion matters. Producing a convincing demonstration and wanting to live
inside its code are different milestones. The later agent, CodeRhapsody, grew
through daily use, correction and arguments about how the machinery should work.
This book preserves those reasons so the next builder need not discover every
one by breaking the same thing.

The intended prize is an agent that fits the way you work. Bill has low vision.
He listens to synthesized speech and steers the agent while it works, before a
mistake turns into another page of edits. Speech, observable progress and the
ability to intervene are part of his working interface. Your priorities may be
different. The chapters distinguish structural decisions from choices you can
make for yourself, and the exercises give those choices somewhere to land.

You should already be comfortable specifying, reading and changing substantial
programs. An AI coding assistant is a working tool throughout the course. You
remain responsible for understanding what it built. A plausible explanation
beside a passing command is an excellent place for a bug to hide.

## Read the first page; keep the reasons

The implementation chapters begin with a TL;DR: what to build, the public
contracts other code depends on, the rules the checks exercise and the decisions
left to you. Read that first to understand the assignment. Read the explanation
when you need to understand why a boundary is there, or why the obvious shortcut
will make the next chapter harder.

The format comes from a programming book Bill found in a university library in
1980. In his account, he read the first page of every chapter in under an hour,
returned the book, and could write BASIC. The rest of each chapter was available
when he needed it. That is a useful bargain for a reader who reads slowly and
wants every sentence to earn its place.

Here the explanation also carries the failures: the setting that was saved but
never applied, the tool result a test never actually checked, the agent that
finished its task while the GUI still said it was busy. These are small stories
with expensive consequences. Knowing the eventual method name would not have
prevented them. Knowing what the method must preserve might.

Most implementation exercises have a local grader. It starts small fake services,
runs your actual program, and checks observable behavior. These checks need no
paid model inference. Live use answers different questions: whether the model
understands the tools, whether the interface lets you intervene, and whether the
result is useful. Keep those results separate. A scripted server can establish
that a file reached the model; it cannot establish that a real model understood
the file.

The accepted chapter solutions are public. Attempt the exercise before reading
its answer. If your design goes sideways, inspect the reference and begin the
next chapter from a working predecessor. Failing one exercise should not end
your course. Keep your failed attempt long enough to understand the difference.

## What changed between editions

The first attempt at this edition went wrong before it had much software to
review. Authors expanded small exercises into elaborate contracts, and reviewers
approved those contracts. More requirements produced more code, which looked like
progress until the comparison returned to what the original exercise had asked
for. That attempt was stopped, and its history was retained.

The restart put the student first. It followed the original chapters with the
architectural lessons supplied up front, without consulting old answers. After
implementation, an independent reviewer compared the result with the first
edition. Genuine bugs were corrected. When a grader depended on an unpublished
spelling or the wrong executable, the exception was documented and the failed
score remained a failed score. Only afterward did author and grader coder revise
the teaching together.

That order changed the editorial work. A chapter could now explain the difficult
step the student had actually encountered, rather than the difficulty its author
imagined. Architecture repairs that the new implementation no longer needed were
folded into the chapters where their lessons and capabilities belonged. Useful
stories and explanations remained. A smaller chapter was never the objective;
a chapter that teaches the right thing was.

The new implementation catches up through the first edition's final chapter,
Gateway, now Chapter 26. Chapter 24 adds an optional workaround for the observed
subscription-route caching problem. Chapter 25 is deliberately reserved for the
actual Codex-to-Ensemble crossover. It is a placeholder because that work is
underway, and an empty space is more useful than a success story written early.

The new Ensemble has already used GPT-6 Astra at high reasoning to read a small
program, reproduce its failures, repair it and pass its unchanged tests. That is
a concrete coding demonstration. Making Ensemble the everyday tool for the next
edition requires more than that bounded task. The chapters preserve both the
working result and the unfinished ambition.

## What to bring

You need Go, a text editor and an AI coding assistant. The exercises introduce
other tools as they become necessary, including the Go debugger. Some later
features have explicit platform limits; the sandbox chapter's implementation is
for macOS. Check a chapter's prerequisites before treating a different environment
as an implementation failure.

Keep three costs separate: running local checks, using real models through the
agent you build, and paying the assistant that helps you build it. Fake-provider
grading avoids inference charges for the first. The other two depend on your
provider, account, model and working habits. The book teaches usage accounting
because a correct answer does not tell you what the request cost.

For a human student, formulate the task yourself before giving it to your coding
assistant. Describing a system precisely enough that another implementer builds
the right thing is part of the exercise. An automated student run instead gives
the agent the chapter and tests whether the teaching transfers. Both are useful;
know which experiment you are doing.

The first edition was written by Bill Cox with CodeRhapsody. Codex is the primary
creator of this second edition, under Bill's direction, using student, author,
grader and review roles. Chapter 24 is the disclosed exception where the
coordinator also implemented and reviewed the work. The first edition remains
available as its own historical artifact, including the claims and failures that
led here.

Build an agent you can understand. Use it long enough to find what the tests
missed. Then leave the next builder the reason, as well as the fix.
