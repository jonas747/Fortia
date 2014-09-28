var Auth = {
	Login: function(event){
		event.preventDefault();

		var user = $("#authLoginUsername").val()
		var pw = $("#authLoginPassword").val()

		if (user === "") {
			Auth.LoginError("Username field is empty")
			return
		};

		if (pw === "") {
			Auth.LoginError("Password field is empty")
			return
		};

		var obj = {
			username: user,
			pw: pw,
		}

		$.ajax({
		    url: 'http://localhost:8080/login',
		    type: 'POST',
		    data: JSON.stringify(obj),
		    dataType: "json",
		    success: function(a) {
		    	Auth.RegisterError("");
		    	$("#authLoginForm")[0].reset()
		    },
		    error: function(a){
		    	var response = a.responseJSON;
		    	if (response === null) {return};
				Auth.LoginError(response.error);
		    },
		    xhrFields: {
		      withCredentials: true
		   },
		});
	},
	LoginError: function(error){
		$("#authLoginErrorMessage").text(error)
	},
	Register: function(event){
		event.preventDefault();
		var user = $("#authRegisterUsername").val();
		var email = $("#authRegisterEmail").val();
		var pass = $("#authRegisterPassword1").val();
		var pass2 = $("#authRegisterPassword2").val();

		if (pass !== pass2) {
			Auth.RegisterError("Passwords don't match");
			return
		};
		
		if (user === "") {
			Auth.RegisterError("User field is empty");
			return
		};

		if (pass === "") {
			Auth.RegisterError("Password field is empty");
			return
		};

		if (pass2 === "") {
			Auth.RegisterError("Password2 field is empty");
			return
		};
		
		if (email === "") {
			Auth.RegisterError("Email field is empty");
			return
		};

		var obj = {
			username: user,
			email: email,
			pw: pass
		}

		$.ajax({
		    url: 'http://localhost:8080/register',
		    type: 'POST',
		    data: JSON.stringify(obj),
		    dataType: "json",
		    success: function(a) {
		    	Auth.RegisterError("");
		    	$("#authRegisterForm")[0].reset()
		    	$("#authRegisterSuccess").modal()
		    },
		    error: function(a){
		    	var response = a.responseJSON;
				Auth.RegisterError(response.error);
		    },
		});
	},
	RegisterError: function(error){
		$("#authRegisterErrorMessage").text(error);
	},
}