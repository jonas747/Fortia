var Fortia = Fortia || {};
Fortia.initAdmWorlds = function(){
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
		template: templates.lobbyservers,

		events: {
			"click .server-join": "join",
		},

		initialize: function() {	    	
	    },

		render: function(){
			this.$el.html(this.template(this));
			return this
		},
		join: function(){
			Fortia.router.navigate("game/"+this.name(), {trigger: true});
		},
		// Getters for template
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
		},
	});

	var AdmWorldsMainView = Backbone.View.extend({
		el: "#main-container",
		template: templates.lobbyadminworlds,
		templateHeader: templates.nav,

		events:{
			"click .world-create": "createWorld",
		},

		initialize: function() {
			this.listenTo(worlds, 'add', this.addWorld);
			this.listenTo(worlds, 'reset', this.addAllWorlds);
		},

		render: function(){
			var header = this.templateHeader(Fortia.userInfo)
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
	    createWorld: function(){

	    },
	    switchTo: function(){
	    	var that = this;
	    	Fortia.authApi.get("me", "", function(info){
	    		if (!info.role) {
	    			info.amenuClass = "hidden";
	    		};
	    		Fortia.userInfo = info;
	    		that.render();
	    	})
	    }
	});
	return new AdmWorldsMainView();
}