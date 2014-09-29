function initHome(){
	var LoginModel = Backbone.Model.extend({
		urlRoot: "/login",
	})

	var HomeView = Backbone.View.extend({
		el: "#main-container",
		template: templates.home,
		events: {
			"submit #login-form": "login",
			"submit #register-form": "register"
		},
		login: function(event){
			event.preventDefault();
			var userField = $("#login-username");
			var passwordField = $("#login-password");

			if(!validateForm(userField, passwordField)){
				return
			}

			var fobj = {
				username: userField.val(),
				pw: passwordField.val(),
			};
			var that = this
			$.ajax({
			    url: 'http://localhost:8080/login',
			    type: 'POST',
			    data: JSON.stringify(fobj),
			    dataType: "json",
			    success: function(a) {
			    	console.log(a)
			    	// Navigate to lobby
			    	localStorage.setItem("username", fobj.username)
			    	router.navigate("lobby", {trigger: true});
			    },
			    error: function(a){
			    	var response = a.responseJSON;
					if (response === "") {
						that.loginError = "Uknown error occured";
					}else{
						that.loginError = response.error;
					}
					that.render()
			    },
			});
		},
		register: function(event){
			event.preventDefault();
			var userField = $("#register-username");
			var emailField = $("#register-email");
			var pw1Field = $("#register-password1");
			var pw2Field = $("#register-password2");
			
			if(!validateForm(userField, emailField, pw1Field, pw2Field)){
				return
			}

			if (pw1Field.val() !== pw2Field.val()) {
				pw2Field.parent().addClass("has-error")
				return
			};

			var fobj = {
				username: userField.val(),
				email: emailField.val(),
				pw: pw1Field.val(),
			}
			var that = this
			$.ajax({
			    url: 'http://localhost:8080/register',
			    type: 'POST',
			    data: JSON.stringify(fobj),
			    dataType: "json",
			    success: function(a) {
			    	console.log(a)
			    },
			    error: function(a){
			    	var response = a.responseJSON;
					if (response === "") {
						that.registerError = "Uknown error occured";
					}else{
						that.registerError = response.error;
					}
					that.render()
			    },
			});
		},
		render: function(){
			this.$el.html(this.template(this));
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