## API server for business rule service


### config

Before bring up the server , you should make neccessary configurations.Under **config** directory, there are three 
sub-directories **dev**,**uat** and **production** which are used in Development, User Acceptance Test and Production 
environment separately. All of them will include following config file:
  
> 1. chaincode.json   *config file for the chaincode this api server will call*
> 2. couchdb.json     *config file for couchdb key store*
> 3. log4js.json      *log config*
> 4. network.json     *config for the fabric network this server will use*
> 5. server.json      *common config for a web server*


### install npm

npm install

### start server

dev env:
> pm2 start ecosystem.config.js

uat env:
> pm2 start ecosystem.config.js --env uat

production env:
> pm2 start ecosystem.config.js --env production