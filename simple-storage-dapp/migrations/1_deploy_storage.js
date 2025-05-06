const SimpleStorage = artifacts.require("./SimpleStorage.sol");

module.exports = function (deployer) {
  deployer.deploy(SimpleStorage, 5); // Initial value of 5
};

