/*

    Copyright 2020 dYdX Trading Inc.

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.

*/

pragma solidity 0.5.16;
pragma experimental ABIEncoderV2;


/**
 * @title I_MakerOracle
 * @author dYdX
 *
 * Interface for the MakerDAO Oracles V2 smart contrats.
 */
interface I_MakerOracle {

    // ============ Getter Functions ============

    /**
     * @notice Returns the current value as a bytes32.
     */
    function peek()
        external
        view
        returns (bytes32, bool);

    /**
     * @notice Requires a fresh price and then returns the current value.
     */
    function read()
        external
        view
        returns (bytes32);

    /**
     * @notice Returns the number of signers per poke.
     */
    function bar()
        external
        view
        returns (uint256);

    /**
     * @notice Returns the timetamp of the last update.
     */
    function age()
        external
        view
        returns (uint32);

    /**
     * @notice Returns 1 if the signer is authorized, and 0 otherwise.
     */
    function orcl(
        address signer
    )
        external
        view
        returns (uint256);

    /**
     * @notice Returns 1 if the address is authorized to read the oracle price, and 0 otherwise.
     */
    function bud(
        address reader
    )
        external
        view
        returns (uint256);

    /**
     * @notice A mapping from the first byte of an authorized signer's address to the signer.
     */
    function slot(
        uint8 signerId
    )
        external
        view
        returns (address);

    // ============ State-Changing Functions ============

    /**
     * @notice Updates the value of the oracle
     */
    function poke(
        uint256[] calldata val_,
        uint256[] calldata age_,
        uint8[] calldata v,
        bytes32[] calldata r,
        bytes32[] calldata s
    )
        external;

    /**
     * @notice Authorize an address to read the oracle price.
     */
    function kiss(
        address reader
    )
        external;

    /**
     * @notice Unauthorize an address so it can no longer read the oracle price.
     */
    function diss(
        address reader
    )
        external;

    /**
     * @notice Authorize addresses to read the oracle price.
     */
    function kiss(
        address[] calldata readers
    )
        external;

    /**
     * @notice Unauthorize addresses so they can no longer read the oracle price.
     */
    function diss(
        address[] calldata readers
    )
        external;
}
