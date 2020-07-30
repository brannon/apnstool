// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package serve

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

var layoutCache *TemplateCache
var viewCache *TemplateCache

func init() {
	layoutCache = new(TemplateCache)
	viewCache = new(TemplateCache)

	layoutCache.RegisterFunc("render_view", RenderView)
	layoutCache.RegisterFunc("has_prefix", strings.HasPrefix)

	layoutCache.MustParse("main", mainLayoutHtml)

	viewCache.MustParse("index", indexViewHtml)
	viewCache.MustParse("error", errorViewHtml)
	viewCache.MustParse("faq", faqViewHtml)
	viewCache.MustParse("send_alert", sendAlertViewHtml)
	viewCache.MustParse("send_background", sendBackgroundViewHtml)
	viewCache.MustParse("send_raw", sendRawViewHtml)
	viewCache.MustParse("send_result", sendResultViewHtml)
}

func LookupLayoutTemplate(name string) (*template.Template, error) {
	tmpl := layoutCache.Lookup(name)
	if tmpl == nil {
		return nil, fmt.Errorf("unable to locate layout template: %s", name)
	}
	return tmpl, nil
}

func RenderView(viewName string, context interface{}) (string, error) {
	viewTmpl := viewCache.Lookup(viewName)
	if viewTmpl == nil {
		return "", fmt.Errorf("unable to locate view template: %s", viewName)
	}

	buf := new(bytes.Buffer)
	err := viewTmpl.Execute(buf, context)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

const mainLayoutHtml string = `
<!doctype HTML>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
	
	<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.0/css/bootstrap.min.css" integrity="sha384-9aIt2nRpC12Uk9gS9baDl411NQApFmC26EwAOH8WgZl5MYYxFfc+NcPb1dKGj7Sk" crossorigin="anonymous">

	<style>
		html {
			position: relative;
			min-height: 100%;
		}
		body {
			margin-bottom: 60px;
		}
		section {
			margin-bottom: 60px;
		}
		.footer {
			position: absolute;
			bottom: 0;
			width: 100%;
			height: 60px;
			line-height: 60px;
			vertical-align: middle;
			/*background-color: #f3f3f3;*/
		}

		code {
			background-color: #f3f3f3;
			/*color: #5a5a5a;*/
			padding: 4px;
			border-radius: 5px;
		}
		pre > code {
			display: block;
			padding: 5px;
			width: 100%
		}

		.navbar-brand {
			font-size: 2em;
			padding-top: 0;
		}
		.form-text kbd {
			background-color: #888;
		}
		.build-info {
			color: #ddd;
		}
		dl > dt {
			margin-top: 30px;
		}
		dl > dd {
			margin-left: 15px;
			padding-top: 15px
		}

		body > .container {
			padding: 100px 15px 0;
		}
		.footer > .container {
			padding-right: 15px;
			padding-left 15px;
		}
	</style>

	<script src="https://code.jquery.com/jquery-3.5.1.slim.min.js" integrity="sha384-DfXdz2htPH0lsSSs5nCTpuj/zy4C+OGpamoFVy38MVBnE+IbbVYUew+OrCXaRkfj" crossorigin="anonymous"></script>
	<script src="https://cdn.jsdelivr.net/npm/popper.js@1.16.0/dist/umd/popper.min.js" integrity="sha384-Q6E9RHvbIyZFJoft+2mJbHaEWldlvI9IOYy5n3zV9zzTtmI3UksdQRVvoxMfooAo" crossorigin="anonymous"></script>
	<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.5.0/js/bootstrap.min.js" integrity="sha384-OgVRvuATP1z7JjHLkuOU7Xw704+h835Lr+6QL9UvYjZE3Ipu6Tp75j7Bh/kR0JKI" crossorigin="anonymous"></script>

	<script src="https://kit.fontawesome.com/1a3d1ce67d.js" crossorigin="anonymous"></script>

	<title>apnstool</title>
</head>

<body>
	<header>
		<nav class="navbar navbar-expand-md fixed-top navbar-dark bg-primary">
			<a class="navbar-brand" href="/">apnstool</a>
			<button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbar-collapse" aria-controls="navbar-collapse" aria-expanded="false" aria-label="Toggle navigation">
				<span class="navbar-toggler-icon"></span>
			</button>
			<div class="collapse navbar-collapse" id="navbar-collapse">
				<ul class="navbar-nav mr-auto">
					<li class="nav-item {{ if eq .Request.Path "/" }}active{{ end }}">
						<a class="nav-link" href="/">Home</a>
					</li>
					<li class="nav-item {{ if eq .Request.Path "/faq" }}active{{ end }}">
						<a class="nav-link" href="/faq">FAQ</a>
					</li>
					<li class="nav-item dropdown {{ if has_prefix .Request.Path "/send/" }}active{{ end }}">
						<a class="nav-link dropdown-toggle" href="#" id="sendDropdown" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
							Send
						</a>
						<div class="dropdown-menu" aria-labelledby="sendDropdown">
							<a class="dropdown-item" href="/send/alert">Alert</a>
							<a class="dropdown-item" href="/send/background">Background</a>
						</div>
					</li>
				</ul>
				<a href="https://github.com/brannon/apnstool"><i class="fab fa-github fa-2x text-white"></i></a>
			</div>
		</nav>
	</header>

	<main role="main" class="container m4">
		{{ render_view .ViewName .ViewContext }}
	</main>

	<footer class="footer">
		<div class="container-fluid">
			<div class="row align-bottom">
				<div class="col-sm-3 text-left">
				</div>
				<div class="col-sm-6 text-center">
				</div>
				<div class="col-sm-3 text-right">
					<small class="build-info">build {{ .Build.BuildDate }}-{{ .Build.CommitHash }}</small>
				</div>
			</div>
		</div>
	</footer>
</body>
</html>
`

const errorViewHtml string = `
<div class="alert alert-danger" role="alert">
	<h4 class="alert-heading">Error sending notification!</h4>
	<hr>
	<div>{{ . }}</div>
</div>
`

const indexViewHtml string = `
<section>
	<div class="jumbotron">
		<h1 class="display-4">apnstool</h1>
		<p>This site is a utility for interacting with APNs (the Apple Push Notification service). It makes it simple to test your APNs credentials and different notification types.</p>
		<hr class="display-4">
		<p>Choose the type of push notification to send:</p>
		<a class="btn btn-primary btn-lg" href="/send/alert" role="button">Alert</a>
		<a class="btn btn-primary btn-lg" href="/send/background" role="button">Background</a> 
	</div>
	<div class="alert alert-warning">
		<p>This tool is not supported for production use!</p>
		<p class="mb-0">See the <a class="alert-link" href="/faq">FAQ</a> for more information.</p>
	</div>
</section>
`

const faqViewHtml string = `
<section>
	<h1 class="display-4">FAQ</h1>

	<dl>
		<dt>Can I use this to send notifications to my app?</dt>
		<dd>
			<p>You may use this tool for testing and troubleshooting notifications.</p>
			<p><strong class="text-danger text-uppercase">This tool is not supported for production use</strong></p>
			<p>To prevent abuse, request rates will be limited and requests may be blocked.</p>
		</dd>

		<dt>What should I use for sending notifications to my app?</dt>
		<dd>
			<p>There are many production-ready options for interacting with APNs that offer cross-platform notifications, advanced device management, and multi-targeting.</p>
			<p>I recommend <a href="https://azure.microsoft.com/en-us/services/notification-hubs/">Azure Notification Hubs</a></p>
		</dd>

		<dt>Are you logging my requests?</dt>
		<dd>
			<p>Yes, when using the web-based version of the tool, requests are logged to track usage and help prevent abuse.</p>
			<p>A one-way hash of the App ID and the Device Token are logged. Any abuse of this tool may result in your App ID being blocked.</p>
		</dd>

		<dt>Are my APNs credentials safe?</dt>
		<dd>
			<p>Yes, credentials are sent over HTTPS and used <em>only</em> to connect to APNs when sending the test notification.</p>
			<p>Credentials are not stored and not logged. Feel free to <a href="https://github.com/brannon/apnstool">review the source code</a>.</p>
		</dd>
		
		<dt>What if I don't feel comfortable sending my APNs credentials?</dt>
		<dd>
			<p>I don't blame you, that's one of the reasons why I created this tool. You can also download the tool and run it locally.</p>
		</dd>

		<dt>Is this site affiliated with Microsoft?</dt>
		<dd>
			<p>No. This site was created by <a href="https://github.com/brannon">@brannon</a> as a personal project and is hosted in a personal Azure account.</p>
			<p>Feel free to use this site for testing, but understand that the site is provided "as is", without warranty of any kind, express or implied.</p>
		</dd>

		<dt>What if I get an error, or something isn't working as expected?</dt>
		<dd>
			<p>Feel free to report problems by <a href="https://github.com/brannon/apnstool/issues">filing an issue</a> in the GitHub project, but understand that there is no official support for this site.</p>
		</dd>
	</dl>
</section>
`

const sendAlertViewHtml string = `
<h1 class="display-4">Send Alert</h1>
<p>
Fill out the fields below to send an APNS alert notification. Alert notifications have the form:</p>
<pre><code>{
  "aps": {
    "alert": "Alert text",
    "badge": 0,
    "sound": "default"
  }
}
</code></pre>
<p>
Alert notifications also set the <code>apns-push-type</code> header to <code>alert</code>.
</p>
<div class="pt-4 pb-4">
<form id="send-alert" action="/send/alert" method="post" enctype="multipart/form-data">
	<div class="form-group row">
		<label for="app-id" class="col-sm-2 col-form-label">App ID</label>
		<div class="col-sm-10">
			<input type="text" class="form-control" name="app-id" id="app-id" required>
			<small class="form-text text-muted">Same as the app's bundle ID. Sets the <code>apns-topic</code> header.</small>
		</div>
	</div>

	<fieldset class="form-group">
		<div class="row">
			<legend class="col-form-label col-sm-2 pt-0">Environment</legend>
			<div class="col-sm-10">
				<div class="form-check">
					<input type="checkbox" class="form-check-input" name="sandbox" id="sandbox" checked>
					<label for="sandbox" class="form-check-label">Use Sandbox</label>
				</div>
				<div class="form-check">
					<small class="form-check-label form-text text-muted">Check this when the <code>aps-environment</code> entitlement is equal to <code>development</code></small>
				</div>
			</div>
		</div>
	</fieldset>


	<fieldset class="form-group">
		<div class="row">
			<legend class="col-form-label col-sm-2 pt-0">Auth Type</legend>
			<div class="col-sm-10">
				<div class="form-check">
					<input type="radio" class="form-check-input auth-type-radio" name="auth-type" id="auth-type-token" value="token">
					<label for="auth-type-token" class="form-check-label">Token</label>
				</div>
				<div class="form-check">
					<input type="radio" class="form-check-input auth-type-radio" name="auth-type" id="auth-type-cert" value="cert">
					<label for="auth-type-cert" class="form-check-label">Certificate</label>
				</div>
			</div>
		</div>
	</fieldset>

	<div class="d-none" id="auth-type-token-container">
		<div class="form-group row">
			<label for="key-file" class="col-sm-2 col-form-label">Key File</label>
			<div class="col-sm-10">
				<input type="file" class="form-control-file" name="key-file" id="key-file">
				<small class="form-text text-muted">Choose the <code>.p8</code> file downloaded from your developer account.</small>
			</div>
		</div>
		<div class="form-group row">
			<label for="key-id" class="col-sm-2 col-form-label">Key ID</label>
			<div class="col-sm-10">
				<input type="text" class="form-control" name="key-id" id="key-id">				
			</div>
		</div>
		<div class="form-group row">
			<label for="team-id" class="col-sm-2 col-form-label">Team ID</label>
			<div class="col-sm-10">
				<input type="text" class="form-control" name="team-id" id="team-id">				
			</div>
		</div>
	</div>

	<div class="d-none" id="auth-type-cert-container">
		<div class="form-group row">
			<label for="cert-file" class="col-sm-2 col-form-label">Certificate File</label>
			<div class="col-sm-10">
				<input type="file" class="form-control-file" name="cert-file" id="cert-file">
				<small class="form-text text-muted">Choose the <code>.p12</code> (or <code>.pfx</code>) file downloaded from your developer account.</small>
			</div>
		</div>
		<div class="form-group row">
			<label for="cert-password" class="col-sm-2 col-form-label">Certificate Password</label>
			<div class="col-sm-10">
				<input type="password" class="form-control" name="cert-password" id="cert-password" placeholder="optional">
			</div>
		</div>			
	</div>

	<div class="form-group row">
		<label for="device-token" class="col-sm-2 col-form-label">Device Token</label>
		<div class="col-sm-10">
			<input type="text" class="form-control" name="device-token" id="device-token" required>
		</div>
	</div>

	<div class="form-group row">
		<label for="device-token" class="col-sm-2 col-form-label">Expiration</label>
		<div class="col-sm-10">
			<input type="number" class="form-control" name="expiration" id="expiration">
			<small class="form-text text-muted">Specifies the time when the notification expires (as a unix timestamp). The value <code>0</code> will cause APNS to attempt immediate delivery, without any retries.</small>
		</div>
	</div>
	<div class="form-group row">
		<label for="device-token" class="col-sm-2 col-form-label">Priority</label>
		<div class="col-sm-10">
			<input type="number" class="form-control" name="priority" id="priority">
			<small class="form-text text-muted">Specifies the priority of the notification. Valid values are <code>5</code> and <code>10</code>.</small>
		</div>
	</div>

	<div class="form-group row">
		<label for="alert-text" class="col-sm-2 col-form-label">Alert Text</label>
		<div class="col-sm-10">
			<input type="text" class="form-control" name="alert-text" id="alert-text">
		</div>
	</div>
	<div class="form-group row">
		<label for="badge-count" class="col-sm-2 col-form-label">Badge Count</label>
		<div class="col-sm-10">
			<input type="number" class="form-control" name="badge-count" id="badge-count" min="0">
			<small class="form-text text-muted">Use the value <code>0</code> to remove the badge from the app icon.</small>
		</div>
	</div>
	<div class="form-group row">
		<label for="sound-name" class="col-sm-2 col-form-label">Sound Name</label>
		<div class="col-sm-10">
			<input type="text" class="form-control" name="sound-name" id="sound-name">
			<small class="form-text text-muted">Use the value <code>default</code> to make the device play a sound.</small>
		</div>
	</div>

	<div>
		<button type="submit" class="btn btn-primary mt-3">Send</button>
	</div>
</form>
</div>

<script type="text/javascript">
	function handleAuthTypeSelected() {
		var authType = $("input[name=auth-type]:checked").val();

		switch (authType) {
		case "token":
			$("#auth-type-token-container").removeClass("d-none");
			$("#auth-type-cert-container").addClass("d-none");
			break;
		case "cert":
			$("#auth-type-token-container").addClass("d-none");
			$("#auth-type-cert-container").removeClass("d-none");
			break;
		}
	}

	$(function() {
		$("input[name=auth-type]:radio").change(function() {
			handleAuthTypeSelected();
		});

		handleAuthTypeSelected();
	})
</script>
`

const sendBackgroundViewHtml string = `
<h1 class="display-4">Send Background</h1>
<p>
Fill out the fields below to send an APNS background notification. Background notifications have the form:</p>
<pre><code>{
  "aps": {
    "content-available": 1
  },
  "custom-data": "value"
}
</code></pre>
<p>
Background notifications also set the <code>apns-push-type</code> header to <code>background</code>.
</p>
<div class="pt-4 pb-4">
<form id="send-background" action="/send/background" method="post" enctype="multipart/form-data">
	<div class="form-group row">
		<label for="app-id" class="col-sm-2 col-form-label">App ID</label>
		<div class="col-sm-10">
			<input type="text" class="form-control" name="app-id" id="app-id" required>
			<small class="form-text text-muted">Same as the app's bundle ID. Sets the <code>apns-topic</code> header.</small>
		</div>
	</div>

	<fieldset class="form-group">
		<div class="row">
			<legend class="col-form-label col-sm-2 pt-0">Environment</legend>
			<div class="col-sm-10">
				<div class="form-check">
					<input type="checkbox" class="form-check-input" name="sandbox" id="sandbox" checked>
					<label for="sandbox" class="form-check-label">Use Sandbox</label>
				</div>
				<div class="form-check">
					<small class="form-check-label form-text text-muted">Check this when the <code>aps-environment</code> entitlement is equal to <code>development</code></small>
				</div>
			</div>
		</div>
	</fieldset>


	<fieldset class="form-group">
		<div class="row">
			<legend class="col-form-label col-sm-2 pt-0">Auth Type</legend>
			<div class="col-sm-10">
				<div class="form-check">
					<input type="radio" class="form-check-input auth-type-radio" name="auth-type" id="auth-type-token" value="token">
					<label for="auth-type-token" class="form-check-label">Token</label>
				</div>
				<div class="form-check">
					<input type="radio" class="form-check-input auth-type-radio" name="auth-type" id="auth-type-cert" value="cert">
					<label for="auth-type-cert" class="form-check-label">Certificate</label>
				</div>
			</div>
		</div>
	</fieldset>

	<div class="d-none" id="auth-type-token-container">
		<div class="form-group row">
			<label for="key-file" class="col-sm-2 col-form-label">Key File</label>
			<div class="col-sm-10">
				<input type="file" class="form-control-file" name="key-file" id="key-file">
				<small class="form-text text-muted">Choose the <code>.p8</code> file downloaded from your developer account.</small>
			</div>
		</div>
		<div class="form-group row">
			<label for="key-id" class="col-sm-2 col-form-label">Key ID</label>
			<div class="col-sm-10">
				<input type="text" class="form-control" name="key-id" id="key-id">				
			</div>
		</div>
		<div class="form-group row">
			<label for="team-id" class="col-sm-2 col-form-label">Team ID</label>
			<div class="col-sm-10">
				<input type="text" class="form-control" name="team-id" id="team-id">				
			</div>
		</div>
	</div>

	<div class="d-none" id="auth-type-cert-container">
		<div class="form-group row">
			<label for="cert-file" class="col-sm-2 col-form-label">Certificate File</label>
			<div class="col-sm-10">
				<input type="file" class="form-control-file" name="cert-file" id="cert-file">
				<small class="form-text text-muted">Choose the <code>.p12</code> (or <code>.pfx</code>) file downloaded from your developer account.</small>
			</div>
		</div>
		<div class="form-group row">
			<label for="cert-password" class="col-sm-2 col-form-label">Certificate Password</label>
			<div class="col-sm-10">
				<input type="password" class="form-control" name="cert-password" id="cert-password" placeholder="optional">
			</div>
		</div>			
	</div>

	<div class="form-group row">
		<label for="device-token" class="col-sm-2 col-form-label">Device Token</label>
		<div class="col-sm-10">
			<input type="text" class="form-control" name="device-token" id="device-token" required>
		</div>
	</div>

	<div class="form-group row">
		<label for="data" class="col-sm-2 col-form-label">Custom Data</label>
		<div class="col-sm-10">
			<textarea class="form-control" name="data" id="data"></textarea>
		</div>
	</div>

	<div>
		<button type="submit" class="btn btn-primary mt-3">Send</button>
	</div>
</form>
</div>

<script type="text/javascript">
	function handleAuthTypeSelected() {
		var authType = $("input[name=auth-type]:checked").val();

		switch (authType) {
		case "token":
			$("#auth-type-token-container").removeClass("d-none");
			$("#auth-type-cert-container").addClass("d-none");
			break;
		case "cert":
			$("#auth-type-token-container").addClass("d-none");
			$("#auth-type-cert-container").removeClass("d-none");
			break;
		}
	}

	$(function() {
		$("input[name=auth-type]:radio").change(function() {
			handleAuthTypeSelected();
		});

		handleAuthTypeSelected();
	})
</script>
`

const sendRawViewHtml string = ``

const sendResultViewHtml string = `
<div class="alert alert-success" role="alert">
	<h4 class="alert-heading">Notification sent!</h4>
	<hr>
	<div>APNS-ID: {{ .ApnsId }}</div>
</div>
`
