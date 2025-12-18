let utils = require("fabric-client/lib/utils.js");
utils.setConfigSetting("crypto-keysize", 256);
utils.setConfigSetting("key-value-store", "fabric-client/lib/impl/CouchDBKeyValueStore.js"); //override

let Client = require("fabric-client");
let logger = require("../utils/log4js").getLogger("credit_service");
let config = require("../utils/config").getConfig();

let chaincodeConfig = require(config.chaincode);
let channelName = chaincodeConfig.channelName;
let chainCodeName = chaincodeConfig.chaincodeName;
//let chainCodeVersion = chaincodeConfig.chaincodeVersion;
let peerOrgID = chaincodeConfig.invokeOrgID;

let tx_id = null;
//let nonce = null;

async function queryByTransactionID(tx_id, channel) {
    try {
        // set target peer
        let peer = peerOrgID
        // let targets = peerOrgID
        // if (targets && targets.length != 0) {
        //     const targetPeers = getTargetPeers(channel, targets);
        //     if (targetPeers.length < targets.length) {
        //         logger.debug('Failed to get all peers for targets: ' + targets);
        //     } else {
        //         peer = targetPeers;
        //     }
        // }
        let transaction_response_payload = await channel.queryTransaction(tx_id, peer)
        if (transaction_response_payload) {
            transaction_response_payload = transaction_response_payload.transactionEnvelope.payload.data.actions[0]
            // logger.debug(transaction_response_payload.data.toString());
            // let data = transaction_response_payload.data
            // logger.debug(data);
            return transaction_response_payload
        } else {
            logger.debug('transaction_response_payload is null');
            return 'transaction_response_payload is null';
        }
    } catch (error) {
        logger.debug("Failed to query due to error:" + error)
        return error.toString()
    }
}


