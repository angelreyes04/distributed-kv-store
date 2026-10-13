# Distributed Key-Value Store in Go — Personal Learning Roadmap

**Goal:** Build real design judgment and distributed-systems skill by building a key-value
store from scratch in Go, evolving it over months into a small replicated, fault-tolerant
database.

**Ground rule:** You are writing the code. Use outside help (AI, docs, books, other people's
source) for *questions, reviews, and comparison after you've attempted something* — never as
a substitute for writing the design or the implementation yourself. If you ever find yourself
about to paste in a generated solution instead of writing your own attempt first, stop —
that's the exact experience this project exists to give you.

---

## How to use this roadmap

- Each stage lists **learning objectives**, a **milestone checklist**, and **rough time**.
  Time estimates are generous on purpose — the point is depth, not speed.
- Don't skip ahead. The design pain in each stage is what produces the lesson in the next one.
- Re-read your own code from two stages ago before starting a new one. You'll usually wince —
  that's growth, and it's worth noting *why* in your decision log.

---

## Stage 0 — Setup & Foundations (1–2 weeks)

**Objective:** Environment, project shape, and Go fundamentals in place before writing store logic.

- [ ] If Go is new to you: work through *A Tour of Go* and *Effective Go* first — you want
      goroutines, channels, and interfaces to feel natural before Stage 1.
- [ ] Initialize a Git repo and push it public from commit one.
- [ ] Lay out an idiomatic Go project structure, e.g.:
  ```
  /cmd/server      — server entrypoint
  /cmd/client      — CLI client entrypoint
  /internal/store  — core key-value logic
  /internal/proto  — wire protocol / codec
  /docs/adr        — decision log
  ```
- [ ] Set up `go test` + `testify` for assertions.
- [ ] Create `docs/adr/0001-project-goals.md` — write one paragraph on what you're
      trying to learn from this project. You'll want this later.

---

## Stage 1 — Single-Node Store + Protocol (3–4 weeks)

**Learning objectives:** goroutines & channels, TCP networking, the discipline of building
the smallest possible working slice, Command pattern (earned, not assumed).

- [ ] Build a raw TCP server (`net` package) that accepts connections, one goroutine per
      connection.
- [ ] Back it with a plain `map[string]string` guarded by `sync.RWMutex`.
- [ ] Get exactly **one** command working end-to-end: `SET key value` → `OK`.
- [ ] Add a test that starts the server, dials a real client, and asserts the behavior.
- [ ] Add `GET` and `DELETE`. Resist adding a `Command` interface until parsing genuinely
      hurts — then introduce it, and write the ADR explaining why.
- [ ] Write a simple CLI client so you can interact with your own server by hand.

**Stage exit criteria:** you can run your server, connect with your own client, and do
SET/GET/DELETE reliably, with tests covering all three.

---

## Stage 2 — Durability (3–4 weeks)

**Learning objectives:** write-ahead logging, crash recovery, serialization tradeoffs.

- [ ] Append every write to a log file before applying it in memory (write-ahead log).
- [ ] On startup, replay the log to rebuild state.
- [ ] Add periodic snapshotting so the log doesn't grow forever, plus log truncation
      after a snapshot.
- [ ] Deliberately kill the process mid-write (`kill -9`) and confirm recovery is correct —
      this is where you'll find your real bugs.
