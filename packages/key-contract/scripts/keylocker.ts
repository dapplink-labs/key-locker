const hre = require("hardhat");

async function main() {
    const factory = await hre.ethers.getContractFactory("KeyLocker");
    const keyLocker = await factory.deploy();
    await keyLocker.initialize()
}

main().then(() => process.exit(0)).catch(error => {
    console.error(error);
    process.exit(1);
});