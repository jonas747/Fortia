function initLobbyMain(){
	var WorldModel = Backbone.Model.extend({
		defaults: function() {
	      return {
	        name: "default server name",
	        started: Date.now(),
	        players: 9999,
	        size: 1,
	      };
	    },
	});

	var WorldList = Backbone.Collection.extend({
		model: WorldModel,
		url: "http://localhost:8080/worlds",
	});

	var worlds = new WorldList(); 

	var WorldView = Backbone.View.extend({
		template: templates.lobbyServer,
		initialize: function() {
	      this.listenTo(this.model, 'change', this.render);
	    },

		render: function(){
			this.$el.html(this.template(this));
			return this
		},
		started: function(){
			var date = new Date(parseInt(this.model.get("started")));
			return date.toString();
		},
		name: function(){
			return this.model.get("name")
		},
		players: function(){
			return this.model.get("players")
		},
		size: function(){
			return this.model.get("size")
		}
	});

	var LobbyMainView = Backbone.View.extend({
		el: "#main-container",
		template: templates.lobbyMain,
		templateHeader: templates.lobbyHeader,

		initialize: function() {
			this.listenTo(worlds, 'add', this.addWorld);
			this.listenTo(worlds, 'reset', this.addAllWorlds);
		},

		render: function(){
			var header = this.templateHeader({username: localStorage.getItem("username")})
			var body = this.template()
			this.$el.html(header + "\n" + body)
			worlds.fetch()
		},

		addWorld: function(model){
			var view = new WorldView({model: model});
			this.$("#lobby-main-worlds").append(view.render().el)
		},

		addAllWorlds: function() {
	      worlds.each(this.addWorld, this);
	    },
	});
	return new LobbyMainView();
}