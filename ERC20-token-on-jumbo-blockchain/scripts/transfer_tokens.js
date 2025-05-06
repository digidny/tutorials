require("dotenv").config(); // Load environment variables
const hre = require("hardhat");

async function main() {
  const [owner, recipient] = await hre.ethers.getSigners();

  // Load token address from environment variables
  const tokenAddress = process.env.TOKEN_ADDRESS;

  if (!tokenAddress) {
    console.error("Error: TOKEN_ADDRESS not found in .env file.");
    return;
  }

  const amountToSend = hre.ethers.parseUnits("100", 18); // Sending 100 tokens (assuming 18 decimals)

  const MyToken = await hre.ethers.getContractAt("MyToken", tokenAddress);

  console.log("Initiating transfer from:", owner.address);
  console.log("Recipient address:", recipient.address);
  console.log("Amount to send:", hre.ethers.formatUnits(amountToSend, 18), "JBT");

  const transferTx = await MyToken.transfer(recipient.address, amountToSend);
  await transferTx.wait();

  console.log("Transfer successful!");
  console.log("Transaction hash:", transferTx.hash);

  const ownerBalance = await MyToken.balanceOf(owner.address);
  const recipientBalance = await MyToken.balanceOf(recipient.address);

  console.log("Owner's new balance:", hre.ethers.formatUnits(ownerBalance, 18), "JBT");
  console.log("Recipient's new balance:", hre.ethers.formatUnits(recipientBalance, 18), "JBT");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
