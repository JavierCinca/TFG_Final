var express = require('express');
var bodyParser = require('body-parser');

var app = express();
app.use(bodyParser.json());

// Setting for Hyperledger Fabric
const { Wallets, FileSystemWallet, Gateway } = require('fabric-network');
const path = require('path');
const fs = require('fs');

const ccpPath = path.resolve(__dirname, '..',  'connection-be1.json');
const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));
ccp.peers['peer0.be1.cinca-network.com'].tlsCACerts.pem = fs.readFileSync(path.resolve(__dirname, '..', ccp.peers['peer0.be1.cinca-network.com'].tlsCACerts.path), 'utf8');
ccp.certificateAuthorities['ca.be1.cinca-network.com'].tlsCACerts.pem = fs.readFileSync(path.resolve(__dirname, '..', ccp.certificateAuthorities['ca.be1.cinca-network.com'].tlsCACerts.path), 'utf8');

app.get('/api/forex', async function (req, res) {
    try {
        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.get('appUser');
        if (!userExists) {
            console.log('An identity for the user "appUser" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
	await gateway.connect(ccp, { wallet, identity: 'appUser', discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('channeldemo');

        // Get the contract from the network.
        const contract = network.getContract('becc');

        // Evaluate the specified transaction.
        const result = await contract.evaluateTransaction('queryAllForexs');
        console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
        res.status(200).json({response: result.toString()});

    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        res.status(500).json({error: error});
        process.exit(1);
    }
});

app.post('/api/addforex/', async function (req, res) {
    console.time('loop');
    try {

        // Create a new file system based wallet for managing identities.
        const walletPath = path.resolve(__dirname, '.', 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.get('appUser');
        if (!userExists) {
            console.log('An identity for the user "appUser" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'appUser', discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('channeldemo');

        // Get the contract from the network.
        const contract = network.getContract('becc');

        // Submit the specified transaction.
        await contract.submitTransaction('createForex', req.body.forexid, req.body.fecha, req.body.base, req.body.oro, req.body.euros, req.body.libra, req.body.bitcoin, req.body.ethereum);
        console.log('Transaction has been submitted');
        res.send('Transaction has been submitted');

        // Disconnect from the gateway.
        await gateway.disconnect();

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
    console.timeEnd('loop');
})



app.listen(9090, 'localhost');
console.log('Running on http://localhost:9090/');