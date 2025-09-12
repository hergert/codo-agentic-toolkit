# Refactor Plan: Command Refinements for Codo Toolkit

## Goal
Standardize and optimize Claude commands following Anthropic's agentic coding guidelines to reduce token usage, improve determinism, and enforce minimal viable changes.

## Trade-offs Analysis

### Benefits
- **Token efficiency**: Hard caps (300-350 tokens) prevent context bloat
- **Predictability**: Explicit argument-hint and minimal allowed-tools
- **Maintainability**: Consistent command structure across all commands
- **Safety**: Fail-fast on unrelated failures prevents scope creep

### Risks
- **Breaking changes**: Commands that depend on current structure may fail
- **User disruption**: Existing workflows need documentation updates
- **Hook dependencies**: Some formatting logic moves to hooks (requires testing)

## Phases

### Phase 1: Command Structure Standardization (Low Risk)
**Files to modify:**
- `pack/dotclaude/commands/map.md`
- `pack/dotclaude/commands/plan.md`
- `pack/dotclaude/commands/prime-context.md`
- `pack/dotclaude/commands/implement.md`
- `pack/dotclaude/commands/tests.md`
- `pack/dotclaude/commands/review.md`

**Changes:**
- Add `argument-hint` metadata to all commands
- Restrict `allowed-tools` to minimal set per command
- Add token limits (≤300-350) to command outputs
- Enforce Read → Process → Report flow

**Checkpoint 1:** Test each command individually with `codo init --no-tui`

### Phase 2: Settings Updates (Low Risk)
**Files to modify:**
- `pack/dotclaude/settings.json`
- `pack/dotclaude/settings.local.json` (if exists)

**Changes:**
- Update notes field to clarify ALLOW_EDITS/FAST_MODE gates
- Remove invalid patterns from settings.local.json
- Align with Claude Code's settings precedence

**Checkpoint 2:** Verify settings load correctly

### Phase 3: Hook Integration (Medium Risk)
**Files to verify (no changes needed):**
- `pack/dotclaude/hooks/pre_tool_use.py`
- `pack/dotclaude/hooks/post_tool_use.py`
- `pack/dotclaude/hooks/user_prompt_submit.py`

**Validation:**
- Ensure hooks can handle formatting/shortstats that commands delegate
- Test determinism flags work as expected

**Checkpoint 3:** Run integration test with hooks enabled

### Phase 4: Doctor & Status Enhancements (Low Risk)
**Files to modify:**
- `cli/cmd/doctor.go`
- `cli/cmd/status.go`

**Changes:**
- Add Python version check to doctor
- Add `--strict` flag to status for CI gates

**Checkpoint 4:** Test CLI commands work correctly

## Dependency Graph
```
Phase 1 (Commands) → Phase 3 (Hook Validation)
        ↓
Phase 2 (Settings) → Phase 4 (CLI Enhancements)
```

## Test Strategy
1. **Unit tests**: Each command tested in isolation with mock inputs
2. **Integration tests**: Full `codo init` → command execution → validation
3. **Regression tests**: Ensure existing projects continue working
4. **Rollback tests**: Verify each phase can be reverted independently

## Rollback Strategy

### Per-Phase Rollback
- **Phase 1**: Keep backup of original commands in `.codo-backup/`
- **Phase 2**: Settings changes are additive; revert via git
- **Phase 3**: No changes, only validation
- **Phase 4**: CLI changes behind feature flags initially

### Emergency Rollback
```bash
# Full rollback to previous pack version
git checkout HEAD~1 -- pack/
codo update --to stable
```

## Acceptance Criteria
- [ ] All 6 commands have argument-hint and minimal allowed-tools
- [ ] Token limits enforced (≤350 tokens per command output)
- [ ] Commands follow Read → Process → Report flow
- [ ] Settings reflect correct gate precedence
- [ ] Doctor shows Python version
- [ ] Status --strict works for CI
- [ ] No regression in existing workflows

## Out of Scope
- Creating new commands
- Modifying hook implementations
- Changing CLI core functionality
- Altering manifest structure

## Implementation Order
1. Start with Phase 1 (commands) - highest impact, lowest risk
2. Phase 2 (settings) in parallel
3. Phase 3 (validation) after Phase 1 complete
4. Phase 4 (CLI) can be done independently

## Risk Mitigation
- Test each command change in `/tmp/test-codo` before committing
- Keep `.codo-backup/` of original files
- Document changes in CHANGELOG
- Provide migration guide for users

---
**Status**: Ready for approval
**Estimated effort**: 2-3 hours
**Risk level**: Low to Medium