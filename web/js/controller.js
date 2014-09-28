$(function(){
	var CurrentView = 1;

	$("#authLoginForm").on("submit", Auth.Login)
	$("#authRegisterForm").on("submit", Auth.Register)

	//$("#authLoginButton").submit(Auth.Login)
	//$("#authRegisterButton").submit(Auth.Register)


})