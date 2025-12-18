var log4js = require('log4js');
var logger = log4js.getLogger('app');


let invokeOrgID = "Org1"
var hfc = require('fabric-client');
hfc.setLogger(logger);

async function getClientForOrg(userorg) {
    logger.debug('getClientForOrg - ****** START %s %s', userorg)
    let config = '-connection-profile-path';
    let client = hfc.loadFromConfig(hfc.getConfigSetting('network' + config));
    logger.debug(">>>>>>>>>>>>> Get client", client)
    client.loadFromConfig(hfc.getConfigSetting(userorg + config));
    await client.initCredentialStores();
    logger.debug('getClientForOrg - ****** END %s\n\n', userorg)

    return client;
}

var getRegisteredUser = async function (username, userpasswd) {
    var userOrg = invokeOrgID
    try {
        var client = await getClientForOrg(userOrg);
        logger.debug('Successfully initialized the credential stores');
        // client can now act as an agent for organization Org1
        // first check to see if the user is already enrolled
        var user = await client.getUserContext(username, true);
        if (user && user.isEnrolled()) {
            logger.info('Successfully loaded member from persistence');
        } else {
            // user was not enrolled, so we will need an admin user object to register
            logger.debug('User %s was not enrolled, so we will need an admin user object to register', username);
            let adminUserObj = await client.setUserContext({
                username: "admin",
                password: "adminpw"
            });
            let caClient = client.getCertificateAuthority();
            // add affiliations
            let affiliationService = caClient.newAffiliationService();
            let registeredAffiliations = await affiliationService.getAll(adminUserObj);
            if (!registeredAffiliations.result.affiliations.some(
                    x => x.name == userOrg.toLowerCase())) {
                let affiliation = userOrg.toLowerCase() + '.department1';
                logger.debug("create affiliation:", affiliation)
                try {
                    await affiliationService.create({
                        name: affiliation,
                        force: true
                    }, adminUserObj);
                } catch (error) {
                    logger.debug("Pass the error:", error)
                }
            }
            // added
            var secret = await caClient.register({
                enrollmentID: username,
                enrollmentSecret: userpasswd,
                affiliation: userOrg.toLowerCase() + '.department1'
            }, adminUserObj);
            logger.debug('Successfully register user %s', username);
        }
        if (secret) {
            var response = {
                success: true,
                secret: secret,
                message: username + ' register Successfully',
            };
            return response;
        } else {
            throw new Error("Secret not found")
        }
    } catch (error) {
        logger.error('Failed to get registered user: %s with error: %s', username, error.toString());
        return 'failed ' + error.toString();
    }

};
exports.getRegisteredUser = getRegisteredUser;