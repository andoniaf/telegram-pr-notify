# Add Event Type

Add support for a new GitHub webhook event type to the notification system.

## Instructions

1. Read the user's request to understand which GitHub event/action to add.
2. Read the existing code to understand current patterns:
   - `pkg/events/events.go` for parsing structs and the `Parse()` switch
   - `pkg/templates/defaults.go` for template constants and the `defaultTemplates` map
   - `pkg/templates/render.go` for template selection logic
3. Create a sample JSON fixture in `testdata/` based on GitHub's webhook documentation. Use sanitized/fake data.
4. Add the event struct and parse function in `pkg/events/events.go` following the existing pattern.
5. Add a template constant in `pkg/templates/defaults.go`:
   - Use `html/template` syntax with the `truncate` function for user-generated content
   - Follow the existing emoji + HTML formatting style
   - Register in the `defaultTemplates` map with `event_name:action` key
6. Add test cases in both `pkg/events/events_test.go` and `pkg/templates/render_test.go`.
7. Run `make test` and `make lint` to verify everything passes.
8. Update README.md's supported events table.
