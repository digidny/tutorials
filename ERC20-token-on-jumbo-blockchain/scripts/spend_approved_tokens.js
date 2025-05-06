require("dotenv").config();

const hre = require("hardhat");

async function main() {
  // We need a different signer here - the one we approved (the spender)
  const [owner, spender, recipient] = await hre.ethers.getSigners();

  const tokenAddress = process.env.TOKEN_ADDRESS;
  const amountToSpend = hre.ethers.parseUnits(
    "25", // lets spend 25 out of allowed tokens
  );

  if (!tokenAddress) {
    console.error("Error: TOKEN_ADDRESS not found in .env file.");
    return;
  }

  const MyToken = await hre.ethers.getContractAt("MyToken", tokenAddress);

  console.log("Spender address (spending):", spender.address);
  console.log("Recipient address:", recipient.address);
  console.log(
    "Amount to spend:",
    hre.ethers.formatUnits(amountToSpend),
    "JBT"
  );

  // The spender needs to call the transferFrom function
  const spendTx = await MyToken.connect(spender).transferFrom(
    owner.address,     // The address whose tokens are being spent
    recipient.address, // The address to receive the tokens
    amountToSpend      // The amount of tokens to transfer
  );
  await spendTx.wait();

  console.log("Spend successful!");
  console.log("Transaction hash:", spendTx.hash);

  const spenderBalance = await MyToken.balanceOf(spender.address);
  const recipientBalance = await MyToken.balanceOf(recipient.address);
  const ownerBalance = await MyToken.balanceOf(owner.address);

  console.log(
    "Spender's new balance:",
    hre.ethers.formatUnits(spenderBalance),
    "JBT"
  );
  console.log(
    "Recipient's new balance:",
    hre.ethers.formatUnits(recipientBalance),
    "JBT"
  );
  console.log(
    "Owner's new balance:",
    hre.ethers.formatUnits(ownerBalance),
    "JBT"
  );
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
