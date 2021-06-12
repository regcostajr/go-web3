// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract ABIEncoding {

    uint[] public uint_dynamic_array;

    constructor(uint[] memory _uint_dynamic_array){
        uint_dynamic_array = _uint_dynamic_array;
    }

    event TestString(string _string, string[] _string_dynamic_array, string[2] _string_fixed_array);

    function testString(string memory _string, string[] memory _string_dynamic_array, string[2] memory _string_fixed_array) public {
        emit TestString(_string, _string_dynamic_array, _string_fixed_array);
    }

}