async function invoke(invokeReq, cb) {
    let channel = null;
    let options = null
    if (invokeReq.hasOwnProperty("options")) {
        options = invokeReq.options
        delete invokeReq.options
    }
    try {
        //let client = Client.loadFromConfig('config/dev/network.yaml');
        let client = Client.loadFromConfig(config.networkYAML);
        logger.info('Successfully loaded a connection profile');
        //client.loadFromConfig('config/dev/'+invokeReq.userOrg+'.yaml');
        client.loadFromConfig(config.orgYAML[invokeReq.userOrg]);

        await client.initCredentialStores();
        logger.info('Successfully created the key value store and crypto store based on the sdk config and connection profile');

        let enrollment = null;
        if (global.enrollment)
            enrollment = global.enrollment;
        else {
            let caService = client.getCertificateAuthority();
            let request = {
                enrollmentID: 'admin',
                enrollmentSecret: 'adminpw',
                profile: 'tls'
            };

            enrollment = await caService.enroll(request);
            enrollment.key = enrollment.key.toBytes();
            logger.info('Successfully called the CertificateAuthority to get the TLS material');
            global.enrollment = enrollment;
        }
        client.setTlsClientCertAndKey(enrollment.certificate, enrollment.key);

        // let caService = client.getCertificateAuthority();
        // let request = {
        //     enrollmentID: 'admin',
        //     enrollmentSecret: 'adminpw',
        //     profile: 'tls'
        // };
        // let enrollment = await caService.enroll(request);
        // logger.info('Successfully called the CertificateAuthority to get the TLS material');
        // let key = enrollment.key.toBytes();
        // let cert = enrollment.certificate;
        // client.setTlsClientCertAndKey(cert, key);

        channel = client.getChannel(channelName);

        //await getUser(client, 'admin', 'adminpw');
        await getUser(client, invokeReq.userName, invokeReq.passwd);
        logger.info('Successfully enrolled user \'admin\' for org1');

        await channel.initialize();

        tx_id = client.newTransactionID(); // get a non admin transaction ID
        let request2 = {
            chaincodeId: chainCodeName,
            fcn: invokeReq.functionName,
            args: invokeReq.args,
            txId: tx_id
        };
        logger.debug("Sending transaction Proposal : transaction id >> " + request2.txId + "; function >>" + request2.fcn + "; args >> " + request2.args);
        let results = await channel.sendTransactionProposal(request2);
        let proposalResponses = results[0];
        let proposal = results[1];
        let all_good = true;
        let response_payload = ""
        logger.debug("Received " + proposalResponses.length + " endorsements totally.");
        for (let i in proposalResponses) {
            let one_good = false;
            let proposal_response = proposalResponses[i];
            if (proposal_response.response && proposal_response.response.status === 200) {
                one_good = channel.verifyProposalResponse(proposal_response);
                if (one_good) {
                    response_payload = proposalResponses[i].response.payload
                    // logger.debug(">>>>>>>>>>>>> proposalResponses:", response_payload.toString("utf8"))
                    logger.debug("proposal was good");
                }
            } else {
                logger.error("invokeChaincode: transaction proposal was bad");
                return cb(proposal_response.response.message, "");
            }
            all_good = all_good & one_good;
        }
        if (all_good) {
            // check all the read/write sets to see if the same, verify that each peer
            // got the same results on the proposal
            all_good = channel.compareProposalResponseResults(proposalResponses);
            logger.debug("compareProposalResponseResults execution did not throw an error");
            if (all_good) {
                logger.debug(" All proposals have a matching read/writes sets");
            } else {
                logger.error(" All proposals do not have matching read/write sets");
            }
        }
        if (all_good) {
            let parse_payload
            if (invokeReq.functionName == "queryAccountHistory") {
                let pageNum = parseInt(options.pageNum)
                let pageSplit = parseInt(options.pageSplit)
                let pageStart = pageNum * pageSplit;
                parse_payload = JSON.parse(response_payload)
                parse_payload = parse_payload.reverse()
                parse_payload = parse_payload.splice(pageStart, pageSplit)
                for (let index = 0; index < parse_payload.length; index++) {
                    let preTx = parse_payload[index]
                    // console.log(">>>>>>>>>>>>>>>>", preTx.tx_id)
                    let transaction_response_payload = await queryByTransactionID(preTx.tx_id, channel)
                    let input_args = transaction_response_payload.payload.chaincode_proposal_payload.input.chaincode_spec.input.args
                    // for (let args_index = 0; args_index < input_args.length; args_index++) {
                    //     let arg = input_args[args_index]
                    //     console.log(arg.toString())
                    // }
                    let args_func_name = input_args[0]
                    let args_body = JSON.parse(input_args[1])
                    // console.log(args_body.toString())
                    if (args_func_name.toString() == "issuePoints") {
                        args_body = args_body.issueProposal

                    }
                    if (args_func_name.toString() == "transferPoints") {
                        args_body = args_body.transferProposal
                    }

                    parse_payload[index].txType = args_func_name.toString()
                    parse_payload[index].txDetails = args_body

                    // console.log(parse_payload)
                }
                parse_payload = JSON.stringify(parse_payload)
            } else {
                parse_payload = response_payload.toString()
            }
            // get the final result

            let request3 = {
                proposalResponses: proposalResponses,
                proposal: proposal,
                admin: true
            };
            // listen event hub
            let promises = [];
            // let eventhub = channel.newChannelEventHub(config.orgYAML[invokeReq.userOrg]["peers"][0])

            // let deployId = tx_id.getTransactionID();
            // let txPromise = new Promise((resolve, reject) => {
            //     let handle = setTimeout(() => {
            //         logger.error('Timeout - Failed to receive the event for commit:  waiting on '+ eventhub.getPeerAddr());
            //         eventhub.disconnect(); // will not be using this event hub
            //         reject('TIMEOUT waiting on '+ eventhub.getPeerAddr());
            //     }, 30000);
            // 	logger.debug("Register event for transaction : >> " + deployId.toString());
            //     eventhub.registerTxEvent(deployId.toString(), (tx, code) => {
            //         clearTimeout(handle);
            //         eventhub.unregisterTxEvent(deployId);

            //         if (code !== 'VALID') {
            //             logger.info('transaction was invalid, code = ' + code);
            //             reject();
            //         } else {
            //             logger.info('transaction has been committed on peer ' + eventhub.getPeerAddr());
            //             resolve();
            //         }
            //     }, (error) => {
            //         clearTimeout(handle);
            //         logger.error('transaction event failed:' + error);
            //         reject(error);
            //     });
            //     eventhub.connect();
            // });

            // promises.push(txPromise);

            logger.debug("Sending endorsements to ordering service -----  endorsement Response >> " + request3.proposalResponses + "; endorsement >>" + request3.proposal);

            let sendPromise = channel.sendTransaction(request3);
            Promise.all([sendPromise].concat(promises))
                .then((results) => {
                    logger.info(" event promise all complete and testing complete");
                    let response = results[0]; // the first returned value is from the 'sendPromise' which is from the 'sendTransaction()' call

                    if (response.status === "SUCCESS") {
                        logger.info("Successfully sent transaction to the orderer.");

                        //close the connections
                        channel.close();
                        logger.debug("Successfully closed all connections");

                        response.txid = tx_id.getTransactionID();
                        response.payload = parse_payload
                        return cb("", response);
                    } else {
                        logger.error("Failed to order the transaction. Error code: " + response.status);
                        return cb(new Error("Failed to order the transaction. Error code: " + response.status), "");
                    }
                }).catch((err) => {
                    logger.error("Failed to send transaction and get notifications within the timeout period.", err);
                    return cb(err, "");
                });
        }
    } catch (err) {
        logger.error("Failed to invoke: " + err.stack ? err.stack : err);
        return cb(err, "");
    }
}

module.exports.invoke = invoke;

/*
 *  credit_service for  query
 * */
