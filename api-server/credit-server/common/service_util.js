let path = require("path");
let fs = require("fs-extra");
let os = require("os");
let util = require("util");

let Client = require("fabric-client");
let copService = require("fabric-ca-client/lib/FabricCAClientImpl.js");
let User = require("fabric-client/lib/User.js");

let logger = require("fabric-client/lib/utils.js").getLogger("service_util");
let config = require("../utils/config").getConfig();


// all temporary files and directories are created under here
let tempdir = path.join(os.tmpdir(), "hfc");

logger.info(util.format(
	"\n\n*******************************************************************************" +
    "\n*******************************************************************************" +
    "\n*                                          " +
    "\n* Using temp dir: %s" +
    "\n*                                          " +
    "\n*******************************************************************************" +
    "\n*******************************************************************************\n", tempdir));

module.exports.getPage = function(){
	
}

module.exports.getTempDir = function () {
	fs.ensureDirSync(tempdir);
	return tempdir;
};

// directory for file based KeyValueStore
module.exports.KVS = path.join(tempdir, "hfc-test-kvs");
module.exports.storePathForOrg = function (org) {
	return module.exports.KVS + "_" + org;
};

// temporarily set $GOPATH to the test fixture folder
module.exports.setupChaincodeDeploy = function () {
	process.env.GOPATH = path.join(__dirname, "../../fabric-sdk-node/test/fixtures");
};

// specifically set the values to defaults because they may have been overridden when
// running in the overall test bucket ('gulp test')
module.exports.resetDefaults = function () {
	global.hfc.config = undefined;
	require("nconf").reset();
};

module.exports.cleanupDir = function (keyValStorePath) {
	let absPath = path.join(process.cwd(), keyValStorePath);
	let exists = module.exports.existsSync(absPath);
	if (exists) {
		fs.removeSync(absPath);
	}
};

module.exports.getUniqueVersion = function (prefix) {
	if (!prefix) prefix = "v";
	return prefix + Date.now();
};

// utility function to check if directory or file exists
// uses entire / absolute path from root
module.exports.existsSync = function (absolutePath /*string*/) {
	try {
		let stat = fs.statSync(absolutePath);
		return (stat.isDirectory() || stat.isFile());
	}
	catch (e) {
		return false;
	}
};

module.exports.readFile = readFile;

Client.addConfigFile(config.network);
let ORGS = Client.getConfigSetting("network");

let tlsOptions = {
	trustedRoots: [],
	verify: false
};

function getMember(username, password, client, userOrg) {
	let caUrl = ORGS[userOrg].ca.url;
	return client.getUserContext(username, true)
		.then((user) => {
			return new Promise((resolve, reject) => {
				if (user && user.isEnrolled()) {
					logger.info("Successfully loaded member from persistence");
					return resolve(user);
				}

				let member = new User(username);
				let cryptoSuite = client.getCryptoSuite();
				if (!cryptoSuite) {
					cryptoSuite = Client.newCryptoSuite();
					if (userOrg) {
						cryptoSuite.setCryptoKeyStore(Client.newCryptoKeyStore({path: module.exports.storePathForOrg(ORGS[userOrg].name)}));
						client.setCryptoSuite(cryptoSuite);
					}
				}
				member.setCryptoSuite(cryptoSuite);

				// need to enroll it with CA server
				let cop = new copService(caUrl, tlsOptions, ORGS[userOrg].ca.name, cryptoSuite);

				return cop.enroll({
					enrollmentID: username,
					enrollmentSecret: password
				}).then((enrollment) => {
					logger.info("Successfully enrolled user '" + username + "'");

					return member.setEnrollment(enrollment.key, enrollment.certificate, ORGS[userOrg].mspid);
				}).then(() => {
					let skipPersistence = false;
					if (!client.getStateStore()) {
						skipPersistence = true;
					}
					return client.setUserContext(member, skipPersistence);
				}).then(() => {
					return resolve(member);
				}).catch((err) => {
					logger.error("Failed to enroll and persist user. Error: " + err.stack ? err.stack : err);
				});
			});
		});
}

