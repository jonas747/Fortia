var Fortia = Fortia || {};
Fortia.apiRequest = function(url, method, success, failed){

}

Fortia.REST = function(options){
	this.urlRoot = ""
	if (options) {
		for(var prop in options){
			this[prop] = options[prop]
		}
	};
	this.responseType = "json";
}

Fortia.REST.prototype = {
	request: function(path, method, data, success, failed){
		var that = this;
		if (failed === undefined) {
			failed = function(req){
				console.log("REST Request failed, api wrapper: ", that, "Request: ", req);
			}
		};

		var url = this.urlRoot + path;
		if (method === "GET") {
			if (data) {
				switch(typeof(data)){
					case "object":
						url += "?";
						for(var key in data){
							url += key+"="+data[key]+"&";
						}
					case "string":
						if (data[0] !== "?") {
							data = "?"+data;
						};
						url += data;
						break;
				}
			};
		};

		var xhr = new XMLHttpRequest();
		xhr.open(method, url, true);
		
		xhr.withCredentials = true;
		xhr.timeout = 10000; // 10 seconds timeout
		xhr.responseType = this.responseType;

		xhr.onload = function(){
			if (xhr.readyState === 4) {
				if (xhr.status === 200) {
					success(xhr.response);
				}else{
					failed(xhr, xhr.response);
				}
			};
		}

		if (method === "GET") {
			xhr.send(null)
		}else{
			if (data) {
				xhr.send(data)
			}else{
				xhr.send(null)
			}
		}
	},
	get: function(path, data, success, failed){
		return this.request(path, "GET", data, success, failed)
	},
	post: function(path, data, success, failed){
		return this.request(path, "POST", data, success, failed)
	},
	patch: function(path, data, success, failed){
		return this.request(path, "PATCH", data, success, failed)
	},
}