/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { FileSystemWallet, Gateway } = require('fabric-network');  
const path = require('path');
const fs = require('fs');
const { Channel } = require('grpc');


const ccpPath = path.resolve(__dirname, '..', 'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

async function main() {
    try {

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists('User1');
        if (!userExists) {
            console.log('does not exist in the wallet');
            console.log(userExists);
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();

        let connectionOptions = {
            identity: 'User1',
            wallet: wallet,
            discovery: { enabled: true, asLocalhost: false }
        }

        // Create a new gateway for connecting to our peer node.
        await gateway.connect(ccp, connectionOptions);
        
        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('general');
        const channel = await network.getChannel('general');

        // Get the contract from the network.
        const contract = network.getContract('asset_management');

        // Submit the specified transaction.
        const arg = process.argv.slice(2);
        await contract.submitTransaction('withdraw', arg[0], arg[1]);
        console.log('Transaction has been submitted');

        // Disconnect from the gateway.
        await gateway.disconnect();

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}

main();
