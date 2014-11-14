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

	var offsetX = 0;
	var offsetY = -50;

	var dirty = false;

	window.addEventListener("mousemove", function(event){
		mousePosX = event.clientX;
		mousePosY = event.clientY;
		dirty = true;
	});
	window.addEventListener("mousedown", function(event){
		buttons[event.button] = true;
	});
	window.addEventListener("mouseup", function(event){
		buttons[event.button] = false;
	});

	Mouse.getX = function(){
		return mousePosX + offsetX;
	};

	Mouse.getY = function(){
		return mousePosY + offsetY;
	};

	Mouse.getMouse = function(){
		return {x: mousePosX+offsetX, y: mousePosY+offsetX};
	}

	Mouse.setOffset = function(x, y){
		offsetX = x;
		offsetY = y;
	}

	Mouse.isButtonDown = function(button){
		return buttons[aliases[button]];
	}

	Mouse.isDirty = function(){
		return dirty;
	}

	Mouse.setDirty = function(d){
		dirty = d;
	}
})()