let path = require("path");
let config = {};
let fs = require("fs");

function init() {
  let configPath;
  if (process.env.NODE_ENV == "dev")
    configPath = path.join(__dirname, "../config/dev");
  if (process.env.NODE_ENV == "uat")
    configPath = path.join(__dirname, "../config/uat");
  if (process.env.NODE_ENV == "production")
    configPath = path.join(__dirname, "../config/production");
  loadConfig(configPath);
  return config;
}

function capitalizeTxt(txt) {
  return txt.charAt(0).toLowerCase() + txt.slice(1);
}

function loadYamlOrg(configPath) {
  fs.readdirSync(configPath).forEach(file => {
    if (file.toUpperCase().indexOf("ORG") == 0) {
      orgName = file.split(".yaml")[0];
      orgYaml = file;
      orgName = capitalizeTxt(orgName);
      config.orgYAML[orgName] = path.join(configPath, orgYaml);
    }
  });
}

function loadConfig(configPath) {
  if (JSON.stringify(config) == "{}") {
    config.chaincode = path.join(configPath, "chaincode.json");
    config.couchdb = path.join(configPath, "couchdb.json");
    config.log4js = path.join(configPath, "log4js.json");
    config.network = path.join(configPath, "network.json");
    config.server = path.join(configPath, "server.json");

    config.networkYAML = path.join(configPath, "network.yaml");
    config.orgYAML = [];
    loadYamlOrg(configPath);
  }
}

function getConfig() {
  if (JSON.stringify(config) == "{}") {
    return init();
  }
  return config;
}

module.exports.init = init;
module.exports.getConfig = getConfig;
