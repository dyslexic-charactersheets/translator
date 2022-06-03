{{ define "live" }}

<div class='panel panel-primary'>
	<div class="panel-heading"><h3 class="panel-title"><span class="glyphicon glyphicon-edit"></span>&nbsp; Live Site</h3></div>
	<div class="panel-body">
		Log into the live site for access to edit & save features on character sheets.
	</div>
	<div class="panel-footer">
		<a class='btn btn-primary btn-sm' href='{{.LiveLoginURL}}' target='_blank'><span class="glyphicon glyphicon-lock"></span>&nbsp; Log in</a> 
	</div>
</div>

{{ end }}