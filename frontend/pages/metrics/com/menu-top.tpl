{{define "MenuTop"}}
<!-- metrics top menu -->
<nav class="nav navbar navbar-expand-sm bg-gradient">
	<div class="container-fluid">
		<!-- logo -->
		<a class="navbar-brand" href="https://poixson.com/">
			<img src="/static/pxn-logo.png" width="100" height="30" alt="pxn"
				style="margin: 2px; margin-left: 20px; margin-right: 10px;" />
		</a>
		<button class="navbar-toggler" type="button"
				data-bs-toggle="collapse"
				data-bs-target="#NavbarSupportedContent"
				aria-controls="NavbarSupportedContent"
				aria-expanded="false"
				aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
		<div class="collapse navbar-collapse" id="NavbarSupportedContent">
			<ul class="navbar-nav me-auto mb-2 mb-sm-0">
				<!-- global -->
				<li class="nav-item">
					<a href="/" class="nav-link {{if eq .Page "global"}}active{{end}}">
						<i class="bi bi-globe2"></i> Global
					</a>
				</li>
			</ul>
			<!-- alpha -->
			<div class="badge fs-4 text-danger text-bg-dark text-uppercase"
				data-bs-toggle="tooltip" data-bs-placement="bottom"
				data-bs-title="Currently in alpha testing and active development. Your feedback is appreciated. Thank you for your support!">Alpha</div>
			<ul class="navbar-nav ms-auto mb-2 mb-sm-0">
				<!-- wiki/docs -->
				<li class="nav-item">
					<a href="/wiki/" class="nav-link {{if eq .Page "wiki"}}active{{end}}">
						<i class="bi bi-file-earmark-text"></i> Wiki/Docs
					</a>
				</li>
				<!-- status -->
				<li class="nav-item">
					<a href="/status/" class="nav-link {{if eq .Page "status"}}active{{end}}">
						<i class="bi bi-clipboard2-pulse"></i> Status
					</a>
				</li>
				<!-- about -->
				<li class="nav-item">
					<a href="/about/" class="nav-link {{if eq .Page "about"}}active{{end}}">
						<i class="bi bi-peace"></i> About Us
					</a>
				</li>
			</ul>
		</div>
	</div>
</nav>
{{end}}
