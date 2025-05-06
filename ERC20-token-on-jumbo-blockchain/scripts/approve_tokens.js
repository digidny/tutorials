require("dotenv").config(); // Load environment variables

const hre = require("hardhat");

async function main() {
  const [owner, spender] = await hre.ethers.getSigners();

  // Load token address from environment variables
  const tokenAddress = process.env.TOKEN_ADDRESS;

  if (!tokenAddress) {
    console.error("Error: TOKEN_ADDRESS not found in .env file.");
    return;
  }

  //Approve 50 tokens
  const amountToApprove = hre.ethers.parseUnits("50", 18);

  const MyToken = await hre.ethers.getContractAt("MyToken", tokenAddress);

  console.log("Owner address (approving):", owner.address);
  console.log("Spender address (to be approved):", spender.address);
  console.log(
    "Amount to approve:",
    hre.ethers.formatUnits(amountToApprove),
    "JBT"
  );

  const approveTx = await MyToken.approve(spender.address, amountToApprove);
  await approveTx.wait();

  console.log("Approval successful!");
  console.log("Transaction hash:", approveTx.hash);

  const allowance = await MyToken.allowance(owner.address, spender.address);
  console.log(
    "Allowance granted to spender:",
    hre.ethers.formatUnits(allowance),
    "JBT"
  );
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
