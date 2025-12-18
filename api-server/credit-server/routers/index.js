module.exports = function (app) {

	let router = require("./router");
	app.use("/api/credit", router);

};