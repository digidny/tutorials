require('dotenv').config({ path: './backend/.env' });
const Web3 = require('web3');
const provider = () => new Web3.providers.HttpProvider(process.env.RPC_URL);

module.exports = {
  networks: {
    jumbo: {
      host: "127.0.0.1",     // Localhost (default: none)
      port: process.env.RPC_PORT,             // Custom port
      provider: provider,
      network_id: 112233,       // Custom network
      gas: 8500000,           // Gas sent with each transaction (default: ~6700000)
      gasPrice: 1000000000,  // 20 gwei (in wei) (default: 100 gwei)
      from: process.env.DEPLOYER_ADDRESS,        // Account to send transactions from (default: accounts[0])
      websocket: true,         // Enable EventEmitter interface for web3 (default: false)
    }
  },

  // Set default mocha options here, use special reporters, etc.
  mocha: {
    // timeout: 100000
  },

  // Configure your compilers
  compilers: {
    solc: {
      version: "0.8.0",      // Fetch exact version from solc-bin (default: truffle's version)
      // docker: true,        // Use "0.5.1" you've installed locally with docker (default: false)
      // settings: {          // See the solidity docs for advice about optimization and evmVersion
      //  optimizer: {
      //    enabled: false,
      //    runs: 200
      //  },
      //  evmVersion: "byzantium"
      // }
    }
  },
};