function getAdmin(client, userOrg) {
	let username = "peer"+userOrg+"Admin";
	return client.getUserContext(username, true)
		.then((user) => {
			return new Promise((resolve, reject) => {
				if (user && user.isEnrolled()) {
					logger.info("Successfully loaded member from persistence");
					return resolve(user);
				}

				let keyPath = path.join(__dirname, util.format(ORGS[userOrg].admin.key));
				let keyPEM = Buffer.from(readAllFiles(keyPath)[0]).toString();
				let certPath = path.join(__dirname, util.format(ORGS[userOrg].admin.cert));
				let certPEM = readAllFiles(certPath)[0];

				//let cryptoSuite = Client.newCryptoSuite();
				//if (userOrg) {
				//	cryptoSuite.setCryptoKeyStore(Client.newCryptoKeyStore({path: module.exports.storePathForOrg(ORGS[userOrg].name)}));
				//	client.setCryptoSuite(cryptoSuite);
				//}

				let cryptoSuite = client.getCryptoSuite();
				if (!cryptoSuite) {
					cryptoSuite = Client.newCryptoSuite();
					if (userOrg) {
						cryptoSuite.setCryptoKeyStore(Client.newCryptoKeyStore({path: module.exports.storePathForOrg(ORGS[userOrg].name)}));
						client.setCryptoSuite(cryptoSuite);
					}
				}

				return resolve(client.createUser({
					username: "peer" + userOrg + "Admin",
					mspid: ORGS[userOrg].mspid,
					cryptoContent: {
						privateKeyPEM: keyPEM.toString(),
						signedCertPEM: certPEM.toString()
					}
				}));
			});
		});
}

function getOrdererAdmin(client) {
	let username = "ordererAdmin";
	return client.getUserContext(username, true)
		.then((user) => {
			return new Promise((resolve, reject) => {
				if (user && user.isEnrolled()) {
					logger.info("Successfully loaded member from persistence");
					return resolve(user);
				}

				let keyPath = path.join(__dirname, ORGS.orderer.keystore);
				let keyPEM = Buffer.from(readAllFiles(keyPath)[0]).toString();
				let certPath = path.join(__dirname, ORGS.orderer.signcerts);
				let certPEM = readAllFiles(certPath)[0];

				return resolve(client.createUser({
					username: "ordererAdmin",
					mspid: "OrdererMSP",
					cryptoContent: {
						privateKeyPEM: keyPEM.toString(),
						signedCertPEM: certPEM.toString()
					}
				}));
			});
		});
}

function readFile(path) {
	return new Promise((resolve, reject) => {
		fs.readFile(path, (err, data) => {
			if (err)
				reject(new Error("Failed to read file " + path + " due to error: " + err));
			else
				resolve(data);
		});
	});
}

function readAllFiles(dir) {
	let files = fs.readdirSync(dir);
	let certs = [];
	files.forEach((file_name) => {
		let file_path = path.join(dir, file_name);
		logger.debug(" looking at file ::" + file_path);
		let data = fs.readFileSync(file_path);
		certs.push(data);
	});
	return certs;
}

module.exports.getOrderAdminSubmitter = function (client) {
	return getOrdererAdmin(client);
};

module.exports.getSubmitter = function (client, peerOrgAdmin, org) {

	let peerAdmin, userOrg;
	if (typeof peerOrgAdmin === "boolean") {
		peerAdmin = peerOrgAdmin;
	} else {
		peerAdmin = false;
	}

	// if the 3rd argument was skipped
	if (typeof peerOrgAdmin === "string") {
		userOrg = peerOrgAdmin;
	} else {
		if (typeof org === "string") {
			userOrg = org;
		} else {
			userOrg = "org1";
		}
	}

	if (peerAdmin) {
		return getAdmin(client, userOrg);
	} else {
		return getMember("admin", "adminpw", client, userOrg);
	}
};
