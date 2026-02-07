# Local Test

Run the action locally to send a real Telegram notification and verify the output.

## Instructions

1. Source credentials from `.envrc`:
   ```bash
   source .envrc
   ```
2. Run the action:
   ```bash
   go run .
   ```
3. Verify the Telegram message was sent successfully (output: "Notification sent successfully").
4. If the user wants to test a specific event type, set `INPUT_EVENT_PAYLOAD` to the content of the appropriate testdata file before running:
   ```bash
   export INPUT_EVENT_PAYLOAD="$(cat testdata/<fixture>.json)"
   go run .
   ```
5. Run the full test suite with race detection:
   ```bash
   go test ./... -v -race
   ```
6. Report results to the user.
