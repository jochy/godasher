package main

var JenkinsViewTemplate = `
<script src="https://momentjs.com/downloads/moment.min.js"></script>
<style>
.jenkins-view-job {
	margin-bottom: 0px !important;
	display: inline-block;
	flex-direction: column;
	width: 100%;
	white-space: nowrap;
	overflow: hidden;
}
.jenkins-active {
    animation-duration: 3000ms;
    animation-name: blink;
    animation-iteration-count: infinite;
    animation-direction: alternate;
    -webkit-animation:blink 3000ms infinite;
}
.jenkins-single-job {
	display: flex !important;
	height: calc(100% - 5px);
	flex-direction: column;
}
.jenkins-single-job .alert {
	display: flex !important;
	height: 100%;
}
.jenkins-single-job .multi-job-container .alert {
	height: unset !important;
	margin-bottom: 5px !important;
}
.jenkins-single-job .alert h4 {
	margin: auto;
}
.jenkins-single-job .alert span {
	margin: auto;
}
@keyframes blink {
    from { opacity: 1 }
    50% { opacity: 0.5 }
	to { opacity: 1 }
}
@-webkit-keyframes blink {
    from { opacity: 1 }
    50% { opacity: 0.5 }
	to { opacity: 1 }
}
</style>

<div style="width: 100%; margin-top:5px;" id="{{ .RandomId }}" class="{{ if eq .View "job" }}jenkins-single-job{{ end }}">
    {{ range $job := (mapValue .Response "jobs") }}
		<div class="alert {{ jobColor $job }} jenkins-view-job" data-width="100">
		    {{ if eq (mapValue $job "_class") "org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject"}}
				<h4>{{ mapValue $job "displayName" }}</h4>
				<div style="width: 100%;word-break: break-all;white-space: normal;" class="multi-job-container">
				{{ range $subjob := (mapValue $job "jobs") }}
					<div class="alert {{ jobColor $subjob }} jenkins-view-job" data-width="100">
						<h6>{{ mapValue $subjob "name" }}</h6>
						{{ if mapValue (mapValue $subjob "lastBuild") "timestamp" }}
							<span>
								<i class="fa fa-clock"></i>
								<span class="jenkins-date" data-timestamp="{{ mapValue (mapValue $subjob "lastBuild") "timestamp" }}"></span>
								- Took {{ msToMin (mapValue (mapValue $subjob "lastBuild") "duration") }} min
							</span>
						{{ end }}
					</div>
				{{ end }}
				</div>
			{{ else }}
				<h4>#{{ mapValue (mapValue $job "lastBuild") "number" }} - {{ mapValue $job "displayName" }}</h4>
				<span>
					<i class="fa fa-clock"></i>
					<span class="jenkins-date" data-timestamp="{{ mapValue (mapValue $job "lastBuild") "timestamp" }}"></span>
					- Took {{ msToMin (mapValue (mapValue $job "lastBuild") "duration") }} min
				</span>
			{{ end }}
		</div>
    {{ end }}
</div>

<script>
	resize(document.getElementById('{{ .RandomId }}'));

	var e = document.getElementsByClassName('jenkins-date');
	for(let i = 0; i < e.length; i++) {
		let elem = e[i];
		if (elem.dataset.timestamp) {
			elem.textContent = moment(new Date(parseFloat(elem.dataset.timestamp)).toISOString()).fromNow();
		}
	}

	function resize(div) {
		let i = 0;
		while (div.parentNode.clientHeight < div.clientHeight) {
			var elements = div.getElementsByClassName('jenkins-view-job');
			for(let i = 0; i < elements.length; i++) {
				let width = elements[i].dataset.width / 2;
				elements[i].dataset.width = width;
				elements[i].style.width = "calc(" + width + "% - 3px)";
			}
			div = document.getElementById('{{ .RandomId }}');
			if (i++ > 20) {
				break;
			}
		}
	}
</script>
`
