const hre = require("hardhat");

async function main() {
    const keyLocker = await hre.ethers.getContractFactory("KeyLocker");
    const locker = await keyLocker.deploy();
    console.log("keyLocker deployed to:", locker.address);
}

main().then(() => process.exit(0)).catch(error => {
    console.error(error);
    process.exit(1);
});