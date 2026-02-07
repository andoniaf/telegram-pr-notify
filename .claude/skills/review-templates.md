# Review Templates

Review and validate all Telegram notification templates for consistency, correctness, and UX.

## Instructions

1. Read all template constants in `pkg/templates/defaults.go`.
2. Read the `Render()` function and `selectDefault()` in `pkg/templates/render.go`.
3. Check each template for:
   - Consistent emoji usage (no duplicates across different event types)
   - Proper `html/template` syntax (not `text/template`)
   - Use of `truncate` function on user-generated content (titles, bodies)
   - Conditional guards (`{{- if .Field}}`) around optional fields
   - Correct `<a href>` links for PR, repo, and actor
   - `<blockquote>` for review/comment bodies
   - Branch lines (`Head.Ref` → `Base.Ref`) where applicable
4. Verify the `defaultTemplates` map has entries for all 11 event+action combinations.
5. Run `go test ./pkg/templates/ -v` to verify all template tests pass.
6. Report any inconsistencies or improvements to the user.