async function query(queryReq, cb) {
    let channel = null;

    try {
        //let client = Client.loadFromConfig('config/dev/network.yaml');
        let client = Client.loadFromConfig(config.networkYAML);
        logger.info('Successfully loaded a connection profile');
        //client.loadFromConfig('config/dev/'+queryReq.userOrg+'.yaml');
        client.loadFromConfig(config.orgYAML[queryReq.userOrg]);

        await client.initCredentialStores();
        logger.info('Successfully created the key value store and crypto store based on the sdk config and connection profile');

        let enrollment = null;
        if (global.enrollment)
            enrollment = global.enrollment;
        else {
            let caService = client.getCertificateAuthority();
            let request = {
                enrollmentID: 'admin',
                enrollmentSecret: 'adminpw',
                profile: 'tls'
            };

            enrollment = await caService.enroll(request);
            enrollment.key = enrollment.key.toBytes();
            logger.info('Successfully called the CertificateAuthority to get the TLS material');
            global.enrollment = enrollment;
        }
        client.setTlsClientCertAndKey(enrollment.certificate, enrollment.key);

        // let caService = client.getCertificateAuthority();
        // let request = {
        //     enrollmentID: 'admin',
        //     enrollmentSecret: 'adminpw',
        //     profile: 'tls'
        // };
        // let enrollment = await caService.enroll(request);
        // logger.info('Successfully called the CertificateAuthority to get the TLS material');
        // let key = enrollment.key.toBytes();
        // let cert = enrollment.certificate;
        // client.setTlsClientCertAndKey(cert, key);

        channel = client.getChannel(channelName);
        //await getUser(client, 'admin','adminpw');
        await getUser(client, queryReq.userName, queryReq.passwd);
        logger.info("Successfully enrolled user 'admin' (e2eUtil 4)");

        tx_id = client.newTransactionID();
        // send query
        let request2 = {
            chaincodeId: chainCodeName,
            txId: tx_id,
            fcn: queryReq.functionName,
            args: queryReq.args
        };
        let targets = [peerOrgID]
        if (targets && targets.length != 0) {
            const targetPeers = getTargetPeers(channel, targets);
            if (targetPeers.length < targets.length) {
                logger.error('Failed to get all peers for targets: ' + targets);
            } else {
                request2.targets = targetPeers;
            }
        }

        let response_payloads = await channel.queryByChaincode(request2);
        try {
            logger.debug(">>>>>>>>>>>>>>>>", response_payloads)
            if (response_payloads) {
                for (let i = 0; i < response_payloads.length; i++) {
                    logger.info("checking query results are ", response_payloads[i].toString("utf8"));
                    if (response_payloads[i]) {
                        var result_index = i
                    }
                    // This is a big bug needs to fix in the near future TODO:
                }
                if (response_payloads[result_index].toString("utf8").indexOf("Error") >= 0) {
                    return cb(response_payloads[result_index].toString("utf8"), "");
                } else {
                    if (!JSON.stringify(response_payloads[result_index])) {
                        return cb("", "Can not stringify response.");
                    } else {
                        logger.debug("The result is: ", response_payloads[result_index].toString("utf8"))
                        if (response_payloads[result_index]) {
                            let result = response_payloads[result_index].toString("utf8");
                            return cb("", result);
                        } else {
                            return cb("", "Emplty Result!")
                        }
                    }
                }
            } else {
                logger.error("response_payloads is null");
                return cb(new Error("Failed to get response on query, response_payloads is null", ""))
            }
        } catch (error) {
            logger.error("response_payloads is null");
            return cb(new Error("Failed to get response on query: " + error), "")
        }
    } catch (err) {
        logger.error("Failed to send query due to error: " + err.stack ? err.stack : err);
        return cb(err, "");
    }
}

module.exports.query = query;

// return an array of peer objects for targets which are a array of peer urls in string (e.g., localhost:7051)
function getTargetPeers(channel, targets) {
    // get all the peers and then find what peer matches a target
    let targetPeers = [];
    if (targets && targets.length != 0) {
        const peers = channel.getPeers();
        for (let i in targets) {
            let found = false;
            for (let j in peers) {
                logger.debug('channel has peer ' + peers[j].getName());
                if (targets[i] === peers[j].getName()) {
                    targetPeers.push(peers[j]);
                    found = true;
                    break;
                }
            }
            if (!found) {
                logger.error('Cannot find the target peer for ' + targets[i]);
            }
        }
    }
    return targetPeers;
}

function getUser(client, username, password, skipPersistence) {

    return client.getUserContext(username, true).then((user) => {
        return new Promise((resolve, reject) => {
            if (user && user.isEnrolled()) {
                logger.info('Successfully loaded member from persistence');
                return resolve(user);
            }
            return resolve(client.setUserContext({
                username: username,
                password: password
            }, skipPersistence));
        });
    });
}