{{ define "dev" }}

<div class='panel panel-info'>
	<div class="panel-heading"><h3 class="panel-title"><span class="glyphicon glyphicon-edit"></span>&nbsp; Development Site</h3></div>
	<div class="panel-body">
		The dev site has previews of material not yet ready for release. It's only available to patrons and contributors, including translators.
	</div>
	<div class="panel-footer">
		<a class='btn btn-primary btn-sm' href='{{.DevLoginURL}}' target='_blank'><span class="glyphicon glyphicon-lock"></span>&nbsp; Log in</a> 
	</div>
</div>

{{ end }}