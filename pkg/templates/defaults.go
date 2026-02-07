package templates

const prOpened = `🔀 <b>New Pull Request</b>
<a href="{{.PR.HTMLURL}}">#{{.PR.Number}}</a> {{truncate .PR.Title 100}}
{{.PR.Head.Ref}} → {{.PR.Base.Ref}}

by <a href="{{.Actor.HTMLURL}}">{{.Actor.Login}}</a> in <a href="{{.Repo.HTMLURL}}"><b>{{.Repo.FullName}}</b></a>`

const prClosed = `❌ <b>Pull Request Closed</b>
<a href="{{.PR.HTMLURL}}">#{{.PR.Number}}</a> {{truncate .PR.Title 100}}

by <a href="{{.Actor.HTMLURL}}">{{.Actor.Login}}</a> in <a href="{{.Repo.HTMLURL}}"><b>{{.Repo.FullName}}</b></a>`

const prMerged = `🟣 <b>Pull Request Merged</b>
<a href="{{.PR.HTMLURL}}">#{{.PR.Number}}</a> {{truncate .PR.Title 100}}
{{.PR.Head.Ref}} → {{.PR.Base.Ref}}

by <a href="{{.Actor.HTMLURL}}">{{.Actor.Login}}</a> in <a href="{{.Repo.HTMLURL}}"><b>{{.Repo.FullName}}</b></a>`

const prReopened = `🔃 <b>Pull Request Reopened</b>
<a href="{{.PR.HTMLURL}}">#{{.PR.Number}}</a> {{truncate .PR.Title 100}}
{{.PR.Head.Ref}} → {{.PR.Base.Ref}}

by <a href="{{.Actor.HTMLURL}}">{{.Actor.Login}}</a> in <a href="{{.Repo.HTMLURL}}"><b>{{.Repo.FullName}}</b></a>`

const prSynchronize = `🔄 <b>Pull Request Updated</b>
<a href="{{.PR.HTMLURL}}">#{{.PR.Number}}</a> {{truncate .PR.Title 100}}
{{.PR.Head.Ref}} → {{.PR.Base.Ref}}

New commits pushed by <a href="{{.Actor.HTMLURL}}">{{.Actor.Login}}</a> in <a href="{{.Repo.HTMLURL}}"><b>{{.Repo.FullName}}</b></a>`

const prReadyForReview = `👀 <b>Pull Request Ready for Review</b>
<a href="{{.PR.HTMLURL}}">#{{.PR.Number}}</a> {{truncate .PR.Title 100}}
{{.PR.Head.Ref}} → {{.PR.Base.Ref}}

by <a href="{{.Actor.HTMLURL}}">{{.Actor.Login}}</a> in <a href="{{.Repo.HTMLURL}}"><b>{{.Repo.FullName}}</b></a>`

const prConvertedToDraft = `📝 <b>Pull Request Converted to Draft</b>
<a href="{{.PR.HTMLURL}}">#{{.PR.Number}}</a> {{truncate .PR.Title 100}}

by <a href="{{.Actor.HTMLURL}}">{{.Actor.Login}}</a> in <a href="{{.Repo.HTMLURL}}"><b>{{.Repo.FullName}}</b></a>`

const reviewApproved = `✅ <b>Pull Request Approved</b>
<a href="{{.PR.HTMLURL}}">#{{.PR.Number}}</a> {{truncate .PR.Title 100}}

by <a href="{{.Actor.HTMLURL}}">{{.Actor.Login}}</a> in <a href="{{.Repo.HTMLURL}}"><b>{{.Repo.FullName}}</b></a>
{{- if .Review.Body}}

<blockquote>{{truncate .Review.Body 500}}</blockquote>
{{- end}}`

const reviewChangesRequested = `🔴 <b>Changes Requested</b>
<a href="{{.PR.HTMLURL}}">#{{.PR.Number}}</a> {{truncate .PR.Title 100}}

by <a href="{{.Actor.HTMLURL}}">{{.Actor.Login}}</a> in <a href="{{.Repo.HTMLURL}}"><b>{{.Repo.FullName}}</b></a>
{{- if .Review.Body}}

<blockquote>{{truncate .Review.Body 500}}</blockquote>
{{- end}}`

const reviewCommented = `💬 <b>Review Submitted</b>
<a href="{{.PR.HTMLURL}}">#{{.PR.Number}}</a> {{truncate .PR.Title 100}}

by <a href="{{.Actor.HTMLURL}}">{{.Actor.Login}}</a> in <a href="{{.Repo.HTMLURL}}"><b>{{.Repo.FullName}}</b></a>
{{- if .Review.Body}}

<blockquote>{{truncate .Review.Body 500}}</blockquote>
{{- end}}`

const reviewCommentCreated = `📝 <b>Inline Comment</b>
<a href="{{.PR.HTMLURL}}">#{{.PR.Number}}</a> {{truncate .PR.Title 100}}

by <a href="{{.Actor.HTMLURL}}">{{.Actor.Login}}</a> in <a href="{{.Repo.HTMLURL}}"><b>{{.Repo.FullName}}</b></a>
{{- if .Comment.Path}}
📄 <code>{{.Comment.Path}}</code>
{{- end}}
{{- if .Comment.Body}}

<blockquote>{{truncate .Comment.Body 500}}</blockquote>
{{- end}}`

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
