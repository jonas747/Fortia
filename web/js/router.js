var Fortia = Fortia || {};
Fortia.initRouter = function(){
	var Router = Backbone.Router.extend({
		routes: {
			"": 			"home",
			"lobby": 		"lobby",
			"adm-worlds": 	"admWorlds",

			"game/:world":				 "game",
			"game/:world/mail": 		 "gameMailInbox",
			"game/:world/mail/inbox": 	 "gameMailInbox",
			"game/:world/mail/sent": 	 "gameMailInbox",
			"game/:world/mail/archived": "gameMailInbox",

			"game/:world/alliance-join": "",
			"game/:world/alliance/:alliance/overview": "",
			"game/:world/alliance/:alliance/members": "",	
		},

		home: function(){
			Fortia.changeView(Fortia.homeView)
			console.log("Im home!")
		},
		lobby: function(){
			Fortia.changeView(Fortia.lobbyView)
			console.log("Im in the lobby!")
		},
		admWorlds: function(){
			Fortia.changeView(Fortia.admWorldsView)
		},

		game: function(world){
			Fortia.changeView(Fortia.gameView)
			Fortia.game.init(Fortia.gameView, world)
			Fortia.game.start();
		},
	});
	Fortia.router = new Router();
}