- [ ] ADR: why you chose your serialization format (`encoding/gob` is a fine starting point;
      note what you'd reconsider if you needed cross-language compatibility).

**Stage exit criteria:** you can kill the process at any point and restart with no data loss
beyond the last unflushed write.

---

## Stage 3 — Client Library & Query Layer (2–3 weeks)

**Learning objectives:** API design for other developers, parsing, gRPC.

- [ ] Extract a real client library (separate Go module or package) with a clean public API —
      pretend someone else has to use it without reading your server code.
- [ ] Write a small recursive-descent parser for a tiny query language on top of raw commands
      (even something minimal like `SET key value IF NOT EXISTS`).
- [ ] Once you've felt what your hand-rolled protocol is missing, rebuild the wire protocol
      on gRPC + protobuf and compare the experience in your ADR log.

**Stage exit criteria:** a friend (or future-you) could use your client library from its
public API alone, without reading the server internals.

---

## Stage 4 — Replication (4–6 weeks)

**Learning objectives:** leader/follower design, Observer and State patterns, reasoning about
concurrent, networked failure.

- [ ] Run multiple instances of your server; designate one as leader.
- [ ] Leader forwards every write to followers over a persistent connection
      (channels are a natural fit here).
- [ ] Followers apply writes in order and reject direct client writes.
- [ ] Kill a follower mid-stream and confirm it catches up correctly on reconnect.
- [ ] ADR: what happens today if the leader dies? (You'll fix this in Stage 5 —
      write down what breaks first.)

**Stage exit criteria:** writes to the leader are visible on all followers, and a follower
that drops and reconnects catches back up without manual intervention.

---

## Stage 5 — Consensus & Fault Tolerance (6–10 weeks)

**Learning objectives:** the hardest and most valuable stage — real distributed state
machine design, leader election, split-brain avoidance.

- [ ] Implement a simplified Raft: leader election, log replication, and commit indexing.
      Build your own before looking at any reference implementation.
- [ ] Test it under real failure: kill the leader, kill a minority of followers, partition
      the network (even crudely, by dropping connections) — confirm the cluster still
      converges on a correct, consistent state.
- [ ] Once your version works, read HashiCorp's `raft` library (used in Consul) and compare
      your design decisions against theirs in an ADR — where did you agree, where did you
      diverge, and why?

**Stage exit criteria:** your cluster survives a leader crash and elects a new one
automatically, with no data loss for committed writes.

---

## Stage 6 — Sharding & Scale-Out (3–5 weeks)

**Learning objectives:** horizontal scaling tradeoffs, routing design.

- [ ] Partition keys across multiple independent replica groups (each running its own
      Stage 4/5 replication set) using consistent hashing.
- [ ] Build a routing layer that directs client requests to the correct shard.
- [ ] Handle the "moving a shard" problem at least conceptually — document the tradeoffs
      even if you don't fully implement live resharding.

**Stage exit criteria:** you can add a new shard, and known keys route to the right place
without touching client code.

---

## Practices to run throughout (not a stage — ongoing)

- **ADR log:** one entry per nontrivial decision, however small it feels at the time.
- **Tests before features:** every stage should be covered by tests before you consider it done.
- **Weekly cadence:** aim for one meaningful commit most weeks rather than a rigid schedule —
  the project is more likely to die from lost momentum than from difficulty.
- **Resist premature abstraction:** if you're building an interface "in case you need it
  later," stop and wait until you actually need it.
- **Write about the hard parts** as you go (even just in the repo's README or a blog) —
  this is what makes the project legible to others later, in interviews or promotion packets.

---

## Reference material (use *after* attempting each stage yourself)

- **"Designing Data-Intensive Applications" by Martin Kleppmann** — maps almost one-to-one
  onto these stages; read each relevant chapter after you've built the stage, not before.
- **MIT 6.824 (Distributed Systems)** — the public labs walk through building Raft in Go;
  great for comparing your Stage 5 design against a rigorous reference after your own attempt.
- **HashiCorp's `raft` library** — a real production Go Raft implementation, useful once
  you've built your own.
- **CodeCrafters' "Build Your Own Redis" track** — good for Stage 1 inspiration on protocol
  design specifically.
- **golang-standards/project-layout** (GitHub) — a reference for idiomatic Go repo structure.

---

## Using AI along the way, without it doing the work

Good uses: asking *why* a pattern applies once you've hit the pain point yourself, getting a
second opinion on an ADR you've already drafted, debugging a specific error message, or
discussing tradeoffs between two designs you've already sketched.

Avoid: asking for the implementation of a stage before you've attempted it, or for "the
correct architecture" up front — the entire value of this roadmap is in feeling the design
pressure before the solution appears.
