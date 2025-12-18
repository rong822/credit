let credit_service = require("../common/service");

module.exports = {

	issuePoints: function (reqString, cb) {
		let invokeRequest = {
			args: [reqString],
			functionName: "issuePoints",
			userOrg: "org1",
			userName: "admin",
			passwd: "adminpw"
		};
		credit_service.invoke(invokeRequest, function (err, result) {
			cb(err, result);
		});
	},
	queryBalance: function (reqString, cb) {
		let invokeRequest = {
			args: [reqString],
			functionName: "queryBalance",
			userOrg: "org1",
			userName: "admin",
			passwd: "adminpw"
		};
		credit_service.invoke(invokeRequest, function (err, result) {
			cb(err, result);
		});
	},
	transferPoints: function (reqString, cb) {
		let invokeRequest = {
			args: [reqString],
			functionName: "transferPoints",
			userOrg: "org1",
			userName: "admin",
			passwd: "adminpw"
		};
		credit_service.invoke(invokeRequest, function (err, result) {
			cb(err, result);
		});
	},
	queryAccountHistory: function (reqString, options, cb) {
		let invokeRequest = {
			args: [reqString],
			functionName: "queryAccountHistory",
			userOrg: "org1",
			userName: "admin",
			passwd: "adminpw",
			options: options
		};
		credit_service.invoke(invokeRequest, function (err, result) {
			cb(err, result);
		});
	},
	registerPoints: function (reqString, cb) {
		let invokeRequest = {
			args: [reqString],
			functionName: "registerPoints",
			userOrg: "org1",
			userName: "admin",
			passwd: "adminpw"
		};
		credit_service.invoke(invokeRequest, function (err, result) {
			cb(err, result);
		});
	},
	// createAccount: function (reqString, cb) {
	// 	let invokeRequest = {
	// 		args: [reqString],
	// 		functionName: "createAccount",
	// 		userOrg: "org1",
	// 		userName: "admin",
	// 		passwd: "adminpw"
	// 	};
	// 	credit_service.invoke(invokeRequest, function (err, result) {
	// 		cb(err, result);
	// 	});
	// },
};