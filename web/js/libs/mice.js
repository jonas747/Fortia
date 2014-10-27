// Simple mice input library
var Mouse = {};

(function(){
	var buttons = [];
	var aliases = {
		"Left":	0,
		"Right": 2,
		"Middle": 1,
	}

	var mousePosX = 0;
	var mousePosY = 0;

	window.addEventListener("mousemove", function(event){
		mousePosX = event.clientX;
		mousePosY = event.clientY;
	});
	window.addEventListener("mousedown", function(event){
		buttons[event.button] = true;
	});
	window.addEventListener("mouseup", function(event){
		buttons[event.button] = false;
	});

	Mouse.getMouseX = function(){
		return mousePosX;
	};

	Mouse.getMouseY = function(){
		return mousePosY;
	};

	Mouse.isButtonDown = function(button){
		return buttons[aliases[button]];
	}
})()