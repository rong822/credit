let express = require("express");
let credit = require("../service/credit");
let service_util = require("../common/service_util");
let router = express.Router();
let logger = require("../utils/log4js").getLogger("credit");


router.post("/issuePoints", issuePoints);
router.post("/queryBalance", queryBalance);
router.post("/transferPoints", transferPoints);
router.post("/queryAccountHistory", queryAccountHistory);
router.post("/registerPoints", registerPoints);
// router.post("/createAccount", createAccount);


function issuePoints(req, res) {
	let reqString = JSON.stringify(req.body);
	credit.issuePoints(reqString, function (err, result) {
		if (result) {
			res.json({
				code: "0",
				results: result,
				msg: "issuePoints succeed "
			});
		} else {
			logger.error(err.toString());
			res.json({
				code: "-1",
				results: err.toString(),
				msg: "issuePoints failed "
			});
		}
	});
}

function queryBalance(req, res) {
	let reqString = JSON.stringify(req.body);
	credit.queryBalance(reqString, function (err, result) {
		if (result) {
			res.json({
				code: "0",
				results: result,
				msg: "queryBalance succeed "
			});
		} else {
			logger.error(err.toString());
			res.json({
				code: "-1",
				results: err.toString(),
				msg: "queryBalance failed "
			});
		}
	});
}

function transferPoints(req, res) {
	let reqString = JSON.stringify(req.body);
	credit.transferPoints(reqString, function (err, result) {
		if (result) {
			res.json({
				code: "0",
				results: result,
				msg: "transferPoints succeed "
			});
		} else {
			logger.error(err.toString());
			res.json({
				code: "-1",
				results: err.toString(),
				msg: "transferPoints failed "
			});
		}
	});
}

function queryAccountHistory(req, res) {
	let options = {
		pageNum: req.body.pageNum,
		pageSplit: req.body.pageSplit
	}
	delete req.body.pageNum
	delete req.body.pageSplit
	let reqString = JSON.stringify(req.body);
	credit.queryAccountHistory(reqString, options, function (err, result) {
		if (result) {
			res.json({
				code: "0",
				results: result,
				msg: "queryAccountHistory succeed "
			});
		} else {
			logger.error(err.toString());
			res.json({
				code: "-1",
				results: err.toString(),
				msg: "queryAccountHistory failed "
			});
		}
	});
}

function registerPoints(req, res) {
	let reqString = JSON.stringify(req.body);
	credit.registerPoints(reqString, function (err, result) {
		if (result) {
			res.json({
				code: "0",
				results: result,
				msg: "registerPoints succeed "
			});
		} else {
			logger.error(err.toString());
			res.json({
				code: "-1",
				results: err.toString(),
				msg: "registerPoints failed "
			});
		}
	});
}
// function createAccount(req, res) {
// 	let reqString = JSON.stringify(req.body);
// 	credit.createAccount(reqString, function (err, result) {
// 		if (result) {
// 			res.json({
// 				code: "0",
// 				results: result,
// 				msg: "createAccount succeed "
// 			});
// 		} else {
// 			logger.error(err.toString());
// 			res.json({
// 				code: "-1",
// 				results: err.toString(),
// 				msg: "createAccount failed "
// 			});
// 		}
// 	});
// }
module.exports = router;