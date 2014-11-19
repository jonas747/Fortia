var Fortia = Fortia || {};

$(function(){
	var homeView, lobbyView, gameView;

	function runFortia(){
		console.log("Sarting everything")

		Fortia.authApi = new Fortia.REST({urlRoot: "http://localhost:8080/"});

		Fortia.homeView = Fortia.initHome();
		Fortia.lobbyView = Fortia.initLobbyMain();
		Fortia.gameView = Fortia.initGameView();
		Fortia.admWorldsView = Fortia.initAdmWorlds();

		Fortia.initRouter();
		Backbone.history.start()	
	}

	function DlTemplate(name, callback){
		$.ajax({
		    url: 'templates/'+name+".tmpl",
		    type: 'GET',
		    dataType: "text",
		    success: function(a) {
		    	var template = Handlebars.compile(a)
		    	callback(template, true);
		    },
		    error: function(a){
		    	callback(null, false);
		    },
		});
	}

	function initTemplates(names, callback){
		var total = names.length
		var templates = {};
		var curProgress = 0;
		console.log("Downloading templates")
		for (var i = 0; i < names.length; i++) {
			(function(){
				var name = names[i]
				DlTemplate(name, function(template, ok){
					console.log("Finnished with "+name, ok)
					templates[name] = template
					curProgress++
					if (curProgress >= total) {
						callback(templates);
					};
				});
			})()
		};
	}

	var templateList = [
		"home",
		"lobbymain",
		"lobbyservers",
		"nav",
		"game",
		"lobbyadminworlds",
		"lobbyadminoverview",
	]

	initTemplates(templateList, function(templates){
		window.templates = templates;
		runFortia();
	});
});

Fortia.changeView = function(newView){
	if (Fortia.currentView) {
		if (Fortia.currentView.switchFrom) {
			Fortia.currentView.switchFrom(newView);
		};
	};
	if (newView.switchTo) {
		newView.switchTo(Fortia.currentView);
	};
	Fortia.currentView = newView
};