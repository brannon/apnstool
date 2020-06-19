// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package serve

import (
	"bytes"
	"fmt"
	"text/template"
)

var layoutCache *TemplateCache
var viewCache *TemplateCache

func init() {
	layoutCache = new(TemplateCache)
	viewCache = new(TemplateCache)

	layoutCache.RegisterFunc("render_view", RenderView)

	layoutCache.MustParse("main", mainLayoutHtml)

	viewCache.MustParse("index", indexViewHtml)
	viewCache.MustParse("error", errorViewHtml)
	viewCache.MustParse("send_alert", sendAlertViewHtml)
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
		.form-text kbd {
			background-color: #888;
		}
	</style>

	<script src="https://code.jquery.com/jquery-3.5.1.slim.min.js" integrity="sha384-DfXdz2htPH0lsSSs5nCTpuj/zy4C+OGpamoFVy38MVBnE+IbbVYUew+OrCXaRkfj" crossorigin="anonymous"></script>
	<script src="https://cdn.jsdelivr.net/npm/popper.js@1.16.0/dist/umd/popper.min.js" integrity="sha384-Q6E9RHvbIyZFJoft+2mJbHaEWldlvI9IOYy5n3zV9zzTtmI3UksdQRVvoxMfooAo" crossorigin="anonymous"></script>
	<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.5.0/js/bootstrap.min.js" integrity="sha384-OgVRvuATP1z7JjHLkuOU7Xw704+h835Lr+6QL9UvYjZE3Ipu6Tp75j7Bh/kR0JKI" crossorigin="anonymous"></script>

	<title>APNS Tool</title>
</head>

<body>
	<div class="container mt-4 mb-4 pl-4 pr-4">
		<div class="text-center mb-4">
			<h1>APNS Tool</h1>
		</div>
		{{ render_view .ViewName .ViewContext }}
	</div>
</body>
</html>
`

const errorViewHtml string = `
<div class="alert alert-danger" role="alert">
	<h4 class="alert-heading">Error sending notification!</h4>
	<div>{{ . }}</div>
</div>
`

const indexViewHtml string = `
<p>Index View</p>
`

const sendAlertViewHtml string = `
<form id="send-alert" action="/send/alert" method="post" enctype="multipart/form-data">
	<div class="form-group row">
		<label for="app-id" class="col-sm-2 col-form-label">App ID</label>
		<div class="col-sm-10">
			<input type="text" class="form-control" name="app-id" id="app-id" required>
			<small class="form-text text-muted">Same as the app's bundle ID. Sets the <kbd>apns-topic</kbd> header.</small>
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
					<small class="form-check-label form-text text-muted">Check this when the <kbd>aps-environment</kbd> entitlement is equal to <kbd>development</kbd></small>
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
				<small class="form-text text-muted">Choose the .p8 file downloaded from your developer account.</small>
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
				<small class="form-text text-muted">Choose the .p12/.pfx file downloaded from your developer account.</small>
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
		<label for="alert-text" class="col-sm-2 col-form-label">Alert Text</label>
		<div class="col-sm-10">
			<input type="text" class="form-control" name="alert-text" id="alert-text">
		</div>
	</div>
	<div class="form-group row">
		<label for="badge-count" class="col-sm-2 col-form-label">Badge Count</label>
		<div class="col-sm-10">
			<input type="number" class="form-control" name="badge-count" id="badge-count" min="0">
		</div>
	</div>
	<div class="form-group row">
		<label for="sound-name" class="col-sm-2 col-form-label">Sound Name</label>
		<div class="col-sm-10">
			<input type="text" class="form-control" name="sound-name" id="sound-name">
			<small class="form-text text-muted">Use the value "default" to make the device play a sound when the notification is received.</small>
		</div>
	</div>

	<div>
		<button type="submit" class="btn btn-primary mt-3">Send</button>
	</div>
</form>

<script type="text/javascript">
	$(function() {
		$("input[name=auth-type]:radio").change(function() {
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
		});
	})
</script>
`

const sendResultViewHtml string = `
<div class="alert alert-success" role="alert">
	<h4 class="alert-heading">Notification sent!</h4>
	<div>APNS-ID: {{ .ApnsId }}</div>
</div>
`
