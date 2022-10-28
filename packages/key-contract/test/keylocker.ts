const hre = require("hardhat")
import { Contract }  from "ethers"
// @ts-ignore
import chai from "chai"
const { expect } = chai
const { ethers } = require("hardhat")


describe("Test KeyLocker", function() {
    let keyLocker: Contract

    before('deploy keylocker contracts', async () => {
        await deployKeyLocker()
    })

    it("test set and get function", async () => {
        const uuid = ethers.utils.solidityKeccak256(['string'], ["0x000000000"]);
        const keys = ["0x1000000000000000000000000000000000000000"]
        keyLocker.setSocialKey(uuid, keys)
        const uuidKey = await keyLocker.getSocialKey(uuid)
        expect(uuidKey[0]).to.eq("0x1000000000000000000000000000000000000000")
    })

    const deployKeyLocker = async () => {
        const factory = await hre.ethers.getContractFactory("KeyLocker");
        keyLocker = await factory.deploy();
        await keyLocker.initialize()
    }
});