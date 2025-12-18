let express = require("express");
let app = express();
let bodyParser = require("body-parser");
let config = require("./utils/config").init();
let serverConfig = require(config.server);
let connectlogger = require("./utils/log4js").getConnectLogger();
let logger = require("./utils/log4js").getLogger("app");
let register = require("./service/register")
let host = process.env.HOST || serverConfig.SERVER_ADDR;
let port = process.env.PORT || serverConfig.LISTEN_PORT;
var cors = require('cors');
var bearerToken = require('express-bearer-token');
var expressJWT = require('express-jwt');
var jwt = require('jsonwebtoken');
require('./app_config.js');
var hfc = require('fabric-client');
var util = require("util")
let fs = require("fs")
// let jwt_decode = require("jwt-decode")

app.use(connectlogger);
app.options('*', cors());
app.use(cors());
//support parsing of application/json type post data
app.use(bodyParser.json());
//support parsing of application/x-www-form-urlencoded post data
app.use(bodyParser.urlencoded({
	extended: false
}));
// set secret variable
app.set('secret', 'thisismysecret');
app.use(expressJWT({
	secret: 'thisismysecret'
}).unless({
	path: ['/api/register']
}));
app.use(bearerToken());
app.use(function (req, res, next) {
	logger.debug(' ------>>>>>> new request for %s', req.originalUrl);
	if (req.originalUrl.indexOf('/api/register') >= 0) {
		return next();
	}

	var token = req.token;
	jwt.verify(token, app.get('secret'), function (err, decoded) {
		if (err) {
			res.send({
				success: false,
				message: err
			});
			return;
		} else {
			// add the decoded user name and org name to the request object
			// for the downstream code to use
			var tempKeyStoragePath = "./token"
			fs.readdirSync(tempKeyStoragePath).forEach(file => {
				fs.readFile(tempKeyStoragePath + "/" + file, "utf8", function read(err, data) {
					if (err) {
						res.send({
							success: false,
							message: err
						});
						return;
					}
					try {
						var decoded_token = JSON.parse(data)
						if (decoded.username == decoded_token.username && decoded.password == decoded_token.password) {
							req.username = decoded.username;
							req.password = decoded.password;
							logger.debug(util.format('Decoded from JWT token: username - %s, password - %s', decoded.username, decoded.password));
							return next();
						}
					} catch (error) {
						console.log("Bypass this key, due to", error)
					}
				});
			})
			// res.send({
			// 	success: false,
			// 	message: "Not find any user"
			// });
			// return;
		}
	});
});
// Register and enroll user
app.post('/api/register', async function (req, res) {
	var username = req.body.username;
	var password = req.body.password
	logger.debug('End point : /register');
	logger.debug('User name : ' + username);
	logger.debug('password name  : ' + password);
	if (!username) {
		res.json(getErrorMessage('\'username\''));
		return;
	}
	if (!password) {
		res.json(getErrorMessage('\'password\''));
		return;
	}
	var token = jwt.sign({
		exp: Math.floor(Date.now() / 1000) + parseInt(360000000),
		username: username,
		password: password
	}, app.get('secret'));
	let response = await register.getRegisteredUser(username, password);
	logger.debug('-- returned from registering the username %s', username, password);
	if (response && typeof response !== 'string') {
		logger.debug('Successfully registered the username %s', username, password);
		fs.writeFile("./token/" + token, JSON.stringify({
			username: username,
			password: password
		}), function (err) {
			if (err) {
				res.json("Restore token error!")
			}
			console.log(token)
		})
		response.token = token;
		res.json(response);
	} else {
		logger.debug('Failed to register the username %s with::%s', username, password, response);
		res.json({
			success: false,
			message: response
		});
	}

});

require("./routers/index")(app);
let server = app.listen(port, host, function () {
	logger.info(`server listening at: http://${host}:${port}/`);
	// register service in zookeeper: to communicate with bigtree java
	//zookeeper.init();

	//Graceful start
	// process.send("ready");
});

//Graceful shutdown
process.on("SIGINT", () => {
	logger.info("SIGINT signal received.");

	// Stops the server from accepting new connections and finishes existing connections.
	server.close(function (err) {
		if (err) {
			logger.error(err);
			process.exit(1);
		}
	});
});

process.on("unhandledRejection", function (err) {
	logger.error("catch unhandled exception:", err.stack);
});

process.on("uncaughtException", function (err) {
	logger.error("catch uncaught exception:", err.stack);
});