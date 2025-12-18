let log4js = require("log4js");
let config = require("../utils/config").getConfig();
log4js.configure(require(config.log4js));

function getLogger(name) {
	let logger = log4js.getLogger(name);
	return logger;
}

function getConnectLogger() {
	let logger = log4js.connectLogger(getLogger("access"), {level: log4js.levels.INFO});
	return logger;
}

module.exports.getLogger = getLogger;
module.exports.getConnectLogger = getConnectLogger;