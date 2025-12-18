let zookeeper = require("node-zookeeper-client");
let server_config = require("../config/server_config");
let client = zookeeper.createClient(server_config.ZK.ADDR, {sessionTimeout: 20000});
let serviceAddr = server_config.SERVER_ADDR + ":" + server_config.LISTEN_PORT;
let servicePath = server_config.ZK.PATH + "/server1";

function init() {
	client.once("connected", function () {
		console.log("Connected to the server.");
		client.mkdirp(server_config.ZK.PATH, function (error, path) {
			if (error) {
				console.log(error.stack);
				return;
			}
			client.exists(servicePath,
				function (event) {
					console.log("Got event: %s.", event);
				},
				function (error, stat) {
					if (error) {
						console.log(error.stack);
						return;
					}
					if (!stat) {
						console.log("Node does not exist.");
						registerSrv();
					}
				});
		});
	});
	client.connect();
}

function registerSrv() {
	client.create(servicePath, new Buffer(serviceAddr), zookeeper.CreateMode.EPHEMERAL, function (error) {
		if (error) {
			if (error.getCode() == zookeeper.Exception.NODE_EXISTS) {
				console.log("Node exists.");
			} else {
				console.log("Failed to create node: %s due to: %s.", servicePath, error.stack);
			}
		} else {
			console.log("Node: %s is successfully created.", servicePath);
		}
	});
}

module.exports.init = init;