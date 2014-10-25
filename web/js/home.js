var Fortia = Fortia || {}; 
Fortia.initHome = function(){
	var HomeView = Backbone.View.extend({
		el: "#main-container",
		template: templates.home,
		events: {
			"submit #login-form": 	 "login",
			"submit #register-form": "register"
		},
		login: function(event){
			event.preventDefault();
			var userField = $("#login-username");
			var passwordField = $("#login-password");

			if(!validateForm(userField, passwordField)){
				return;
			}

			var fobj = {
				username: userField.val(),
				pw: passwordField.val(),
			};
			var that = this
			
			Fortia.authApi.post("login", fobj, function(response){
				console.log(response);
				// Navigate to lobby
				localStorage.setItem("username", fobj.username);
				Fortia.router.navigate("lobby", {trigger: true});			
			}, function(req){
				var response = req.responseJSON;
				if (response === "") {
					that.loginError = "Unknown error occured";
				}else{
					that.loginError = response.error;
				}
				that.render()
			})
		},
		register: function(event){
			event.preventDefault();
			var userField = 	$("#register-username");
			var emailField = 	$("#register-email");
			var pw1Field = 		$("#register-password1");
			var pw2Field = 		$("#register-password2");
			
			if(!validateForm(userField, emailField, pw1Field, pw2Field)){
				this.registerError = "Check your fields again!"
				this.render();
				return;
			}

			if (pw1Field.val() !== pw2Field.val()) {
				pw2Field.parent().addClass("has-error")
				pw1Field.parent().addClass("has-error")
				return;
			};

			var fobj = {
				username: userField.val(),
				email: emailField.val(),
				pw: pw1Field.val(),
			}
			var that = this
			
			Fortia.authApi.post("register", fobj, function(reponse){
				that.render();
				$("#register-success").modal();
				$("#login-username").val(fobj.username);
			}, function(req){
		    	var response = req.responseJSON;
				if (response === "") {
					that.registerError = "Uknown error occured";
				}else{
					that.registerError = response.error;
				}
				that.render()
			});
		},
		render: function(){
			this.$el.html(this.template(this));
		},
		switchTo: function(){
			this.render();
		}
	})

	function validateForm(){
		var ok = true
		for (var i = 0; i < arguments.length; i++) {
			var field = arguments[i]
			var value = field.val()
			var parent = field.parent()
			if (value === "") {
				parent.addClass("has-error")
				ok = false
			}else{
				parent.removeClass("has-error")
			}
		};
		return ok
	}

	var homeView = new HomeView();
	return homeView;
}