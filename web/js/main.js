var router;

$(function(){

	$.ajaxPrefilter(function( options, originalOptions, jqXHR ) {
		if (!options.xhrFields) {
			options.xhrFields = {
		      withCredentials: true
			}
		}else{
			options.xhrFields.withCredentials = true;
		}
	});

	window.templates = {
		home:  			initTemplate("#home-template"),
		lobbyServer:  	initTemplate("#lobby-server-template"),
		lobbyMain:  	initTemplate("#lobby-main-template"),
		lobbyHeader:  	initTemplate("#lobby-header-template"),
	}

	var homeView = initHome();
	var lobbyView = initLobbyMain();

	var Router = Backbone.Router.extend({
		routes: {
			"": 		"home",
			"lobby": 	"lobby",
		},

		home: function(){
			homeView.render();
			console.log("Im home!")
		},
		lobby: function(){
			lobbyView.render();
			console.log("Im in the lobby!")
		}
	});

	router = new Router();
	Backbone.history.start()
})

function initTemplate(id){
	var source   = $(id).html();
	var template = Handlebars.compile(source);
	return template;
}