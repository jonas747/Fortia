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
	this.dType = "json";
}

Fortia.REST.prototype = {
	request: function(path, method, data, success, failed){
		var d = "";
		if (data) {
			d = JSON.stringify(data);
		};
		var that = this;
		if (failed === undefined) {
			failed = function(req){
				console.log("REST Request failed, api wrapper: ", that, "Request: ", req);
			}
		};
		$.ajax({
			    url: this.urlRoot+path,
			    type: method,
			    data: d,
			    dataType: this.dType,
			    success: success,
			    error: failed,
			    xhrFields: {
		      		withCredentials: true
				}
			});
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