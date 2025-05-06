require("@nomicfoundation/hardhat-toolbox");

const owner = "2a3569dbc2f6afb8ad94eb65ac23d8530538b3ed153d8dfae26e163c80fcbbad";
const spender = "55782da8d496f78c228cb15afcaa6f222f6d17c231d72bb12ac97153394e7074";
const recipient = "4a0665930111bf18c0c6618db483ec7d58ee4fcb9caef790f2177d37b06b7e9f"

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: "0.8.24",

  networks: {
    jumbochain: {
      url: `https://testnode.jumbochain.org`, // RPC URL which you can get from
      //the chainlist, as mentioned in the documentation.
      accounts: [`${owner}`, `${spender}`, `${recipient}`],
    },
  },
};
