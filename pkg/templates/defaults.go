package templates

const prOpened = `🔀 <b>New Pull Request</b>
<a href="{{.PR.HTMLURL}}">#{{.PR.Number}}</a> {{.PR.Title}}
{{.PR.Head.Ref}} → {{.PR.Base.Ref}}

by <a href="{{.Actor.HTMLURL}}">{{.Actor.Login}}</a> in <b>{{.Repo.FullName}}</b>`

const prClosed = `❌ <b>Pull Request Closed</b>
<a href="{{.PR.HTMLURL}}">#{{.PR.Number}}</a> {{.PR.Title}}

by <a href="{{.Actor.HTMLURL}}">{{.Actor.Login}}</a> in <b>{{.Repo.FullName}}</b>`

const prMerged = `🟣 <b>Pull Request Merged</b>
<a href="{{.PR.HTMLURL}}">#{{.PR.Number}}</a> {{.PR.Title}}
{{.PR.Head.Ref}} → {{.PR.Base.Ref}}

by <a href="{{.Actor.HTMLURL}}">{{.Actor.Login}}</a> in <b>{{.Repo.FullName}}</b>`

const prReopened = `🔄 <b>Pull Request Reopened</b>
<a href="{{.PR.HTMLURL}}">#{{.PR.Number}}</a> {{.PR.Title}}

by <a href="{{.Actor.HTMLURL}}">{{.Actor.Login}}</a> in <b>{{.Repo.FullName}}</b>`

const prSynchronize = `🔄 <b>Pull Request Updated</b>
<a href="{{.PR.HTMLURL}}">#{{.PR.Number}}</a> {{.PR.Title}}

New commits pushed by <a href="{{.Actor.HTMLURL}}">{{.Actor.Login}}</a> in <b>{{.Repo.FullName}}</b>`

const prReadyForReview = `✅ <b>Pull Request Ready for Review</b>
<a href="{{.PR.HTMLURL}}">#{{.PR.Number}}</a> {{.PR.Title}}

by <a href="{{.Actor.HTMLURL}}">{{.Actor.Login}}</a> in <b>{{.Repo.FullName}}</b>`

const prConvertedToDraft = `📝 <b>Pull Request Converted to Draft</b>
<a href="{{.PR.HTMLURL}}">#{{.PR.Number}}</a> {{.PR.Title}}

by <a href="{{.Actor.HTMLURL}}">{{.Actor.Login}}</a> in <b>{{.Repo.FullName}}</b>`

const reviewApproved = `✅ <b>Pull Request Approved</b>
<a href="{{.PR.HTMLURL}}">#{{.PR.Number}}</a> {{.PR.Title}}

by <a href="{{.Actor.HTMLURL}}">{{.Actor.Login}}</a> in <b>{{.Repo.FullName}}</b>
{{- if .Review.Body}}

💬 {{truncate .Review.Body 500}}
{{- end}}`

const reviewChangesRequested = `🔴 <b>Changes Requested</b>
<a href="{{.PR.HTMLURL}}">#{{.PR.Number}}</a> {{.PR.Title}}

by <a href="{{.Actor.HTMLURL}}">{{.Actor.Login}}</a> in <b>{{.Repo.FullName}}</b>
{{- if .Review.Body}}

💬 {{truncate .Review.Body 500}}
{{- end}}`

const reviewCommented = `💬 <b>Review Comment</b>
<a href="{{.PR.HTMLURL}}">#{{.PR.Number}}</a> {{.PR.Title}}

by <a href="{{.Actor.HTMLURL}}">{{.Actor.Login}}</a> in <b>{{.Repo.FullName}}</b>
{{- if .Review.Body}}

💬 {{truncate .Review.Body 500}}
{{- end}}`

const reviewCommentCreated = `💬 <b>Review Comment</b>
<a href="{{.PR.HTMLURL}}">#{{.PR.Number}}</a> {{.PR.Title}}

by <a href="{{.Actor.HTMLURL}}">{{.Actor.Login}}</a> in <b>{{.Repo.FullName}}</b>
📄 {{.Comment.Path}}

💬 {{truncate .Comment.Body 500}}`

// defaultTemplates maps event_name + action to a default template string.
var defaultTemplates = map[string]string{
	"pull_request:opened":            prOpened,
	"pull_request:closed":            prClosed,
	"pull_request:merged":            prMerged,
	"pull_request:reopened":          prReopened,
	"pull_request:synchronize":       prSynchronize,
	"pull_request:ready_for_review":  prReadyForReview,
	"pull_request:converted_to_draft": prConvertedToDraft,

	"pull_request_review:approved":          reviewApproved,
	"pull_request_review:changes_requested": reviewChangesRequested,
	"pull_request_review:commented":         reviewCommented,

	"pull_request_review_comment:created": reviewCommentCreated,
}
